package user

import (
	"github.com/dadamssg/commandbus"
	commandpkg "github.com/dadamssg/starterapp/app/command"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
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

type TokenCommand struct {
	commandpkg.Command
	accessToken  *Token
	refreshToken *Token
}

func (t *TokenCommand) SetAccessToken(token *Token) {
	t.accessToken = token
}

func (t *TokenCommand) SetRefreshToken(token *Token) {
	t.refreshToken = token
}

func (t *TokenCommand) AccessToken() *Token {
	return t.accessToken
}

func (t *TokenCommand) RefreshToken() *Token {
	return t.refreshToken
}

type IssueAccessTokenCommand struct {
	TokenCommand
	Username string
	Password string
}

type RenewAccessTokenCommand struct {
	TokenCommand
	Token string
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
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

func registerUserHandler(app *commandbus.CommandBus, users UserRepository, mailer UserMailer) {

	app.RegisterHandler(&RegisterUserCommand{}, func(cmd interface{}) {
		command, _ := cmd.(*RegisterUserCommand)
		uid, _ := uuid.NewV4()
		token, _ := uuid.NewV4()

		password := []byte(command.PlainPassword)

		hashedPassword, _ := bcrypt.GenerateFromPassword(password, 10)

		user := &User{
			Id:                uid.String(),
			CreatedAt:         time.Now().Local(),
			Username:          command.Username,
			Email:             command.Email,
			Password:          string(hashedPassword),
			Enabled:           false,
			ConfirmationToken: token.String(),
		}

		if err := users.Add(user); err != nil {
			commandpkg.AddCommandError(command, 500, "Internal server error.")
			return
		}

		go mailer.SendRegistrationEmail(user)

		command.User = user
	})
}

func issueAuthTokenHandler(app *commandbus.CommandBus, users UserRepository, accessTokens TokenRepository, refreshTokens TokenRepository) {

	app.RegisterHandler(&IssueAccessTokenCommand{}, func(cmd interface{}) {
		command, _ := cmd.(*IssueAccessTokenCommand)

		user, _ := users.ByUsername(command.Username)

		if user == nil {
			commandpkg.AddCommandError(command, 400, "Username not found.")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(command.Password)); err != nil {
			commandpkg.AddCommandError(command, 400, "Invalid password.")
			return
		}

		t := time.Now().UTC()

		aToken := &Token{
			ExpiresAt: t.Add(time.Duration(3600) * time.Second),
			Token:     randSeq(75),
			UserId:    user.Id,
		}

		rToken := &Token{
			ExpiresAt: t.Add(time.Duration(3600) * time.Second * 2),
			Token:     randSeq(75),
			UserId:    user.Id,
		}

		accessTokens.Add(aToken)
		refreshTokens.Add(rToken)

		command.SetAccessToken(aToken)
		command.SetRefreshToken(rToken)
	})
}

func renewAccessTokenHandler(app *commandbus.CommandBus, accessTokens TokenRepository, refreshTokens TokenRepository) {

	app.RegisterHandler(&RenewAccessTokenCommand{}, func(cmd interface{}) {
		command, _ := cmd.(*RenewAccessTokenCommand)

		oldRefreshToken, _ := refreshTokens.ByToken(command.Token)

		if oldRefreshToken == nil {
			commandpkg.AddCommandError(command, 400, "Invalid refresh token.")
			return
		}

		t := time.Now().UTC()

		if !t.Before(oldRefreshToken.ExpiresAt) {
			commandpkg.AddCommandError(command, 400, "Token has expired.")
			return
		}

		aToken := &Token{
			ExpiresAt: t.Add(time.Duration(3600) * time.Second),
			Token:     randSeq(75),
			UserId:    oldRefreshToken.UserId,
		}

		rToken := &Token{
			ExpiresAt: t.Add(time.Duration(1209600) * time.Second),
			Token:     randSeq(75),
			UserId:    oldRefreshToken.UserId,
		}

		accessTokens.Add(aToken)
		refreshTokens.Add(rToken)

		command.SetAccessToken(aToken)
		command.SetRefreshToken(rToken)
	})
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
