package platform

import "wave/internal/service"

//~~~~~~~~~~~~~~~~~~~~~~| API

// API - platform api structure
type API struct {
	sv service.IService
}

// NewAPI - cretae a new API
func NewAPI(sv service.IService) *API {
	api := &API{}
	return api
}
