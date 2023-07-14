FROM golang:1.18 AS builder
WORKDIR /go/src/user-repository
COPY idp-repository/user-repository .
COPY idp-repository/pkg /go/src/pkg
COPY idp-repository/env /go/src/env
COPY idp-repository/protos /go/src/protos
COPY idp-repository/go.mod /go/src
COPY idp-repository/go.sum /go/src
RUN cd /go/src && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-repository ./cmd/main.go

FROM scratch
COPY --from=builder /go/src/user-repository/user-repository ./user-repository
ENV GRPC_PORT=:81 \
    DB_HOST=localhost \
    DB_PORT=5432 \
    DB_USER=user \
    DB_PASS=pass \
    DB_NAME=membership \
    REDIS_HOST=localhost \
    REDIS_PORT=6379 \
    REDIS_PASS=pass \
    MODE=Development

CMD ["./user-repository"]