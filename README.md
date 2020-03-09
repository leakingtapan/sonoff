# Sonoff
sonoff server and similuated switch device. This work is based of [simple-sonoff-server](https://github.com/mdopp/simple-sonoff-server) and blog posts http://blog.nanl.de/2017/05/sonota-flashing-itead-sonoff-devices-via-original-ota-mechanism/ and https://blog.ipsumdomus.com/sonoff-switch-complete-hack-without-firmware-upgrade-1b2d6632c01. Aside from the sonoff server, there are several extra feature implemented. The features including:
* a simulated sonoff switch that can connect to both eWeLink cloud and local sonoff server.
* a cli tool `sonoff` that bundles both client and server
* a RESTfull API based off OpenAPI for the sonoff server.
* a docker image that works for both x86_64 and arm platform

This project is built in golang. 

# Server
Start the sonoff server using:
```sh
./sonoff server --server-ip {sonoffServerIp} --server-port {sonoffServerPort}
```

# Switch
Start the simulated switch using:
```sh
./sonoff switch --device-spec-path {deviceSpec}
```
