package main

import (
	"bufio"
	"crypto/tls"

	"fmt"
	"net"
	"os"
	"time"

	"github.com/quic-go/quic-go"
)

const addr = "100.0.0.1:30000"

func main() {
	fmt.Println("hello")
	err := Client()
	if err != nil {
		fmt.Println("client err")
		fmt.Println(err)
	}

}

func Client() error {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	conn, err := quic.DialAddr(addr, tlsConf, &quic.Config{
		KeepAlivePeriod: time.Minute * 5,
		EnableDatagrams: true,
		EnableMPQuic:    true,
	})
	if err != nil {
		fmt.Println("conn err")
		return err
	}

	mpquic, ok := conn.(quic.MPConnection)
	if ok {
		fmt.Println("convert quic conn to MPQUIC interface!!!")
		mpquic.Hello()
		conn2, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 6666})
		if err != nil {
			fmt.Println("create UDPConn error!")
		}
		mpquic.SetSecondConn(conn2)
		go test(mpquic)
	} else {
		fmt.Println("mpquic convert error")
	}
	time.Sleep(2 * time.Second)
	sendData := []byte("hello! It is a test!\n")
	conn.SendMessage(sendData)
	time.Sleep(5 * time.Second)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("What message do you want to send?")
		fmt.Print("message: ")
		text, _ := reader.ReadString('\n')
		sendData := []byte(text)
		err = conn.SendMessage(sendData)
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(1 * time.Second)
	}
}

func test(mp quic.MPConnection) {
	time.Sleep(10 * time.Second)
	err := mp.InitiatePathValidation()
	if err != nil {
		fmt.Printf("encounter err:%v\n", err)
	}
	time.Sleep(1 * time.Second)
	mp.Migration()

	time.Sleep(30 * time.Second)
	//Do the second connection migration
	conn2, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(11, 0, 0, 1), Port: 45000})
	if err != nil {
		fmt.Println("create UDPConn error!")
	}
	mp.SetSecondConn(conn2)
	time.Sleep(1 * time.Second)
	err = mp.InitiatePathValidation()
	if err != nil {
		fmt.Printf("encounter err:%v\n", err)
	}
	time.Sleep(1 * time.Second)
	mp.Migration()

}
