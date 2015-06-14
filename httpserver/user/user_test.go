package user

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/dadamssg/starterapp/app"
	"github.com/dadamssg/starterapp/app/user"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var db *sql.DB

var app *starterapp.App

func init() {
	config := starterapp.ReadConfig("/config/starterapp.yml")
	db = starterapp.SetupDatabase(config.DbConfig)
	app = starterapp.New(config)
	//users := user.NewPSQLUserRepository(db)
}

func TestMain(m *testing.M) {
	db.Query("select truncate_tables('go_app')")
	os.Exit(m.Run())
}

func TestCanRegisterUser(t *testing.T) {
	var jsonStr = []byte(`{
        "user": {
            "username": "1a232344",
            "email": "123a3232da@d93g.comd",
            "plainPassword": "adadgdsa3"
        }
    }`)
	req, _ := http.NewRequest("POST", "http://localhost:8080/register", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fail(t, "Couldn't make request.")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		fail(t, "Couldn't register user!")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func TestGetUserToken(t *testing.T) {

	cmd := &user.RegisterUserCommand{
		Username:      "johndoe",
		Email:         "jdoe@example.org",
		PlainPassword: "s3cr3t123",
	}

	app.CommandBus.Handle(cmd)

	var jsonStr = []byte(`{
        "authenticate": {
            "username": "johndoe",
            "password": "s3cr3t123"
        }
    }`)
	req, _ := http.NewRequest("POST", "http://localhost:8080/authenticate", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fail(t, "Couldn't make request.")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		fail(t, "Couldn't authenticate!")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func fail(t *testing.T, msg string) {
	t.Log(msg)
	t.Fail()
}
