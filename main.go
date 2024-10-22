package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/lucsky/cuid"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Creds  string
	Port   string
	Dns_wg string
}

type ReqBody struct {
	Creds string `json:"creds"`
}

var (
	_, b, _, _   = runtime.Caller(0)
	basepath     = filepath.Dir(b)
	serverConfig Config
)

func ValidateRequest(w http.ResponseWriter, r *http.Request, f func(http.ResponseWriter, *http.Request)) {
	var reqBody ReqBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s\nOutput: %s", err.Error(), "Invalid credentials were provided."), http.StatusBadRequest)
		return
	}
	if reqBody.Creds != serverConfig.Creds {
		http.Error(w, fmt.Sprintf("Error: %s", "Wrong credentials were provided."), http.StatusBadRequest)
		return
	}
	f(w, r)
}

func CreateNewUser(w http.ResponseWriter, req *http.Request) {
	cuid := "user-" + cuid.Slug()
	b, err := exec.Command("bash", "/root/vpn.sh", "--addclient", cuid, "--dns1", serverConfig.Dns_wg).CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s\nOutput: %s", err.Error(), string(b)), http.StatusInternalServerError)
		return
	}
	filePath := "/root/" + cuid + ".conf"
	conf, err := exec.Command("cat", filePath).CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error! %s\nOutput: %s", err.Error(), string(conf)), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, fmt.Sprintf("Success! Output: %s@#$%s", filePath, conf))
}

func ReadConfig(filePath string) {
	attempts := 25
	var tmpCfg Config

	for attempts > 0 {
		rawCfg, err := os.ReadFile(filePath)
		if err == nil {
			if err := toml.Unmarshal([]byte(rawCfg), &tmpCfg); err != nil {
				panic(fmt.Sprintf("Error reading configuration file: %v", filePath))
			}
			serverConfig = tmpCfg
			break
		} else {
			fmt.Printf("File '%v' was not found! Attempts left: %d\n", filePath, attempts)
		}
		time.Sleep(5 * time.Second)
		attempts--
	}

	if attempts == 0 {
		panic(fmt.Sprintf("File '%v' cannot be found.", filePath))
	}
}

func main() {
	ReadConfig(basepath + "/serverConfig.toml")
	http.Handle("/", http.FileServer(http.Dir(basepath+"/static")))

	http.HandleFunc("/create_new_user", func(w http.ResponseWriter, r *http.Request) {
		ValidateRequest(w, r, CreateNewUser)
	})

	http.ListenAndServe(":"+serverConfig.Port, nil)
}
