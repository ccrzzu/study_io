package goroutine_per_conn

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func main_server() {
	ln, err := net.Listen("tcp", ":8972")
	if err != nil {
		panic(err)
	}

	// import _ "net/http/pprof"  这个 才有用
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()
	var connections []net.Conn
	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()
	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}
			log.Printf("accept err: %v", e)
			return
		}
		go func() {
			handleConn(conn)
		}()
		connections = append(connections, conn)
		if len(connections)%100 == 0 {
			log.Printf("total number of connections: %v", len(connections))
		}
	}
}
func handleConn(conn net.Conn) {
	io.Copy(ioutil.Discard, conn)
}
