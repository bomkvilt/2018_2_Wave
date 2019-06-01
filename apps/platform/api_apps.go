package platform

import (
	"net/http"
	"wave/apps/platform/models"

	"github.com/gorilla/mux"
)

// easyjson:json
type appsCrutch struct {
	Apps models.Apps `json:"apps"`
}

//easyjson:json
type ctgsCrutch struct {
	Categories models.Categories `json:"categories"`
}

func (ap API) GetApps(rw http.ResponseWriter, r *http.Request) {
	apps, err := ap.db.GetApps()
	if err != nil {
		ap.Logger().Errorf("unable to get apps: %s", err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	data, err := appsCrutch{apps}.MarshalJSON()
	if err != nil {
		ap.Logger().Errorf("serialisation fail: %s", err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}

func (ap API) GetAppCategories(rw http.ResponseWriter, r *http.Request) {
	ctgs, err := ap.db.GetAppCategories()
	if err != nil {
		ap.Logger().Errorf("unable to get app categories: %s", err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	data, err := ctgsCrutch{ctgs}.MarshalJSON()
	if err != nil {
		ap.Logger().Errorf("serialisation fail: %s", err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}

func (ap API) GetMyApps(rw http.ResponseWriter, r *http.Request) {
	apps, err := ap.db.GetUserApps(ap.getUID(r))
	if err != nil {
		ap.Logger().Errorf("unable to get apps: %s", err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	data, err := appsCrutch{apps}.MarshalJSON()
	if err != nil {
		ap.Logger().Errorf("serialisation fail: %s", err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}

func (ap API) AddMyApp(rw http.ResponseWriter, r *http.Request) {
	appName := r.FormValue("name")
	if err := ap.db.AddMyApp(ap.getUID(r), appName); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (ap API) GetCategoryApps(rw http.ResponseWriter, r *http.Request) {
	apps, err := ap.db.GetCategoryApps(mux.Vars(r)["category"])
	if err != nil {
		ap.Logger().Errorf("unable to get apps: %s", err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	data, err := appsCrutch{apps}.MarshalJSON()
	if err != nil {
		ap.Logger().Errorf("serialisation fail: %s", err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}

func (ap API) GetAppInfo(rw http.ResponseWriter, r *http.Request) {
	app, err := ap.db.GetApp(mux.Vars(r)["name"])
	if err != nil {
		ap.Logger().Errorf("unable to get an app: %s", err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	data, err := app.MarshalJSON()
	if err != nil {
		ap.Logger().Errorf("serialisation fail: %s", err.Error())
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}
