-------RUN:----------
# will uncoment main sections and then:
nats-server
go run reply.go
go run req.go

