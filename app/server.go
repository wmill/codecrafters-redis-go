package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
)

var config = struct{
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}



func main() {

	dir := flag.String("dir", "", "The directory where RDB files are stored")
	dbfilename := flag.String("dbfilename", "", "The name of the RDB file")

	flag.Parse()

	config.Lock()
	config.m["dir"] = *dir
	config.m["dbfilename"] = *dbfilename
	config.Unlock()

	
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go newHandleConnection(conn)
	}
}

