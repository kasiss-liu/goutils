package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8989")
	if err != nil {
		fmt.Println("Dial Error :", err.Error())
	}

	for i := 0; i < 1000; i++ {
		n, err := conn.Write([]byte("app log " + time.Now().Format("2006-01-02 15:05:05") + " testlog for server\n"))

		if err != nil {
			fmt.Println("Write Error:", err.Error())
			break
		} else {
			fmt.Println("Written bytes:", strconv.Itoa(n))
		}

		rsp := make([]byte, 64)
		n, err = conn.Read(rsp)
		if err != nil {
			fmt.Println("Response Error:", err.Error())
		} else {
			fmt.Println("Response:", string(rsp[:n]))
		}
	}
	conn.Close()
}
