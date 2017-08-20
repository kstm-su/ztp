package main

import (
	"fmt"
	"os"

	"github.com/bgpat/dhcpd/server"
	"github.com/parnurzeal/gorequest"
)

func main() {
	apiURL := os.Getenv("API_URL")
	s, err := server.New(func(lease *server.Lease) server.Reply {
		fmt.Printf("lease: %+v\n", lease)
		_, body, err := gorequest.New().Get(apiURL + "/nodes").End()
		if err != nil {
			fmt.Println(err)
			return &server.NAKReply{}
		}
		fmt.Printf("nodes: %+v\n", body)
		/*
			for _, node := range body {
			}
		*/
		return nil
	})
	if err == nil {
		fmt.Printf("%+v\n", s.Handler)
		err = s.Listen()
	}
	fmt.Fprintln(os.Stderr, err.Error())
}
