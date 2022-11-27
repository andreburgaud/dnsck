package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	textVersion = "0.3.0"
	appName     = "dnsck"
)

const textUsage = `
%[1]s performs DNS lookups and validates the connection with one
or more hosts. It returns the results in JSON format allowing
parsing with tools like jq (https://stedolan.github.io/jq/).

USAGE:
  %[1]s [FLAGS] <host_names> | <hosts_file>

FLAGS:
  -h, --help       displays %[1]s help
  -V, --version    displays %[1]s version
  -f, --file       reads input from a CSV file

EXAMPLES:
  %[1]s -h
  %[1]s -V
  %[1]s -f some_hosts.txt
  %[1]s google.com
  %[1]s example.com:443
  %[1]s google.com example.com:80

`

// Result holds the list of servers and count summaries sent back to the main function
type Result struct {
	Servers        []Server `json:"servers"`
	Count          int      `json:"count"`
	CountConnError int      `json:"count_connection_error"`
	CountDnsError  int      `json:"count_dns_error"`
}

// Server holds a server information
type Server struct {
	Hostname   string     `json:"hostname"`
	Connection Connection `json:"connection"`
	Domain     Domain     `json:"domain"`
}

// Domain holds the result of a domain lookups
type Domain struct {
	Name        string   `json:"name"`
	IPAddresses []string `json:"ip_addresses,omitempty"`
	CName       string   `json:"canonical_name,omitempty"`
	TextRecords []string `json:"dns_text_records,omitempty"`
	Error       string   `json:"error,omitempty"`
}

// Connection holds the result of a connection attempt to a server
type Connection struct {
	Address string `json:"address"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// usage displays the help
func usage() {
	fmt.Printf(textUsage, appName)
}

// connect performs lookup for the host, cname and other DNS records
func validateDomain(domain string) Domain {
	d := Domain{Name: domain}
	addrs, err := net.LookupHost(domain)
	if err != nil {
		d.Error = err.Error()
		return d
	}

	d.IPAddresses = addrs

	cname, err := net.LookupCNAME(domain)
	if err != nil {
		d.Error = err.Error()
		return d
	}

	d.CName = cname

	records, err := net.LookupTXT(domain)
	if err == nil {
		d.TextRecords = records
	}

	return d
}

// connect attempts to connect to a server via a TCP connect
func connect(addr string) Connection {
	c := Connection{Address: addr}

	duration, _ := time.ParseDuration("2s")
	conn, err := net.DialTimeout("tcp", addr, duration)
	if err != nil {
		c.Success = false
		c.Error = err.Error()
		return c
	}
	defer conn.Close()
	c.Success = true
	return c
}

// validateServer perform the validation of the connection and
// domain information for one server
func validateServer(server string, ch chan<- Server) {
	var (
		port   string
		domain string
		err    error
	)

	if domain, port, err = net.SplitHostPort(server); err != nil {
		domain = server
		port = "443"
	}

	d := validateDomain(domain)

	addr := net.JoinHostPort(domain, port)
	c := connect(addr)

	ch <- Server{
		Hostname:   domain,
		Connection: c,
		Domain:     d,
	}
}

// validateServers triggers the validation of the connection and domain information
// for a list of server.
func validateServers(servers []string) Result {
	r := Result{}

	ch := make(chan Server)
	defer close(ch)

	for _, server := range servers {
		fmt.Fprintf(os.Stderr, "processing: %s\n", server)
		go validateServer(server, ch)
	}

	count := 0
	countConnErr := 0
	countDnsErr := 0

	for range servers {
		count++
		s := <-ch
		if !s.Connection.Success {
			countConnErr++
		}
		if len(s.Domain.Error) > 0 {
			countDnsErr++
		}
		r.Servers = append(r.Servers, s)
	}
	r.Count = count
	r.CountConnError = countConnErr
	r.CountDnsError = countDnsErr
	return r
}

// readServers returns the first column of a CSV containing servers as first fields
// The server might be a domain name / host name followed by an optional :<port>
// Example: google.com:443
func readServers(file string) []string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Skip first line (headers)
	reader := csv.NewReader(f)
	reader.Comment = '#'

	if _, err := reader.Read(); err != nil {
		log.Fatal(err)
	}

	var hosts []string
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		hosts = append(hosts, record[0])
	}
	return hosts
}

func buildJsonDomains(result Result) (string, error) {
	buf, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func printResults(result Result) {
	j, err := buildJsonDomains(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while converting struct to json")
	}
	fmt.Println(j)
}

func main() {

	var help bool
	var version bool
	var file string

	flag.BoolVar(&help, "help", false, "help")
	flag.BoolVar(&version, "version", false, "version")
	flag.BoolVar(&version, "V", false, "version")
	flag.StringVar(&file, "file", "", "file")
	flag.StringVar(&file, "f", "", "file")
	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}

	if version {
		fmt.Printf("%s version %s\n", appName, textVersion)
		os.Exit(0)
	}

	if len(file) > 0 {
		servers := readServers(file)
		printResults(validateServers(servers))
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	printResults(validateServers(args))
}
