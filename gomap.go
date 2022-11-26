package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/varuuntiwari/gomap/scan"
)

// Check errors and panic if found
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

// Metadata to print on startup
const Version = "1.0.1"
const Source = "github.com/varuuntiwari/gomap"
const Author = "varuuntiwari@github"

var Host string
var Arr string
var Timeout int
var Ports []int

// Function welcomeText gives a decorated output to the terminal about the tool to the user.
func welcomeText() {
	// Green-colored output
	fmt.Print("\033[32m")
	fmt.Println("+---------------------------------+")
	fmt.Printf("|   Welcome to GoMap %v\t  |\n", Version)
	fmt.Printf("|   Author: %v\t  |\n", Author)
	fmt.Printf("|   Source code: %v\t\t  |\n", fmt.Sprint("\x1b]8;;" + Source + "\x07" + "GitHub" + "\x1b]8;;\x07"))
	fmt.Println("+---------------------------------+")
	// End of colored output
	fmt.Print("\033[0m")
}

func main() {
	welcomeText()

	// Parse parameters
	flag.StringVar(&Host, "h", "", "specify host to scan")
	flag.StringVar(&Arr, "p", "def", "specify ports to scan separated by commas")
	flag.IntVar(&Timeout, "t", 2, "timeout for scanning every port")
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
	s := scan.Scanner{Host: Host, Ports: Ports, Timeout: Timeout}
	s.PrettyRun()
}