package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	mapUsers     map[string]*User
	broadMessage chan string

	mapLock sync.RWMutex
}

func NewServer(ip string, port int) *Server {
	server := &Server{ip, port, make(map[string]*User),
		make(chan string), sync.RWMutex{}}
	return server
}

func (this *Server) ListenBroadMessage() {
	for {
		message := <-this.broadMessage
		this.mapLock.Lock()
		for _, user := range this.mapUsers {
			user.ch <- message
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) BroadCastUserOnline(user *User, msg string){
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	
	this.broadMessage <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	// fmt.Println("链接已建立")
	user := newUser(conn)
	this.mapLock.Lock()	
	this.mapUsers[user.Name] = user
	this.mapLock.Unlock()

	this.BroadCastUserOnline(user,"已上线")

	go func(){
		buf := make([]byte,4096)
		for{
			n, err := conn.Read(buf)
			if n == 0{
				this.BroadCastUserOnline(user,"已下线")
				return
			}
			if err != nil && err != io.EOF{
				fmt.Println("出错",err)
			}
			msg := string(buf[:n-1])
			this.BroadCastUserOnline(user,msg)	

		}
	}()
}

func (this *Server) Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	defer listen.Close()

	if err != nil {
		fmt.Println("连接错误,", err)
		return
	}
	go this.ListenBroadMessage()

	for {
		coon, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go this.Handler(coon)
	}
}
