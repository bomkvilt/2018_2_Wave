package api

import (
	psql "Wave/internal/database"
	lg "Wave/internal/logger"
	mc "Wave/internal/metrics"
	"Wave/internal/services/auth/proto"
	"Wave/internal/models"
	"Wave/internal/misc"

	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"os"
	"io"
	"log"

	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	"golang.org/x/net/context"

	"Wave/application/manager"
	"Wave/application/room"
	"Wave/application/snake"

	_ "github.com/lib/pq"
)

type Handler struct {
	DB psql.DatabaseModel
	LG *lg.Logger
	Prof *mc.Profiler
	AuthManager auth.AuthClient
}

func (h *Handler) uploadHandler(r *http.Request) (created bool, path string) {
	log.Println("hey")
	file, _, err := r.FormFile("avatar")
	log.Println("gurl")
	defer file.Close()

	if err != nil {

		h.LG.Sugar.Infow("upload failed, unable to read from FormFile, default avatar set",
		"source", "api.go",
		"who", "uploadHandler",)

        return true, "/img/avatars/default"
	}

	prefix := "/img/avatars/"
	hash := ksuid.New()
	fileName := hash.String()

	createPath := "." + prefix + fileName
	path = prefix + fileName

	out, err := os.Create(createPath)
	defer out.Close()

    if err != nil {

		h.LG.Sugar.Infow("upload failed, file couldn't be created",
		"source", "api.go",
		"who", "uploadHandler",)

        return false, ""
    }

    _, err = io.Copy(out, file)
    if err != nil {

        h.LG.Sugar.Infow("upload failed, couldn't copy data",
		"source", "api.go",
		"who", "uploadHandler",)

		return false, ""
    }

	h.LG.Sugar.Infow("upload succeded",
		"source", "api.go",
		"who", "uploadHandler",)

	return true, path
}

func (h *Handler) SlashHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	h.Prof.Hits.Add(1)
	return
}

func (h *Handler) RegisterPOSTHandler(rw http.ResponseWriter, r *http.Request) {
	user := models.UserEdit{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	isCreated, avatarPath := h.uploadHandler(r)

	if isCreated && avatarPath != "" {
		user.Avatar = avatarPath
	} else if !isCreated {
		fr := models.ForbiddenRequest{
			Reason: "Bad avatar.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users failed, bad avatar.",
		"source", "api.go",
		"who", "RegisterPOSTHandler",)

		return
	}

	cookie, err := h.AuthManager.Create(
			context.Background(),
			&auth.Credentials{
			Username: user.Username,
			Password: user.Username,
			Avatar: user.Avatar,
		})

	if err != nil {

		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/users failed",
		"source", "api.go",
		"who", "RegisterPOSTHandler",)

		return
	}
	// or validation failed
	if cookie.CookieValue == "" {
		fr := models.ForbiddenRequest{
			Reason: "Username already in use.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users failed, username already in use.",
		"source", "api.go",
		"who", "RegisterPOSTHandler",)

		return
	}

	sessionCookie := misc.MakeSessionCookie(cookie.CookieValue)
	http.SetCookie(rw, sessionCookie)
	rw.WriteHeader(http.StatusCreated)

	h.LG.Sugar.Infow("/users succeded",
		"source", "api.go",
		"who", "RegisterPOSTHandler",)

	return
}

func (h *Handler) MeGETHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)

	profile, err := h.DB.GetMyProfile(cookie)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := profile.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/users/me succeded",
		"source", "api.go",
		"who", "MeGETHandler",)

	return
}

func (h *Handler) EditMePUTHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)

	user := models.UserEdit{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	isCreated, avatarPath := h.uploadHandler(r)

	if isCreated && avatarPath != "" {
		user.Avatar = avatarPath
	} else if !isCreated {
		fr := models.ForbiddenRequest{
			Reason: "Update didn't happend, shitty avatar.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users/me failed, bad avatar.",
		"source", "api.go",
		"who", "EditMePUTHandler",)

		return
	}

	_, err := h.DB.UpdateProfile(user, cookie)

	if err != nil {
		fr := models.ForbiddenRequest{
			Reason: "Update didn't happend, shitty username and/or password.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/users/me failed",
		"source", "api.go",
		"who", "EditMePUTHandler",)

		return
	}

	h.LG.Sugar.Infow("/users/me succeded, user profile updated",
	"source", "api.go",
	"who", "EditMePUTHandler",)

	rw.WriteHeader(http.StatusOK)

	return
}

func (h *Handler) UserGETHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profile, err := h.DB.GetProfile(vars["name"])

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/users/{name} failed",
		"source", "api.go",
		"who", "UserGETHandler",)

		return
	}

	if reflect.DeepEqual(models.UserExtended{}, profile) {
		rw.WriteHeader(http.StatusNotFound)

		h.LG.Sugar.Infow("/users/{name} failed",
		"source", "api.go",
		"who", "UserGETHandler",)

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := profile.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/users/{name} succeded",
		"source", "api.go",
		"who", "UserGETHandler",)

	return
}

func (h *Handler) LeadersGETHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, _ := strconv.Atoi(vars["count"])
	p, _ := strconv.Atoi(vars["page"])
	leaders, err := h.DB.GetTopUsers(c, p)

	if err != nil || reflect.DeepEqual(models.Leaders{}, leaders) {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/users/leaders failed",
		"source", "api.go",
		"who", "LeadersGETHandler",)

		return
	}

	rw.WriteHeader(http.StatusOK)
	payload, _ := leaders.MarshalJSON()
	fmt.Fprintln(rw, string(payload))

	h.LG.Sugar.Infow("/users/leaders succeded",
	"source", "api.go",
	"who", "LeadersGETHandler",)

	return
}

func (h *Handler) LoginPOSTHandler(rw http.ResponseWriter, r *http.Request) {
	user := models.UserCredentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	cookie, err := h.DB.LogIn(user)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/session failed",
		"source", "api.go",
		"who", "LoginPOSTHandler",)

		return
	}

	if cookie == "" {
		fr := models.ForbiddenRequest{
			Reason: "Incorrect password/username.",
		}

		payload, _ := fr.MarshalJSON()
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(rw, string(payload))

		h.LG.Sugar.Infow("/session failed, incorrect password/username",
		"source", "api.go",
		"who", "LoginPOSTHandler",)

		return
	}

	sessionCookie := misc.MakeSessionCookie(cookie)
	http.SetCookie(rw, sessionCookie)
	rw.WriteHeader(http.StatusOK)

	h.LG.Sugar.Infow("/session succeded",
		"source", "api.go",
		"who", "LoginPOSTHandler",)

	return
}

func (h *Handler) LogoutDELETEHandler(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)
	success, _ := h.AuthManager.Delete(
		context.Background(),
		&auth.Cookie{
			CookieValue: cookie,
		})

	if success.Resp != true {
		rw.WriteHeader(http.StatusInternalServerError)

		h.LG.Sugar.Infow("/session failed",
		"source", "api.go",
		"who", "LogoutDELETEHandler",)

		return
	}

	http.SetCookie(rw, misc.MakeSessionCookie(""))
	rw.WriteHeader(http.StatusOK)

	h.LG.Sugar.Infow("/session succeded",
		"source", "api.go",
		"who", "LogoutDELETEHandler",)

	return
}

func (h *Handler) EditMeOPTHandler(rw http.ResponseWriter, r *http.Request) {

	h.LG.Sugar.Infow("/users/me succeded",
		"source", "api.go",
		"who", "EditMeOPTHandler",)

}

func (h *Handler) LogoutOPTHandler(rw http.ResponseWriter, r *http.Request) {

	h.LG.Sugar.Infow("/session succeded",
		"source", "api.go",
		"who", "LogoutOPTHandler",)

<<<<<<< HEAD
}

func (h *Handler) IsLoggedIn(rw http.ResponseWriter, r *http.Request) {
	cookie := misc.GetSessionCookie(r)
	success, _ := h.SessManager.Check(
		context.Background(),
		&session.SessionID{
			ID: cookie,
		})
	if success.Dummy == true {
		rw.WriteHeader(http.StatusOK)

		return
	}

	rw.WriteHeader(http.StatusUnauthorized)
}

/************************* websocket block ************************************/

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	waitTime = 15 * time.Second
)

//action_uid uniqely generated on the front
//action_id : 
// 1 - add user to the room
// 2 - remove user from the room
// 3 - start
// 4 - rollback

type lobbyReq struct {
	actionID 	string `json:"action_id"`
	actionUID 	string `json:"action_uid"`
	username 	string `json:"username"`
}

type lobbyRespGenereic struct {
	actionUID string `json:"action_id"`
	status 	 string `json:"status"`
}

func contains(sl []string, str string) bool {
    for _, cur := range sl {
        if str == cur {
            return true
        }
    }
    return false
}

func (h *Handler) LobbyHandler(rw http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()
/*
	lobby := []string{}

	go func(client *websocket.Conn, lb []string){
		ticker := time.NewTicker(waitTime)
		defer func() {
			ticker.Stop()
			client.Close()
		}()
			for {
					in := lobbyReq{}
				
					err := client.ReadJSON(&in)
					if err != nil {
						break
					}
			
					fmt.Printf("Got message: %#v\n", in)

					out := lobbyRespGenereic{}

					switch in.actionID {
						case "1": 
							if in.username == "" {
								break
							}
							out.actionUID = in.actionUID
							lb = append(lb, in.username)
							out.status = "success" 

							if err = client.WriteJSON(out); err != nil {
								break
							}
		
						case "2":
							if in.username == "" {
								break
							}
							out.actionUID = in.actionUID
							if contains(lb, in.username) {
								for _, cur := range lb {
									if cur == in.username {
										cur = ""
									}
								}
								out.status = "success"
								if err = client.WriteJSON(out); err != nil {
									break
								}
							} else {
								out.status = "failure"
								if err = client.WriteJSON(out); err != nil {
									break
								}			
							}
						case "3":
							out.actionUID = in.actionUID
							out.status = "success"
							if err = client.WriteJSON(out); err != nil {
								break
							}		
							
						case "4":
							out.actionUID = in.actionUID
							out.status = "success"
							if err = client.WriteJSON(out); err != nil {
								break
							}		
					}

					<-ticker.C
		}
	}(ws, lobby)
	*/
/*
	go func(client *websocket.Conn) {
	}(ws)
	*/
	log.Println("ws@")
	return
}


// ----------------| ws

// TODO:: get the value from configuration files
const wsAppTickRate = 500 * time.Millisecond

// TODO:: place the variables into the handler struct
// TODO:: add function-constructor for the handler
// TODO:: initialize the variables in the constructor
// TODO:: pass throw a db connection instead of nil
var wsApp *app.App = func() *app.App {
	wsApp := app.New("app", wsAppTickRate, nil)
	wsApp.CreateLobby(snake.RoomType, "snake")
	go wsApp.Run()
	return wsApp
}()
var wsUpgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) WSHandler(rw http.ResponseWriter, r *http.Request) {
	ws, err := wsUpgrader.Upgrade(rw, r, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		user := room.NewUser(wsApp.GetNextUserID(), ws)
		user.LG = h.LG
		user.AddToRoom(wsApp)
		user.Listen()
	}()
=======
>>>>>>> deploy
}