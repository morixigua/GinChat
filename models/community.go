package models

import (
	"ginchat/utils"
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}

func (table *Community) TableName() string {
	return "community"
}

func CreateCommunity(community Community) (int, string) {
	tx := utils.DB.Begin()
	//事务一旦开始，不论什么异常，最终都会 Rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if len(community.Name) == 0 {
		return -1, "群名称不能为空"
	}
	if community.OwnerId == 0 {
		return -1, "环境异常，请先登录"
	}
	if err := tx.Create(&community).Error; err != nil {
		tx.Rollback()
		return -1, "建群失败"
	}
	contact := Contact{}
	contact.OwnerId = community.OwnerId
	contact.TargetId = community.ID
	contact.Type = 2
	if err := tx.Create(&contact).Error; err != nil {
		tx.Rollback()
		return -1, "添加关系失败"
	}
	tx.Commit()
	return 0, "建群成功"
}

func LoadCommunity(ownerId uint) ([]*Community, string) {
	contact := make([]Contact, 0)
	contactIds := make([]uint, 0)
	utils.DB.Where("owner_id= ? and type=2", ownerId).Find(&contact)
	for _, v := range contact {
		contactIds = append(contactIds, v.TargetId)
	}
	data := make([]*Community, 0)
	utils.DB.Where("id in ?", contactIds).Find(&data)
	//for _, v := range data {
	//	fmt.Println(v)
	//}
	return data, "查询成功"
}

func JoinGroup(userId uint, comId string) (int, string) {
	contact := Contact{}
	community := Community{}

	utils.DB.Where("id=? or name=?", comId, comId).Find(&community)
	if community.Name == "" {
		return -1, "没有找到群"
	}
	utils.DB.Where("owner_id= ? and target_id=? and type =2", userId,
		community.ID).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1, "已加过此群"
	} else {
		contact = Contact{}
		contact.OwnerId = userId
		contact.TargetId = community.ID
		contact.Type = 2
		utils.DB.Create(&contact)
		return 0, "加群成功"
	}
}
