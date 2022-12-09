package main

import (
	"encoding/json"
	"fmt"
	"github.com/garfeng/n2n_user_edge/model"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	b := &model.ServerConfig{
		Port:                 ":8973",
		Data:                 "./data/data.db",
		SuperNodeAddr:        ":5645",
		CommunityTemplate:    "./templates/community.tmpl",
		CommunityDestination: "./supernode/community.list",
	}
	buff, _ := json.MarshalIndent(b, "", "  ")
	fmt.Println(string(buff))
}

func TestAddUser(t *testing.T) {
	*cfgPath = "../etc/server.json"
}
