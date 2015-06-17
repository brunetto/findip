package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
    "regexp"
	"time"
)

const Debug bool = false

func main() {
	
	
	
	var usage string = `
	fip
	===
	
	Display your ip address as seen from the Internet.
	Use like:
	
	fip once
	
	to run it one time and get back the ip on the STDOUT 
	(useful to start a file server for example)
	
	Use like:
	
	fip loop
	
	to have it run indefinitely and updating the ip.txt file
	(useful to have an update address to connect to).
	
	You can also query the ip address as seen from the local network 
	with:
	
	fip local <interface>
	
	`
	
	var (
		ipRes string
		err error
	)
	
	if len(os.Args) < 2 {
		log.Fatal(usage)
	}
	
	if os.Args[1] == "once" {
		if ipRes, err = retrieveIp(); err != nil {
			log.Fatal("Error retrieving ip: ", err)
		}
		// Exit gracefully after printing the ip
		os.Stdout.WriteString(ipRes)
		os.Exit(0)
	} else if os.Args[1] == "loop" {
		retrieveIpLoop()
	} else if os.Args[1] == "local" {
		if len(os.Args) != 3 {
			log.Fatal(usage)
		}
		
		if ipRes, err = retrieveLocalIp(os.Args[2]); err != nil {
			log.Fatal("Error retrieving local ip: ", err)
		}
		// Exit gracefully after printing the ip
		os.Stdout.WriteString(ipRes)
		os.Exit(0)		
	} else {
		log.Println("Wrong argument: ", os.Args[1])
		log.Fatal(usage)
	}
}

// Retrieve ip address as seen from the Internet
func retrieveIp () (ipRes string, err error) {
	var (
		ipService string = "http://checkip.dyndns.org"
		ipString string
		response *http.Response
		resp []byte
	)
	
	if Debug{log.Println("New check")}
		if response, err = http.Get(ipService); err != nil {
			if Debug{log.Println("Error while downloading ip form ", ipService)}
			return "", err
		}
		
	if Debug{log.Println("Read downloaded data")}
	if resp, err = ioutil.ReadAll(response.Body); err != nil {
		if Debug{log.Println("Can't read data from ", ipService)}
		return "", err
	}
	
	if Debug{log.Println("Close response and outFile")}
	response.Body.Close()
		
	if Debug{log.Println("bytes to string")}
	ipString = string(resp[:])
	if Debug{log.Println("String is: ", ipString)}
	
	if ipRes, err = regIP(ipString); err != nil {
		return "", err
	}
	
	if Debug{log.Println("IP is: ", ipRes)}
	
	return ipRes, nil
}
	
// Retrieve ip address as seen from the Internet every "delay"
func retrieveIpLoop () () {
	log.Println("Starting in loop mode")
	
	var (
		ipRes string
		outFile *os.File
		err error
		delay time.Duration = 15 * time.Minute
	)
	
    for ;; {
		
		if ipRes, err = retrieveIp(); err != nil {
			log.Println("Error retrieving ip: ", err)
			log.Println("Wait ", delay)
			time.Sleep(delay)
			continue
		}
		
		if Debug{log.Println("Write to file")}
		if outFile, err = os.Create("ip.txt"); err != nil {
			log.Println("Error while creating ip.txt: ", err)
			log.Println("Wait ", delay)
			time.Sleep(delay)
			continue
		}
		outFile.WriteString(ipRes)
		outFile.Close()
		
		if Debug{log.Println("Wait ", delay)}
		time.Sleep(delay)
    }
}

// Retrieve ip address as seen from the local network
func retrieveLocalIp (ifaceSelected string) (ipRes string, err error) {
		var (
			addrs []net.Addr
			addr net.Addr
			ifaces []net.Interface
			iface net.Interface
		)
		
		if ifaces, err = net.Interfaces(); err != nil {
			return "", err
		}
		
		for _, iface = range ifaces {
			if iface.Name == ifaceSelected {
				if addrs, err = iface.Addrs(); err != nil {
					return "", err
				}
				for _, addr = range addrs {
					if ipRes, err = regIP(addr.String()); err != nil {
// 						return "", err // If the addr does not correspond, maybe is the ipv6 one
					}
					fmt.Printf(" %v", ipRes)
				}
			}
		}
	return ipRes, nil
}

func regIP (ipIn string) (ipOut string, err error) {
	var (
		ipReg *regexp.Regexp = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
		ipRes []string
	)
		
	if Debug{log.Println("Reg string")}
	if ipRes = ipReg.FindStringSubmatch(ipIn); ipRes == nil {
		return "", fmt.Errorf("Can't extract ip from %v", ipIn)
	}
	ipOut = ipRes[0]
	return ipOut, nil
}









