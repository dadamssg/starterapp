package user

import (
	"github.com/dadamssg/commandbus"
)

func Connect(bus *commandbus.CommandBus, users UserRepository) {
	registerUserHandler(bus, users)
	findUserByIdHandler(bus, users)
}
