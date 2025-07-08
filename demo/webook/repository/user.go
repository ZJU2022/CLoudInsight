package repository

import (
	"CloudInsight/demo/webook/domain"
	"CloudInsight/demo/webook/repository/dao"
	"context"
	"database/sql"
)

// 别名机制，向上返回，
var ErrUserDuplicate = dao.ErrUserDuplicate
var ErrUserNotFound = dao.ErrDataNotFound

//go:generate mockgen -source=./user.go -package=repomocks -destination=mocks/user.mock.go UserRepository
type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	// Update 更新数据，只有非 0 值才会更新
	Update(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	// FindByWechat 暂时可以认为按照 openId来查询
	// 将来可能需要按照 unionId 来查询
	FindByWechat(ctx context.Context, openId string) (domain.User, error)
}

// CachedUserRepository 使用了缓存的 repository 实现
type CachedUserRepository struct {
	dao dao.UserDAO
}

// func NewCachedUserRepository() UserRepository {}构造函数

// 抽象存储，代表数据在程序中的存储
func (ur *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	// 将业务对象（domain.User）转换为数据对象（dao.User），并调用 dao 层的方法进行数据库操作
	return ur.dao.Insert(ctx, dao.User{
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		WechatUnionId: sql.NullString{
			String: u.WechatInfo.UnionId,
			Valid:  u.WechatInfo.UnionId != "",
		},
		WechatOpenId: sql.NullString{
			String: u.WechatInfo.OpenId,
			Valid:  u.WechatInfo.OpenId != "",
		},
	})
}
