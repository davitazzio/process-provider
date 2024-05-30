package processservice

// import (
// 	"strconv"

// 	"github.com/crossplane/crossplane-runtime/pkg/logging"
// 	"golang.org/x/crypto/ssh"
// )

// type FLCConfig struct {
// 	ClientID  string `json:"client_id"`
// 	Topic     string `json:"topic"`
// 	NumEpochs int64  `json:"num_epochs"`
// 	Batch     int64  `json:"batch"`
// 	IBUS      int64  `json:"ibus"`
// }

// type FLCMetrics struct {
// 	Loss          int64 `json:"loss"`
// 	AverageLoss5  int64 `json:"average_loss_5"`
// 	AverageLoss10 int64 `json:"average_loss_10"`
// }

// func ObserveProcess(nodeAddress string, logger logging.Logger) (int, error) {
// 	client, session, err := connectToHost("dtazzioli", nodeAddress)
// 	if err != nil {
// 		logger.Debug("errore CONNESSIONE")
// 		logger.Debug(err.Error())
// 		return -1, err
// 	}
// 	defer client.Close()

// 	output, err := session.CombinedOutput("pgrep -f mqttsub.py")
// 	if err != nil {
// 		logger.Debug("errore grep")
// 		logger.Debug(err.Error())
// 		return -1, err
// 	}

// 	process_pid, _ := strconv.Atoi(string(output))
// 	return process_pid, nil

// 	// 	_, err = session.CombinedOutput(fmt.Sprintf("kill %d", process_pid))
// 	// 	if err.Error() == "Process exited with status 143 from signal TERM" {
// 	// 		client.Close()

// 	//			return nil
// 	//		} else {
// 	//			logger.Debug("errore kill")
// 	//			logger.Debug(err.Error())
// 	//			client.Close()
// 	//			return err
// 	//		}
// 	//	}
// }

// func connectToHost(user, host string) (*ssh.Client, *ssh.Session, error) {

// 	sshConfig := &ssh.ClientConfig{
// 		User: user,
// 		Auth: []ssh.AuthMethod{ssh.Password("dtazzioli")},
// 	}
// 	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

// 	client, err := ssh.Dial("tcp", host+":22", sshConfig)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	session, err := client.NewSession()
// 	if err != nil {
// 		client.Close()
// 		return nil, nil, err
// 	}

// 	return client, session, nil
// }

// // func ObserveFLC(nodeAddress string, logger logging.Logger) (*FLCMetrics, error) {

// // 	resp, err := http.Get(fmt.Sprintf("http://%s:9399/", nodeAddress))

// // 	if err != nil {
// // 		logger.Debug("errore nella connection")
// // 		logger.Debug(err.Error())
// // 		return nil, err
// // 	}
// // 	defer resp.Body.Close()
// // 	body, _ := io.ReadAll(resp.Body)
// // 	// manage the body
// // 	logger.Debug(fmt.Sprint(string(body)))

// // 	metrics := &FLCMetrics{}
// // 	err = json.Unmarshal(body, metrics)
// // 	if err != nil {
// // 		logger.Debug("errore nell'unmarshal")
// // 		logger.Debug(err.Error())
// // 		return nil, err
// // 	}
// // 	return metrics, nil

// // }

// // func SetConfig(nodeAddress string, config FLCConfig, logger logging.Logger) (bool, error) {

// // 	config_json, err := json.Marshal(config)
// // 	if err != nil {
// // 		logger.Debug("errore nella codifica json")
// // 	}

// // 	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/config", nodeAddress), bytes.NewReader(config_json))
// // 	if err != nil {
// // 		logger.Debug("errore nella creazione della richiesta")
// // 		logger.Debug(err.Error())
// // 		return false, err
// // 	}
// // 	// create http client
// // 	// do not forget to set timeout; otherwise, no timeout!
// // 	client := http.Client{Timeout: 10 * time.Second}
// // 	// send the request
// // 	res, err := client.Do(req)
// // 	if err != nil {
// // 		logger.Debug("impossible to send request: %s", err)
// // 		return false, err
// // 	}
// // 	logger.Debug("status Code: %d", res.StatusCode)

// // 	return true, nil

// // }

// // func StartFLC() {

// // }
