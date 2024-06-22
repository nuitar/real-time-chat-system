package main

import (
	"fmt"
	"io"
	"net"
)

type User struct{
	Name string
	Addr string
	ch chan string
	conn net.Conn

	server *Server
}

func newUser(conn net.Conn,server *Server) *User{
	addr := conn.RemoteAddr().String()
	user := &User{addr,addr,make(chan string),conn,server}
	
	user.Online()
	go user.ListenMessage()
	go user.sendMessage()
	return user
}

func (this *User)sendMessage(){
	buf := make([]byte,4096)
		for{
			n, err := this.conn.Read(buf)
			if n == 0{
				this.server.BroadCastUserOnline(this,"已下线")
				return
			}
			if err != nil && err != io.EOF{
				fmt.Println("出错",err)
			}
			msg := string(buf[:n-1])
			this.server.BroadCastUserOnline(this,msg)
		}
}

func (this *User)Online(){
	this.server.mapLock.Lock()	
	this.server.mapUsers[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.BroadCastUserOnline(this,"已上线")
}

func (this *User)OffLine(){
	this.server.mapLock.Lock()	
	delete(this.server.mapUsers,this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCastUserOnline(this,"已下线")
}

func (this *User)ListenMessage(){
	for{
		msg := <- this.ch
		this.conn.Write([]byte(msg + "\n"))
	}
}