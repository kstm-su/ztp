package main

import (
	"bytes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/moby/tool/src/moby"
)

func main() {
	router := gin.Default()

	moby.MobyDir = "/tmp/moby"
	router.POST("/create", func(c *gin.Context) {

		log.Printf("%v", c.PostForm("config"))
		config := []byte(c.PostForm("config"))

		m, err := moby.NewConfig(config)
		if err != nil {
			log.Fatalf("Invalid config: %v", err)
		}
		buf := new(bytes.Buffer)
		if err := moby.Build(m, buf, false, ""); err != nil {
			log.Fatalf("%v", err)
		}
		image := buf.Bytes()
		if err := moby.Outputs("/tmp/tftpboot/output", image, []string{"iso-bios", "iso-efi"}, 1024, false); err != nil {
			log.Fatalf("Error writing outputs: %v", err)
		}

		c.JSON(200, gin.H{
			"path": "/tftpboot/output.iso",
		})
	})
	router.Run(":3000")
}
