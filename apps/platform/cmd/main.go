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
	s.Logger().Infof("platform server has been started on a port %s", p)
	http.ListenAndServe(p, r)
}
