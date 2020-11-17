FROM alpine:3.12 as alpine

# Install ca certificates to enable https requests
RUN apk --no-cache add ca-certificates


FROM golang:1.15.5 as build
ARG APP_NAME

WORKDIR /go/src/github.com/Sharsie/tv-status-rpio/

COPY ["./cmd", "./cmd"]
COPY ["./go.mod", "./go.mod"]
COPY ["./go.sum", "./go.sum"]

ENV GOOS=linux
ENV GOARCH=amd64
ARG COMMAND
ARG DOCKER_TAG

RUN CGO_ENABLED=0 go build -mod readonly -ldflags "-s -w \
    		-X github.com/Sharsie/$APP_NAME/cmd/$COMMAND/version.Tag=$DOCKER_TAG" -o ./bin/$COMMAND ./cmd/$COMMAND

FROM scratch

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ARG COMMAND
ARG APP_NAME

COPY --from=build /go/src/github.com/Sharsie/$APP_NAME/bin/$COMMAND /app

ENTRYPOINT ["/app"]
