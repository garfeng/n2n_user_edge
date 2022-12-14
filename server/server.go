package main

import (
	"bytes"
	"flag"
	"github.com/garfeng/n2n_user_edge/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
)

var (
	cfgPath    = flag.String("c", "./etc/server.json", "config path")
	cfg        *model.ServerConfig
	globalConn *gorm.DB
)

func init() {
	flag.Parse()
}

func main() {
	if (*cfgPath) == "" {
		*cfgPath = "./etc/server.json"
	}

	var err error
	cfg, err = model.LoadJSON[model.ServerConfig](*cfgPath)
	if err != nil {
		panic(err)
	}

	if *isManage {
		adminCommand()
		return
	}

	connectDB()
	connectSuperNode()
	initTemplate()
	refreshCommunityList()

	r := gin.Default()
	r.POST("/auth/changePassword", ChangePasswordHandler)
	g := r.Group("/admin", AdminAuth)

	g.POST("/addUser", AddUser)
	g.POST("/disableUser", DisableUser)
	r.Run(cfg.Port)

}

func ChangePasswordHandler(c *gin.Context) {
	req := new(model.ChangePasswordParam)
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ChangePasswordResp{
			Status:  http.StatusBadRequest,
			Message: "bad request",
		})
		return
	}

	existedUser := &model.UserTable{}
	err = globalConn.Debug().Model(existedUser).First(existedUser,
		"username = ? and password = ?", req.Username, req.OldPassword).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.ChangePasswordResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if existedUser.ID == 0 {
		c.JSON(http.StatusBadRequest, &model.ChangePasswordResp{
			Status:  http.StatusBadRequest,
			Message: "username or password error",
		})
		return
	}

	udpMutex.Lock()
	defer udpMutex.Unlock()
	if gopher == nil || udpConnToSuperNode == nil {
		c.JSON(http.StatusBadRequest, &model.ChangePasswordResp{
			Status:  http.StatusInternalServerError,
			Message: "super node unavailable",
		})
		return
	}

	err = globalConn.Debug().Model(&model.UserTable{}).
		Where("id = ?", existedUser.ID).
		Update("password", req.NewPassword).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.ChangePasswordResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	err = refreshCommunityList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.ChangePasswordResp{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &model.ChangePasswordResp{
		Status:  0,
		Message: "ok",
	})
}

type CommunityData struct {
	Users []model.UserTable
}

func refreshCommunityList() error {
	// Reload all users and write to community.list
	users := []model.UserTable{}
	globalConn.Model(&model.UserTable{}).Find(&users, "disabled = false")
	w := bytes.NewBuffer(nil)
	communityTemplate.Execute(w, &CommunityData{Users: users})
	buff := w.Bytes()
	ioutil.WriteFile(cfg.CommunityDestination, buff, 0755)

	// reload command
	_, err := udpConnToSuperNode.Write([]byte("reload_communities"))
	return err
}
