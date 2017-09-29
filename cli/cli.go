package main

import (
	"net"
	"fmt"
	"flag"
	"bufio"
)

var (
	bind    = flag.String("b", "localhost:9999", "Address to bind on")
	conn_type = "tcp"
	//write container id
	containerId = "39111edf6286"
)


func main()  {
	conn, err := net.Dial(conn_type, *bind)
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