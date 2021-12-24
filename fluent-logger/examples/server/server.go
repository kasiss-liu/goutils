package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	fluentLogger "github.com/kasiss-liu/goutils/fluent-logger"
	loadConfig "github.com/kasiss-liu/goutils/load-config"
)

func main() {
	//读取配置 获取到writers
	config := loadConfig.New("loggers", "./config/server.json")
	writers, err := config.Get("writers").Array()

	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	//注册每个配置里的writer
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
	//启动日志处理器
	err = fluentLogger.StartLogger()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	//准备启动一个socket server 用来接收数据
	var addrPort string
	var sType string
	//读取配置
	server, err := config.Get("server").MapString()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	if addr, ok := (server["AddrPort"]).(string); ok {
		addrPort = addr
	} else {
		fmt.Println("Error:", "Wrong AddrPort Type, Need string")
		return
	}
	if stype, ok := (server["type"]).(string); ok {
		sType = stype
	} else {
		fmt.Println("Error:", "Wrong Type, Need string")
		return
	}
	//启动tcp server
	netServer, err := net.Listen(sType, addrPort)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	//监听接入
	for {
		conn, err := netServer.Accept()
		if err != nil {
			fmt.Println("AcceptError:", err.Error())
			continue
		}
		//处理数据
		go setDataToChan(conn)
	}

}

//goroutiue 处理client输入数据
func setDataToChan(conn net.Conn) (err error) {

	conn.SetReadDeadline(time.Now().Add(time.Minute * 1))
	defer conn.Close()

	lineReader := bufio.NewReader(conn)
	for {
		input, err := lineReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Conn Closed")
			}
			fmt.Println("read line:", err.Error())
			break
		}
		//判断接入是否需要退出
		if isClose(input) {
			conn.Write([]byte("bye bye"))
			break
		}
		//读取处理获取到的string
		reader := bufio.NewReader(strings.NewReader(input))
		//然后读取第一个空格之前的内容
		markBytes, err := reader.ReadBytes(byte(' '))
		if err != nil {
			fmt.Println("mark:", err.Error())
			break
		}
		mark := strings.Trim(string(markBytes), " ")
		//然后读取第二个空格之前的内容
		typeBytes, err := reader.ReadBytes(byte(' '))
		if err != nil {
			fmt.Println("type:", err.Error())
			break
		}
		stype := strings.Trim(string(typeBytes), " ")
		//读取接下来的20个字符 作为时间
		timeStampBytes := make([]byte, 20)
		_, err = reader.Read(timeStampBytes)
		if err != nil {
			fmt.Println("time:", err.Error())
			break
		}
		timestamp := strings.Trim(string(timeStampBytes), " ")
		timeNow, err := time.Parse("2006-01-02 15:04:05", timestamp)
		if err != nil {
			timeNow = time.Now()
		}
		//然后读取剩余内容
		content, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("surplus", err.Error())
			break
		}
		//删除最后一位的换行 \r\n
		content = strings.TrimRight(content, "\r\n")
		content = strings.TrimRight(content, "\n")

		switch strings.ToUpper(stype) {
		case "LOG":
			err = fluentLogger.Log(mark, content, timeNow)
		case "NOTICE":
			err = fluentLogger.Notice(mark, content, timeNow)
		case "WARNING":
			err = fluentLogger.Warning(mark, content, timeNow)
		case "ERROR":
			err = fluentLogger.Error(mark, content, timeNow)
		default:
			err = fluentLogger.Write(mark, stype, content, timeNow)
		}

		if err != nil {
			conn.Write([]byte(err.Error()))
		} else {
			conn.Write([]byte("success\n"))
		}
	}

	return err
}

//验证两个简单的关闭规则
func isClose(content string) bool {
	if !strings.Contains(content, " ") {
		order := strings.Trim(content, "\n")
		order = strings.Trim(content, "\r\n")
		order = strings.TrimSpace(content)
		if order == "quit" || order == "q" {
			return true
		}
	}
	return false
}
