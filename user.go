package main

import "net"

type User struct{
	Name string
	Addr string
	ch chan string
	conn net.Conn
}

func newUser(conn net.Conn) *User{
	addr := conn.RemoteAddr().String()
	user := &User{addr,addr,make(chan string),conn}

	go user.ListenMessage()
	return user
}

func (this *User)ListenMessage(){
	for{
		msg := <- this.ch
		this.conn.Write([]byte(msg + "\n"))
	}
}