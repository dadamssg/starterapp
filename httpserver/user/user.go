package user

import (
	"errors"
	"github.com/dadamssg/commandbus"
	"github.com/dadamssg/jsonhttp"
	"github.com/dadamssg/starterapp/app"
	commandpkg "github.com/dadamssg/starterapp/app/command"
	"github.com/dadamssg/starterapp/app/user"
	httpcommandpkg "github.com/dadamssg/starterapp/httpserver/command"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type TokenCommand interface {
	AccessToken() *user.Token
	RefreshToken() *user.Token
}

func Connect(app *starterapp.App, router *mux.Router) {
	router.HandleFunc("/register", registerUser(app.CommandBus)).Methods("POST")
	router.HandleFunc("/users/{id}", findUserById(app.CommandBus)).Methods("GET")

	router.HandleFunc("/authenticate", issueAuthToken(app.CommandBus)).Methods("POST")
}

func userRepresentation(user *user.User) interface{} {
	u := make(map[string]interface{})
	u["id"] = user.Id
	u["created_at"] = user.CreatedAt.Unix()
	u["username"] = user.Username
	u["enabled"] = user.Enabled
	u["email"] = user.Email

	return struct {
		User map[string]interface{} `json:"user"`
	}{
		u,
	}
}

type registerUserRequest struct {
	Command user.RegisterUserCommand `json:"user"`
}

type issueAuthTokenRequest struct {
	Command user.IssueAccessTokenCommand `json:"authenticate"`
}

func findUserById(bus *commandbus.CommandBus) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		command := user.FindUserByIdCommand{
			Id: id,
		}

		bus.Handle(&command)

		if command.HasErrors() {
			httpcommandpkg.SendCommandErrors(w, &command)
			return
		}

		if command.User == nil {
			jsonhttp.SendError(w, 404, errors.New("No user found."))
			return
		}

		jsonhttp.SendJSON(w, 200, userRepresentation(command.User))
	}
}

func registerUser(bus *commandbus.CommandBus) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request registerUserRequest

		if err := jsonhttp.MapOrSendError(w, r, &request); err != nil {
			return
		}

		command := request.Command

		bus.Handle(&command)

		if command.HasErrors() {
			httpcommandpkg.SendCommandErrors(w, &command)
			return
		}

		jsonhttp.SendJSON(w, 201, userRepresentation(command.User))
	}
}

func issueAuthToken(bus *commandbus.CommandBus) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		grantType := r.FormValue("grant_type")

		var cmd TokenCommand

		if grantType == "password" {
			cmd = &user.IssueAccessTokenCommand{
				Username: r.FormValue("username"),
				Password: r.FormValue("password"),
			}
		}

		if grantType == "refresh_token" {
			cmd = &user.RenewAccessTokenCommand{
				Token: r.FormValue("token"),
			}
		}

		if cmd == nil {
			jsonhttp.SendError(w, 400, errors.New("Invalid token request."))
			return
		}

		bus.Handle(cmd)

		errorable, _ := cmd.(commandpkg.Errorable)
		if errorable.HasErrors() {
			httpcommandpkg.SendCommandErrors(w, errorable)
			return
		}

		command := cmd

		expires_in := command.AccessToken().ExpiresAt.Unix() - time.Now().Unix()

		t := make(map[string]interface{})
		t["access_token"] = command.AccessToken().Token
		t["expires_in"] = expires_in
		t["token_type"] = "Bearer"
		t["scope"] = nil
		t["refresh_token"] = command.RefreshToken().Token

		jsonhttp.SendJSON(w, 201, t)
	}
}
