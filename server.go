package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var pingoo = func() []byte {
	pingu, _ := os.ReadFile("pingoo.txt")
	return pingu
}()

func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Message struct {
	UserName    string
	SendingTime string
	Content     []byte
}

type User struct {
	Name string
	Adrr net.Addr
	Conn net.Conn
}

type Server struct {
	Adrr     string
	Listener net.Listener
	LogFile  *os.File
	Quitch   chan struct{}
	Msgch    chan Message
	UserList []User
}

func CreatServer(listenAdrr string) *Server {
	filleName := "ServerLog_" + time.Now().Format("02-01-2006")
	logFile, err := os.Create(filleName)
	LogFatal(err)
	return &Server{
		Adrr:     listenAdrr,
		LogFile:  logFile,
		Quitch:   make(chan struct{}),
		Msgch:    make(chan Message),
		UserList: make([]User, 0),
	}
}

// [2020-01-20 15:48:41][client.name]:[client.message]
func logNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func serverAlert(s *Server, content string) {
	msg := Message{
		UserName:    "Server",
		SendingTime: logNow(),
		Content:     []byte(content),
	}
	s.Msgch <- msg
}

func formateMsg(msg Message) string {
	return fmt.Sprintf("[%s][%s]: %s\n", msg.SendingTime, msg.UserName, msg.Content)
}

func listenUser(s *Server, user User) {
	for {
		buf := make([]byte, 2048)
		n, err := user.Conn.Read(buf)
		if err != nil {
			user.Conn.Write([]byte("You have been disconected !!!!"))
			user.Conn = nil
			serverAlert(s, fmt.Sprintf("%s has left our chat...", user.Name))
			for i, userF := range s.UserList {
				if userF.Adrr == user.Adrr && userF.Name == user.Name {
					// user.Conn.Close()
					s.UserList = RemoveIndex(s.UserList, i)
					return
				}
			}
		}
		if len(strings.Fields(string(buf))) == 1 {
			user.Conn.Write([]byte("\x1B[A\r\033[K"))
			continue
		}
		s.Msgch <- Message{
			UserName:    user.Name,
			SendingTime: logNow(),
			Content:     buf[:n-1],
		}
	}
}

func StartServer(s *Server) {
	ln, err := net.Listen("tcp", ":"+s.Adrr)
	LogFatal(err)
	s.Listener = ln
	// Listen New Conection
	go func() {
		for {
			conn, err := s.Listener.Accept()
			if len(s.UserList) == 2 {
				conn.Write([]byte("Server is full !!!!"))
				conn.Close()
				continue
			}
			buf := make([]byte, 2048)
			if err != nil {
				serverAlert(s, "Conn err: "+fmt.Sprint(err))
				return
			}
			// conn.SetWriteDeadline(time.Now())
			conn.Write(pingoo)
			n, err := conn.Read(buf)
			if err != nil {
				serverAlert(s, "Conn err: "+fmt.Sprint(err))
				return
			}
			fmt.Println(conn.RemoteAddr())
			name := string(buf[:n-1])
			s.UserList = append(s.UserList, User{
				Adrr: conn.LocalAddr(),
				Name: name,
				Conn: conn,
			})
			stat, _ := s.LogFile.Stat()
			history := make([]byte, stat.Size())
			s.LogFile.ReadAt(history, 0)
			// LogFatal(err)
			conn.Write(history)
			serverAlert(s, fmt.Sprint(name, " has joined our chat..."))
			// listen User
			go listenUser(s, s.UserList[len(s.UserList)-1])
		}
	}()

	// listenMsg and log
	go func() {
		for {
			msg := <-s.Msgch
			fMsg := formateMsg(msg)
			s.LogFile.Write([]byte(fMsg))
			fmt.Print(fMsg)
			for _, user := range s.UserList {
				if user.Name == msg.UserName {
					user.Conn.Write([]byte("\x1B[A\r\033[K"))
				}
				user.Conn.Write([]byte(fMsg))
			}
		}
	}()
	serverAlert(s, "ServerStart")
	<-s.Quitch
	close(s.Msgch)
}

func main() {
	port := "8989"
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	Server := CreatServer(port)
	StartServer(Server)
}

func RemoveIndex(s []User, index int) []User {
	return append(s[:index], s[index+1:]...)
}
