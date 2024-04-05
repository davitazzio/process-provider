package processservice

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
)

// A ProcessService does nothing.
type ProcessService struct {
	ProcessName string
	processPID  int
	Active      bool
}

type ProcessServiceList struct {
	items []*ProcessService
}

func (pl *ProcessServiceList) GetItems() []*ProcessService {
	return pl.items
}

func (pl *ProcessServiceList) AddService(service *ProcessService) {
	pl.items = append(pl.items, service)

}

var instances *ProcessServiceList

func GetInstance(processName string) *ProcessService {

	// Pattern singleton for each "process" created
	if instances == nil {
		instances = &ProcessServiceList{}
	}
	for _, process := range instances.GetItems() {
		if process.GetName() == processName {
			return process
		}
	}

	newProcess := &ProcessService{ProcessName: processName, processPID: 0, Active: false}
	instances.AddService(newProcess)
	return newProcess
}

func (p *ProcessService) StartProcess(nodeAddress string, nodePort string, programPath string, remoteUser string, logger logging.Logger) (int, error) {

	app := "scp"
	arg0 := programPath                                                        //"./codice_python/start.sh"
	arg1 := fmt.Sprintf("%s@%s:/home/%s", remoteUser, nodeAddress, remoteUser) //"datavix@dtazzioli-kubernetes.cloudmmwunibo.it:/home/datavix/"
	cmd := exec.Command(app, arg0, arg1)
	_, err := cmd.Output()

	if err != nil {
		logger.Debug("errore nell'invio del comando di copia")
		p.Active = false
		return -1, err
	}

	app = "ssh"
	arg0 = fmt.Sprintf("%s@%s", remoteUser, nodeAddress)
	arg1 = "'python3 app/start.py &'"
	cmd = exec.Command(app, arg0, arg1)
	output, err := cmd.Output()
	logger.Debug(string(output))
	if err != nil {
		logger.Debug("errore nell'avvio dell'applicazione")
		p.Active = false
		return -1, err
	}
	process_pid, _ := strconv.Atoi(string(output))

	p.processPID = process_pid
	p.Active = true

	return process_pid, nil

}

func (p *ProcessService) UpdateActive(condition bool) {
	p.Active = condition
}
func (p *ProcessService) GetActive() bool {
	return p.Active
}

func (p *ProcessService) GetName() string {
	return p.ProcessName
}

func (p *ProcessService) SetName(name string) {
	p.ProcessName = name
}

func DeleteProcessService(processName string) {
	if instances == nil {
		return
	}
	for i, process := range instances.GetItems() {
		if process.GetName() == processName {
			remove(instances.items, i)
		}
	}
}
func remove(s []*ProcessService, i int) []*ProcessService {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (p *ProcessService) ObserveProcess(nodeAddress string, nodePort string, remoteUser string, logger logging.Logger) (bool, error) {
	if !p.Active {
		return false, fmt.Errorf("processo non attivo")
	}
	app := "ssh"
	arg0 := fmt.Sprintf("%s@%s", remoteUser, nodeAddress) //"datavix@dtazzioli-kubernetes.cloudmmwunibo.it"
	arg1 := fmt.Sprintf("ps -A | grep %d", p.processPID)
	cmd := exec.Command(app, arg0, arg1)

	// fmt.Println(cmd.Err)
	_, err := cmd.Output()

	if err != nil {
		fmt.Print("errore nell'avvio dell'applicazione")
		fmt.Print(err)
		return false, err
	}

	p.Active = true
	return true, nil
}

func (p *ProcessService) TerminateProcess(nodeAddress string, nodePort string, remoteUser string, logger logging.Logger) error {
	app := "ssh"
	arg0 := fmt.Sprintf("%s@%s", remoteUser, nodeAddress) //"datavix@dtazzioli-kubernetes.cloudmmwunibo.it"
	arg1 := fmt.Sprintf("kill %d", p.processPID)
	cmd := exec.Command(app, arg0, arg1)

	// fmt.Println(cmd.Err)
	_, err := cmd.Output()

	if err != nil {
		fmt.Print("errore nell'avvio dell'applicazione")
		fmt.Print(err)
		return err
	}

	return nil
}
