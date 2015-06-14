package starterapp

import (
	"database/sql"
	"fmt"
	"github.com/dadamssg/commandbus"
	"github.com/dadamssg/starterapp/app/command"
	"github.com/dadamssg/starterapp/app/user"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type App struct {
	CommandBus *commandbus.CommandBus
}

func (app *App) RegisterHandler(t interface{}, function func(cmd interface{})) {
	app.CommandBus.RegisterHandler(t, function)
}

type Config struct {
	DbConfig       DatabaseConfig
	MandrillApiKey string
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

	validator := command.NewValidator()
	app.CommandBus.AddMiddleware(1, command.RegisterMiddleware(validator))

	db := SetupDatabase(config.DbConfig)
	users := user.NewPSQLUserRepository(db)
	accessTokens := user.NewPSQLAccessTokenRepository(db)
	refreshTokens := user.NewPSQLRefreshTokenRepository(db)
	mailer := user.NewMandrillMailer(config.MandrillApiKey)

	user.Connect(app.CommandBus, validator, users, accessTokens, refreshTokens, mailer)

	return app
}

func ReadConfig(configPath string) Config {
	data, err := ioutil.ReadFile(configPath)
	panicIf(err)

	config := make(map[string]string)

	err = yaml.Unmarshal([]byte(data), &config)
	panicIf(err)

	return Config{
		DbConfig: DatabaseConfig{
			User:     config["database_user"],
			Password: config["database_password"],
			Host:     config["database_host"],
			Database: config["database_name"],
		},
		MandrillApiKey: config["mandrill_api_key"],
	}
}

func SetupDatabase(config DatabaseConfig) *sql.DB {
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
