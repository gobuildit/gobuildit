# Guessing Game

There are two applications here: 1) a command line application, and 2) a web
server. Both are implementations of a guessing game, just with different user
interfaces.

The guessing game picks a random number and then asks a player to guess that
number. When the user guesses the number, the game is over.

## Command line application

To start the command line application, run the following command:

```
go run cli/main.go
```

## Web Server

To start the web server, run the following command:

```
go run web/main.go
```

In a separate session, send your guess to the web server with:

```
curl -i localhost:8080/guesses -H "Content-Type: application/json" -d '{"number": $YOUR_GUESS}'
```

If your guess is correct, the web server will return a `201 Created` status and a
success message:

```
HTTP/1.1 201 Created
Date: Tue, 20 Mar 2018 20:29:31 GMT
Content-Length: 55
Content-Type: application/json; charset=utf-8

{"message": "Congratulation! You win!"}
```

Otherwise, the server will return a `418 I'm a teapot` status, which means your
guess is not correct and the game is still on.
