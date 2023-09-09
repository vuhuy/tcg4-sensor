# pcan-client-example

Real-world client example receiving CAN and requesting REST data data. This example is specifically written for Windows x64 with PEAK-System PCAN-USB adapters for bus communication. It should work on Mac too. Developed using Node.js 18.17.1 and tested on Windows 11 / Firefox 117.0. You can download a precompiled binary from the [GitHub release page](https://github.com/vuhuy/tcg4-sensor/releases), or build one from source. 

## TCG4 configuration

```shell
$ sudo vi /etc/sensor.conf
```
```
SENSOR_ENABLE=ON
SENSOR_ARGS="--rest-api-key=0b239503-a040-4eec-acc6-c4fca7534b5b --motion-frequency 10 --can-frequency 10"
```

```shell
$ sudo vi /etc/gps.conf
```
```
GPS_ENABLE=ON
GPS_RATE_MS=100
```

## Build from source

**Prerequisite software**: [Node.js 18 LTS](https://nodejs.org/en/download)

Download and install modules (needs to be executed only once).

```shell
PS > npm i
```

Run the example.

```shell
PS > node index.js
```
