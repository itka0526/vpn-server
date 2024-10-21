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
)

// Add your custom DNS :)
var (
	// DNS        = "147.45.231.11" # was used for WireGuard
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

type ReqBody struct {
	Creds string `json:"creds"`
}

func CreateNewUser(w http.ResponseWriter, req *http.Request) {
	var reqBody ReqBody

	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s\nOutput: %s", err.Error(), "Invalid credentials were provided."), http.StatusBadRequest)
		return
	}
	if reqBody.Creds != os.Getenv("creds") {
		http.Error(w, fmt.Sprintf("Error: %s", "Wrong credentials were provided."), http.StatusBadRequest)
		return
	}
	cuid := "user-" + cuid.Slug()
	// b, err := exec.Command("bash", "/root/wireguard.sh", "--addclient", cuid, "--dns1", DNS).CombinedOutput() # was used for WireGuard
	b, err := exec.Command("bash", "/root/vpn.sh", "--addclient", cuid).CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s\nOutput: %s", err.Error(), string(b)), http.StatusInternalServerError)
		return
	}
	conf, err := exec.Command("cat", "/root/"+cuid+".conf").CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error! %s\nOutput: %s", err.Error(), string(conf)), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, fmt.Sprintf("Success! Output: %s", conf))
}

func Init() {
	attempts := 25
	for attempts > 0 {
		cred, err := os.ReadFile(basepath + "/creds.txt")
		if err == nil {
			os.Setenv("creds", string(cred))
			break
		} else {
			fmt.Printf("File 'creds.txt' was not found! Attempts left: %d\n", attempts)
		}
		time.Sleep(5 * time.Second)
		attempts--
	}
	if attempts == 0 {
		panic("File 'creds.txt' cannot be found.")
	}
}

func main() {
	Init()
	fmt.Println(basepath)
	http.Handle("/", http.FileServer(http.Dir(basepath+"/static")))
	http.HandleFunc("/create_new_user", CreateNewUser)
	http.ListenAndServe(":80", nil)
}
