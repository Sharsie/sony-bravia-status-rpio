ARG ALPINE_IMAGE
ARG BUILD_IMAGE

FROM $ALPINE_IMAGE as alpine

# Install ca certificates to enable https requests
RUN apk --no-cache add ca-certificates


FROM $BUILD_IMAGE as build
ARG APP_NAME

WORKDIR /go/src/github.com/Sharsie/tv-status-rpio/

COPY ["./cmd", "./cmd"]
COPY ["./go.mod", "./go.mod"]
COPY ["./go.sum", "./go.sum"]

ARG GOOS
ARG GOARCH
ARG GOARM
ARG COMMAND
ARG DOCKER_TAG


RUN CGO_ENABLED=0 go build -mod readonly -ldflags "-s -w \
    		-X github.com/Sharsie/$APP_NAME/cmd/$COMMAND/version.Tag=$DOCKER_TAG" -o ./bin/$COMMAND ./cmd/$COMMAND

FROM scratch

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ARG COMMAND

COPY --from=build /go/src/github.com/Sharsie/tv-status-rpio/bin/$COMMAND /app

ENTRYPOINT ["/app"]
