package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
    "regexp"
	"time"
)

const Debug bool = false

func main() {
	
	var (
		ipRes string
		outFile *os.File
		err error
		delay time.Duration = 15 * time.Minute
	)
	
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
	
	`
	
	if len(os.Args) != 2 {
		log.Fatal(usage)
	}
	
	if os.Args[1] == "once" {
		if ipRes, err = retrieveIp(); err != nil {
			log.Fatal("Error retrieving ip: ", err)
		}
		os.Stdout.WriteString(ipRes)
		os.Exit(0)
	} else if os.Args[1] != "loop" {
		log.Println("Wrong argument: ", os.Args[1])
		log.Fatal(usage)
	}
	
	log.Println("Starting in loop mode")
	
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

func retrieveIp () (string, error) {
	var (
		ipService string = "http://checkip.dyndns.org"
		ipReg *regexp.Regexp = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
		ipRes []string
		ipString string
		err error
		response *http.Response
		resp []byte
	)
	
	if Debug{log.Println("New check")}
		if response, err = http.Get(ipService); err != nil {
			log.Println("Error while downloading ip form ", ipService)
			return "", err
		}
		
	if Debug{log.Println("Read downloaded data")}
	if resp, err = ioutil.ReadAll(response.Body); err != nil {
		log.Println("Can't read data from ", ipService)
		return "", err
	}
	
	if Debug{log.Println("Close response and outFile")}
	response.Body.Close()
		
	if Debug{log.Println("bytes to string")}
	ipString = string(resp[:])
	if Debug{log.Println("String is: ", ipString)}
	
	if Debug{log.Println("Reg string")}
	if ipRes = ipReg.FindStringSubmatch(ipString); ipRes == nil {
		log.Println("Can't extract ip from ", ipString)
		return "", err
	}
	
	if Debug{log.Println("IP is: ", ipRes[0])}
	
	return ipRes[0], nil
}
