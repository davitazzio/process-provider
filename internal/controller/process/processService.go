package process

import (
	"fmt"
	"strconv"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"golang.org/x/crypto/ssh"
)

func ObserveProcess(nodeAddress string, logger logging.Logger) (int64, error) {
	client, session, err := connectToHost("dtazzioli", nodeAddress)
	if err != nil {
		logger.Debug("errore CONNESSIONE")
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

func CreateProcess(nodeAddress string, logger logging.Logger) error {

	client, session, err := connectToHost("dtazzioli", nodeAddress)
	if err != nil {
		logger.Debug("errore CONNESSIONE")
		logger.Debug(err.Error())
		return err
	}

	command_str := "scp -r dtazzioli@dtazzioli-processprovider.cloudmmwunibo.it:/home/dtazzioli/process-provider/script.sh /home/dtazzioli"

	// exec.Command("scp", "./script.sh", fmt.Sprintf("dtazzioli@%s:/home/dtazzioli", nodeAddress))
	output, err := session.CombinedOutput(command_str)
	logger.Debug(string(output))
	if err != nil {
		logger.Debug("errore copia")
		logger.Debug(err.Error())
		client.Close()
		return err
	}
	client.Close()
	client, session, err = connectToHost("dtazzioli", nodeAddress)
	if err != nil {
		logger.Debug("errore CONNESSIONE")
		logger.Debug(err.Error())
		return err
	}

	output, err = session.CombinedOutput("bash /home/dtazzioli/script.sh")
	logger.Debug(string(output))
	if err != nil {
		logger.Debug("errore script")
		logger.Debug(err.Error())
		client.Close()
		return err
	}
	client.Close()
	return nil

}
func KillProcess(nodeAddress string, logger logging.Logger) (bool, error) {

	client, session, err := connectToHost("dtazzioli", nodeAddress)
	if err != nil {
		logger.Debug("errore CONNESSIONE")
		logger.Debug(err.Error())
		return false, err
	}
	defer client.Close()

	output, err := session.CombinedOutput("pgrep -f mqttsub.py")
	if err != nil {
		logger.Debug("errore grep")
		logger.Debug(err.Error())
		return false, err
	}

	process_pid, _ := strconv.ParseInt(string(output), 10, 64)

	client, session, err = connectToHost("dtazzioli", nodeAddress)
	if err != nil {
		logger.Debug("errore CONNESSIONE")
		logger.Debug(err.Error())
		return false, err
	}

	output, err = session.CombinedOutput(fmt.Sprintf("kill %d", process_pid))
	logger.Debug(string(output))
	defer client.Close()
	if err != nil {
		logger.Debug("errore kill")
		logger.Debug(err.Error())
		return false, err
	}
	return true, nil

}

func connectToHost(user, host string) (*ssh.Client, *ssh.Session, error) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password("dtazzioli")},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host+":22", sshConfig)
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
