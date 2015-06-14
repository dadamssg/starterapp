package user

import (
	"github.com/dadamssg/commandbus"
	"github.com/dadamssg/starterapp/app/command"
)

func Connect(bus *commandbus.CommandBus, validator *command.Validator, users UserRepository, accessTokens TokenRepository, refreshTokens TokenRepository, mailer UserMailer) {
	registerUserHandler(bus, users, mailer)
	findUserByIdHandler(bus, users)
	issueAuthTokenHandler(bus, users, accessTokens, refreshTokens)
	renewAccessTokenHandler(bus, accessTokens, refreshTokens)
	registerUserValidator(validator, users)
}
