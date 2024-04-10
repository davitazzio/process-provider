package processservice

import (
	"fmt"
	"strconv"

	"golang.org/x/crypto/ssh"

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

	// app := "scp"
	// arg0 := programPath                                                        //"./codice_python/start./bin/sh"
	// bash_command := fmt.Sprintf("scp %s %s@%s:/home/%s", programPath, remoteUser, nodeAddress, remoteUser) //"datavix@dtazzioli-kubernetes.cloudmmwunibo.it:/home/datavix/"
	// cmd := exec.Command("/bin/sh", "-c", bash_command)
	// logger.Debug(bash_command)
	// _, err := cmd.Output()

	// if err != nil {
	// 	logger.Debug("errore nell'invio del comando di copia")
	// 	logger.Debug(err.Error())
	// 	p.Active = false
	// 	return -1, err

	// }

	nodeAddress = nodeAddress + ":22"

	client, session, err := connectToHost(remoteUser, nodeAddress)
	if err != nil {
		logger.Debug("errore CONNESSIONE")
		logger.Debug(err.Error())
		p.Active = false
		return -1, err
	}

	_, err = session.CombinedOutput(fmt.Sprintf("scp -r dtazzioli@192.168.17.107:/home/dtazzioli/codice_python /home/%s/", remoteUser))
	if err != nil {
		logger.Debug("errore COPIA")
		logger.Debug(err.Error())
		p.Active = false
		return -1, err
	}

	// app = "ssh"
	// arg0 = fmt.Sprintf("%s@%s", remoteUser, nodeAddress)
	// bash_command = fmt.Sprintf("ssh %s@%s 'python3 app/start.py &'", remoteUser, nodeAddress)
	// cmd = exec.Command("/bin/sh", "-c", bash_command)
	// output, err := cmd.Output()
	// logger.Debug(string(output))
	// if err != nil {
	// 	logger.Debug("errore nell'avvio dell'applicazione")
	// 	logger.Debug(err.Error())
	// 	p.Active = false
	// 	return -1, err
	// }
	client.Close()

	client, session, err = connectToHost(remoteUser, nodeAddress)
	if err != nil {
		logger.Debug("errore CONNESSIONE")
		logger.Debug(err.Error())
		p.Active = false
		return -1, err
	}

	output, err := session.CombinedOutput("python3 app/start.py &")
	if err != nil {
		logger.Debug("errore nell'avvio dell'applicazione")
		logger.Debug(err.Error())
		p.Active = false
		return -1, err
	}

	client.Close()

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
		return false, nil
	}
	nodeAddress = nodeAddress + ":22"
	// app := "ssh"
	// bash_command := fmt.Sprintf("ssh %s@%s ps -A | grep %d", remoteUser, nodeAddress, p.processPID) //"datavix@dtazzioli-kubernetes.cloudmmwunibo.it"
	// // arg1 := fmt.Sprintf("ps -A | grep %d", p.processPID)
	// cmd := exec.Command("/bin/sh", "-c", bash_command)

	// // fmt.Println(cmd.Err)
	// _, err := cmd.Output()

	// if err != nil {
	// 	fmt.Print("errore nell'avvio dell'applicazione")
	// 	fmt.Print(err)
	// 	return false, err
	// }

	client, session, err := connectToHost(remoteUser, nodeAddress)
	if err != nil {
		logger.Debug("errore nell'CONNESSIONE OBSERVE")
		logger.Debug(err.Error())
		p.Active = false
		return false, err
	}
	_, err = session.CombinedOutput(fmt.Sprintf("ps -A | grep %d", p.processPID))
	if err != nil {
		logger.Debug("errore PS GREP ")
		logger.Debug(err.Error())
		p.Active = false
		return false, err
	}

	client.Close()

	p.Active = true
	return true, nil
}

func (p *ProcessService) TerminateProcess(nodeAddress string, nodePort string, remoteUser string, logger logging.Logger) error {
	// app := "ssh"
	// bash_command := fmt.Sprintf("ssh %s@%s kill %d", remoteUser, nodeAddress, p.processPID) //"datavix@dtazzioli-kubernetes.cloudmmwunibo.it"
	// // arg1 := fmt.Sprintf("kill %d", p.processPID)
	// cmd := exec.Command("/bin/sh", "-c", bash_command)

	// // fmt.Println(cmd.Err)
	// _, err := cmd.Output()

	// if err != nil {
	// 	fmt.Print("errore nell'avvio dell'applicazione")
	// 	fmt.Print(err)
	// 	return err
	// }
	nodeAddress = nodeAddress + ":22"

	client, session, err := connectToHost(remoteUser, nodeAddress)
	if err != nil {
		logger.Debug("errore CONNESSIONE TERMINA")
		logger.Debug(err.Error())
		p.Active = true
		return err
	}
	_, err = session.CombinedOutput(fmt.Sprintf("kill %d", p.processPID))
	if err.Error() == "Process exited with status 143 from signal TERM" {
		client.Close()
		p.Active = false

		return nil
	} else {
		logger.Debug("errore kill")
		logger.Debug(err.Error())
		p.Active = true
		client.Close()
		return err
	}

}

func connectToHost(user, host string) (*ssh.Client, *ssh.Session, error) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password("dtazzioli")},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}
