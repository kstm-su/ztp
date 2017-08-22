package main

import (
	"bytes"
	"log"

	"github.com/moby/tool/src/moby"
)

func main() {
	moby.MobyDir = "/tmp/moby"
	config := []byte(`
kernel:
  image: linuxkit/kernel:4.9.44
  cmdline: "console=tty0 console=ttyS0 console=ttyAMA0"
init:
  - linuxkit/init:2122f8b7202b383c1be0a91a02122b0c078ca6ac
  - linuxkit/runc:a1b564248a0d0b118c11e61db9f84ecf41dd2d2a
  - linuxkit/containerd:8e4aa6c09e9bceee8300a315c23e0333e187f5fa
onboot:
  - name: dhcpcd
    image: linuxkit/dhcpcd:f3f5413abb78fae9020e35bd4788fa93df4530b7
    command: ["/sbin/dhcpcd", "--nobackground", "-f", "/dhcpcd.conf", "-1"]
services:
  - name: getty
    image: linuxkit/getty:797cb79e0a229fcd16ebf44a0da74bcec03968ec
    env:
     - INSECURE=true
trust:
  org:
    - linuxkit
`)
	m, err := moby.NewConfig(config)
	if err != nil {
		log.Fatalf("Invalid config: %v", err)
	}
	buf := new(bytes.Buffer)
	if err := moby.Build(m, buf, false, ""); err != nil {
		log.Fatalf("%v", err)
	}
	image := buf.Bytes()
	if err := moby.Outputs("./output", image, []string{"iso-bios", "iso-efi"}, 1024, false); err != nil {
		log.Fatalf("Error writing outputs: %v", err)
	}
}
