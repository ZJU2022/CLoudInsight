package webook

import "context"

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context,
		biz, phone, inputCode string) (bool, error)
}
