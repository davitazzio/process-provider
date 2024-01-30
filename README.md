# Process-provider: modifica al template

Il manifest [examples/dummypod/ngix-server-depl.yaml] dichiara il pod "ricevente", crea un servizio e lo associa a una porta 
e infine crea un tipo `Process` che comunica con il pod appena creato. 

Nel manifest [examples/provider/config.yaml] viene definito il `ProviderConfig` che serve per associare i tipi `Process` al `process-provider`

Nella cartella [packages/crds] sono definiti i tipi come `CustomResourceDefinition`

Il file [packages/crossplane.yaml] contiene la configurazione per il provider `process-provider` necessario per il `make build`.

Il file [package/manifest-kubernetes.yaml] contiene la configurazione per il provider `process-provider` necessario il deployment su kubernetes.


Refer to Crossplane's [CONTRIBUTING.md] file for more information on how the
Crossplane community prefers to work. The [Provider Development][provider-dev]
guide may also be of use.

[examples/dummypod/ngix-server-depl.yaml]: https://github.com/davitazzio/process-provider/blob/main/examples/dummypod/ngix-server-depl.yaml
[examples/provider/config.yaml]: https://github.com/davitazzio/process-provider/blob/main/examples/provider/config.yaml
[packages/crds]: https://github.com/davitazzio/process-provider/tree/main/package/crds
[packages/crossplane.yaml]: https://github.com/davitazzio/process-provider/blob/main/package/crossplane.yaml
[package/manifest-kubernetes.yaml]: https://github.com/davitazzio/process-provider/blob/main/package/manifest-kubernetes.yaml
[CONTRIBUTING.md]: https://github.com/crossplane/crossplane/blob/master/CONTRIBUTING.md
[provider-dev]: https://github.com/crossplane/crossplane/blob/master/contributing/guide-provider-development.md

