# TODO APP USING GRPC

fun project of a production ready todo application utilizing diffrent features of gRPC.

In this project we utilize unary request, server streaming , client streaming and bidirectional streaming.

Running docker Todo -server container

```
docker run -p 8081:8081 paschalolo/grpcserver:1.0.0
```

Running server locally

```
make server
```

Running client locally

```
make client
```
