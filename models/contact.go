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

// AddFriend
// 添加好友
func AddFriend(userId uint, targetName string) (int, string) {
	targetUser := FindUserByName(targetName)
	if targetUser.Salt != "" {
		if targetUser.ID == userId {
			return -1, "不能加自己"
		}
		cont := Contact{}
		utils.DB.Where("owner_id= ? and target_id = ? and type=1",
			userId, targetUser.ID).Find(&cont)
		if cont.ID != 0 {
			return -1, "不能重复添加"
		}
		//互相加好友
		//开启事务tx，双向控制，只有一个成功会回滚
		tx := utils.DB.Begin()
		//事务一旦开始，不论什么异常，最终都会 Rollback
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		if err := tx.Error; err != nil {
			return -1, "添加好友失败"
		}
		contact := Contact{}
		contact.OwnerId = userId
		contact.TargetId = targetUser.ID
		contact.Type = 1
		//utils.DB.Create(&contact)
		if err := tx.Create(&contact).Error; err != nil {
			tx.Rollback()
			return -1, "添加好友失败"
		}
		contact = Contact{}
		contact.OwnerId = targetUser.ID
		contact.TargetId = userId
		contact.Type = 1
		//utils.DB.Create(&contact)
		if err := tx.Create(&contact).Error; err != nil {
			tx.Rollback()
			return -1, "添加好友失败"
		}
		tx.Commit()
		return 0, "添加好友成功"
	}
	return -1, "没有找到此用户"
}

// SearchUserByGroupId
// 通过群找到群人员的方法
func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("target_id= ? and type=2", communityId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, v.OwnerId)
	}
	return objIds
}
