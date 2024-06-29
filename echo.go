package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
)

var sigChannel = make(chan os.Signal)

func main() {
	go echoUDPServer()
	go echoTCPServer()
	go echoHTTPServer()
	
	log.Println("echo server start")
	flag := true
	for flag {
		sig := <- sigChannel
		log.Printf("receive signal: %d desc:%s",sig , sig.String())
		if sig == syscall.SIGTERM || sig == syscall.SIGINT {
			flag = false
			break
		}
	}
	log.Println("echo server stop")
}

func echoUDPServer() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 40000})
	if err != nil {
		log.Printf("[UDP]listen err: %s", err.Error())
		return
	}
	buffer := make([]byte, 2048)
	var length int
	var remoteAddr *net.UDPAddr
	for {
		length, remoteAddr, err = conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}
		_, _ = conn.WriteToUDP(buffer[:length], remoteAddr)
		log.Println("udp client:", remoteAddr.IP.String())
	}
}

func echoTCPServer() {
	conn, err := net.Listen("tcp", "0.0.0.0:40000")
	if err != nil {
		log.Printf("[TCP] listen err: %s", err.Error())
		return
	}
	for {
		tcpConn, err := conn.Accept()
		if err != nil {
			log.Printf("[TCP] accept err: %s", err.Error())
			continue
		}
		go handleTCP(tcpConn)
	}

}

func handleTCP(tcpConn net.Conn) {
	//log.Printf("[TCP] handleTCP: %s", tcpConn.RemoteAddr().String())
	var length int
	var err error
	var buffer []byte
	for {
		buffer = make([]byte, 1024)
		length, err = tcpConn.Read(buffer)
		if err != nil {
			log.Printf("[TCP] Read err: %s", err.Error())
			break
		}
		_, err = tcpConn.Write(buffer[:length])
		if err != nil {
			log.Printf("[TCP] Write err: %s", err.Error())
			break
		}
		log.Println("[TCP] client:", tcpConn.RemoteAddr().String())
	}
	if tcpConn != nil {
		_ = tcpConn.Close()
	}
}

func echoHTTPServer() {
	http.HandleFunc("/", httpIndexHandler)
	http.HandleFunc("/exportip/", httpExportIpHandler)
	err := http.ListenAndServe("0.0.0.0:41000", nil)
	if err != nil {
		log.Printf("[HTTP] listen err: %s", err.Error())
		return
	}
}

func httpIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	_, err := w.Write([]byte("hello"))
	if err != nil {
		log.Printf("[HTTP] listen err: %s", err.Error())
	}
}

func httpExportIpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	host := strings.Split(r.RemoteAddr, ":")[0]
	_, err := w.Write([]byte(fmt.Sprintf("Your export ip is: %s", host)))
	if err != nil {
		log.Printf("[HTTP] listen err: %s", err.Error())
	}
}
