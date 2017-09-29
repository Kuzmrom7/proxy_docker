package main

import (
	"net"
	"fmt"
	"flag"
	"bufio"
)

var (
	bind    = flag.String("b", ":9999", "Address to bind on")
	conn_type = "tcp"
	containerId = "39111edf6286"
)
const(
	CONN_TYPE = "tcp"
	CONN_HOST = "localhost:8080"
)

func main()  {
	conn, err := net.Dial(conn_type, "localhost:9999")
	if err != nil{
		fmt.Print("Can't connect to server")
	}

	for {
		// Send container id
		fmt.Fprintf(conn, containerId + "\n")

		//Localhost and Stream
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: "+message)
	}
}