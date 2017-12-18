package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	dhcp "github.com/krolaw/dhcp4"
	"github.com/kstm-su/ztp/dhcp/server"
	"github.com/parnurzeal/gorequest"
)

type node struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	MACAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
	Image      image  `json:"image"`
}

type image struct {
	ID          int    `json:"id"`
	Path        string `json:"path"`
	Config      string `json:"config"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var (
	client *gorequest.SuperAgent
	apiURL string
)

func main() {
	apiURL = os.Getenv("API_URL")
	client = gorequest.New()
	if socket := os.Getenv("API_SOCKET"); socket != "" {
		client.Transport.Dial = func(_, _ string) (net.Conn, error) {
			return net.Dial("unix", socket)
		}
	}
	s, err := server.New(func(lease *server.Lease) server.Reply {
		macAddr := lease.CHAddr.String()
		nodes := []node{}
		_, _, err := client.Get(apiURL + "/nodes").EndStruct(&nodes)
		if err != nil {
			fmt.Println(err)
			return &server.NAKReply{}
		}
		images := []image{}
		_, _, err = client.Get(apiURL + "/images").EndStruct(&images)
		if err != nil {
			fmt.Println(err)
			return &server.NAKReply{}
		}
		for _, node := range nodes {
			if strings.ToLower(node.MACAddress) == macAddr {
				fmt.Printf("node: %+v\n", node)
				if lease.IPAddr == nil {
					nodeIP := net.ParseIP(node.IPAddress)
					if nodeIP == nil {
						if err := lease.Find(); err != nil {
							fmt.Println(err)
							return &server.NAKReply{}
						}
					} else {
						lease.IPAddr = nodeIP
						lease.Update()
					}
				}
				if node.Image.Path == "" {
					node.Image.Path = "/default"
					node.Name = ""
				}
				reply := &server.ACKReply{
					Lease: lease,
					Options: dhcp.Options{
						dhcp.OptionBootFileName: []byte(node.Image.Path + "/syslinux/pxelinux.0"),
						dhcp.OptionHostName:     []byte(node.Name),
					},
				}
				fmt.Printf("reply: %+v\n", reply)
				node.IPAddress = lease.IPAddr.String()
				go client.Put(fmt.Sprintf("%s/nodes/%d", apiURL, node.ID)).Send(node).End()
				return reply
			}
		}
		fmt.Println("unknown MAC address: ", macAddr)
		if lease.IPAddr == nil {
			if err := lease.Find(); err != nil {
				fmt.Println(err)
				return &server.NAKReply{}
			}
		}
		reply := &server.ACKReply{
			Lease: lease,
			Options: dhcp.Options{
				dhcp.OptionBootFileName: []byte("/default/syslinux/pxelinux.0"),
			},
		}
		fmt.Printf("reply: %+v\n", reply)
		go client.Post(fmt.Sprintf("%s/nodes", apiURL)).Send(node{
			MACAddress: lease.CHAddr.String(),
			IPAddress:  lease.IPAddr.String(),
		}).End()
		return reply
	})
	if err := restoreIPs(s); err != nil {
		log.Fatalf("failed to restore leased IP addresses: %v\n", err)
	}
	if err == nil {
		fmt.Printf("%+v\n", s.Handler)
		err = s.Listen()
	}
	fmt.Fprintln(os.Stderr, err.Error())
}

func restoreIPs(s *server.Server) []error {
	nodes := []node{}
	if _, _, err := client.Get(apiURL + "/nodes").EndStruct(&nodes); err != nil {
		return err
	}
	for _, node := range nodes {
		ip := net.ParseIP(node.IPAddress)
		if ip == nil {
			return []error{fmt.Errorf("failed to parse IP addr of node#%d", node.ID)}
		}
		mac, err := net.ParseMAC(node.MACAddress)
		if err != nil {
			return []error{err}
		}
		s.Handler.Leases.Use(ip, mac).Update()
	}
	return nil
}
