package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

type User struct {
	Name   string
	Output chan Message
}

type Message struct {
	Username string
	Text     string
}

type ChatServer struct {
	Users map[string]User
	Join  chan User
	Leave chan User
	Input chan Message
}

func (cs *ChatServer) Run() {
	for {
		select {
		case user := <-cs.Join:
			cs.Users[user.Name] = user
			go func() {
				cs.Input <- Message{
					Username: "SYSTEM",
					Text:     fmt.Sprintf("%s joined", user.Name),
				}
			}()
		case user := <-cs.Leave:
			delete(cs.Users, user.Name)
			go func() {
				cs.Input <- Message{
					Username: "SYSTEM",
					Text:     fmt.Sprintf("%s left", user.Name),
				}
			}()
		case msg := <-cs.Input:
			for _, user := range cs.Users {
				user.Output <- msg
			}
		}
	}
}

func handleConn(chatServer *ChatServer, conn net.Conn) {
	defer conn.Close()

	user := User{
		Name:   "TEST",
		Output: make(chan Message),
	}
	chatServer.Join <- user
	defer func() {
		chatServer.Leave <- user
	}()

	// Read from conn
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			ln := scanner.Text()
			chatServer.Input <- Message{user.Name, ln}
		}
	}()

	// write to conn
	for msg := range user.Output {
		io.WriteString(conn, msg.Username+": "+msg.Text)
	}
}

func main() {
	server, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer server.Close()

	chatServer := &ChatServer{
		Users: make(map[string]User),
		Join:  make(chan User),
		Leave: make(chan User),
		Input: make(chan Message),
	}
	go chatServer.Run()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleConn(chatServer, conn)
	}
}