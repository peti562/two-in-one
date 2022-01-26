package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	Id      uint32 `gorm:"column:c_id;primary_key:true" json:"id"`
	Body    string `gorm:"column:c_body" json:"body"`
	Deleted bool   `gorm:"column:c_deleted" json:"deleted"`
	UserId  uint32 `gorm:"column:fk_user_id" json:"userId"`
}

func (comment *Comment) TableName() string {
	return "comments"
}

func (comment *Comment) FindById(gormDb *gorm.DB, commentId uint32) error {
	return gormDb.Model(&comment).
		Find(&comment, commentId).
		Where("c_deleted", false).
		Error
}

func (comment *Comment) GetByUserId(gormDb *gorm.DB, userId uint32) ([]*Comment, error) {
	var comments []*Comment
	exception := gormDb.Model(&comment).
		Where("fk_user_id", userId).
		Where("c_deleted", false).
		Find(&comments).Error

	return comments, exception
}

func (comment *Comment) UpdateBody(gormDb *gorm.DB, commentId uint32, body string) error {
	return gormDb.Model(&comment).
		Limit(1).
		Where("c_id", commentId).
		Where("c_deleted", false).
		Update("c_body", body).
		Error
}

func (comment *Comment) Delete(gormDb *gorm.DB, commentId uint32) error {
	return gormDb.Model(&comment).
		Limit(1).
		Where("c_id", commentId).
		Update("c_deleted", true).
		Error
}
