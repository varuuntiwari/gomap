package scan

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