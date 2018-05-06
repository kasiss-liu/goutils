package main

import (
	"fmt"
	"net"

	"github.com/kasiss-liu/go-tools/fluent-logger"
	"github.com/kasiss-liu/go-tools/load-config"
)

func main() {
	config := loadConfig.New("loggers", "./config/server.json")
	writers, err := config.Get("writers").Array()

	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	for _, m := range writers {
		if writer, ok := m.(map[string]interface{}); ok {
			dir := writer["dir"].(string)
			wType := writer["type"].(string)
			mode := writer["mode"].(float64)
			mark := writer["mark"].(string)
			prefix := writer["prefix"].(string)

			switch wType {
			case "file":
				res := fluentLogger.RegisterFileLogger(mark, dir, prefix, int(mode))
				if !res {
					fmt.Println("Error:", "writer `"+mark+"` register failed")
					return
				}
			}
		}
	}
	err = fluentLogger.StartLogger()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	server, err := config.Get("server").MapString()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	if addrPort, ok := (server["AddrPort"]).(string); !ok {
		fmt.Println("Error:", "Wrong AddrPort Type, Need string")
		return
	}
	if sType, ok := (server["type"]).(string); !ok {
		fmt.Println("Error:", "Wrong Type, Need string")
		return
	}
	server, err := net.Listen(stype, addrPort)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("AcceptError:", err.Error())
		}
		var data = make([]byte, 1024)
		conn.Read(data)

	}

}
