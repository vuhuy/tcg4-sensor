# TCG4 sensor

[![Go Report Card](https://goreportcard.com/badge/github.com/vuhuy/tcg4-sensor?style=flat-square)](https://goreportcard.com/report/github.com/vuhuy/tcg4-sensor) [![Release](https://img.shields.io/github/release/vuhuy/tcg4-sensor?style=flat-square)](https://github.com/vuhuy/tcg4-sensor/releases/latest)

Provides the `sensor` application to read the built-in TCG4 motion sensor and GPS data via GPSd. This data is then spammed over CAN for everyone to hear. An HTTP REST server is also availble with additional satellite data. Developed using Go 1.21.1 and tested on Koala v3.06r0.

## Usage

Several options can be configured when running the application.

- Set an API key to require authentication for REST. If set, each request will be checked for a valid X-API-Key header.
- Enable or disable sending certain CAN frames. CAN frames with estimated errors are disabled by default.
- Configure the used CAN message ID of certain CAN frames. The flag accepts decimal representations of the CAN message ID.
- The polling frequency of the motion (accelerometer) sensor. The GPS polling frequency is configured by the `/etc/gps.conf`.

```shell
$ sensor -help
```
```
Usage of sensor:
  -alt-frame-id uint
        CAN frame ID for GPS altitude data [m] (float64 LE). Set frame ID to 0 to disable. (default 202)
  -can-extended
        Use extended CAN. Set to true to enable.
  -can-frequency float
        Message frequency [Hz] of all sensor data on CAN bus. Set frame ID to 0 to disable. (default 1)
  -can-interface string
        CAN interface name to send sensor data. (default "can0")
  -epc-frame-id uint
        CAN frame ID for the GPS estimated climb error [m/s] (float64 LE). Set frame ID to enable.
  -epd-frame-id uint
        CAN frame ID for the GPS estimated track (direction) error [°] (float64 LE). Set frame ID to enable.
  -eph-frame-id uint
        CAN frame ID for the GPS estimated horizontal position (2D) error [m] (float64 LE). Set frame ID to enable.
  -eps-frame-id uint
        CAN frame ID for the GPS estimated speed error [m/s] (float64 LE). Set frame ID to enable.
  -ept-frame-id uint
        CAN frame ID for the GPS estimated time stamp error [s] (float64 LE). Set frame ID to enable.
  -epv-frame-id uint
        CAN frame ID for the GPS estimated vertical error [m] (float64 LE). Set frame ID to enable.
  -epx-frame-id uint
        CAN frame ID for the GPS estimated longitude error [m] (float64 LE). Set frame ID to enable.
  -epy-frame-id uint
        CAN frame ID for the GPS estimated latitude error [m] (float64 LE). Set frame ID to enable.
  -gdop-frame-id uint
        CAN frame ID for the GPS geometric (hyperspherical) dilution of precision (float64 LE). Set frame ID to enable.
  -gps-frame-id uint
        CAN frame ID for GPS mode (1:uint8), status (2:uint8), visible satellites (3+4:uint16 LE), used satellites (5+6:uint16 LE), and quality data (7: uint8) data (8: not used). Set frame ID to enable. (default 204)
  -gpsd-host string
        Hostname of the device that runs the GPSd TCP feed. (default "localhost")
  -gpsd-port uint
        Port running the GPSd TCP feed. (default 2947)
  -hdop-frame-id uint
        CAN frame ID for the GPS horizontal dilution of precision (float64 LE). Set frame ID to enable.
  -lat-frame-id uint
        CAN frame ID for GPS latitude data [°] (float64 LE). Set frame ID to 0 to disable. (default 200)
  -lon-frame-id uint
        CAN frame ID for GPS longitude data [°] (float64 LE). Set frame ID to 0 to disable. (default 201)
  -motion-frame-id uint
        CAN frame ID for motion X [mg] (1+2:int16 LE), Y [mg] (3+4:int16 LE), Z [mg] (5+6:int16 LE), and scale [g] (7:uint8) data (8: not used). Set frame ID to 0 to disable. (default 205)
  -motion-frequency float
        Polling frequency [Hz] of motion sensor data (default 1)
  -pdop-frame-id uint
        CAN frame ID for the GPS position (spherical/3D) dilution of precision dilution of precision (float64 LE). Set frame ID to enable.
  -rest-api-key string
        Expected X-API-Key header value to authenticate HTTP requests. Set a key to enable.
  -rest-port uint
        Port used to serve the HTTP REST API. (default 8081)
  -sep-frame-id uint
        CAN frame ID for the GPS estimated spherical (3D) position error [m] (float64 LE). Set frame ID to enable.
  -speed-frame-id uint
        CAN frame ID for GPS speed data [m/s] (float64 LE). Set frame ID to 0 to disable. (default 203)
  -tdop-frame-id uint
        CAN frame ID for the GPS time dilution of precision (float64 LE). Set frame ID to enable.
  -vdop-frame-id uint
        CAN frame ID for the GPS vertical (altitude) dilution of precision dilution of precision (float64 LE). Set frame ID to enable.
  -verbose
        More verbose output for debugging purposes. Set to true to enable.
  -version
        Print the current application version.
  -xdop-frame-id uint
        CAN frame ID for the GPS longitudinal dilution of precision (float64 LE). Set frame ID to enable.
  -ydop-frame-id uint
        CAN frame ID for the GPS latitudinal dilution of precision (float64 LE). Set frame ID to enable.
```

Some CAN frames and REST endpoints return enumerator data.

```
mode (uint8)
0 = unknown
1 = no fix
2 = 2D
3 = 3D

status (uint8)
0 = unknown
1 = normal
2 = DGPS
3 = RTK fixed
4 = RTK floating
5 = DR
6 = GNSSDR
7 = time (surveyed)
8 = simulated
9 = P(Y)
		
quality (uint8)
0 = no signal
1 = searching signal
2 = signal acquired
3 = signal detected but unusable
4 = code locked and time-synchronized

gnssid (uint8)
0 = GPS (GP)
1 = SBAS (SB)
2 = Galileo (GA)
3 = BeiDou (BD)
4 = IMES (IM)
5 = QZSS (QZ)
6 = GLONASS (GL)
7 = NavIC (IRNSS) (IR)
```

## Installation

Make sure `ypgpsd` and `gps.conf` are configured correctly on your TCG4. Some `GPS_RATE_MS` values, like 500, can result in weird behavior with missing GPS data. This is most likely a bug in the GPS driver.

```shell
$ sudo mv /etc/taf/_ygpsd.config /etc/taf/ygpsd.config
$ sudo vi /etc/gps.conf
```
```
GPS_ENABLE=ON
GPS_RATE_MS=1000
```

Copy the `sensor` binary to the `/usr/local/bin` directory on your TCG4 using tools like `scp`, `wget`, or WinSCP. You can download a precompiled binary from the [GitHub release page](https://github.com/vuhuy/tcg4-sensor/releases), or build one from source. Make sure it's owned by `root` and has the correct execution permissions.

```shell
$ sudo chmod 755 /usr/local/bin/sensor
$ sudo chown root:root /usr/local/bin/sensor
```

Create an `init.d` script and configuration file on your TCG4 for auto startup. An example is provided in the `init` directory of this project. `S50sensor` goes in `/etc/init.d/` and `sensor.conf` goes in `/etc/`. Make sure both files are owned by `root` with the correct execution permissions for `S50sensor`.

```shell
$ sudo chmod 744 /etc/init.d/S50sensor
$ sudo chmod 644 /etc/sensor.conf
$ sudo chown root:root /etc/init.d/S50sensor
$ sudo chown root:root /etc/sensor.conf
```

## Build from source

**Prerequisite software**: [Go](https://go.dev/doc/install)

Build `go-sensor` for TCG4 (armv7).

```
env GOOS=linux GOARCH=arm go build -o build/sensor ./cmd/sensor
```