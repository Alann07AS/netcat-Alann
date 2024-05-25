package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestCreatServer(t *testing.T) {
	server := CreatServer("8000")
	if server == nil {
		t.Fatal("CreatServer() returned a nil value")
	}
	if server.Adrr != "8000" {
		t.Errorf("Incorrect Adrr value, got: %s, want: %s", server.Adrr, "8000")
	}
	if server.Quitch == nil {
		t.Errorf("Quitch channel is nil")
	}
	if server.Msgch == nil {
		t.Errorf("Msgch channel is nil")
	}
	if server.UserList == nil {
		t.Errorf("UserList slice is nil")
	}
	if server.LogFile == nil {
		t.Errorf("LogFile is nil")
	}
}

func TestStartServer(t *testing.T) {
	server := CreatServer("8000")
	if server == nil {
		t.Fatal("CreatServer() returned a nil value")
	}
	go StartServer(server)

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		t.Errorf("Error dialing server: %s", err)
	}
	if conn == nil {
		t.Errorf("conn is nil")
	}
	defer conn.Close()
}

func TestServerAlert(t *testing.T) {
	server := CreatServer("8000")
	if server == nil {
		t.Fatal("CreatServer() returned a nil value")
	}
	go serverAlert(server, "Test message")

	select {
	case msg := <-server.Msgch:
		if msg.UserName != "Server" {
			t.Errorf("Incorrect username, got: %s, want: %s", msg.UserName, "Server")
		}
		if string(msg.Content) != "Test message" {
			t.Errorf("Incorrect message content, got: %s, want: %s", string(msg.Content), "Test message")
		}
		layout := "2006-02-01 15:04:05"
		_, err := time.Parse(layout, msg.SendingTime)
		if err != nil {
			t.Errorf("Invalid sending time: %s", msg.SendingTime)
		}
		if fmt.Sprintf("[%s][%s]: %s\n", msg.SendingTime, msg.UserName, msg.Content) != formateMsg(msg) {
			t.Errorf("Incorrect formatted message")
		}
	case <-time.After(time.Second):
		t.Errorf("Timed out")
	}
}
