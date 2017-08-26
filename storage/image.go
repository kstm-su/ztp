package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/moby/tool/src/moby"
	"github.com/parnurzeal/gorequest"
)

var (
	imageTypes = []string{"iso-bios", "iso-efi"}
	outputDir  = ""
	buildQueue = []*Image{}
	building   bool
	apiURL     = ""
)

type Image struct {
	ID     int     `json:"id"`
	Config string  `json:"config"`
	Size   int     `json:"size"`
	Path   *string `json:"path"`
	Error  *string `json:"error"`
}

func init() {
	moby.MobyDir = os.Getenv("TEMPORARY_DIR")
	outputDir = os.Getenv("OUTPUT_DIR")
	apiURL = os.Getenv("API_URL")

}

func (i *Image) MD5() string {
	h := md5.New()
	h.Write([]byte(i.Config))
	return hex.EncodeToString(h.Sum(nil))
}

func (i *Image) Append() {
	buildQueue = append(buildQueue, i)
	if !building {
		go buildNext()
	}
}

func (i *Image) Build() error {
	m, err := moby.NewConfig([]byte(i.Config))
	if err != nil {
		errStr := err.Error()
		i.Error = &errStr
		return err
	}
	buf := new(bytes.Buffer)
	if err := moby.Build(m, buf, false, ""); err != nil {
		errStr := err.Error()
		i.Error = &errStr
		return err
	}
	if i.Size == 0 {
		i.Size = 1024
	}
	path := filepath.Join(outputDir, i.MD5())
	if err := moby.Outputs(path, buf.Bytes(), imageTypes, i.Size, false); err != nil {
		errStr := err.Error()
		i.Error = &errStr
		return err
	}
	i.Path = &path
	return nil
}

func buildNext() {
	if len(buildQueue) == 0 {
		return
	}
	image := buildQueue[0]
	buildQueue = buildQueue[1:]
	building = true
	if err := image.Build(); err != nil {
		log.Printf("[ERROR] failed to build: %v\n", err)
	}
	_, body, err := gorequest.New().Put(fmt.Sprintf("%s/images/%d", apiURL, image.ID)).Send(image).End()
	if err != nil {
		log.Printf("[ERROR] failed to update image: %v\n", err)
	} else {
		log.Printf("[DEBUG] update image: %#v\n", body)
	}
	building = false
	buildNext()
}
