package auth

import (
	"context"
	"wave/apps/auth/iauth"
	"wave/internal/service"
)

//go:generate protoc --go_out=plugins=grpc:. ./iauth/auth.proto
//go:generate easyjson -output_filename jsons.gen.go .

// --------------------------| authService

// IAuthServer interface
type IAuthServer interface {
	iauth.IAuthServer
	service.IService
	GetPort() string
}

type authService struct {
	service.IService
	cf Config
	db *db
}

// NewAuthService - create a new auth service
func NewAuthService(config service.FServiceConfig) IAuthServer {
	au := &authService{IService: service.NewService(config)}
	service.PanicIf(au.Config(&au.cf, "auth"))
	au.db = newDB(au.Logger(), au.cf)
	return au
}

func (au *authService) GetPort() string {
	return au.cf.Port
}

// --------------------------|

func (au *authService) CreateAccaunt(_ context.Context, u *iauth.User) (*iauth.Cookie, error) {
	cookie, uid, err := au.db.SignUp(u.Username, u.Password)
	if err != nil {
		return cookieFail, nil
	}
	return &iauth.Cookie{
		Status: &iauth.Status{
			BOK: true,
			Uid: uid,
		},
		Cookie: cookie,
	}, nil
}

func (au *authService) LogIn(_ context.Context, u *iauth.User) (*iauth.Cookie, error) {
	cookie, uid, err := au.db.LogIn(u.Username, u.Password)
	if err != nil {
		return cookieFail, nil
	}
	return &iauth.Cookie{
		Status: &iauth.Status{
			BOK: true,
			Uid: uid,
		},
		Cookie: cookie,
	}, nil
}

func (au *authService) LogOut(_ context.Context, c *iauth.Cookie) (*iauth.Status, error) {
	err := au.db.LogOut(c.Cookie)
	if err != nil {
		return statusFail, nil
	}
	return statusOK, nil
}

func (au *authService) IsLoggedIn(_ context.Context, c *iauth.Cookie) (*iauth.Status, error) {
	uid, err := au.db.GetSession(c.Cookie)
	if err != nil {
		return statusFail, nil
	}
	return &iauth.Status{
		BOK: true,
		Uid: uid,
	}, nil
}

// --------------------------|

var (
	statusOK   = &iauth.Status{BOK: true}
	statusFail = &iauth.Status{BOK: false}
	cookieFail = &iauth.Cookie{Status: statusFail}
)
