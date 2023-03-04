package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string //用户名
	PassWord      string //密码
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"` //电话
	Email         string `valid:"email"`                      //邮箱
	Identity      string
	ClientIP      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time //登录时间
	HeartbeatTime time.Time //心跳时间
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"` //登出时间
	IsLogout      bool
	DeviceInfo    string
	Avatar        string //头像
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}
func FindUserByNameAndPwd(name, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and pass_word= ?", name, password).First(&user)
	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id= ?", user.ID).Update("identity", temp)
	return user
}

// FindUserByName FindUserByPhone FindUserByEmail 重复注册校验
func FindUserByName(name string) UserBasic { //查找用户是否存在
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}
func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("phone = ?", phone).First(&user)
}
func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("email = ?", email).First(&user)
}
func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}
func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}
func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{
		Name:     user.Name,
		PassWord: user.PassWord,
		Phone:    user.Phone,
		Email:    user.Email,
		Avatar:   user.Avatar,
	})
}

// FindByID
// 查找某个用户
func FindByID(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id=?", id).First(&user)
	return user
}
