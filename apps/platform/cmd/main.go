package main

import (
	"net/http"

	mw "wave/internal/middlewares"
	sv "wave/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	var (
		c = sv.LoadConfig("configs/landing.json")
		s = sv.NewService(c)
		h = func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("gg wp, mafaker!")) }
		r = mux.NewRouter()
	)
	r.HandleFunc("/", mw.Pipe(h, mw.Cors(s)))
	http.ListenAndServe(":8080", r)
}
