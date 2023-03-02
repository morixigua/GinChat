package models

import (
	"ginchat/utils"
	"gorm.io/gorm"
)

// Contact 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系信息
	TargetId uint //对应的谁
	Type     int  //对应类型 1好友 2群组 3
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

// SearchFriend
// 查询好友列表
func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("owner_id=? and type=1", userId).Find(&contacts)
	for _, v := range contacts {
		//fmt.Println(v)
		objIds = append(objIds, v.TargetId)
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	//fmt.Println(users)
	return users
}

func AddFriend(userId uint, targetId uint) int {
	user := UserBasic{}
	if targetId != 0 {
		user = FindByID(targetId)
		if user.Salt != "" {
			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetId = targetId
			contact.Type = 1
			utils.DB.Create(&contact)
			//互相加
			contact = Contact{}
			contact.OwnerId = targetId
			contact.TargetId = userId
			contact.Type = 1
			utils.DB.Create(&contact)
			return 0
		}
	}
	return -1
}
