package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	//user.Name = c.Query("name")
	//password := c.Query("password")
	//repassword := c.Query("repassword")
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")
	if user.Name == "" || password == "" {
		c.JSON(-1, gin.H{
			"code":    -1, //0成功,-1失败
			"message": "用户名和密码不能为空!",
			"data":    user,
		})
		return
	}
	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1, //0成功,-1失败
			"message": "两次密码不一致!",
			"data":    user,
		})
		return
	}

	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1, //0成功,-1失败
			"message": "用户名已注册!",
			"data":    data,
		})
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	salt := fmt.Sprintf("%010d", r.Int31())            //加随机数
	user.PassWord = utils.MakePassword(password, salt) //加密
	user.Salt = salt
	//user.PassWord = password
	models.CreateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功,-1失败
		"message": "新增用户成功!",
		"data":    data,
	})
}

// FindUserByNameAndPwd
// @Summary 登录
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	//name := c.Query("name")
	//password := c.Query("password")
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功,-1失败
			"message": "该用户不存在",
			"data":    user,
		})
		return
	}
	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功,-1失败
			"message": "密码不正确",
			"data":    user,
		})
		return
	}
	data = models.FindUserByNameAndPwd(name, user.PassWord)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功,-1失败
		"message": "登录成功",
		"data":    data,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功,-1失败
		"message": "删除用户成功",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	fmt.Println("update:", user)
	_, err := govalidator.ValidateStruct(user) //govalidator校验电话号码和邮箱
	if err != nil {
		fmt.Println(err)
		c.JSON(-1, gin.H{
			"code":    -1, //0成功,-1失败
			"message": "修改参数不匹配!",
			"data":    user,
		})
		return
	}
	models.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功,-1失败
		"message": "修改用户成功",
		"data":    user,
	})
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(c, ws)
}
func MsgHandler(c *gin.Context, ws *websocket.Conn) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println("MsgHandler 发送失败:", err)
		}
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	//id, err := strconv.Atoi(c.PostForm("userId"))//body的form-data里查
	id, err := strconv.Atoi(c.Request.FormValue("userId"))
	if err != nil {
		log.Fatal(err)
	}
	users := models.SearchFriend(uint(id))
	utils.RespOKList(c.Writer, users, len(users))
	//c.JSON(http.StatusOK, gin.H{
	//	"code":    0, //0成功,-1失败
	//	"message": "查询好友列表成功!",
	//	"data":    users,
	//})
}

func AddFriend(c *gin.Context) {
	userId, err := strconv.Atoi(c.Request.FormValue("userId"))
	if err != nil {
		log.Fatal(err)
	}
	targetId, err := strconv.Atoi(c.Request.FormValue("targetId"))
	if err != nil {
		log.Fatal(err)
	}
	code := models.AddFriend(uint(userId), uint(targetId))
	if code == 0 {
		utils.RespOK(c.Writer, code, "添加好友成功")
	} else {
		utils.RespFail(c.Writer, "添加好友失败")
	}
}
