package main

import (
	"flag"
	"fmt"
	"goEnum/portscan"
	"goEnum/webdiscovery"
	"os"
	"strconv"
)

func usage() {
	fmt.Fprintf(os.Stderr, "A basic use of goEnum is something like this: goEnum -t example.com -pi 0 -po 1024 -p tcp")
	os.Exit(1)
}

func main() {
	var target string
	var pi int
	var po int

	flag.StringVar(&target, "t", "", "Target")
	flag.IntVar(&pi, "pi", 0, "Starting Port")
	flag.IntVar(&po, "po", 1024, "Ending Port")

	flag.Parse()

	head := `                            
         _____               
 ___ ___|   __|___ _ _ _____ 
| . | . |   __|   | | |     |
|_  |___|_____|_|_|___|_|_|_|
|___|                        
  `
	fmt.Println(head)
	fmt.Println("[!] Phase 1 started: Port recon")
	fmt.Printf("[!] Running scan on target %s on ports %d to %d\n\n", target, pi, po)
	var results []portscan.Results = portscan.LoopScan("tcp", target, pi, po)
	var webservers []webdiscovery.WebPages

	for i := 0; i < len(results); i++ {
		res := results[i]
		if res.State == "Open" {
			port, _ := strconv.Atoi(res.Port)
			fmt.Printf("[+] Port %d is Open\n", port)
		}
	}
	fmt.Println("\n[!] Phase 2 started: Web server discovery\n")
	for i := 0; i < len(results); i++ {
		res := results[i]
		if res.State == "Open" {
			port, _ := strconv.Atoi(res.Port)
			http, httpe := webdiscovery.Resolve(target, port, "http")
			https, httpse := webdiscovery.Resolve(target, port, "https")
			if httpe == 0 {
				fmt.Printf("[+] Web server available on port %d -> http://%s:%d\n", port, target, port)
				webservers = append(webservers, http)
			}
			if httpse == 0 {
				fmt.Printf("[+] Web server available on port %d -> https://%s:%d\n", port, target, port)
				webservers = append(webservers, https)
			}
		}
	}

	//fmt.Println("\n[!] Phase 3 started: Web Fuzzing with ffuf (light)\n")

}
