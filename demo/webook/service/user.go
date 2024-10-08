package service

import (
	"CloudInsight/demo/webook/domain"
	"CloudInsight/demo/webook/repository"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// 别名机制，向上返回，避免跨层依赖
var (
	ErrDuplicateEmail        = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

//go:generate mockgen -source=./user.go -package=svcmocks -destination=./mocks/user.mock.go UserService
type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	//Login(ctx context.Context, email string, password string) (domain.User, error)
	//UpdateNonSensitiveInfo(ctx context.Context,user domain.User) error
	//FindById(ctx context.Context, uid int64) (domain.User, error)
	//FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	//FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
	//logger *zap.Logger
}

// 构造函数
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// 用户注册的核心业务逻辑，包括密码的安全加密和用户数据的持久化存储，确保数据的正确性和安全性
func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}
