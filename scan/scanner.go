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

// Map ports to commonly used services
var ServicePorts = map[int]string{
	21: "ftp",
	22: "ssh",
	23: "telnet",
	25: "smtp",
	43: "whois",
	53: "dns",
	69: "tftp",
	80: "http",
	123: "ntp",
	135: "msrpc",
	389: "ldap",
	443: "https",
	512: "rexec",
	513: "rlogin",
	514: "syslog",
	520: "rip",
	587: "smtp",
	1433: "mssql",
	3306: "mysql",
	5432: "postgres",
	8080: "http-proxy",	
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
	fmt.Println("Ports Scanned:\nPort\tService\t\tStatus")
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