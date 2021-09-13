package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func main() {
	fmt.Println("hello world!")

	fi := "nmap-services.csv"

	f, err := os.Open(fi)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	r := bufio.NewReader(f)
	s, e := Readln(r)
	for e == nil {
		fmt.Println(s)
		s, e = Readln(r)

		split := strings.Fields(s)
		fmt.Println("length: ", len(split))
		fmt.Println(split[0])
		fmt.Println(split[1])
		if "1/udp" == split[1] {
			fmt.Println(split[0])
		}

		break
	}

}
