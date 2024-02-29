/*
Copyright 2022 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package process

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane/provider-processprovider/apis/process/v1alpha1"
	apisv1alpha1 "github.com/crossplane/provider-processprovider/apis/v1alpha1"
	processservice "github.com/crossplane/provider-processprovider/internal/controller/processService"
	"github.com/crossplane/provider-processprovider/internal/features"
)

const (
	errNotProcess   = "managed resource is not a Process custom resource"
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errGetCreds     = "cannot get credentials"
	errNewClient    = "cannot create new Service"
)

var (
	newProcessService = func(_ []byte, processName string) (*processservice.ProcessService, error) {
		return processservice.GetInstance(processName), nil
	}
)

// Setup adds a controller that reconciles Process managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.ProcessGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.ProcessGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			newServiceFn: newProcessService,
			logger:       o.Logger}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithPollInterval(o.PollInterval),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		WithEventFilter(resource.DesiredStateChanged()).
		For(&v1alpha1.Process{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	logger       logging.Logger
	kube         client.Client
	usage        resource.Tracker
	newServiceFn func(creds []byte, processName string) (*processservice.ProcessService, error)
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Process)
	if !ok {
		return nil, errors.New(errNotProcess)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	cd := pc.Spec.Credentials
	data, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.kube, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	svc, err := c.newServiceFn(data, mg.GetName())
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	return &external{service: svc, logger: c.logger}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	// A 'client' used to connect to the external resource API. In practice this
	// would be something like an AWS SDK client.
	logger  logging.Logger
	service *processservice.ProcessService
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Process)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotProcess)
	}
	cr.Status.SetConditions(v1.ReconcileSuccess())
	cr.SetConditions(v1.Available())

	// These fmt statements should be removed in the real implementation.
	// fmt.Printf("Observing: %+v", cr)
	// c.logger.Debug("SONO DENTRO A OBSERVE")
	// c.logger.Debug(strconv.FormatBool(c.service.GetExecuted()))
	// c.logger.Debug(strconv.FormatBool(c.service.Executed))

	return managed.ExternalObservation{
		// Return false when the external resource does not exist. This lets
		// the managed resource reconciler know that it needs to call Create to
		// (re)create the resource, or that it has successfully been deleted.
		ResourceExists: c.service.GetExecuted(),

		// Return false when the external resource exists, but it not up to date
		// with the desired managed resource state. This lets the managed
		// resource reconciler know that it needs to call Update.
		ResourceUpToDate: true,

		// Return any details that may be required to connect to the external
		// resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Process)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotProcess)
	}

	c.logger.Debug(fmt.Sprintf("Creating resource %s with parameter 'node_address' %s.", cr.Name, cr.Spec.ForProvider.NodeAddress))

	// an example of resource creatino connecting to an HTTP server
	sendHTTPReq(cr.Spec.ForProvider.NodeAddress, cr.Spec.ForProvider.NodePort, cr.Spec.ForProvider.Service, c.logger)

	// newCondition := true
	mg.SetConditions(v1.Available())

	meta.SetExternalCreateSucceeded(mg, time.Now())
	return managed.ExternalCreation{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Process)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotProcess)
	}

	c.logger.Debug(fmt.Sprintf("Updating: %+v", cr))

	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Process)
	if !ok {
		return errors.New(errNotProcess)
	}
	processservice.DeleteProcessService(mg.GetName())

	c.logger.Debug(fmt.Sprintf("Deleting resource: %+v", cr))

	return nil
}

func sendHTTPReq(nodeAddress string, nodePort string, service *string, logger logging.Logger) {

	address := "http://" + nodeAddress + ":" + nodePort + "/" + *service

	logger.Debug(fmt.Sprintf("Sending HTTP request to node %s", address))

	resp, err := http.Get(address)
	if err != nil {
		logger.Debug("Error decoding request.")
		logger.Debug(err.Error())
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Debug(err.Error())
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Debug("Error reading response.")
		logger.Debug(err.Error())
	}
	logger.Debug(string(body))
}
