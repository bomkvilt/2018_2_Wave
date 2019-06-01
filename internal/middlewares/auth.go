package middlewares

import (
	"context"
	"net/http"
	"wave/apps/auth"
	"wave/apps/auth/iauth"
	"wave/internal/cookies"
	"wave/internal/service"

	"google.golang.org/grpc"
)

type EAuthFlags int

const (
	ENoAuth = 0
	EAuth   = 1 << iota
	EAdmin
)

// --------------------------|

type authKey string

const (
	authUID  authKey = "uid"
	authAuth authKey = "auth"
)

// easyjson:json
type authConfig struct {
	Adress string `json:"adress"`
}

// --------------------------|

type IAuth interface {
	MW(EAuthFlags) IMiddleware
}

type fAuth struct {
	iauth.IAuthClient
}

// -----------|

func NewAuth(sv service.IService) IAuth {
	mw := &fAuth{}
	{
		config := auth.Config{}
		service.PanicIf(sv.Config(&config, "auth"))
		con, err := grpc.Dial(
			"127.0.0.1"+config.Port,
			grpc.WithInsecure(),
		)
		if err != nil {
			sv.Logger().Errorf("tcp port %s cannot be listened: %s", config.Port, err.Error())
			service.Panic(err)
		}
		// defer con.Close()
		mw.IAuthClient = iauth.NewIAuthClient(con)
	}
	return mw
}

func (mw fAuth) MW(flags EAuthFlags) IMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ok := true
			r = mw.injectAuth(r)
			if flags&EAuth != 0 {
				if r, ok = mw.authCheck(w, r); !ok {
					return
				}
			}
			if flags&EAdmin != 0 {
				if r, ok = mw.authCheck(w, r); !ok {
					return
				}
			}
			next(w, r)
		}
	}
}

// inject an auth service to the request
func (mw fAuth) injectAuth(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), authAuth, mw.IAuthClient)
	return r.WithContext(ctx)
}

// check an authentification and store uid to the context
func (mw fAuth) authCheck(w http.ResponseWriter, r *http.Request) (*http.Request, bool) {
	cookie := mw.makeCookie(r)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		//TODO:: metrics
		return r, false
	}
	status, err := mw.IsLoggedIn(r.Context(), cookie)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		//TODO:: metrics
		return r, false
	}
	if !status.BOK {
		w.WriteHeader(http.StatusUnauthorized)
		//TODO:: metrics
		return r, false
	}
	ctx := context.WithValue(r.Context(), authUID, status.Uid)
	return r.WithContext(ctx), true
}

// check administrator permissions
func (mw fAuth) adminCheck(w http.ResponseWriter, r *http.Request) (*http.Request, bool) {
	cookie := mw.makeCookie(r)
	if cookie == nil {
		return r, false
	}
	status, err := mw.IsAdmin(r.Context(), cookie)
	if err != nil {
		return r, false
	}
	return r, status.BOK
}

// make iauth cookie
func (mw fAuth) makeCookie(r *http.Request) *iauth.Cookie {
	cookie := cookies.GetSessionCookie(r)
	if cookie == "" {
		return nil
	}
	return &iauth.Cookie{Cookie: cookie}
}

// --------------------------|

func GetAuth(r *http.Request) (iauth.IAuthClient, bool) {
	v := r.Context().Value(authAuth)
	if v == nil {
		return nil, false
	}
	a, ok := v.(iauth.IAuthClient)
	return a, ok
}

func GetUID(r *http.Request) (int64, bool) {
	v := r.Context().Value(authUID)
	if v == nil {
		return 0, false
	}
	u, ok := v.(int64)
	return u, ok
}
