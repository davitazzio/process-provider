package process

import (
	"fmt"
	"os"
	"strconv"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/provider-processprovider/apis/process/v1alpha1"
	"golang.org/x/crypto/ssh"
)

func ObserveProcess(cr_specs v1alpha1.ProcessParameters, logger logging.Logger) (int64, error) {
	logger.Debug(fmt.Sprintf("Osservazione processo su %s", cr_specs.NodeAddress))
	client, session, err := connectToHost(cr_specs, logger)
	if err != nil {
		logger.Debug("errore CONNESSIONE observing")
		logger.Debug(err.Error())
		return -1, err
	}
	defer client.Close()

	output, err := session.CombinedOutput("pgrep -f mqttsub.py")
	if err != nil {
		logger.Debug("errore grep")
		logger.Debug(err.Error())
		return -1, err
	}

	process_pid, _ := strconv.ParseInt(string(output), 10, 64)
	return process_pid, nil

	// 	_, err = session.CombinedOutput(fmt.Sprintf("kill %d", process_pid))
	// 	if err.Error() == "Process exited with status 143 from signal TERM" {
	// 		client.Close()

	//			return nil
	//		} else {
	//			logger.Debug("errore kill")
	//			logger.Debug(err.Error())
	//			client.Close()
	//			return err
	//		}
	//	}
}

func CreateProcess(cr_specs v1alpha1.ProcessParameters, logger logging.Logger) error {

	client, session, err := connectToHost(cr_specs, logger)
	if err != nil {
		logger.Debug("errore CONNESSIONE creation")
		logger.Debug(err.Error())
		return err
	}
	defer client.Close()
	//downloading code from gitlab on the remote machine running script.sh on the remote machine
	file, err := os.Open("script.sh")
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}
	defer file.Close()
	//run remotely test.sh
	session.Stdin = file
	output, err := session.CombinedOutput("bash")
	if err != nil {
		fmt.Println("Error executing command: ", err)
	}
	logger.Debug(string(output))

	return nil
}
func KillProcess(cr_specs v1alpha1.ProcessParameters, logger logging.Logger) (bool, error) {

	client, session, err := connectToHost(cr_specs, logger)
	if err != nil {
		logger.Debug("errore CONNESSIONE")
		logger.Debug(err.Error())
		return false, err
	}

	output, err := session.CombinedOutput("pgrep -f mqttsub.py")
	if err != nil {
		logger.Debug("errore grep")
		logger.Debug(err.Error())
		return false, err
	}

	process_pid, _ := strconv.ParseInt(string(output), 10, 64)
	client.Close()
	client, session, err = connectToHost(cr_specs, logger)
	if err != nil {
		logger.Debug("errore CONNESSIONE")
		logger.Debug(err.Error())
		return false, err
	}
	defer client.Close()
	output, err = session.CombinedOutput(fmt.Sprintf("kill %d", process_pid))
	logger.Debug(string(output))
	if err != nil {
		logger.Debug("errore kill")
		logger.Debug(err.Error())
		return false, err
	}
	return true, nil

}

func connectToHost(cr_specs v1alpha1.ProcessParameters, logger logging.Logger) (*ssh.Client, *ssh.Session, error) {
	logger.Debug(fmt.Sprintf("Connessione a %s, username %s, password %s", cr_specs.NodeAddress, cr_specs.Username, cr_specs.Password))
	sshConfig := &ssh.ClientConfig{User: cr_specs.Username, Auth: []ssh.AuthMethod{ssh.Password(cr_specs.Password)}, HostKeyCallback: ssh.InsecureIgnoreHostKey()}

	client, err := ssh.Dial("tcp", cr_specs.NodeAddress+":22", sshConfig)
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
