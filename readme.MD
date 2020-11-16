# TV Status RPIO

Golang project aimed at switching GPIO-pins on/off based on the Status of the TV, using
[go-rpio](https://github.com/stianeikeland/go-rpio) package

## Important note
The GPIO numbers do not

## Debugging
To see the status without touch the GPIO pins, use the following environemnt variables:

- TV_ACTIVE_GPIO_NUMBER=0
- DEBUG=true

## Settings
Using environment variables, the following settings are available:

- TV_HOSTNAME [required]
  - Hostname where the TV API is available, including scheme, e.g. http://192.168.0.1
- TV_ACTIVE_GPIO_NUMBER [required, int >= 0]
  - GPIO number you want to switch based on TV activation status
  - if 0 is provided the switching will not occur
- TV_CHECK_PERIOD [default: 5s]
  - In case of API that requires polling, the time between requests. `time.ParseDuration` is used to parse it from string
- DEBUG
  - Will report the status even without having access to the GPIO


## On/Off status
Watches the TV status and reports it's on/off state, switching desired GPIO to HIGH or LOW

### Binary
- Build the binary using `make build-arm32v7`
- run TV_HOSTNAME={HOSTNAME} TV_ACTIVE_GPIO_NUMBER={NUMBER} ./bin/arm32v7/linux/is-on

### Dockerfile
Use prebuild docker image at docker.io/docksee/tv-status-rpio-is-on:latest-arm32v7
or docker.io/docksee/tv-status-rpio-is-on:latest for testing on amd64 without GPIO access

Memory devices need to be mounted to the docker container, or see go-rpio package for other ways to use the GPIO pins

```
docker run \
    -d \
    --restart always \
    --device /dev/mem \
    --device /dev/gpiomem
    -e TV_HOSTNAME={HOSTNAME}\
    -e TV_ACTIVE_GPIO_NUMBER={NUMBER}
    docker.io/docksee/tv-status-rpio-is-on:latest-arm32v7
```


## Supported models
- Sony Bravia

### Tested models
- Sony Bravia XF9005
