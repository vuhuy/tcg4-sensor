package can

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/vuhuy/tcg4-sensor/internal/global"
	"github.com/vuhuy/tcg4-sensor/internal/gps"
	"github.com/vuhuy/tcg4-sensor/internal/helper"
	"github.com/vuhuy/tcg4-sensor/internal/motion"
	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

// Struct to store CAN frame data.
type Can struct{}

// Send a GPS float64 value in a single CAN frame.
func sendFloatFrame(id uint32, value float64, tx *socketcan.Transmitter) {
	if id == 0 {
		return
	}

	frame := can.Frame{}
	frame.ID = id
	frame.Length = 8
	frame.IsExtended = *global.CanExtended

	binary.LittleEndian.PutUint64(frame.Data[:], math.Float64bits(value))
	_ = tx.TransmitFrame(context.Background(), frame)

	if *global.Verbose {
		fmt.Printf("[%v] CAN frame sended: %#v\n", time.Now().UTC(), frame)
	}
}

// Send additional GPS mode, status, nsat, usat, and qual values in a single CAN frame
func sendGpsFrame(mode uint8, status uint8, nsat uint16, usat uint16, qual uint8, tx *socketcan.Transmitter) {
	if *global.GpsFrameId == 0 {
		return
	}

	frame := can.Frame{}
	frame.ID = uint32(*global.GpsFrameId)
	frame.Length = 8
	frame.IsExtended = *global.CanExtended

	var nsatBytes [2]byte
	var usatBytes [2]byte

	binary.LittleEndian.PutUint16(nsatBytes[:], nsat)
	binary.LittleEndian.PutUint16(usatBytes[:], usat)

	frame.Data[0] = mode
	frame.Data[1] = status
	frame.Data[2] = nsatBytes[0]
	frame.Data[3] = nsatBytes[1]
	frame.Data[4] = usatBytes[0]
	frame.Data[5] = usatBytes[1]
	frame.Data[6] = qual
	frame.Data[7] = 0x0

	_ = tx.TransmitFrame(context.Background(), frame)

	if *global.Verbose {
		fmt.Printf("[%v] CAN frame sended: %#v\n", time.Now().UTC(), frame)
	}
}

// Send the measured motion sensor X, Y, Z values and the used scale in a single CAN frame.
func sendMotionFrame(x, y, z int16, scale uint8, tx *socketcan.Transmitter) {
	if *global.MotionFrameId == 0 {
		return
	}

	frame := can.Frame{}
	frame.ID = uint32(*global.MotionFrameId)
	frame.Length = 8
	frame.IsExtended = *global.CanExtended

	var xBytes [2]byte
	var yBytes [2]byte
	var zBytes [2]byte

	binary.LittleEndian.PutUint16(xBytes[:], uint16(x))
	binary.LittleEndian.PutUint16(yBytes[:], uint16(y))
	binary.LittleEndian.PutUint16(zBytes[:], uint16(z))

	frame.Data[0] = xBytes[0]
	frame.Data[1] = xBytes[1]
	frame.Data[2] = yBytes[0]
	frame.Data[3] = yBytes[1]
	frame.Data[4] = zBytes[0]
	frame.Data[5] = zBytes[1]
	frame.Data[6] = scale
	frame.Data[7] = 0x0

	_ = tx.TransmitFrame(context.Background(), frame)

	if *global.Verbose {
		fmt.Printf("[%v] CAN frame sended: %#v\n", time.Now().UTC(), frame)
	}
}

// Start sending CAN messages periodically.
func (data *Can) Start(gpsData *gps.Gps, motionData *motion.Motion, done chan struct{}) error {
	fmt.Printf("Opening CAN interface... ")

	conn, err := socketcan.DialContext(context.Background(), "can", *global.CanInterface)

	if err != nil {
		fmt.Printf("Fail\n")
		fmt.Fprintf(os.Stderr, "Cannot connect to interface %s: %s\n", *global.CanInterface, err)

		close(done)

		return err
	}

	fmt.Printf("OK\n")

	tx := socketcan.NewTransmitter(conn)
	ticker := time.NewTicker(time.Duration(helper.ConvertHzToMilliseconds(*global.CanFrequency)) * time.Millisecond)

	global.Wg.Add(1)

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				conn.Close()
				global.Wg.Done()
				return
			case <-ticker.C:
				sendGpsFrames(gpsData, tx)
				sendMotionFrames(motionData, tx)
			}
		}
	}()

	return nil
}

// Send GPS related CAN frames.
func sendGpsFrames(gpsData *gps.Gps, tx *socketcan.Transmitter) {
	lastTpvUpdate, lat, lon, alt, speed, mode, status, epc, epd, eph, eps, ept, epx, epy, epv, sep := gpsData.GetTpv()
	lastSkyUpdate, qual, xdop, ydop, vdop, tdop, hdop, pdop, gdop, nsat, usat, _ := gpsData.GetSky()

	if lastTpvUpdate.IsZero() || lastSkyUpdate.IsZero() {
		return
	}

	sendGpsFrame(mode, status, nsat, usat, qual, tx)
	sendFloatFrame(uint32(*global.LatFrameId), lat, tx)
	sendFloatFrame(uint32(*global.LonFrameId), lon, tx)
	sendFloatFrame(uint32(*global.AltFrameId), alt, tx)
	sendFloatFrame(uint32(*global.SpeedFrameId), speed, tx)
	sendFloatFrame(uint32(*global.EpcFrameId), epc, tx)
	sendFloatFrame(uint32(*global.EpdFrameId), epd, tx)
	sendFloatFrame(uint32(*global.EphFrameId), eph, tx)
	sendFloatFrame(uint32(*global.EpsFrameId), eps, tx)
	sendFloatFrame(uint32(*global.EptFrameId), ept, tx)
	sendFloatFrame(uint32(*global.EpxFrameId), epx, tx)
	sendFloatFrame(uint32(*global.EpyFrameId), epy, tx)
	sendFloatFrame(uint32(*global.EpvFrameId), epv, tx)
	sendFloatFrame(uint32(*global.SepFrameId), sep, tx)
	sendFloatFrame(uint32(*global.XdopFrameId), xdop, tx)
	sendFloatFrame(uint32(*global.YdopFrameId), ydop, tx)
	sendFloatFrame(uint32(*global.VdopFrameId), vdop, tx)
	sendFloatFrame(uint32(*global.TdopFrameId), tdop, tx)
	sendFloatFrame(uint32(*global.HdopFrameId), hdop, tx)
	sendFloatFrame(uint32(*global.PdopFrameId), pdop, tx)
	sendFloatFrame(uint32(*global.GdopFrameId), gdop, tx)
}

// Send motion related CAN frames.
func sendMotionFrames(motionData *motion.Motion, tx *socketcan.Transmitter) {
	lastUpdate, x, y, z, scale := motionData.Get()

	if lastUpdate.IsZero() {
		return
	}

	sendMotionFrame(x, y, z, scale, tx)
}
