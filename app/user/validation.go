package user

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/dadamssg/starterapp/app/command"
)

func registerUserValidator(validator *command.Validator, users UserRepository) {
	validator.Register(&RegisterUserCommand{}, func(cmd interface{}) []command.CommandError {
		c, _ := cmd.(*RegisterUserCommand)

		errs := []command.CommandError{}

		if !govalidator.IsEmail(c.Email) {
			errs = append(errs, newError("Invalid email."))
		}

		if !govalidator.IsAlphanumeric(c.Username) {
			errs = append(errs, newError("Username must be alphanumeric."))
		}

		if len(c.PlainPassword) < 6 {
			errs = append(errs, newError("Password must be at least 6 characters."))
		}

		if len(c.PlainPassword) >= 30 {
			errs = append(errs, newError("Password must be fewer than 30 characters."))
		}

		if len(c.Username) < 3 {
			errs = append(errs, newError("Username must be at least 3 characters."))
		}

		if len(c.Username) >= 20 {
			errs = append(errs, newError("Username must be fewer than 20 characters."))
		}

		// return before hitting the database
		if len(errs) > 0 {
			return errs
		}

		if user, _ := users.ByEmail(c.Email); user != nil {
			errs = append(errs, newError("Email already registered."))
		}

		if user, _ := users.ByUsername(c.Username); user != nil {
			errs = append(errs, newError("Username already registered."))
		}

		return errs
	})
}

func newError(err string) command.CommandError {
	return command.CommandError{Code: 400, Err: errors.New(err)}
}
