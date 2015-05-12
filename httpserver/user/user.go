package user

import (
	"errors"
	"github.com/dadamssg/commandbus"
	"github.com/dadamssg/jsonhttp"
	"github.com/dadamssg/starterapp/app"
	"github.com/dadamssg/starterapp/app/user"
	commandpkg "github.com/dadamssg/starterapp/httpserver/command"
	"github.com/gorilla/mux"
	"net/http"
)

func Connect(app *starterapp.App, router *mux.Router) {
	router.HandleFunc("/register", registerUser(app.CommandBus)).Methods("POST")
	router.HandleFunc("/users/{id}", findUserById(app.CommandBus)).Methods("GET")
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

func findUserById(bus *commandbus.CommandBus) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		command := user.FindUserByIdCommand{
			Id: id,
		}

		bus.Handle(&command)

		if command.HasErrors() {
			commandpkg.SendCommandErrors(w, &command)
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
			commandpkg.SendCommandErrors(w, &command)
			return
		}

		jsonhttp.SendJSON(w, 201, userRepresentation(command.User))
	}
}
