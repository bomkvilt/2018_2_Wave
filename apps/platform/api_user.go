package platform

import (
	"net/http"
	"wave/apps/auth/iauth"
	"wave/apps/platform/models"
	"wave/internal/cookies"
	"wave/internal/service"
)

func (ap API) SugnUp(rw http.ResponseWriter, r *http.Request) {
	var (
		user = models.UserProfile{
			Username: r.FormValue("username"),
			Avatar:   r.FormValue("avatar"),
		}
		password = r.FormValue("password")
	)

	cookie, err := ap.getAuth(r).CreateAccaunt(r.Context(), &iauth.User{
		Username: user.Username,
		Password: password,
	})
	if err != nil {
		ap.Logger().Infof("Unexpected error with auth service: %s", err.Error())
		return
	}

	if cookie.Status.BOK {
		user.UID = cookie.Status.Uid
		if err := ap.db.CreatePrifile(user); err != nil {
			ap.Logger().Errorf("Error during user creation: %s", err.Error())
			service.Panic(err)
		}
		session := cookies.MakeSessionCookie(cookie.Cookie)
		cookies.SetCookie(rw, session)
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusForbidden)
	}
}

func (ap API) LogIn(rw http.ResponseWriter, r *http.Request) {
	var (
		username = r.FormValue("username")
		password = r.FormValue("password")
	)
	cookie, err := ap.getAuth(r).LogIn(r.Context(), &iauth.User{
		Username: username,
		Password: password,
	})
	if err != nil {
		ap.Logger().Infof("Unexpected error with auth service: %s", err.Error())
		return
	}

	if cookie.Status.BOK {
		session := cookies.MakeSessionCookie(cookie.Cookie)
		cookies.SetCookie(rw, session)
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusForbidden)
	}
}

func (ap API) LogOut(rw http.ResponseWriter, r *http.Request) {
	cookie := cookies.GetSessionCookie(r)
	_, err := ap.getAuth(r).LogOut(r.Context(), &iauth.Cookie{Cookie: cookie})
	if err != nil {
		rw.WriteHeader(http.StatusForbidden)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}

func (ap API) MyProfile(rw http.ResponseWriter, r *http.Request) {
	user, err := ap.db.GetProfile(ap.getUID(r))
	payload, err := user.MarshalJSON()
	if err != nil {
		ap.Logger().Errorf("Unexpected error during marshaling: %s", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(payload)
}
