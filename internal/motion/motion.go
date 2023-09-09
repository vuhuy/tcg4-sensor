package motion

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vuhuy/tcg4-sensor/internal/global"
	"github.com/vuhuy/tcg4-sensor/internal/helper"
)

const motionSensorXyzPath = "/sys/bus/i2c/devices/0-0018/xyz"
const motionSensorControlPath = "/sys/bus/i2c/devices/0-0018/control"

var cachedScaleValue uint8 = 0

// Struct to store motion data from I2C Sysfs.
type Motion struct {
	Mutex      sync.RWMutex
	LastUpdate time.Time
	X          int16
	Y          int16
	Z          int16
	Scale      uint8
}

// Store motion data with mutex lock.
func (data *Motion) Store(x, y, z int16, scale uint8) {
	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	data.LastUpdate = time.Now()
	data.X = x
	data.Y = y
	data.Z = z
	data.Scale = scale
}

// Get motion data with mutex lock.
func (data *Motion) Get() (time.Time, int16, int16, int16, uint8) {
	data.Mutex.RLock()
	defer data.Mutex.RUnlock()

	return data.LastUpdate, data.X, data.Y, data.Z, data.Scale
}

// Check if file exists.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// Start updating the motion data by reading the motion sensor periodically.
func (data *Motion) Start(done chan struct{}) error {
	fmt.Printf("Starting motion monitor... ")

	ticker := time.NewTicker(time.Duration(helper.ConvertHzToMilliseconds(*global.MotionFrequency)) * time.Millisecond)

	if !fileExists(motionSensorXyzPath) || !fileExists(motionSensorControlPath) {
		fmt.Printf("Fail\n")
		fmt.Fprintf(os.Stderr, "Failed to read motion sensor: snsor i2c Sysfs file does not exists")

		close(done)

		return errors.New("sensor i2c Sysfs file does not exists")
	}

	fmt.Print("OK\n")

	global.Wg.Add(1)

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				global.Wg.Done()
				return
			case <-ticker.C:
				xyzSuccess, x, y, z := readXyz()

				if cachedScaleValue == 0 {
					scaleSuccess, scale := readScale()

					if xyzSuccess && scaleSuccess {
						cachedScaleValue = scale
						data.Store(x, y, z, scale)
					}
				} else if xyzSuccess {
					data.Store(x, y, z, cachedScaleValue)
				}
			}
		}
	}()

	return nil
}

// Read xyz values from I2C Sysfs.
func readXyz() (bool, int16, int16, int16) {
	file, openErr := os.Open(motionSensorXyzPath)

	if openErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to open motion sensor: %s\n", openErr)

		return false, 0, 0, 0
	}

	defer file.Close()

	xyz, readErr := io.ReadAll(file)

	if readErr != nil {
		fmt.Fprintf(os.Stderr, "Cannot read motion sensor: %s\n", readErr)

		return false, 0, 0, 0
	}

	xyzString := string(xyz[:])
	xyzSplit := strings.Split(xyzString, ";")

	for i := range xyzSplit {
		xyzSplit[i] = strings.TrimSpace(xyzSplit[i])
	}

	if len(xyzSplit) != 3 {
		fmt.Fprintf(os.Stderr, "Error parsing motion sensor data: %s\n", xyzString)

		return false, 0, 0, 0
	}

	x, xParseErr := strconv.ParseInt(xyzSplit[0], 10, 16)
	y, yParseErr := strconv.ParseInt(xyzSplit[1], 10, 16)
	z, zParseErr := strconv.ParseInt(xyzSplit[2], 10, 16)

	if xParseErr != nil || yParseErr != nil || zParseErr != nil {
		fmt.Fprintf(os.Stderr, "Error parsing motion sensor data: %s\n", xyzString)

		return false, 0, 0, 0
	}

	if *global.Verbose {
		fmt.Printf("[%v] Read motion sensor XYZ: (%d; %d; %d)\n", time.Now().UTC(), x, y, z)
	}

	return true, int16(x), int16(y), int16(z)
}

// Read scale from I2C Sysfs.
func readScale() (bool, uint8) {
	file, openErr := os.Open(motionSensorControlPath)

	if openErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to open motion settings: %s\n", openErr)

		return false, 0
	}

	defer file.Close()

	control, readErr := io.ReadAll(file)

	if readErr != nil {
		fmt.Fprintf(os.Stderr, "Cannot read motion settings: %s\n", readErr)

		return false, 0
	}

	scaleRegex := regexp.MustCompile(`scale=(\d*)g`)
	scaleMatches := scaleRegex.FindAllStringSubmatch(string(control[:]), -1)

	if len(scaleMatches) != 1 || len(scaleMatches[0]) != 2 {
		fmt.Fprintf(os.Stderr, "Error parsing motion settings: cannot find scale\n")

		return false, 0
	}

	scale, parseErr := strconv.Atoi(scaleMatches[0][1])

	if parseErr != nil {
		fmt.Fprintf(os.Stderr, "Error parsing motion settings: %s\n", scaleMatches[0][1])

		return false, 0
	}

	if *global.Verbose {
		fmt.Printf("[%v] Read motion sensor scale: %d\n", time.Now().UTC(), scale)
	}

	return true, uint8(scale)
}
