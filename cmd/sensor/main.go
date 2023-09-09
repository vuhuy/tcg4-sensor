package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/vuhuy/tcg4-sensor/internal/can"
	"github.com/vuhuy/tcg4-sensor/internal/global"
	"github.com/vuhuy/tcg4-sensor/internal/gps"
	"github.com/vuhuy/tcg4-sensor/internal/helper"
	"github.com/vuhuy/tcg4-sensor/internal/motion"
	"github.com/vuhuy/tcg4-sensor/internal/rest"
)

var AppVersion string = "development"

// Main.
func main() {
	flag.Parse()
	helper.PrintVersion(AppVersion)

	// Handle graceful shutdown using SIGINT or SIGTERM.
	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	global.Wg.Add(2)

	go func() {
		select {
		case <-sigs:
			fmt.Printf("\nReceived interrupt or termination signal\n")
			close(done)
			global.Wg.Done()
		case <-done:
			global.Wg.Done()
		}
	}()

	go func() {
		// Initialize GPS.
		gpsData := gps.Gps{
			Tpv: gps.TpvReport{
				LastUpdate: time.Time{},
			},
			Sky: gps.SkyReport{
				LastUpdate: time.Time{},
			},
		}

		gpsErr := gpsData.Start(done)

		if gpsErr != nil {
			os.Exit(2)
		}

		// Initialize motion.
		motion := motion.Motion{
			LastUpdate: time.Time{},
		}

		motionErr := motion.Start(done)

		if motionErr != nil {
			os.Exit(3)
		}

		// Initialize CAN.
		canData := can.Can{}

		canErr := canData.Start(&gpsData, &motion, done)

		if canErr != nil {
			os.Exit(4)
		}

		// Initialize HTTP REST server.
		restData := rest.Rest{}

		restErr := restData.Start(&gpsData, &motion, done)

		if restErr != nil {
			os.Exit(5)
		}

		// Ready.
		fmt.Printf("CAN sensor data is currently being send on %s at a frequency of %f Hz\n", *global.CanInterface, *global.CanFrequency)
		fmt.Printf("An HTTP REST server running on %s is available for requesting sensor data\n", ":"+strconv.FormatUint(*global.RestPort, 10))

		global.Wg.Done()
	}()

	// Application end.
	<-done

	fmt.Printf("Exiting, waiting for background processes to finish... ")

	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	global.Wg.Wait()

	fmt.Printf("Bye!\n")

	os.Exit(1)
}
