package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	//connection with ssh to the host
	hostAddress := "nodeserf.cloudmmwunibo.it"
	password := "lTLm#191"
	config := &ssh.ClientConfig{User: "lucaserf", Auth: []ssh.AuthMethod{ssh.Password(password)}, HostKeyCallback: ssh.InsecureIgnoreHostKey()}
	client, err := ssh.Dial("tcp", hostAddress+":22", config)
	if err != nil {
		fmt.Println("Error connecting to host: ", err)
		return
	}
	fmt.Println("Connected to host: ", hostAddress)
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Error creating session: ", err)
	}
	defer session.Close()
	file, err := os.Open("test.sh")
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

	fmt.Println(string(output))

}
