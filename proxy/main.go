package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/tidwall/redcon"
)

var injectedDelay time.Duration
var delayActive bool

func main() {
	listenAddr := "localhost:20000"
	envAddr := os.Getenv("SERVER_ADDR")
	if len(envAddr) > 0 {
		listenAddr = envAddr
	}
	upstreamAddr := "localhost:6379"
	envRedisAddr := os.Getenv("REDIS_ADDR")
	if len(envRedisAddr) > 0 {
		upstreamAddr = envRedisAddr
	}

	fmt.Printf("Proxy listening on %s, forwarding to %s\n", listenAddr, upstreamAddr)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleConnection(conn, upstreamAddr)
	}
}

func handleConnection(clientConn net.Conn, upstreamAddr string) {
	defer clientConn.Close()

	upstreamConn, err := net.Dial("tcp", upstreamAddr)
	if err != nil {
		fmt.Printf("Could not connect to upstream Redis: %v\n", err)
		return
	}
	defer upstreamConn.Close()

	reader := redcon.NewReader(clientConn)

	for {
		cmd, err := reader.ReadCommand()
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("Error reading command from client: %v\n", err)
			return
		}

		command := strings.ToUpper(string(cmd.Args[0]))

		if command == "FAULT.START" && len(cmd.Args) == 2 {
			delayStr := string(cmd.Args[1])
			delay, parseErr := time.ParseDuration(delayStr + "ms")
			if parseErr == nil {
				injectedDelay = delay
				delayActive = true
				fmt.Printf("Started injecting %v delay.\n", injectedDelay)
				clientConn.Write([]byte("+OK\r\n"))
			} else {
				clientConn.Write([]byte("-ERR invalid delay argument\r\n"))
			}
			continue
		}

		if command == "FAULT.STOP" && len(cmd.Args) == 1 {
			injectedDelay = 0
			delayActive = false
			fmt.Printf("Stopped injecting delay.\n")
			clientConn.Write([]byte("+OK\r\n"))
			continue
		}

		if delayActive && injectedDelay > 0 {
			fmt.Printf("Applying %v injected delay for command: %s\n", injectedDelay, command)
			time.Sleep(injectedDelay)
		}

		_, err = upstreamConn.Write(cmd.Raw)
		if err != nil {
			fmt.Printf("Error writing command to upstream: %v\n", err)
			return
		}

		upstreamRespBuf := make([]byte, 4096)
		n, err := upstreamConn.Read(upstreamRespBuf)
		if err != nil {
			fmt.Printf("Error reading response from upstream: %v\n", err)
			clientConn.Write([]byte(fmt.Sprintf("-ERR %v\r\n", err)))
			return
		}

		_, err = clientConn.Write(upstreamRespBuf[:n])
		if err != nil {
			fmt.Printf("Error writing response to client: %v\n", err)
			return
		}
	}
}
