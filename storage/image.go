package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/moby/tool/src/moby"
	"github.com/parnurzeal/gorequest"
)

var (
	imageTypes = []string{"kernel+initrd"}
	outputDir  = ""
	buildQueue = []*Image{}
	building   bool
	apiURL     = ""
)

type Image struct {
	ID          int     `json:"id"`
	Path        *string `json:"path"`
	Config      string  `json:"config"`
	Size        int     `json:"size"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Error       *string `json:"error"`
}

func init() {
	moby.MobyDir = os.Getenv("TEMPORARY_DIR")
	outputDir = os.Getenv("OUTPUT_DIR")
	apiURL = os.Getenv("API_URL")

}

func (i *Image) Append() {
	i.Path = nil
	i.Error = nil
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
	path := filepath.Join(outputDir, strconv.Itoa(i.ID))

	if err := os.RemoveAll(path); err != nil {
		log.Printf("[ERROR] faild to initialize %s", path)
		return err
	}

	if err := os.MkdirAll(path+"/syslinux/pxelinux.cfg", 0755); err != nil {
		return err
	}
	if err := os.Symlink("/usr/share/syslinux/pxelinux.0", path+"/syslinux/pxelinux.0"); err != nil {
		return err
	}
	if err := os.Symlink("/usr/share/syslinux/ldlinux.c32", path+"/syslinux/ldlinux.c32"); err != nil {
		return err
	}

	if err := moby.Formats(path+"/linuxkit", buf.Bytes(), imageTypes, i.Size, false); err != nil {
		errStr := err.Error()
		i.Error = &errStr
		return err
	}
	i.Path = &path

	config, err := ioutil.ReadFile("/usr/share/syslinux/pxelinux.cfg/default")
	if err != nil {
		return err
	}
	cmdline, err := ioutil.ReadFile(path + "/linuxkit-cmdline")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(
		path+"/syslinux/pxelinux.cfg/default",
		[]byte(string(config)+"\tAPPEND "+string(cmdline)),
		0644,
	)
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
