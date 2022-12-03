package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	host      = flag.String("host", "localhost:8080", "The host and port of load balancer")
	connCount = 0
)

func main() {
	servers, err := getServers()
	if err != nil {
		panic(err)
	}

	if host == nil {
		panic("addr can  not be nil")
	}

	listener, err := net.Listen("tcp", *host)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection: %v", err)
		}

		server := chooseBackend(servers)
		fmt.Printf("count:%d, server:%s", connCount, server)
		go func() {
			err := proxy(conn, server)
			if err != nil {
				log.Printf("error proxying connection: %v", err)
			}
		}()
	}

}

func proxy(clientConn net.Conn, server string) error {
	serverConn, err := net.Dial("tcp", server)
	if err != nil {
		return fmt.Errorf("error dialing server: %v", err)
	}

	go io.Copy(serverConn, clientConn)
	go io.Copy(clientConn, serverConn)

	return nil
}

func chooseBackend(servers []string) string {
	server := servers[connCount%len(servers)]
	connCount++

	return server
}

func getServers() ([]string, error) {
	file, err := os.Open("servers.txt")
	if err != nil {
		return nil, fmt.Errorf("error opening servers.txt: %v", err)
	}
	defer file.Close()

	var servers []string
	for {
		var server string
		_, err := fmt.Fscanln(file, &server)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading servers.txt: %v", err)
		}

		servers = append(servers, server)
	}

	return servers, nil
}
