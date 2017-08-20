package main

import (
	"fmt"
	"os"

	"github.com/bgpat/dhcpd/server"
	"github.com/parnurzeal/gorequest"
)

type node struct {
	id         int
	name       string
	macAddress string `json:"mac_address"`
	ipAddress  string `json:"ip_address"`
	image      image
}

type image struct {
	id          int
	path        string
	config      string
	name        string
	description string
}

func main() {
	apiURL := os.Getenv("API_URL")
	s, err := server.New(func(lease *server.Lease) server.Reply {
		fmt.Printf("lease: %+v\n", lease)
		nodes := []node{}
		_, _, err := gorequest.New().Get(apiURL + "/nodes").EndStruct(&nodes)
		if err != nil {
			fmt.Println(err)
			return &server.NAKReply{}
		}
		macAddr := lease.CHAddr.String()
		for _, node := range nodes {
			if node.macAddress == macAddr {
				fmt.Printf("node: %+v\n", node)
				return nil
			}
		}
		return nil
	})
	if err == nil {
		fmt.Printf("%+v\n", s.Handler)
		err = s.Listen()
	}
	fmt.Fprintln(os.Stderr, err.Error())
}
