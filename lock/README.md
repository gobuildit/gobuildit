# lock

This code demonstrates the problem with locking around I/O. There are two parts:

1. a server which is the main focus, and
2. a client which initiates an HTTP request, but reads no bytes.

For a detailed discussion of this code, see [here][blog].

## Running the demo

First, start the server:

```
go run server/main.go

# Starts server on localhost:8080
```

Then, confirm the server is running with:

```
curl localhost:8080/; echo
```

The response will be a large number of "1" characters, indicating only one
request has been served.

Now, to observe the locking problem, start the slow client:

```
go run client/main.go
```

With the slow client connected to the server, try to send a second request:

```
curl localhost:8080/; echo
```

The second request will hang on account of the poorly implemented locking
strategy. Let's fix that.

Stop both the server and the client. Then, open `server/main.go` and comment out the
first `root` handler and remove the comments around the second `root` handler.

Next, let's rerun our test.

Start the server again:

```
go run server/main.go

# Starts server on localhost:8080
```

Start the slow client again:

```
go run client/main.go
```

And finally, send a second request:

```
curl localhost:8080/; echo
```

Notice this time how the second request immediately returns, even though the
slow client is still connected.

## Note

If you did not see the behavior described above, you may need to adjust the
number of bytes written to each client to ensure the kernel's TCP socket buffer
is filled. See the [blog post][blog] for more details.

[blog]: https://commandercoriander.net/blog/2018/04/10/dont-lock-around-io/
