package main

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	ReadConfig(basepath + "/serverConfig.toml")
	if len(serverConfig.Creds) <= 4 {
		t.Error("Creds is too short or does not exist")
	} else if len(serverConfig.Port) <= 0 {
		t.Error("No port")
	} else if len(serverConfig.Srv) <= 0 {
		t.Error("Server type is unknown")
	} else if len(serverConfig.Dns_wg) <= 6 {
		t.Error("No WireGuard DNS was found")
	}
}
