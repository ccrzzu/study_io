package epoll

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

var epoller *epoll

func main() {
	setLimit()
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
	
	// 创建epoll
	epoller, err = MkEpoll()
	if err != nil {
		panic(err)
	}

	// epoll wait
	go epollWait()

	// epoll ctl add
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
		if err := epoller.Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
	}
}

func setLimit() {
	panic("unimplemented")
}

func epollWait() {
	var buf = make([]byte, 8)
	for {
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}
			if _, err := conn.Read(buf); err != nil {
				// epoll ctl del
				if err := epoller.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}
				conn.Close()
			}
		}
	}
}
