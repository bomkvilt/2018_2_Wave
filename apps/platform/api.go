package platform

import (
	"wave/apps/auth"
	"wave/apps/auth/iauth"
	"wave/internal/service"

	"google.golang.org/grpc"
)

//go:generate easyjson -output_filename ./models/jsons.gen.go ./models
//go:generate easyjson -output_filename jsons.gen.go .

//~~~~~~~~~~~~~~~~~~~~~~| API

// API - platform api structure
type API struct {
	service.IService

	db     *db
	conf   Config
	authCf auth.Config
	authMs iauth.IAuthClient
}

// NewAPI - cretae a new API
func NewAPI(config service.FServiceConfig) *API {
	api := &API{IService: service.NewService(config)}
	service.PanicIf(api.Config(&api.conf, "platform"))
	service.PanicIf(api.Config(&api.authCf, "auth"))
	api.db = newDB(api.Logger(), api.conf)

	// connect to an auth microservice
	{
		con, err := grpc.Dial(
			"127.0.0.1"+api.authCf.Port,
			grpc.WithInsecure(),
		)
		if err != nil {
			api.Logger().Errorf("tcp port %s cannot be rised: %s", api.authCf.Port, err.Error())
			service.Panic(err)
		}
		// defer con.Close()
		api.authMs = iauth.NewIAuthClient(con)
	}
	return api
}

func (ap API) GetPort() string {
	return ap.conf.Port
}
