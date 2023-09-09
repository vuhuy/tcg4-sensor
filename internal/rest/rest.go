package rest

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vuhuy/tcg4-sensor/internal/global"
	"github.com/vuhuy/tcg4-sensor/internal/gps"
	"github.com/vuhuy/tcg4-sensor/internal/motion"
	"github.com/vuhuy/tcg4-sensor/pkg/gpsd"
)

// Struct to store REST server data.
type Rest struct{}

// Struct to store sensor data.
type SensorsJson struct {
	Gps    GpsJson    `json:"gps"`
	Motion MotionJson `json:"motion"`
}

// Struct to store GPS data.
type GpsJson struct {
	Tpv TpvJson `json:"tpv"`
	Sky SkyJson `json:"sky"`
}

// Struct store GPS TPV report data.
type TpvJson struct {
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Alt    float64 `json:"alt"`
	Speed  float64 `json:"speed"`
	Mode   uint8   `json:"mode"`
	Status uint8   `json:"status"`
	Epc    float64 `json:"epc"`
	Epd    float64 `json:"epd"`
	Eph    float64 `json:"eph"`
	Eps    float64 `json:"eps"`
	Ept    float64 `json:"ept"`
	Epx    float64 `json:"epx"`
	Epy    float64 `json:"epy"`
	Epv    float64 `json:"epv"`
	Sep    float64 `json:"sep"`
}

// Struct to store GPS SKY report data.
type SkyJson struct {
	Qual       uint8           `json:"qual"`
	Xdop       float64         `json:"xdop"`
	Ydop       float64         `json:"ydop"`
	Vdop       float64         `json:"vdop"`
	Tdop       float64         `json:"tdop"`
	Hdop       float64         `json:"hdop"`
	Pdop       float64         `json:"pdop"`
	Gdop       float64         `json:"gdop"`
	Nsat       uint16          `json:"nsat"`
	Usat       uint16          `json:"usat"`
	Satellites []SatelliteJson `json:"satellites"`
}

// Struct to store GPS satellite data.
type SatelliteJson struct {
	Prn    float64 `json:"prn"`
	Az     float64 `json:"az"`
	El     float64 `json:"el"`
	Ss     float64 `json:"ss"`
	Gnssid uint8   `json:"gnssid"`
	Used   bool    `json:"used"`
}

// Struct to store motion data.
type MotionJson struct {
	X     int16 `json:"x"`
	Y     int16 `json:"y"`
	Z     int16 `json:"z"`
	Scale uint8 `json:"scale"`
}

// Prepare a request.
func prepareRequest(w http.ResponseWriter, r *http.Request, lastUpdate time.Time) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)

		return false
	} else if r.Method == "GET" {
		if *global.Verbose {
			fmt.Printf("[%v] HTTP request: %s \"%s %s\" \"%s\"\n", time.Now().UTC(), r.RemoteAddr, r.Method, r.URL.Path, r.UserAgent())
		}

		if lastUpdate.IsZero() {
			w.WriteHeader(http.StatusNoContent)

			return false
		}

		if *global.RestApiKey != "" {
			apiKeyHeader := r.Header.Get("X-API-Key")

			if apiKeyHeader != *global.RestApiKey {
				w.WriteHeader(http.StatusUnauthorized)

				fmt.Printf("[%v] Unauthorized HTTP request: %s \"%s %s\" \"%s\"\n", time.Now().UTC(), r.RemoteAddr, r.Method, r.URL.Path, r.UserAgent())
				fmt.Fprintf(w, "Unauthorized: Missing or wrong X-API-Key header")

				return false
			}
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Last-Modified", lastUpdate.Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)

		return true
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return false
	}
}

// Handle sensors request.
func handleSensorsRequest(w http.ResponseWriter, r *http.Request, gpsData *gps.Gps, motion *motion.Motion) {
	lastTpvUpdate, lat, lon, alt, speed, mode, status, epc, epd, eph, eps, ept, epx, epy, epv, sep := gpsData.GetTpv()
	lastSkyUpdate, qual, xdop, ydop, vdop, tdop, hdop, pdop, gdop, nsat, usat, satellites := gpsData.GetSky()
	lastMotionUpdate, x, y, z, scale := motion.Get()

	lastUpdate := lastTpvUpdate

	if lastSkyUpdate.After(lastUpdate) {
		lastUpdate = lastSkyUpdate
	}

	if lastMotionUpdate.After(lastUpdate) {
		lastUpdate = lastMotionUpdate
	}

	if !prepareRequest(w, r, lastUpdate) {
		return
	}

	data := SensorsJson{
		Gps: GpsJson{
			Tpv: TpvJson{
				Lat:    lat,
				Lon:    lon,
				Alt:    alt,
				Speed:  speed,
				Mode:   mode,
				Status: status,
				Epc:    epc,
				Epd:    epd,
				Eph:    eph,
				Eps:    eps,
				Ept:    ept,
				Epx:    epx,
				Epy:    epy,
				Epv:    epv,
				Sep:    sep,
			},
			Sky: SkyJson{
				Qual:       qual,
				Xdop:       xdop,
				Ydop:       ydop,
				Vdop:       vdop,
				Tdop:       tdop,
				Hdop:       hdop,
				Pdop:       pdop,
				Gdop:       gdop,
				Nsat:       nsat,
				Usat:       usat,
				Satellites: createSatelliteJson(satellites),
			},
		},
		Motion: MotionJson{
			X:     x,
			Y:     y,
			Z:     z,
			Scale: scale,
		},
	}

	jsonString, _ := json.Marshal(&data)
	fmt.Fprint(w, string(jsonString))
}

// Handle sensor/gps request.
func handleSensorsGpsRequest(w http.ResponseWriter, r *http.Request, data *gps.Gps) {
	lastTpvUpdate, lat, lon, alt, speed, mode, status, epc, epd, eph, eps, ept, epx, epy, epv, sep := data.GetTpv()
	lastSkyUpdate, qual, xdop, ydop, vdop, tdop, hdop, pdop, gdop, nsat, usat, satellites := data.GetSky()

	lastUpdate := lastTpvUpdate

	if lastSkyUpdate.After(lastUpdate) {
		lastUpdate = lastSkyUpdate
	}

	if !prepareRequest(w, r, lastUpdate) {
		return
	}

	jsonData := GpsJson{
		Tpv: TpvJson{
			Lat:    lat,
			Lon:    lon,
			Alt:    alt,
			Speed:  speed,
			Mode:   mode,
			Status: status,
			Epc:    epc,
			Epd:    epd,
			Eph:    eph,
			Eps:    eps,
			Ept:    ept,
			Epx:    epx,
			Epy:    epy,
			Epv:    epv,
			Sep:    sep,
		},
		Sky: SkyJson{
			Qual:       qual,
			Xdop:       xdop,
			Ydop:       ydop,
			Vdop:       vdop,
			Tdop:       tdop,
			Hdop:       hdop,
			Pdop:       pdop,
			Gdop:       gdop,
			Nsat:       nsat,
			Usat:       usat,
			Satellites: createSatelliteJson(satellites),
		},
	}

	jsonString, _ := json.Marshal(&jsonData)
	fmt.Fprint(w, string(jsonString))
}

// Handle sensor/gps/tpv request.
func handleSensorsGpsTpvRequest(w http.ResponseWriter, r *http.Request, data *gps.Gps) {
	lastUpdate, lat, lon, alt, speed, mode, status, epc, epd, eph, eps, ept, epx, epy, epv, sep := data.GetTpv()

	if !prepareRequest(w, r, lastUpdate) {
		return
	}

	jsonData := TpvJson{
		Lat:    lat,
		Lon:    lon,
		Alt:    alt,
		Speed:  speed,
		Mode:   mode,
		Status: status,
		Epc:    epc,
		Epd:    epd,
		Eph:    eph,
		Eps:    eps,
		Ept:    ept,
		Epx:    epx,
		Epy:    epy,
		Epv:    epv,
		Sep:    sep,
	}

	jsonString, _ := json.Marshal(&jsonData)
	fmt.Fprint(w, string(jsonString))
}

// Handle sensor/gps/sky request.
func handleSensorsGpsSkyRequest(w http.ResponseWriter, r *http.Request, data *gps.Gps) {
	lastUpdate, qual, xdop, ydop, vdop, tdop, hdop, pdop, gdop, nsat, usat, satellites := data.GetSky()

	if !prepareRequest(w, r, lastUpdate) {
		return
	}

	jsonData := SkyJson{
		Qual:       qual,
		Xdop:       xdop,
		Ydop:       ydop,
		Vdop:       vdop,
		Tdop:       tdop,
		Hdop:       hdop,
		Pdop:       pdop,
		Gdop:       gdop,
		Nsat:       nsat,
		Usat:       usat,
		Satellites: createSatelliteJson(satellites),
	}

	jsonString, _ := json.Marshal(&jsonData)
	fmt.Fprint(w, string(jsonString))
}

// Handle sensor/motion request.
func handleSensorsMotionRequest(w http.ResponseWriter, r *http.Request, data *motion.Motion) {
	lastUpdate, x, y, z, scale := data.Get()

	if !prepareRequest(w, r, lastUpdate) {
		return
	}

	jsonData := MotionJson{
		X:     x,
		Y:     y,
		Z:     z,
		Scale: scale,
	}

	jsonString, _ := json.Marshal(&jsonData)
	fmt.Fprint(w, string(jsonString))
}

// Create the satellite JSON struct.
func createSatelliteJson(satellites []gpsd.Satellite) []SatelliteJson {
	satelliteJson := []SatelliteJson{}

	for _, satellite := range satellites {
		satelliteJson = append(satelliteJson, SatelliteJson{
			Prn:    satellite.PRN,
			Az:     satellite.Az,
			El:     satellite.El,
			Ss:     satellite.Ss,
			Gnssid: satellite.Gnssid,
			Used:   satellite.Used,
		})
	}

	return satelliteJson
}

// Start the HTTP REST server and listen for incoming requests.
func (data *Rest) Start(gpsData *gps.Gps, motionData *motion.Motion, done chan struct{}) error {
	fmt.Printf("Starting HTTP server... ")

	http.HandleFunc("/sensors", func(w http.ResponseWriter, r *http.Request) {
		handleSensorsRequest(w, r, gpsData, motionData)
	})

	http.HandleFunc("/sensors/gps", func(w http.ResponseWriter, r *http.Request) {
		handleSensorsGpsRequest(w, r, gpsData)
	})

	http.HandleFunc("/sensors/gps/tpv", func(w http.ResponseWriter, r *http.Request) {
		handleSensorsGpsTpvRequest(w, r, gpsData)
	})

	http.HandleFunc("/sensors/gps/sky", func(w http.ResponseWriter, r *http.Request) {
		handleSensorsGpsSkyRequest(w, r, gpsData)
	})

	http.HandleFunc("/sensors/motion", func(w http.ResponseWriter, r *http.Request) {
		handleSensorsMotionRequest(w, r, motionData)
	})

	listen, err := net.Listen("tcp", ":"+strconv.FormatUint(*global.RestPort, 10))

	if err != nil {
		fmt.Printf("Fail\n")
		fmt.Fprintf(os.Stderr, "Cannot start HTTP server: %s\n", err)

		close(done)

		return err
	}

	fmt.Printf("OK\n")

	serve := make(chan struct{})

	global.Wg.Add(2)

	go func() {
		err := http.Serve(listen, nil)

		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			fmt.Fprintf(os.Stderr, "Error serving content: %s\n", err)
			close(serve)
		}

		global.Wg.Done()
	}()

	go func() {
		select {
		case <-serve:
			close(done)
			listen.Close()
		case <-done:
			listen.Close()
		}

		global.Wg.Done()
	}()

	return nil
}
