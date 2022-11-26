// Package scan defines the struct Scanner and methods associated with it.
//
// This package is made to faciliate scanning for the tool shodan-go at
// github.com/varuuntiwari/shodan-go, thus it provides an abstract way
// of scanning a particular host and given ports along with various options.
//
// Other than being used for another tool, it can also be compiled and used as
// a standalone scanner for a single host at a time.
//
// The program accepts flag parameters as follows:
//
// -h : specify host to scan, eg. 101.43.75.11, scanme.nmap.org
//
// -p : specify ports to scan, eg. 80, 1-100, 22,21,53. defaults to 1-1000
//
// -t : specify timeout for receiving response from every port. defaults to 2 seconds.
package scan

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// Scanner is a class which stores the host and ports to be
// scanned. This is to provide abstraction for users importing this
// package for individually using the scanning operations given.
//
// The function Run() initializes the scan and declares the
// instance scanned by changing Scanned variable to true. It stores the
// open ports in OpenPorts to be used for further operations such as logging
// or simply printing it out.
//
// The function ShowPorts() prints out the OpenPorts variable
// along with the corresponding service in ServicePorts, if exists.
// Else it is declared as an open port running an unknown service.
type Scanner struct {
	Host  		string
	Ports 		[]int
	Timeout 	int
	OpenPorts 	[]int
	Scanned 	bool
}

// Function ShowPorts is a struct function of Scanner struct which
// prints out the open ports of the struct after checking if the Run function
// has scanned the host.
//
// The function gets the service running on the port from the ServicePorts mapping
// and prints it. Incase not found, it declares it unknown.
//
// The function currently does not print out closed or filtered ports.
func (sc Scanner) ShowPorts() (error) {
	if !sc.Scanned {
		return errors.New("host not scanned")
	}
	fmt.Println("\nPorts Scanned:\nPort\tService\t\tStatus")
	for _, port := range sc.OpenPorts {
		serv, ok := ServicePorts[port]
		if ok {
			fmt.Printf("%v\t%v\t\topen\n", port, serv)
		} else {
			fmt.Printf("%v\tunknown\t\topen\n", port)
		}
	}
	return nil
}

// Function Run is a struct function belonging to Scanner struct which initializes the
// scan of the Host. It iterates through the Ports slice, assigning a goroutine to every
// port and adding the port number to OpenPorts if the TCP connection succeeds.
//
// The function uses sync.WaitGroup for checking if every port has been scanned before
// declaring the Scanner struct as Scanned. It also uses sync.Mutex to lock the area where
// a port is added to the OpenPorts slice to ensure no operation is skipped or done twice.
//
// The function ends with setting the Scanned variable as true.
func (sc *Scanner) Run() (err error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, port := range sc.Ports {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			addr := sc.Host + ":" + fmt.Sprint(port)
			conn, err := net.DialTimeout("tcp", addr, time.Second * time.Duration(sc.Timeout))
			mu.Lock()
			if err == nil && conn != nil {
				sc.OpenPorts = append(sc.OpenPorts, port)
			}
			defer mu.Unlock()
		}(port)
	}
	wg.Wait()
	sc.Scanned = true
	return nil
}

// Function PrettyRun runs the scan, along with timing it and printing the details
// of the scan to the console.
func (sc *Scanner) PrettyRun() (err error) {
	startTime := time.Now()
	fmt.Printf("\nStarting scan at %v\n", startTime.Format("02-01-2006 15:04:05"))
	err = sc.Run()
	if err != nil { return }
	// Print result of scan
	err = sc.ShowPorts()
	if err != nil { return }
	endTime := time.Now()
	fmt.Printf("Scan ended at %v\n", endTime.Format("02-01-2006 15:04:05"))
	fmt.Printf("Scanned %v ports in %v\n", len(sc.Ports), endTime.Sub(startTime))
	return
}

// Function Refresh clears the OpenPorts and changes Scanned status to false.
// Then it proceeds to scan the Host again if the parameter is set to true.
func (sc *Scanner) Refresh(runAgain bool) (err error) {
	sc.Scanned = false
	sc.OpenPorts = nil
	if runAgain {
		err = sc.Run()
	}
	return
}