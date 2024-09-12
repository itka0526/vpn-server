package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/lucsky/cuid"
)

func Home(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "<!DOCTYPE html><html><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width\"><title>Wireguard - VPN</title></head><body><h1 style=\"text-align:center\">Welcome!</h1><h2 style=\"text-align:center\">To generate a new user go to \"/create_new_user\"</h2></body></html>")
}

type ReqBody struct {
	Creds string `json:"creds"`
}

func CreateNewUser(w http.ResponseWriter, req *http.Request) {
	var reqBody ReqBody

	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil || reqBody.Creds != os.Getenv("creds") {
		http.Error(w, fmt.Sprintf("Error: %s\nOutput: %s", err.Error(), "Invalid credentials were provided."), http.StatusBadRequest)
		return
	}
	cuid := "user-" + cuid.Slug()
	b, err := exec.Command("bash", "/root/wireguard.sh", "--addclient", cuid).CombinedOutput()
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
		cred, err := os.ReadFile("/root/vpn-server/creds.txt")
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
		panic("'creds.txt' cannot be found.")
	}
}

func main() {
	Init()

	http.HandleFunc("/", Home)
	http.HandleFunc("/create_new_user", CreateNewUser)
	http.ListenAndServe(":80", nil)
}
