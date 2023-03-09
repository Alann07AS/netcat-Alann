package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	openConnection  = make(map[net.Conn]bool)
	newConnections  = make(chan net.Conn)
	deadConnections = make(chan net.Conn)
	testCH          = make(chan string)
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	argrs := os.Args
	if len(argrs) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(0)
	} 
	ln, err := net.Listen("tcp", ":"+argrs[1])
	logFatal(err)
	pingu, _ := os.ReadFile("pingu.txt")
	fmt.Println(string(pingu))
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			logFatal(err)
			openConnection[conn] = true
			newConnections <- conn
		}
	}()
	go func() {
		time.Sleep(time.Second)
		testCH <- "pouette"
	}()
	// testCH <- "po"
	fmt.Println(<-newConnections)
}
