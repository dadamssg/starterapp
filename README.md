starterapp
==========

This is an attempt to come up with a modular application architecture in Go. The application 
makes use of my command bus implementation. The idea being that the application/command bus struct
can be passed to several modules and they can plug themselves in. 

## Goal
The primary goal of this is to create a client agnostic application. The client could be an http server, CLI client, or anything else
that can take input, create command structs, and pass them to the application to be handled. 

## Structure
The meat of the application lives in `/app`. The application modules live there. There is only one domain module
so far, `user`. The application is created in `/app/app.go` and is passed to the user module to be connected
in `/app/user/app.go`. 

Currently there only two commands/handlers, a `RegisterUserCommand` and `FindUserByIdCommand`. They live in `app/user/command.go`.

## HTTP Server Client
There is only one client to the application, the http server. It takes the same modular approach. The 
application and http router are created in `/httpserver/main.go`. They are then passed to the user module in 
`/httpserver/user/user.go` to register the http routes to http handlers that are responsible for creating command
structs, passing them to the application/command bus to be handled, then writing a response to the `http.ResponseWriter`.

The http server application relies on a yml configuration file to keep sensitive configuration data. It looks like this:

```yml
# /config/starterapp.yml

database_user:     go_app
database_password: s3cr3t123
database_host:     localhost
database_port:     5432
database_name:     startapp_dev
```

The http server currently has two endpoints. One to register a user:

### POST /register
```json
{
    "user": {
        "username": "johndoe",
        "email": "jdoe@example.org",
        "password": "s3cr3t123"
    }
}
```

### Response - 201
```json
{
    "user": {
        "created_at": 1431397823,
        "email": "jdoe@example.org",
        "enabled": false,
        "id": "615164bc-8129-468d-5c11-f6b352dc740e",
        "username": "johndoe"
    }
}
```

The other endpoint is to retrieve an already registered user:

### GET /user/{id}

### Response - 200
```json
{
    "user": {
        "created_at": 1431397823,
        "email": "jdoe@example.org",
        "enabled": false,
        "id": "615164bc-8129-468d-5c11-f6b352dc740e",
        "username": "johndoe"
    }
}
```

I'm planning on tackling a CLI client next. 
