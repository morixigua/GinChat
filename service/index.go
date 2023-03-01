package service

import (
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")
	//c.JSON(http.StatusOK, gin.H{
	//	"message": "welcome!!",
	//})
}

// ToRegister
// 注册跳转
func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
	//c.JSON(http.StatusOK, gin.H{
	//	"message": "welcome!!",
	//})
}

// ToChat
// 聊天主页面跳转
func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/main.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/foot.html",
	)
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	var user = models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	//fmt.Println("tochat>>>>>>>>>", user)
	ind.Execute(c.Writer, user)
	//c.JSON(http.StatusOK, gin.H{
	//	"message": "welcome!!",
	//})
}

//func Chat(c *gin.Context) {
//	models.Chat(c.Writer, c.Request)
//}
