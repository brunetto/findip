package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
    "regexp"
	"time"
)

func main() {
	const Debug bool = false
	var (
		ipService string = "http://checkip.dyndns.org"
		ipReg *regexp.Regexp = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
		ipRes []string
		ipString string
		outFile *os.File
		response *http.Response
		resp []byte
		err error
		delay time.Duration = 15 * time.Minute
	)

    for ;;{
		
		if Debug{log.Println("New check")}
		if response, err = http.Get(ipService); err != nil {
			log.Println("Error while downloading ip form ", ipService, ": ", err)
			log.Println("Wait ", delay)
			time.Sleep(delay)
			continue
		}
		
		if Debug{log.Println("Read downloaded data")}
		if resp, err = ioutil.ReadAll(response.Body); err != nil {
			log.Println("Can't read data from ", ipService, " with error: ", err)
			log.Println("Wait ", delay)
			time.Sleep(delay)
			continue
		}
		
		if Debug{log.Println("bytes to string")}
		ipString = string(resp[:])
		if Debug{log.Println("String is: ", ipString)}
		
		if Debug{log.Println("Reg string")}
		if ipRes = ipReg.FindStringSubmatch(ipString); ipRes == nil {
			log.Println("Can't extract ip from ", ipString)
			log.Println("Wait ", delay)
			time.Sleep(delay)
			continue
		}
		
		if Debug{log.Println("IP is: ", ipRes[0])}
		
		if Debug{log.Println("Write to file")}
		if outFile, err = os.Create("ip.txt"); err != nil {
			log.Println("Error while creating ip.txt: ", err)
			log.Println("Wait ", delay)
			time.Sleep(delay)
			continue
		}
		outFile.WriteString(ipRes[0])
		
		if Debug{log.Println("Close response and outFile")}
		response.Body.Close()
		outFile.Close()		
		
		if Debug{log.Println("Wait ", delay)}
		time.Sleep(delay)
    }
}

