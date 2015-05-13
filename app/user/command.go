package user

import (
	"github.com/dadamssg/commandbus"
	commandpkg "github.com/dadamssg/starterapp/app/command"
	"github.com/nu7hatch/gouuid"
	"time"
)

type RegisterUserCommand struct {
	commandpkg.Command
	Username      string
	Email         string
	PlainPassword string
	User          *User
}

type FindUserByIdCommand struct {
	commandpkg.Command
	Id   string
	User *User
}

func findUserByIdHandler(app *commandbus.CommandBus, users UserRepository) {

	app.RegisterHandler(&FindUserByIdCommand{}, func(cmd interface{}) {
		command, _ := cmd.(*FindUserByIdCommand)

		user, _ := users.ById(command.Id)

		if user == nil {
			commandpkg.AddCommandError(command, 404, "User not found.")
			return
		}

		command.User = user
	})
}
func registerUserHandler(app *commandbus.CommandBus, users UserRepository) {

	app.RegisterHandler(&RegisterUserCommand{}, func(cmd interface{}) {
		command, _ := cmd.(*RegisterUserCommand)
		uid, _ := uuid.NewV4()
		token, _ := uuid.NewV4()

		user := &User{
			Id:                uid.String(),
			CreatedAt:         time.Now().Local(),
			Username:          command.Username,
			Email:             command.Email,
			Password:          command.PlainPassword,
			Enabled:           false,
			ConfirmationToken: token.String(),
		}

		if err := users.Add(user); err != nil {
			commandpkg.AddCommandError(command, 500, "Internal server error.")
			return
		}

		command.User = user
	})
}
