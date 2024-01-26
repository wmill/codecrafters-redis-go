package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var db = struct{
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}


type RedisCommand struct {
	name string
	args []string
}


// a function that parses redis commands of the form "*<number of arguments>\r\n$<length of command>\r\n<command>\r\n"
// and returns a RedisCommand struct
func parseCommand(command string) RedisCommand {
	cmd := RedisCommand{}
	chunks := strings.Split(command, "\r\n")
	args := make([]string, 0)
	cmd.name = strings.ToUpper(chunks[2])
	for i := 4; i < len(chunks); i += 2 {
		args = append(args, chunks[i])
	}
	cmd.args = args
	return cmd
}

func newHandleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		readLen, err := conn.Read(buffer)
		// fmt.Println(string(buffer))
		if readLen == 0 || err != nil && err.Error()  == "EOF" {
			continue
		}
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			continue
		}
		command := parseCommand(string(buffer))
		switch command.name {
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		case "ECHO":
			conn.Write([]byte("+" + command.args[0] + "\r\n"))
		case "SET":
			db.Lock()
			db.m[command.args[0]] = command.args[1]
			db.Unlock()
			if (len(command.args) > 2) {
				var timeout time.Duration
				baseTimeout, _ := strconv.Atoi(command.args[3])
				if (strings.ToUpper(command.args[2]) == "PX") {
					timeout = time.Duration(baseTimeout) * time.Millisecond
				} else if (strings.ToUpper(command.args[2]) == "EX") {
					timeout = time.Duration(baseTimeout) * time.Second
				}
				timer := time.NewTimer(timeout)

				go func()  {
					<-timer.C
					db.Lock()
					delete(db.m, command.args[0])
					db.Unlock()
				}()
			}
			conn.Write([]byte("+OK\r\n"))
		case "GET":
			db.RLock()
			if val, ok := db.m[command.args[0]]; ok {
				conn.Write([]byte("$" + fmt.Sprint(len(val)) + "\r\n" + val + "\r\n"))
			} else {
				conn.Write([]byte("$-1\r\n"))
			}
			db.RUnlock()
		case "CONFIG":
			switch strings.ToUpper(command.args[0]) {
			case "GET":
				config.RLock()
				key := command.args[1]
				if val, ok := config.m[key]; ok {
					conn.Write([]byte("*2\r\n"+ "$" + fmt.Sprint(len(key)) + "\r\n" + key + "\r\n" + "$" + fmt.Sprint(len(val)) + "\r\n" + val + "\r\n"))
				} else {
					conn.Write([]byte("$-1\r\n"))
				}
				config.RUnlock()
			}
		} 
	}
}