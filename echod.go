package main

import (
	"container/list"
	"log"
	"net"
	"strconv"
)

const (
	port   = 6667
	bufmax = 256
)

func main() {
	// Create a listener
	clients := list.New()
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for { // Main accept loop
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		// add new connection to list
		clients.PushBack(conn)
		log.Print("client connected")

		// spawn go routine to handle connection
		go func(c net.Conn) {
			// close and remove connection from list when finished
			defer func() {
				c.Close()
				for e := clients.Front(); e != nil; e = e.Next() {
					if e.Value.(net.Conn) == c {
						clients.Remove(e)
						log.Print("client disconnected")
					}
				}
			}()

			b := make([]byte, bufmax)
			for { // keep reading from client
				n, err := c.Read(b)
				if err != nil {
					break
				} else if n > 0 { // we got something, send it to everyone else
					for e := clients.Front(); e != nil; e = e.Next() {
						c := e.Value.(net.Conn)
						n, err = c.Write(b[:n])
					}
				}
			}
		}(conn)
	}
}
