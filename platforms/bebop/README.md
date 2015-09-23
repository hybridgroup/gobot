# Bebop 

The Bebop from Parrot is an inexpensive quadcopter that is controlled using WiFi. It includes a built-in front-facing HD video camera, as well as a second lower resolution bottom facing video camera.


## How to Install
```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/bebop
```
## How to Connect

The Bebop is a WiFi device, so there is no additional work to establish a connection to a single drone. However, in order to connect to multiple drones, you need to perform some configuration steps on each drone via SSH.
