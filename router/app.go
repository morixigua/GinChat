package router

import (
	"ginchat/docs"
	"ginchat/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	//swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//静态资源
	r.Static("/asset", "asset/")
	r.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	r.LoadHTMLGlob("views/**/*")
	//首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)

	//用户主界面
	//r.GET("/chat", service.Chat)
	r.POST("/searchFriends", service.SearchFriends)

	//用户模块 增删改查
	//r.GET("/user/getUserList", service.GetUserList)
	//r.GET("/user/createUser", service.CreateUser)
	//r.GET("/user/deleteUser", service.DeleteUser)
	r.POST("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)

	//发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	//上传文件
	r.POST("/attach/upload", service.Upload)
	//添加好友
	r.POST("/contact/addFriend", service.AddFriend)
	//创建群
	r.POST("/contact/createCommunity", service.CreateCommunity)
	//群列表
	r.POST("/contact/loadCommunity", service.LoadCommunity)
	r.POST("/contact/joinGroup", service.JoinGroup)
	//心跳续命 不合适 因为Node 所以前端发过来的消息再receProc里面处理
	//r.POST("/user/heartbeat", service.Heartbeat)

	r.POST("/user/redisMsg", service.RedisMsg)
	r.POST("/user/find", service.FindByID)

	return r
}
