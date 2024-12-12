package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"github.com/lucsky/cuid"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Creds  string
	Port   string
	Dns_wg string
	Srv    string
}

type ValidateRequestReqBody struct {
	Creds string `json:"creds"`
}

var (
	_, b, _, _   = runtime.Caller(0)
	basepath     = filepath.Dir(b)
	serverConfig Config
)

func ValidateRequest(w http.ResponseWriter, r *http.Request, f func(http.ResponseWriter, *http.Request)) {
	var reqBody ValidateRequestReqBody
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

func CreateNewUserWG(w http.ResponseWriter, req *http.Request) {
	cuid := "user-" + cuid.Slug()
	b, err := exec.Command("sudo", "bash", "/root/vpn.sh", "--addclient", cuid, "--dns1", serverConfig.Dns_wg).CombinedOutput()
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

func CreateNewUserOV(w http.ResponseWriter, req *http.Request) {
	cuid := "user-" + cuid.Slug()
	// CANNOT ADD CUSTOM DNS! HAD TO ADD DNS DURING INITIAL SETUP.
	b, err := exec.Command("sudo", "bash", "/root/vpn.sh", "--addclient", cuid).CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s\nOutput: %s", err.Error(), string(b)), http.StatusInternalServerError)
		return
	}
	filePath := "/root/" + cuid + ".ovpn"
	conf, err := exec.Command("cat", filePath).CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error! %s\nOutput: %s", err.Error(), string(conf)), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, fmt.Sprintf("Success! Output: %s@#$%s", filePath, conf))
}

type DelUserReqBody struct {
	Creds       string   `json:"creds"`
	ClientNames []string `json:"clientNames"`
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var rb DelUserReqBody
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s\nOutput: %s", err.Error(), "Invalid credentials were provided."), http.StatusBadRequest)
		return
	}
	if rb.Creds != serverConfig.Creds {
		http.Error(w, fmt.Sprintf("Error: %s", "Wrong data was provided."), http.StatusBadRequest)
		return
	}
	re := regexp.MustCompile(`(?mi)user-[^\.]+`)
	for _, rawClientName := range rb.ClientNames {
		cn := re.FindString(rawClientName)
		b, err := exec.Command("sudo", "bash", "/root/vpn.sh", "--revokeclient", cn, "-y").CombinedOutput()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %s\nOutput: %s", err.Error(), string(b)), http.StatusInternalServerError)
			return
		}
	}
	io.WriteString(w, "Success!")
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

	switch serverConfig.Srv {
	case "wg":
		http.HandleFunc("/create_new_user", func(w http.ResponseWriter, r *http.Request) {
			ValidateRequest(w, r, CreateNewUserWG)
		})
	case "ov":
		http.HandleFunc("/create_new_user", func(w http.ResponseWriter, r *http.Request) {
			ValidateRequest(w, r, CreateNewUserOV)
		})
	}

	http.HandleFunc("/delete_user", DeleteUser)
	http.ListenAndServe(":"+serverConfig.Port, nil)
}
