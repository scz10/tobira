package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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
	fmt.Println("Start forwarding port " + tunnel.Local.String() + " -> " + listener.Addr().String())
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

	copyConn := func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()

		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(local, remote)
	go copyConn(remote, local)
}
func connected() (ok bool) {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	return true
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
	if !connected() {
		fmt.Println("No internet connection, waiting for internet connection")
	}
	for !connected() {
		time.Sleep(60 * time.Second)
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	localPortPtr := flag.Int("local", 80, "Specify local port that want be forwarded, default is port 80")
	remotePortPtr := flag.Int("remote", 0, "Specify remote port for destination local forwarded port, default is random available port")
	flag.Parse()

	serverPort, _ := strconv.Atoi(os.Getenv("REMOTE_PORT"))

	localEndpoint := &Endpoint{
		Host: "localhost",
		Port: *localPortPtr,
	}

	serverEndpoint := &Endpoint{
		Host: os.Getenv("REMOTE_SERVER"),
		Port: serverPort,
	}

	remoteEndpoint := &Endpoint{
		Host: os.Getenv("REMOTE_SERVER"),
		Port: *remotePortPtr,
	}

	sshConfig := &ssh.ClientConfig{
		User: os.Getenv("REMOTE_USERNAME"),
		Auth: []ssh.AuthMethod{
			ssh.Password(os.Getenv("REMOTE_PASSWORD")),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if val, _ := strconv.ParseBool(os.Getenv("PASSWORDLESS")); val == true {
		sshConfig = &ssh.ClientConfig{
			User: os.Getenv("REMOTE_USERNAME"),
			Auth: []ssh.AuthMethod{
				publicKeyFile(os.Getenv("SSH_KEY")),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	}

	tunnel := &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	tunnel.Start()
}
