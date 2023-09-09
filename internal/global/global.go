package global

import (
	"flag"
	"sync"
)

// Define some global vars.
var (
	Wg              sync.WaitGroup
	GpsdHost        = flag.String("gpsd-host", "localhost", "Hostname of the device that runs the GPSd TCP feed.")
	GpsdPort        = flag.Uint64("gpsd-port", 2947, "Port running the GPSd TCP feed.")
	RestPort        = flag.Uint64("rest-port", 8081, "Port used to serve the HTTP REST API.")
	RestApiKey      = flag.String("rest-api-key", "", "Expected X-API-Key header value to authenticate HTTP requests. Set a key to enable.")
	CanInterface    = flag.String("can-interface", "can0", "CAN interface name to send sensor data.")
	CanExtended     = flag.Bool("can-extended", false, "Use extended CAN. Set to true to enable.")
	CanFrequency    = flag.Float64("can-frequency", 1, "Message frequency [Hz] of all sensor data on CAN bus. Set frame ID to 0 to disable.")
	LatFrameId      = flag.Uint64("lat-frame-id", 200, "CAN frame ID for GPS latitude data [°] (float64 LE). Set frame ID to 0 to disable.")
	LonFrameId      = flag.Uint64("lon-frame-id", 201, "CAN frame ID for GPS longitude data [°] (float64 LE). Set frame ID to 0 to disable.")
	AltFrameId      = flag.Uint64("alt-frame-id", 202, "CAN frame ID for GPS altitude data [m] (float64 LE). Set frame ID to 0 to disable.")
	SpeedFrameId    = flag.Uint64("speed-frame-id", 203, "CAN frame ID for GPS speed data [m/s] (float64 LE). Set frame ID to 0 to disable.")
	GpsFrameId      = flag.Uint64("gps-frame-id", 204, "CAN frame ID for GPS mode (1:uint8), status (2:uint8), visible satellites (3+4:uint16 LE), used satellites (5+6:uint16 LE), and quality data (7: uint8) data (8: not used). Set frame ID to enable.")
	EpcFrameId      = flag.Uint64("epc-frame-id", 0, "CAN frame ID for the GPS estimated climb error [m/s] (float64 LE). Set frame ID to enable.")
	EpdFrameId      = flag.Uint64("epd-frame-id", 0, "CAN frame ID for the GPS estimated track (direction) error [°] (float64 LE). Set frame ID to enable.")
	EphFrameId      = flag.Uint64("eph-frame-id", 0, "CAN frame ID for the GPS estimated horizontal position (2D) error [m] (float64 LE). Set frame ID to enable.")
	EpsFrameId      = flag.Uint64("eps-frame-id", 0, "CAN frame ID for the GPS estimated speed error [m/s] (float64 LE). Set frame ID to enable.")
	EptFrameId      = flag.Uint64("ept-frame-id", 0, "CAN frame ID for the GPS estimated time stamp error [s] (float64 LE). Set frame ID to enable.")
	EpxFrameId      = flag.Uint64("epx-frame-id", 0, "CAN frame ID for the GPS estimated longitude error [m] (float64 LE). Set frame ID to enable.")
	EpyFrameId      = flag.Uint64("epy-frame-id", 0, "CAN frame ID for the GPS estimated latitude error [m] (float64 LE). Set frame ID to enable.")
	EpvFrameId      = flag.Uint64("epv-frame-id", 0, "CAN frame ID for the GPS estimated vertical error [m] (float64 LE). Set frame ID to enable.")
	SepFrameId      = flag.Uint64("sep-frame-id", 0, "CAN frame ID for the GPS estimated spherical (3D) position error [m] (float64 LE). Set frame ID to enable.")
	XdopFrameId     = flag.Uint64("xdop-frame-id", 0, "CAN frame ID for the GPS longitudinal dilution of precision (float64 LE). Set frame ID to enable.")
	YdopFrameId     = flag.Uint64("ydop-frame-id", 0, "CAN frame ID for the GPS latitudinal dilution of precision (float64 LE). Set frame ID to enable.")
	VdopFrameId     = flag.Uint64("vdop-frame-id", 0, "CAN frame ID for the GPS vertical (altitude) dilution of precision dilution of precision (float64 LE). Set frame ID to enable.")
	TdopFrameId     = flag.Uint64("tdop-frame-id", 0, "CAN frame ID for the GPS time dilution of precision (float64 LE). Set frame ID to enable.")
	HdopFrameId     = flag.Uint64("hdop-frame-id", 0, "CAN frame ID for the GPS horizontal dilution of precision (float64 LE). Set frame ID to enable.")
	PdopFrameId     = flag.Uint64("pdop-frame-id", 0, "CAN frame ID for the GPS position (spherical/3D) dilution of precision dilution of precision (float64 LE). Set frame ID to enable.")
	GdopFrameId     = flag.Uint64("gdop-frame-id", 0, "CAN frame ID for the GPS geometric (hyperspherical) dilution of precision (float64 LE). Set frame ID to enable.")
	MotionFrameId   = flag.Uint64("motion-frame-id", 205, "CAN frame ID for motion X [mg] (1+2:int16 LE), Y [mg] (3+4:int16 LE), Z [mg] (5+6:int16 LE), and scale [g] (7:uint8) data (8: not used). Set frame ID to 0 to disable.")
	MotionFrequency = flag.Float64("motion-frequency", 1, "Polling frequency [Hz] of motion sensor data")
	Verbose         = flag.Bool("verbose", false, "More verbose output for debugging purposes. Set to true to enable.")
	Version         = flag.Bool("version", false, "Print the current application version.")
)
