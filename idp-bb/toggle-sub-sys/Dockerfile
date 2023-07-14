FROM golang:1.20 AS builder
WORKDIR /app
COPY ./login-throttles-sub-system .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o login-throttles-sub-system ./cmd/main.go

FROM scratch
COPY --from=builder /app/login-throttles-sub-system ./login-throttles-sub-system
CMD ["./login-throttles-sub-system"]
