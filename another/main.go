package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
)

var (
	listenAddr  = flag.String("l", "localhost:9999", "Local address")
	connectAddr = flag.String("r", "localhost:5222", "NATS Server address")
)

//var listenAddr = flag.String("listen", ":8025", "Address to listen on.")
//var connectAddr = flag.String("connect", "mail:25", "Address to connect to.")
var logFile = flag.String("logfile", "", "file to log hex junk into")

var logWriter = log.New(ioutil.Discard, "", 0)

func myCopy(dst io.Writer, src io.Reader, ch chan bool) {
	defer func() { ch <- true }()
	io.Copy(dst, src)
}

func maybeFatal(msg string, e error) {
	if e != nil {
		log.Printf(msg, e)
		runtime.Goexit()
	}
}

type logHexer string

func (lh logHexer) Write(b []byte) (int, error) {
	//logWriter.Printf("%v\n%v\n", string(lh), hex.Dump(b))
	fmt.Printf(string(b))
	return len(b), nil // logger doesn't error
}

func handleConn(c net.Conn, destAddr string) {
	defer c.Close()
	client, err := net.Dial("tcp", destAddr)
	maybeFatal("Error connecting to the other side:  %v", err)
	defer client.Close()

	log.Printf("Connected a new conn")

	ch := make(chan bool)

	go myCopy(c, io.TeeReader(client, logHexer("server")), ch)
	go myCopy(client, io.TeeReader(c, logHexer("client")), ch)

	<-ch
	c.Close()
	client.Close()
	<-ch

	log.Printf("Lost a connection")
}

func main() {
	flag.Parse()

	if *logFile != "" {
		lf, err := os.Create(*logFile)
		if err != nil {
			log.Fatalf("Unable to create logfile: %v", err)
		}
		defer lf.Close()
		logWriter = log.New(lf, "", log.Lmicroseconds)
	}

	addr, err := net.ResolveTCPAddr("tcp", *listenAddr)
	maybeFatal("Error resolving listen address:  %v", err)
	l, err := net.ListenTCP("tcp", addr)
	maybeFatal("Error listening: %v", err)

	for {
		c, err := l.Accept()
		maybeFatal("Error connecting:  %v", err)
		go handleConn(c, *connectAddr)
	}
}
