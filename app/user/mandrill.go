package user

import (
	"fmt"
	"github.com/mattbaird/gochimp"
)

type MandrillMailer struct {
	mandrill *gochimp.MandrillAPI
}

func (mailer *MandrillMailer) SendRegistrationEmail(user *User) error {

	text := `<a href="https://jsonstub.com/register/confirm/%s">
        Click this link to finish activating your account.
        </a>`

	content := []gochimp.Var{
		gochimp.Var{"greeting", "Howdy and welcome!"},
		gochimp.Var{"main", fmt.Sprintf(text, user.ConfirmationToken)},
	}

	recipients := []gochimp.Recipient{
		gochimp.Recipient{Email: user.Email},
	}

	message := gochimp.Message{
		Subject:   "Account Confirmation",
		FromEmail: "robot@jsonstub.com",
		FromName:  "JsonStub",
		To:        recipients,
	}

	_, err := mailer.mandrill.MessageSendTemplate("basic", content, message, false)

	return err
}

func NewMandrillMailer(apiKey string) *MandrillMailer {
	mandrillApi, err := gochimp.NewMandrill(apiKey)

	if err != nil {
		panic(err)
	}

	return &MandrillMailer{
		mandrill: mandrillApi,
	}
}
