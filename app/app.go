package starterapp

import (
	"database/sql"
	"fmt"
	"github.com/dadamssg/commandbus"
	"github.com/dadamssg/starterapp/app/user"
	_ "github.com/lib/pq"
)

type App struct {
	CommandBus *commandbus.CommandBus
}

func (app *App) RegisterHandler(t interface{}, function func(cmd interface{})) {
	app.CommandBus.RegisterHandler(t, function)
}

type Config struct {
	DbConfig DatabaseConfig
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func New(config Config) *App {
	app := &App{
		CommandBus: commandbus.New(),
	}

	db := setupDatabase(config.DbConfig)
	users := user.NewPSQLUserRepository(db)

	user.Connect(app.CommandBus, users)

	return app
}

func setupDatabase(config DatabaseConfig) *sql.DB {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	db, err := sql.Open("postgres", conn)
	panicIf(err)
	panicIf(db.Ping())

	return db
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
