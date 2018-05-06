package main

import (
	"bufio"
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

	var addrPort string
	var sType string

	server, err := config.Get("server").MapString()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	if addr, ok := (server["AddrPort"]).(string); !ok {
		fmt.Println("Error:", "Wrong AddrPort Type, Need string")
		return
	} else {
		addrPort = addr
	}
	if stype, ok := (server["type"]).(string); !ok {
		fmt.Println("Error:", "Wrong Type, Need string")
		return
	} else {
		sType = stype
	}
	netServer, err := net.Listen(sType, addrPort)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	for {
		conn, err := netServer.Accept()
		if err != nil {
			fmt.Println("AcceptError:", err.Error())
			continue
		}

		//		var data = make([]byte, 1024)
		//		conn.Read(data)
		fmt.Println("start exec")
		go setDataToChan(conn)
		//		reader := bufio.NewReader(conn)
		//		n := reader.Size()
		//		strBytes := make([]byte, 0, n)
		//		reader.Read(strBytes)
		//		fmt.Printf("%d,%#v", n, string(strBytes))
	}

}

func setDataToChan(conn net.Conn) (int, error) {

	defer conn.Close()
	//然后读取第一个空格之前的内容
	reader := bufio.NewReader(conn)
	for {

		//删除最后一位的换行 \r\n

		typeBytes, err := reader.ReadBytes(byte(' '))
		if err != nil {
			fmt.Println(err.Error())
			return 0, err
		}
		//然后读取第二个空格之前的内容
		timeBytes, err := reader.ReadBytes(byte(' '))
		if err != nil {
			fmt.Println(err.Error())
			return 0, err
		}
		//然后读取剩余内容
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			return 0, err
		}
		fmt.Println(string(typeBytes), string(timeBytes), str)
	}

	return 0, nil
}
