# Nokeval temperature sensor reader
This is a utility which reads the measurements of a wireless thermometer Measurinator.com service. This works with Nokeval FTR970B and compatible devices.

It can be run as a daemon and an example launchd configuration file is included.

## How to use	
### Easy install for Debian/Raspbian 9 (Stretch)

Open apt-get file in editor:

    vim /etc/apt/sources.list.d/measurinator.list

Add:

    deb [trusted=yes] http://koti.kapsi.fi/hkroger/debs/stretch ./

Save & run:

	apt update
	apt-get install nokeval-reader

Edit configs:

	cd /opt/nokeval_reader
	cp config.yaml.example config.yaml
	vim config.yaml
	
Add key where it says `<key here>` and client id where it says `<client id here>`.

And start the service

	systemctl start nokeval_reader

## Development

### Mac (for development)

```
$ brew install golang
$ mkdir ~/go
```

### Raspberry pi

Run all of the following steps on your raspberry pi.

Install:

```
$ apt update
$ apt install -y git
$ wget https://dl.google.com/go/go1.13.1.linux-armv6l.tar.gz
$ tar xzvf go1.13.1.linux-armv6l.tar.gz
$ sudo mv go /usr/local
$ mkdir ~/go


```

To run stuff:

```
$ export PATH=/usr/local/go/bin:$PATH
$ export GOROOT=/usr/local/go
```

## Setup your environment

Create directory for your Golang build environment:

```
$ mkdir ~/go
```


Remember to execute this (or add this to your `.bash_profile` or so) every time you login:

```
$ export GOPATH=~/go
```

## Get the code
The reader itself:

```
$ go get -v github.com/hkroger/nokeval-reader-go/...
```

## Build
```
$ cd ~/go/src/github.com/hkroger/nokeval-reader-go/
$ go build -o nokeval-reader cmd/reader/main.go
```

## Run

Verbose mode:

```
$ ./nokeval-reader -v -c  /opt/nokeval_reader/config.yaml
```

Production mode:

```
$ ./nokeval-reader -c  /opt/nokeval_reader/config.yaml
```

## Build the deb (has to be on Raspberry pi)

```
$ cd ~/go/src/github.com/hkroger/nokeval-reader-go/
$ ./build_deb.sh
```