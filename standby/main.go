package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func ipToNumber(ip net.IP) uint {
	p := ip
	var ipNumber uint = 0
	p4 := p.To4()
	if len(p4) != 4 {
		return 0
	}
	for _, p := range p4 {
		ipNumber = ipNumber*256 + uint(p)
	}
	return ipNumber
}

func isInLeaseRange(start net.IP, leaseRange uint, target net.IP) bool {
	fmt.Printf("Check %s <= %s\n", start.String(), target.String())
	fmt.Printf("Check %d <= %d <= %d\n", ipToNumber(start), ipToNumber(target), ipToNumber(start)+leaseRange)
	return ipToNumber(start) <= ipToNumber(target) && ipToNumber(target) <= ipToNumber(start)+leaseRange
}

func main() {
	managerURL := os.Getenv("MANAGER_URL")

	startIP := net.ParseIP(os.Getenv("DHCP_START_IP_ADDR"))
	leaseRange, err := strconv.Atoi(os.Getenv("DHCP_LEASE_RANGE"))
	if err != nil {
		fmt.Println(err)
		log.Fatal("Failed strconv.Atoi")
	}

	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal("Cannot fetch hostname")
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal("Cannot fetch list of interfaces")
	}
	macAddress := ""
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal("Cannot fetch address of interface")
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
				break
				fmt.Println(v.IP)
			}
			if isInLeaseRange(startIP, uint(leaseRange), ip) && bytes.Compare(i.HardwareAddr, nil) != 0 {
				fmt.Println("found")
				macAddress = i.HardwareAddr.String()
			}
		}
	}
	if macAddress == "" {
		log.Fatal("Not found such a network interface")
	}

	fmt.Println(macAddress)
	m := map[string]interface{}{
		"name":        hostName,
		"mac_address": macAddress,
	}
	mJson, _ := json.Marshal(m)
	contentReader := bytes.NewReader(mJson)
	req, _ := http.NewRequest("POST", managerURL+"/api/nodes", contentReader)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	println(string(body))

}
