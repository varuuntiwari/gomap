package scan

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// Scanner is a class which initializes the host and ports to be
// scanned. It has the functions Run() which initializes the scan
// and declares the instance scanned by changing Scanned variable to
// true. The function ShowPorts() prints out the OpenPorts variable
// along with the corresponding service in ServicePorts, if exists.
// Else it is declared as an open port running an unknown service.
type Scanner struct {
	Host  		string
	Ports 		[]int
	OpenPorts 	[]int
	Scanned 	bool
}

// Map ports to their services
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

func (sc *Scanner) Run() (err error) {
	var wg sync.WaitGroup
	for _, port := range sc.Ports {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			addr := sc.Host + ":" + fmt.Sprint(port)
			conn, err := net.DialTimeout("tcp", addr, time.Second * 2)
			if err == nil && conn != nil {
				sc.OpenPorts = append(sc.OpenPorts, port)
			}
		}(port)
	}
	wg.Wait()
	sc.Scanned = true
	return nil
}