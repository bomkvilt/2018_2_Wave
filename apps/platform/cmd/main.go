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
		cors   = mw.Cors(s)
		authMW = mw.NewAuth(s)
		noAuth = authMW.MW(mw.ENoAuth)
		auth   = authMW.MW(mw.EAuth)
		// admin  = authMW.MW(mw.EAuth | mw.EAdmin)

		public = mw.Combine(noAuth, cors)
		common = mw.Combine(auth, cors)
		// private = mw.Combine(admin, cors)
	)

	r.HandleFunc("/users", mw.Pipe(s.SugnUp, public)).Methods("POST")
	r.HandleFunc("/users/me", mw.Pipe(s.MyProfile, common)).Methods("GET")
	r.HandleFunc("/session", mw.Pipe(s.LogIn, common)).Methods("POST")
	r.HandleFunc("/session", mw.Pipe(s.LogOut, common)).Methods("DELETE")

	r.HandleFunc("/apps", mw.Pipe(s.GetApps, common)).Methods("GET")
	r.HandleFunc("/apps/categories", mw.Pipe(s.GetAppCategories, common)).Methods("GET")
	r.HandleFunc("/apps/category/{category}", mw.Pipe(s.GetCategoryApps, common)).Methods("GET")
	r.HandleFunc("/apps/{name}", mw.Pipe(s.GetAppInfo, common)).Methods("GET")
	r.HandleFunc("/me/apps", mw.Pipe(s.GetMyApps, common)).Methods("GET")
	r.HandleFunc("/me/apps", mw.Pipe(s.AddMyApp, common)).Methods("POST")

	s.Logger().Infof("platform server has been started on a port %s", p)
	http.ListenAndServe(p, r)
}
