package user

import (
	"context"
	"errors"

	"github.com/marcelo-rocha/task-service-challenge/domain"
	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
)

type UserInfo struct {
	Id   int64
	Kind entities.UserKind
}

var ErrContextWithoutUserInfo = errors.New("context without user info")

func GetAuthenticatedUserInfo(ctx context.Context) (*UserInfo, error) {
	v := ctx.Value(domain.UserInfoKey)
	info, ok := v.(*UserInfo)
	if !ok {
		return nil, ErrContextWithoutUserInfo
	}
	return info, nil
}

func ContextWithUserInfo(ctx context.Context, userId int64, kind entities.UserKind) context.Context {
	return context.WithValue(ctx, domain.UserInfoKey, &UserInfo{Id: userId, Kind: kind})
}
