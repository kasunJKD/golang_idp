FROM golang:1.18-alpine AS builder

#RUN git config --global url."https://x-token-auth:ATCTT3xFfGN0xri0u3ZUP4ES4yHeR_K8uX8cDzW0GbOFyB85UlgMJpWQdQDhFwpw_0T6FRuewNWzsmHmSUbhXhyOZpTRkhfUxP23Oz7bYMzsBaUk2NPeuPMGUIxCsJEbr9H-nk26oAKyhH45H_bL4vjtBqambARQqosIxWxuti6axobnLBesPLo=0BCF3EC1@bitbucket.org/".insteadOf "https://bitbucket.org/"
#RUN go env -w GOPRIVATE="bitbucket.org/project-99-games/otp_model"
ARG ACCESS_TOKEN_USR
ARG ACCESS_TOKEN_PWD

RUN apk add --no-cache ca-certificates git

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Create a netrc file using the credentials specified using --build-arg
RUN printf "machine bitbucket.org\n\
    login ${ACCESS_TOKEN_USR}\n\
    password ${ACCESS_TOKEN_PWD}\n\
    \n\
    machine api.bitbucket.org\n\
    login ${ACCESS_TOKEN_USR}\n\
    password ${ACCESS_TOKEN_PWD}\n"\
    >> /root/.netrc
RUN chmod 600 /root/.netrc

RUN go env -w GOPRIVATE="bitbucket.org/project-99-games/otp_model"


WORKDIR /go/src/gateway
COPY idp-service/go.mod /go/src
COPY idp-service/go.sum /go/src
RUN cd /go/src && go mod download
COPY idp-service/gateway .
COPY idp-service/protos /go/src/protos
COPY idp-service/pkg /go/src/pkg
COPY idp-service/env /go/src/env
RUN cd /go/src && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway ./main.go

FROM scratch
COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/gateway/gateway ./gateway
EXPOSE 48080
ENV HTTP_PORT=:80 \
    GRPC_PORT=:81 \
    HOST=localhost \
    USER_SERVICE=localhost \
    TOKEN_SERVICE=localhost \
    LOGIN_SERVICE=localhost \
    MODE=Development

CMD ["./gateway"]