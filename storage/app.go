package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const configPath = "./tftpboot/"

type ResponceCreateImage struct {
	Error	string	`json:"error"`
	Path	string	`json:"path"`
}

func FileHash(config string) string{
	h := md5.New()
	h.Write([]byte(config))
	return hex.EncodeToString(h.Sum(nil))
}

func CreateConfigFile(path, config string) error{
	file, err := os.Create(path);
	if err != nil {
		log.Fatal(err)
		return errors.New("Do not create config/yml")
	}
	defer file.Close()

	file.Write(([]byte)(config))
	return nil
}

func GetFilenameWithoutExtension(path string) string{
	fmt.Println(path)
	pos := strings.LastIndex(path,".")
	fmt.Println(path[:pos])
	return path[:pos]
}

func BuildImage(yamlPath string) (string, error){
	fmt.Println("Building Image")
	name := GetFilenameWithoutExtension(yamlPath)
	cmd := exec.Command("moby", "build", "-output", "iso-bios", "-name", name, yamlPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Println("Building...")
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	file := name + ".iso"
	fmt.Println("path: " + file)
	return file, nil
}

func CreateImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var imagePath string
	var err error

	w.Header().Set("Content-Type", "application/json")

	err = r.ParseForm()
	if err != nil {
		fmt.Println("error")
	}

	config := r.PostFormValue("config")
	hash := FileHash(config)
	yamlPath := configPath + hash + ".yml"

	fmt.Println(config)

	if err := CreateConfigFile(yamlPath, config); err == nil {
		imagePath, err = BuildImage(yamlPath)
	}

	fmt.Println(imagePath)

	res := ResponceCreateImage{"", imagePath[1:]}

	fmt.Printf("%+v\n", res)

	uj, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(uj)

	w.WriteHeader(200)
	/*if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Fatal(err)
	}*/
	fmt.Fprintf(w, "%s", uj)
}

func main() {
	fmt.Println("Welcome to Linuxkit Image Storage!");

	router := httprouter.New()
	router.POST("/create", CreateImage)

	log.Fatal(http.ListenAndServe(":8080",router))
}
