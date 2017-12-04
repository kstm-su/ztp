package main

import (
	"fmt"
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

func main() {
	apiURL := os.Getenv("API_URL")
	client := gorequest.New()
	if socket := os.Getenv("API_SOCKET"); socket != "" {
		client.Transport.Dial = func(_, _ string) (net.Conn, error) {
			return net.Dial("unix", socket)
		}
	}
	s, err := server.New(func(lease *server.Lease) server.Reply {
		fmt.Printf("lease: %+v\n", lease)
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
						lease.Find()
						if lease == nil {
							fmt.Println(err)
							return &server.NAKReply{}
						}
					} else {
						lease.IPAddr = nodeIP
					}
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
		lease.Find()
		if lease == nil {
			fmt.Println(err)
			return &server.NAKReply{}
		}
		reply := &server.ACKReply{
			Lease: lease,
			Options: dhcp.Options{
				dhcp.OptionBootFileName: []byte(images[0].Path + "/syslinux/pxelinux.0"),
				dhcp.OptionHostName:     []byte(images[0].Name),
			},
		}
		fmt.Printf("reply: %+v\n", reply)
		return reply
	})
	if err == nil {
		fmt.Printf("%+v\n", s.Handler)
		err = s.Listen()
	}
	fmt.Fprintln(os.Stderr, err.Error())
}
