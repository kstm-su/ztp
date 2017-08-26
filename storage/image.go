package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/moby/tool/src/moby"
)

var (
	imageTypes = []string{"iso-bios", "iso-efi"}
	outputDir  = ""
)

func init() {
	moby.MobyDir = os.Getenv("TEMPORARY_DIR")
	outputDir = os.Getenv("OUTPUT_DIR")
}

type Image struct {
	Config string `json:"config"`
	Size   int    `json:"size"`
	Path   string `json:"path"`
}

func (i *Image) ID() string {
	h := md5.New()
	h.Write([]byte(i.Config))
	return hex.EncodeToString(h.Sum(nil))
}

func (i *Image) Build() error {
	m, err := moby.NewConfig([]byte(i.Config))
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := moby.Build(m, buf, false, ""); err != nil {
		return err
	}
	if i.Size == 0 {
		i.Size = 1024
	}
	i.Path = filepath.Join(outputDir, i.ID())
	if err := moby.Outputs(i.Path, buf.Bytes(), imageTypes, i.Size, false); err != nil {
		return err
	}
	return nil
}
