// The client sends a GET request and reads nothing from the response, causing a
// poorly implemented server to block.
package main

import (
	"log"
	"net"
)

func main() {
	println("dialing")
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	println("sending GET request")
	_, err = conn.Write([]byte("GET / HTTP/1.1\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write([]byte("Host: localhost\r\n\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	println("blocking and never reading")
	select {}
}
