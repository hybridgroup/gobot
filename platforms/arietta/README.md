# Arietta

This package provides the Gobot adaptor for the [Acme Arietta](http://acmesystems.it/arietta)

## Getting Started

## Installing
```
go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/platforms/gobot/arietta
```

## Cross compiling for the Arietta
You must first configure your Go environment for arm linux cross compiling

```bash
$ cd $GOROOT/src
$ GOOS=linux GOARCH=arm GOARM=5 ./make.bash --no-clean
```

Then compile your Gobot program with
```bash
$ GOOS=linux GOARCH=arm GOARM=5 go build examples/arietta_blink.go
```

Then upload using SCP and execute it with
``` bash
$ scp arietta_blink root@192.168.7.2:/home/root/
$ ssh -t root@192.168.7.2 "./arietta_blink"
```

## Example
