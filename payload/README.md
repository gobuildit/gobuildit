# JSON vs Protobuf Payloads

When interacting with remote servers using HTTP, JSON stands out as a prominent
and convenient way to transfer data. In fact, given the prevalence of mobile and
JavaScript clients, JSON has become a fantastic choice when designing a remote
server's API. In turn, as backend systems have grown in capability and often
multiplied into a distributed collection of servers, JSON has enjoyed even more
use as the primary means for one server to transfer data to another server.

In spite of its wide use, JSON has some significant drawbacks. Nothing about the
JSON format guarantees a fixed structure -- keys may be present in one payload
but not in another and yet in both cases the payload may still be valid JSON.
Further, assuming there exists a shared representation of data, one must often
repeatedly write the same boilerplate code to serialize and deserialize that
representation, which in addition to being tedious, may also be error-prone.

Given the state of JSON, protocol buffers, protobuf for short, have been
enjoying a rise in use recently. Protobufs are structured data, and provide a
fixed representation shared between client and server. Finally, there are tools
that generate code to serialize and deserialize protobufs, alleviating the
burden of writing any boilerplate. Although gzipped payloads may remove any
meaningful difference, protobuf is also a more compact representation than plain
JSON.

The code in this package demonstrates the details of working with protobufs in
comparison to JSON. It is structured as a client and a server.

To run the example, first start the server:

```
go run cmd/server/main.go
```

Next, in another session, run the client:

```
go run cmd/client/main.go
```

By default, the client will send a JSON representation to the server as part of
a POST request. The server will print out the content length in addition to a
string representation of the data transferred. If all goes well, the client will
exit silently.

Next, run the client again, this time passing a flag to use protobuf
serialization:

```
go run cmd/client/main.go -proto
```

The server will again print out the content length of the protobuf payload and a
string representation of the data transferred.

The key takeaway is the protobuf transfer takes less than half as many bytes as
the JSON transfer.
