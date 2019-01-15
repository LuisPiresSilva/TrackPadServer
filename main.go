package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/go-vgo/robotgo"
	"net"
	"os"
	"strconv"
	"stringutil"
	"time"
)

var addr = flag.String("addr", "192.168.1.71", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")

func main() {

	fmt.Printf(stringutil.Reverse("test") + "\n")

	robotgo.MoveMouse(200, 300)
	time.Sleep(500 * time.Millisecond)
	robotgo.MoveMouse(400, 400)
	time.Sleep(500 * time.Millisecond)
	robotgo.MoveMouse(600, 400)

	fmt.Printf("hello, mouse\n")

	// initServer(addr)

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
	}

	// var currentIP
	// var currentNetworkHardwareName string

	for _, address := range addrs {

		// check the address type and if it is not a loopback the display it
		// = GET LOCAL IP ADDRESS
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("Current IP address : ", ipnet.IP.String())
				initServer(ipnet.IP.String())
				// currentIP = ipnet.IP.String()
			}
		}
	}

	ifaces, err := net.Interfaces()
	// handle err
	if err == nil {
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			// handle err
			if err == nil {
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
						fmt.Println("ip net: %s", ip)
					case *net.IPAddr:
						ip = v.IP
						fmt.Println("ip adrr: %s", ip)
					}

					_ = ip //avoid build error, use it somewhere

					// process IP address
				}
			} else {
				fmt.Println("error: %s", err)
			}
		}
	} else {
		fmt.Println("error: %s", err)
	}
}

func initServer(addr string) {
	flag.Parse()

	fmt.Println("Starting server...")

	src := addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s.\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)

	scanner := bufio.NewScanner(conn)

	for {
		ok := scanner.Scan()

		if !ok {
			break
		}

		handleMessage(scanner.Text(), conn)
	}

	fmt.Println("Client at " + remoteAddr + " disconnected.")
}

func handleMessage(message string, conn net.Conn) {
	fmt.Println("> " + message)

	if len(message) > 0 && message[0] == '/' {
		switch {
		case message == "/time":
			resp := "It is " + time.Now().String() + "\n"
			fmt.Print("< " + resp)
			conn.Write([]byte(resp))

		case message == "/quit":
			fmt.Println("Quitting.")
			conn.Write([]byte("I'm shutting down now.\n"))
			fmt.Println("< " + "%quit%")
			conn.Write([]byte("%quit%\n"))
			os.Exit(0)

		default:
			conn.Write([]byte("Unrecognized command.\n"))
		}
	}
}
