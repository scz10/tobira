package main

import (
	"fmt"
    "net"
    "log"
    "io"
    "io/ioutil"

	"golang.org/x/crypto/ssh"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

type SSHtunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

func (tunnel *SSHtunnel) Start() error {
    serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
    }
    
	listener, err := serverConn.Listen("tcp", tunnel.Remote.String())
	if err != nil {
		log.Fatalln(fmt.Printf("Listen open port ON remote server error: %s", err))
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHtunnel) forward(remote net.Conn) {
	local, err := net.Dial("tcp", tunnel.Local.String())
	if err != nil {
		log.Fatalln(fmt.Printf("Dial INTO local service error: %s", err))
	}	

	copyConn:=func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()
		
		_, err:= io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(local, remote)
	go copyConn(remote, local)
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot read SSH public key file %s", file))
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot parse SSH public key file %s", file))
		return nil
	}
	return ssh.PublicKeys(key)
}

func main() {
	localEndpoint := &Endpoint{
		Host: "localhost",
		Port: 3306,
	}

	serverEndpoint := &Endpoint{
		Host: "192.168.100.9",
		Port: 22,
	}

	remoteEndpoint := &Endpoint{
		Host: "localhost",
		Port: 6000,
	}

	sshConfig := &ssh.ClientConfig{
		User: "coba",
		Auth: []ssh.AuthMethod{
			//publicKeyFile("C:/Users/Xian/.ssh/id_rsa"),
			ssh.Password("aja"),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	tunnel := &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	tunnel.Start()
}