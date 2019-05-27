# Nokeval temperature sensor reader

## Build

### Mac (for development)

```
$ brew install golang
```

### Raspberry pi

#### Prerequisites
Install:

```
$ apt update
$ apt install -y golang git
$ mkdir ~/go
```

Remember to execute this (or add this to your `.bash_profile` or so) every time you login:

```
$ export GOPATH=~/go
```

#### Get the code
The reader itself:

```
$ go get -v github.com/hkroger/nokeval-reader-go
```

Get deps:

```
$ go get -v github.com/bvinc/go-sqlite-lite
$ go get -v github.com/sirupsen/logrus
$ go get -v github.com/jacobsa/go-serial/serial
$ go get -v gopkg.in/yaml.v2
```

#### Build
```
$ cd ~/go/src/github.com/hkroger/nokeval-reader-go/
$ go build -o nokeval-reader cmd/reader/main.go
```

#### Run

Verbose mode:

```
$ ./nokeval-reader -v -c  /opt/nokeval_reader/config.yaml
```

Production mode:

```
$ ./nokeval-reader -c  /opt/nokeval_reader/config.yaml
```
