package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
)

var globalThread int = 100
var globalports = []int{}

func worker(ports chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range ports {
		PORT := ":" + strconv.Itoa(p)
		l, err := net.Listen("tcp4", PORT)
		if err != nil {
			fmt.Println(err)
			//return
			continue
		}
		defer l.Close()

		for {
			c, err := l.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			go handleConnection(c)
		}
	}
}

func handleConnection(c net.Conn) {
	fmt.Printf("From %s to %s\n", c.RemoteAddr().String(), c.LocalAddr().String())
	c.Close()
}

func main() {
	/*
		var portsRange string
		flag.StringVar(&portsRange, "p", "", "ports range, split by ',' (default is top1k)")

		flag.Parse()

		if portsRange != "" {
			tmp := strings.Split(portsRange, ",")
			for _, p := range tmp {
				port, _ := strconv.Atoi(p)
				globalports = append(globalports, port)
			}
		} else {
			//fmt.Println("no -p")
			globalports = top1kports[0:]
		}
	*/
	// usually set p from 1 to 65535 to test full port range. need to test on shellbox to see if resource is allowed.
	for p := 10; p <= 100; p++ {
		globalports = append(globalports, p)
	}

	fmt.Println(globalports)

	ports := make(chan int, globalThread)

	var wg sync.WaitGroup

	for i := 0; i < cap(ports); i++ {
		wg.Add(1)
		go worker(ports, &wg)
	}

	go func() {
		for _, p := range globalports {
			ports <- p
		}
	}()

	wg.Wait()
	close(ports)
	log.Println("[+] Finished")
}
