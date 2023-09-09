package gps

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/vuhuy/tcg4-sensor/internal/global"
	"github.com/vuhuy/tcg4-sensor/pkg/gpsd"
)

// Struct to store GPS data from GPSd.
type Gps struct {
	Mutex sync.RWMutex
	Tpv   TpvReport
	Sky   SkyReport
}

// Struct to store TPV report data.
type TpvReport struct {
	LastUpdate time.Time
	Lat        float64
	Lon        float64
	Alt        float64
	Speed      float64
	Mode       uint8
	Status     uint8
	Epc        float64
	Epd        float64
	Eph        float64
	Eps        float64
	Ept        float64
	Epx        float64
	Epy        float64
	Epv        float64
	Sep        float64
}

// Struct to store TPV report data.
type SkyReport struct {
	LastUpdate time.Time
	Xdop       float64
	Ydop       float64
	Vdop       float64
	Tdop       float64
	Hdop       float64
	Pdop       float64
	Gdop       float64
	Nsat       uint16
	Usat       uint16
	Qual       uint8
	Satellites []gpsd.Satellite
}

// Store TPV report with mutex lock.
func (data *Gps) StoreTpv(lat, lon, alt, speed float64, mode, status uint8, epc, epd, eph, eps, ept, epx, epy, epv, sep float64) {
	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	data.Tpv.LastUpdate = time.Now()
	data.Tpv.Lat = lat
	data.Tpv.Lon = lon
	data.Tpv.Alt = alt
	data.Tpv.Speed = speed
	data.Tpv.Mode = mode
	data.Tpv.Status = status
	data.Tpv.Epc = epc
	data.Tpv.Epd = epd
	data.Tpv.Eph = eph
	data.Tpv.Eps = eps
	data.Tpv.Ept = ept
	data.Tpv.Epx = epx
	data.Tpv.Epy = epy
	data.Tpv.Epv = epv
	data.Tpv.Sep = sep
}

// Store SKY report with mutex lock.
func (data *Gps) StoreSky(qual uint8, xdop, ydop, vdop, tdop, hdop, pdop, gdop float64, nsat, usat uint16, satellites []gpsd.Satellite) {
	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	data.Sky.LastUpdate = time.Now()
	data.Sky.Qual = qual
	data.Sky.Xdop = xdop
	data.Sky.Ydop = ydop
	data.Sky.Vdop = vdop
	data.Sky.Tdop = tdop
	data.Sky.Hdop = hdop
	data.Sky.Pdop = pdop
	data.Sky.Gdop = gdop
	data.Sky.Nsat = nsat
	data.Sky.Usat = usat
	data.Sky.Satellites = nil
	data.Sky.Satellites = append(data.Sky.Satellites, satellites...)
}

// Get TPV report with mutex lock.
func (data *Gps) GetTpv() (time.Time, float64, float64, float64, float64, uint8, uint8, float64, float64, float64, float64, float64, float64, float64, float64, float64) {
	data.Mutex.RLock()
	defer data.Mutex.RUnlock()

	return data.Tpv.LastUpdate, data.Tpv.Lat, data.Tpv.Lon, data.Tpv.Alt, data.Tpv.Speed, data.Tpv.Mode, data.Tpv.Status, data.Tpv.Epc, data.Tpv.Epd, data.Tpv.Eph, data.Tpv.Eps, data.Tpv.Ept, data.Tpv.Epx, data.Tpv.Epy, data.Tpv.Epv, data.Tpv.Sep
}

// Get SKY report with mutex lock.
func (data *Gps) GetSky() (time.Time, uint8, float64, float64, float64, float64, float64, float64, float64, uint16, uint16, []gpsd.Satellite) {
	data.Mutex.RLock()
	defer data.Mutex.RUnlock()

	return data.Sky.LastUpdate, data.Sky.Qual, data.Sky.Xdop, data.Sky.Ydop, data.Sky.Vdop, data.Sky.Tdop, data.Sky.Hdop, data.Sky.Pdop, data.Sky.Gdop, data.Sky.Nsat, data.Sky.Usat, data.Sky.Satellites
}

// Start updating the GPS data periodically using GPSd.
func (data *Gps) Start(done chan struct{}) error {
	fmt.Printf("Starting GPS monitor... ")

	gps, err := gpsd.Dial(*global.GpsdHost + ":" + strconv.FormatUint(*global.GpsdPort, 10))

	if err != nil {
		fmt.Printf("Fail\n")
		fmt.Fprintf(os.Stderr, "Failed to connect to GPSd on %s:%d %s\n", *global.GpsdHost, *global.GpsdPort, err)

		close(done)

		return err
	}

	gps.AddFilter("SKY", func(r interface{}) {
		sky := r.(*gpsd.SKYReport)
		data.StoreSky(sky.Qual, sky.Xdop, sky.Ydop, sky.Vdop, sky.Tdop, sky.Hdop, sky.Pdop, sky.Gdop, sky.Nsat, sky.Usat, sky.Satellites)

		if *global.Verbose {
			fmt.Printf("[%v] Read GPS SKY report: %#v\n", time.Now().UTC(), sky)
		}
	})

	gps.AddFilter("TPV", func(r interface{}) {
		tpv := r.(*gpsd.TPVReport)
		data.StoreTpv(tpv.Lat, tpv.Lon, tpv.Alt, tpv.Speed, uint8(tpv.Mode), tpv.Status, tpv.Epc, tpv.Epd, tpv.Eph, tpv.Eps, tpv.Ept, tpv.Epx, tpv.Epy, tpv.Epv, tpv.Sep)

		if *global.Verbose {
			fmt.Printf("[%v] Read GPS TPV report: %#v\n", time.Now().UTC(), tpv)
		}
	})

	fmt.Printf("OK\n")

	watch := gps.Watch()

	global.Wg.Add(1)

	go func() {
		select {
		case <-watch:
			close(done)
			gps.Close()
			global.Wg.Done()
		case <-done:
			gps.Close()
			global.Wg.Done()
		}
	}()

	return nil
}
