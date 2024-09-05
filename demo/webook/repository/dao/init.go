package dao

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	// gorm自动创建，更新数据库表结构
	// 可以使用sql转go工具帮助理解
	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	//GoRM将结构体字段映射为Sql
	result := db.Create(&User{
		Email:         sql.NullString{String: "test@example.com", Valid: true},
		Password:      "hashed_password_here", // 通常应该是已经加密过的密码
		Phone:         sql.NullString{String: "1234567890", Valid: true},
		Birthday:      sql.NullInt64{Int64: 946684800000, Valid: true}, // 生日时间戳（例如 2000-01-01）
		Nickname:      sql.NullString{String: "testuser", Valid: true},
		AboutMe:       sql.NullString{String: "Hello, I am a test user.", Valid: true},
		WechatOpenId:  sql.NullString{String: "wechat_open_id_here", Valid: true},
		WechatUnionId: sql.NullString{String: "wechat_union_id_here", Valid: true},
		Ctime:         time.Now().Unix(), // 当前时间作为创建时间
		Utime:         time.Now().Unix(), // 当前时间作为更新时间
	})

	if result.Error != nil {
		return result.Error
	}
	return nil
}
