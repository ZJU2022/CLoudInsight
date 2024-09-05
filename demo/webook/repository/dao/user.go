package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// ErrDataNotFound 通用的数据没找到
var ErrDataNotFound = gorm.ErrRecordNotFound

// ErrUserDuplicate 这个算是 user 专属的
var ErrUserDuplicate = errors.New("用户邮箱或者手机号冲突")

//go:generate mockgen -source=./user.go -package=daomocks -destination=mocks/user.mock.go UserDAO
type UserDAO interface {
	Insert(ctx context.Context, u User) error
	//UpdateNonZeroFields(ctx context.Context, u User) error
	//FindByPhone(ctx context.Context, phone string) (User, error)
	//FindByEmail(ctx context.Context, email string) (User, error)
	//FindByWechat(ctx context.Context, openId string) (User, error)
	//FindById(ctx context.Context, id int64) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewGORMUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

// dao中操作的不是domain.User，因为domain.User是业务概念，不一定和数据库中的表或者列完全对应的上，而dao.User则是直接映射到表里面的
func (ud *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := ud.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrNo uint16 = 1062
		if me.Number == uniqueIndexErrNo {
			//尽可能给前端返回准确的错误信息
			return ErrUserDuplicate
		}
	}
	return err
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 设置为唯一索引
	Email    sql.NullString `gorm:"unique"`
	Password string

	//Phone *string
	Phone sql.NullString `gorm:"unique"`

	// 这三个字段表达为 sql.NullXXX 的意思，
	// 就是希望使用的人直到，这些字段在数据库中是可以为 NULL 的
	// 这种做法好处是看到这个定义就知道数据库中可以为 NULL，坏处就是用起来没那么方便
	// 大部分公司不推荐使用 NULL 的列
	// 所以你也可以直接使用 string, int64，那么对应的意思是零值就是每填写
	// 这种做法的好处是用起来好用，但是看代码的话要小心空字符串的问题
	// 生日。一样是毫秒数
	Birthday sql.NullInt64
	// 昵称
	Nickname sql.NullString
	// 自我介绍
	// 指定是 varchar 这个类型，并且长度是 1024
	// 因此你可以看到在 web 里面有这个校验
	AboutMe sql.NullString `gorm:"type:varchar(1024)"`

	// 微信有关数据。有些公司会尝试把这些数据分离出去做一个单独的表
	// 从而避免这个表有过多的列，但是暂时来说
	// 我们还没到这个地步
	WechatOpenId  sql.NullString `gorm:"type:varchar(256);unique"`
	WechatUnionId sql.NullString `gorm:"type:varchar(256)"`

	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
