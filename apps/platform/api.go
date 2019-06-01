package platform

import (
	"errors"
	"net/http"
	"wave/apps/auth/iauth"
	"wave/internal/middlewares"
	"wave/internal/service"
)

//go:generate easyjson -output_filename ./models/jsons.gen.go ./models
//go:generate easyjson -output_filename jsons.gen.go .

//~~~~~~~~~~~~~~~~~~~~~~| API

// API - platform api structure
type API struct {
	service.IService

	conf Config
	db   *db
}

// NewAPI - cretae a new API
func NewAPI(config service.FServiceConfig) *API {
	api := &API{IService: service.NewService(config)}
	service.PanicIf(api.Config(&api.conf, "platform"))
	api.db = newDB(api.Logger(), api.conf)
	return api
}

func (ap API) GetPort() string {
	return ap.conf.Port
}

func (ap API) getAuth(r *http.Request) iauth.IAuthClient {
	ms, ok := middlewares.GetAuth(r)
	if !ok {
		service.Panic(errors.New("Cannot get auth service"))
	}
	return ms
}

func (ap API) getUID(r *http.Request) int64 {
	uid, ok := middlewares.GetUID(r)
	if !ok {
		service.Panic(errors.New("Cannot get user id"))
	}
	return uid
}
