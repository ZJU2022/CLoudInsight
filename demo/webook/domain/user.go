package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string
	Name     string
	Nickname string
	// YYYY-MM-DD
	Birthday time.Time
	AboutMe  string

	Phone string

	// UTC 0 的时区
	Ctime time.Time

	WechatInfo WechatInfo

	//Addr Address
}

// WechatInfo 微信的授权信息
type WechatInfo struct {
	// OpenId 是应用内唯一
	OpenId string
	// UnionId 是整个公司账号内唯一
	UnionId string
}

/*
type Product struct {
	ID    uint
	Code  string
	Price uint
}*/

// TodayIsBirthday 判定今天是不是我的生日
func (u User) TodayIsBirthday() bool {
	now := time.Now()
	return now.Month() == u.Birthday.Month() && now.Day() == u.Birthday.Day()
}

//type Address struct {
//	Province string
//	Region   string
//}

//func (u User) ValidateEmail() bool {
// 在这里用正则表达式校验
//return u.Email
//}
