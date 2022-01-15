package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	localAddr  = flag.String("l", "localhost:9999", "Local address")
	remoteAddr = flag.String("r", "localhost:5222", "NATS Server address")
)

func main() {

	flag.Parse()
	fmt.Printf("Listening: %v\nProxying: %v\n\n", *localAddr, *remoteAddr)

	listener, err := net.Listen("tcp", *localAddr)
	if err != nil {
		panic(err)
	}
	for {
		localConn, err := listener.Accept()
		log.Println("New connection", localConn.RemoteAddr())
		if err != nil {
			log.Println("error accepting connection", err)
			continue
		}
		go func() {
			defer localConn.Close()
			natsConn, err := net.Dial("tcp", *remoteAddr)
			if err != nil {
				log.Println("error dialing remote addr", err)
				return
			}
			defer natsConn.Close()
			closer := make(chan struct{}, 2)

			go copy(closer, natsConn, io.TeeReader(localConn, logger("server")))
			go copy(closer, localConn, io.TeeReader(localConn, logger("client")))

			//go copy(closer, natsConn, localConn)
			//go copy(closer, localConn, natsConn)

			<-closer
			log.Println("Connection complete", localConn.RemoteAddr())
		}()

	}

}

type logger string

func (lh logger) Write(b []byte) (int, error) {
	fmt.Printf(string(b))
	return len(b), nil
}

func copy(closer chan struct{}, dst io.Writer, src io.Reader) {
	_, _ = io.Copy(dst, src)
	closer <- struct{}{} // connection is closed, send signal to stop proxy
}
