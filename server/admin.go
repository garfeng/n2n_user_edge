package main

import (
	"bytes"
	"changeme/lib"
	"changeme/model"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
)

func AdminAuth(c *gin.Context) {
	u, p, _ := c.Request.BasicAuth()

	if p != cfg.AdminPassword || u != "admin" {
		c.JSON(http.StatusForbidden, model.NewCommonResp(http.StatusForbidden, "forbidden"))
		c.Abort()
		return
	}
	c.Next()
}

func AddUser(c *gin.Context) {
	newUser := &model.AddUserReq{}
	err := c.BindJSON(newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewCommonResp(http.StatusBadRequest, err.Error()))
		return
	}

	user := &model.UserTable{}
	conn := globalConn.Debug()
	conn.Model(user).First(user, "username = ?", newUser.Username)
	if user.ID != 0 {
		c.JSON(http.StatusBadRequest, model.NewCommonResp(http.StatusBadRequest, "user existed"))
		return
	}

	hash, err := lib.Keygen(newUser.Username, newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewCommonResp(http.StatusInternalServerError, err.Error()))
		return
	}

	err = conn.Model(&model.UserTable{}).Create(&model.UserTable{
		Model:    gorm.Model{},
		Username: newUser.Username,
		Password: hash,
		Disabled: false,
		UserId:   newUser.UserId,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewCommonResp(http.StatusInternalServerError, err.Error()))
		return
	}

	err = refreshCommunityList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewCommonResp(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewCommonResp(0, "ok"))
}

func DisableUser(c *gin.Context) {
	userDisabled := &model.DisableUserReq{}
	err := c.BindJSON(userDisabled)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewCommonResp(http.StatusBadRequest, err.Error()))
		return
	}

	user := model.UserTable{}
	conn := globalConn.Debug()
	conn.Model(user).First(user, "username = ? and disabled = ?", userDisabled.Username, false)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, model.NewCommonResp(http.StatusBadRequest, "user not existed"))
		return
	}

	err = conn.Model(user).Where("id = ?", user.ID).Update("disabled", true).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewCommonResp(http.StatusInternalServerError, err.Error()))
		return
	}

	err = refreshCommunityList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewCommonResp(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewCommonResp(0, "ok"))
}

var (
	isManage       = flag.Bool("m", false, "is manage cmd")
	manageUsername = flag.String("u", "", "username to set")
	managePassword = flag.String("p", "", "user password to set")
	manageUserId   = flag.String("uid", "", "userId to set")
	manageCmd      = flag.String("cmd", "addUser", "addUser|disableUser")
)

func adminCommand() {
	url := cfg.ServerUrl + "/admin/" + *manageCmd

	buff, _ := json.Marshal(&model.AddUserReq{
		Username: *manageUsername,
		Password: *managePassword,
		UserId:   *manageUserId,
	})
	w := bytes.NewBuffer(buff)

	req, err := http.NewRequest(http.MethodPost, url, w)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.SetBasicAuth("admin", cfg.AdminPassword)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	cr := &model.CommonResp{}
	json.Unmarshal(body, cr)
	if cr.Status != 0 {
		fmt.Println("Error:", cr.Message)
	} else {
		fmt.Println("succeed")
	}
}
