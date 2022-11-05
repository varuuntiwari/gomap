package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/varuuntiwari/gomap/scan"
)

// Handle errors
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

const Version = "1.0"

var Host string
var Arr string
var Ports []int

func main() {
	fmt.Printf("Welcome to GoMap %v\n----------------\n", Version)
	// Parse parameters
	flag.StringVar(&Host, "h", "", "specify host to scan")
	flag.StringVar(&Arr, "p", "def", "specify ports to scan separated by commas")
	flag.Parse()
	if flag.Parsed() {
		if Host == "" {
			panic("[-] host not specified")
		}
		fmt.Println("[+] Parameters parsed")
	} else {
		panic("[-] could not parse parameters")
	}
	// Get ports from string parameter
	re := regexp.MustCompile(`^\d+[-]\d+$`)
	
	if Arr == "def" {
		Arr = "1-1000"
	}
	if re.MatchString(Arr) {
		arr := strings.Split(Arr, "-")
		i, _ := strconv.Atoi(arr[0])
		j, _ := strconv.Atoi(arr[1])
		if i < 1 || i > j || j > 65535 {
			panic("[-] port range specified out of range")
		}
		for ; i <= j; i++ {
			Ports = append(Ports, i)
		}
		fmt.Printf("[+] Scanning ports in the range %v\n", Arr)
	} else {
		tmp := strings.Split(Arr, ",")
		for _, port := range tmp {
			if t, err := strconv.Atoi(port); err == nil {
				Ports = append(Ports, t)
			} else {
				panic(err.Error())
			}
		}
		fmt.Printf("[+] The ports to be scanned are %+v\n", Ports)
	}

	// Initialize Scanner instance
	s := scan.Scanner{Host: Host, Ports: Ports}
	startTime := time.Now()
	fmt.Printf("Starting scan at %v\n", startTime.Format("02-01-2006 15:04:05"))
	err := s.Run()
	HandleError(err)
	// Print result of scan
	err = s.ShowPorts()
	HandleError(err)
	endTime := time.Now()
	fmt.Printf("Scan ended at %v\n", endTime.Format("02-01-2006 15:04:05"))
	fmt.Printf("Scanned %v ports in %v", len(s.Ports), endTime.Sub(startTime))
}