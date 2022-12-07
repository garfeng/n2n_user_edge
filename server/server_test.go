package main

import (
	"changeme/model"
	"encoding/json"
	"fmt"
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
