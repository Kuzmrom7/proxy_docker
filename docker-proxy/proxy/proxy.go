package proxy

import (
	"net"
	"log"
	"io"
	"context"
	"encoding/json"
	"bufio"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
)


//Proxy instance containing the list of servers.
type Proxy struct {
	bind     string
	listener net.Listener
}

//New proxy instance
func New(bind string) *Proxy {
	return &Proxy{bind: bind}
}

//Start the proxy server
func (p *Proxy) Start() {
	listener, err := net.Listen("tcp", p.bind)
	if err != nil {
		return
	}
	p.listener = listener
	log.Println("Started proxy bound on", p.bind)
	for {
		if conn, err := listener.Accept(); err == nil {
			containerID, _ := bufio.NewReader(conn).ReadString('\n')
			containerID = containerID[:len(containerID)-1]

			p.handle(conn, containerID)
		} else {
			log.Fatal("Nothing to read in connection? ", err)
			p.Close()
		}
	}
}

//Close the proxy instance
func (p *Proxy) Close() {
	err := p.listener.Close()
	var resp string
	if err == nil {
		resp = "Success."
	} else {
		resp = "Failed."
	}
	log.Println("Closing proxy server:", resp)
}

// Handles incoming requests.
func (p Proxy) handle(up net.Conn, containerID string) {
	defer up.Close()

	//Attach docker client
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}


	//Attach docker container
	resp, err := cli.ContainerAttach(ctx, containerID, types.ContainerAttachOptions{

		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: false,
	})
	if err != nil {
		panic(err)
	}


	// Container inspect
	info,err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		panic(err)

	}
	//parsing in json
	prJson, _ := json.Marshal(info.NetworkSettings.Ports)

	//Send information about host and port container
	up.Write([]byte(prJson))
	up.Write([]byte("\n"))

	//defer
	defer resp.Close()

	//Demultiplex stream
	pipe(up, resp.Conn)
}

//Pipe docker stream on proxy server
func pipe(a, b net.Conn) error {
	errors := make(chan error, 1)
	copy := func(write, read net.Conn) {
		_, err := io.Copy(write, read)
		errors <- err
	}
	go copy(a, b)
	go copy(b, a)
	return <-errors
}