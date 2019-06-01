package main

import (
	"net/http"

	ap "wave/apps/platform"
	mw "wave/internal/middlewares"
	sv "wave/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	var ( // common
		c = sv.LoadConfig("configs/ms_platform.json")
		s = ap.NewAPI(c)
		r = mux.NewRouter()
		p = s.GetPort()
	)
	var ( // middlewares
		cors = mw.Cors(s)
	)

	r.HandleFunc("/users", mw.Pipe(s.SugnUp, cors)).Methods("POST")
	r.HandleFunc("/users/me", mw.Pipe(s.MyProfile, cors)).Methods("GET")
	r.HandleFunc("/session", mw.Pipe(s.LogIn, cors)).Methods("POST")
	r.HandleFunc("/session", mw.Pipe(s.LogOut, cors)).Methods("DELETE")

	r.HandleFunc("/apps", mw.Pipe(s.GetApps, cors)).Methods("GET")
	r.HandleFunc("/apps/categories", mw.Pipe(s.GetAppCategories, cors)).Methods("GET")
	r.HandleFunc("/apps/category/{category}", mw.Pipe(s.GetCategoryApps, cors)).Methods("GET")
	r.HandleFunc("/apps/{name}", mw.Pipe(s.GetAppInfo, cors)).Methods("GET")
	r.HandleFunc("/me/apps", mw.Pipe(s.GetMyApps, cors)).Methods("GET")
	r.HandleFunc("/me/apps", mw.Pipe(s.AddMyApp, cors)).Methods("POST")

	s.Logger().Infof("platform server has been started on a port %s", p)
	http.ListenAndServe(p, r)
}
