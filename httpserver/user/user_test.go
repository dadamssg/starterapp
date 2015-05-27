package user

import (
	"github.com/dadamssg/starterapp/app"
	//"github.com/dadamssg/starterapp/app/user"
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var db *sql.DB

func init() {
	config := starterapp.ReadConfig("/config/starterapp.yml")
	db = starterapp.SetupDatabase(config.DbConfig)
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
	req, err := http.NewRequest("POST", "http://localhost:8080/register", bytes.NewBuffer(jsonStr))
	if err != nil {
		fail(t, "Couldn't create request.")
	}
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

	// fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func fail(t *testing.T, msg string) {
	t.Log(msg)
	t.Fail()
}

// func config() starterapp.Config {

// 	config := readConfig()

// 	return starterapp.Config{
// 		DbConfig: starterapp.DatabaseConfig{
// 			User:     config["database_user"],
// 			Password: config["database_password"],
// 			Host:     config["database_host"],
// 			Database: config["database_name"],
// 		},
// 	}
// }

// func readConfig() map[string]string {

// 	data, err := ioutil.ReadFile("/config/starterapp.yml")
// 	panicIf(err)

// 	config := make(map[string]string)

// 	err = yaml.Unmarshal([]byte(data), &config)
// 	panicIf(err)

// 	return config
// }

// func setupDatabase(config DatabaseConfig) *sql.DB {
// 	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
// 		config.User,
// 		config.Password,
// 		config.Host,
// 		config.Port,
// 		config.Database)

// 	db, err := sql.Open("postgres", conn)
// 	panicIf(err)
// 	panicIf(db.Ping())

// 	return db
// }

// func panicIf(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
