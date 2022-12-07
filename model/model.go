package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"io/ioutil"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SetupEdgeParam struct {
	AutoRun    bool              `json:"autoRun"`
	Server     string            `json:"server"`
	EdgeParams map[string]string `json:"edgeParams"`
}

type ChangePasswordParam struct {
	Username    string `json:"username"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ChangePasswordResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type UserTable struct {
	gorm.Model

	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	Disabled bool   `gorm:"column:disabled"`

	// user defined id
	UserId string `gorm:"column:user_id"`
}

func (u *UserTable) TableName() string {
	return "users"
}

type ServerConfig struct {
	Port                 string `json:"port"`
	Data                 string `json:"data"`
	SuperNodeAddr        string `json:"superNodeAddr"`
	CommunityTemplate    string `json:"communityTemplate"`
	CommunityDestination string `json:"communityDestination"`
}

func LoadJSON[T any](name string) (*T, error) {
	buff, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	v := new(T)
	err = json.Unmarshal(buff, v)
	return v, err
}

func SaveJSON(name string, v any) error {
	buff, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil
	}
	return ioutil.WriteFile(name, buff, 0755)
}
