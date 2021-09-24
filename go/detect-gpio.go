package main

import (
	"fmt"
	"os"
	"os/signal"

	bbhw "github.com/btittelbach/go-bbhw"
	"github.com/gordonklaus/portaudio"
)

// set gpio pin to high after detecting incoming sound
func main() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	pin50 := bbhw.NewMMappedGPIO(50, bbhw.OUT)
	defer pin50.Close()

	enablePin(pin50, false)

	portaudio.Initialize()
	defer portaudio.Terminate()

	buffer := make([]float32, 8)
	stream, err := openStream("Merging RAVENNA", 1, 0, 48000, len(buffer), buffer)
	chk(err)

	fmt.Println("listening to stream, waiting for data ...")
	chk(stream.Start())
	defer stream.Close()

	nSamples := 0
	for {
		chk(stream.Read())
		nSamples += len(buffer)
		if nSamples > 0 {

			for _, v := range buffer {
				// check sound-level
				if v > 0 {
					enablePin(pin50, true)
					fmt.Println("received audio data", v)
					nSamples = 0
					return
				}
			}

		}
		select {
		case <-sig:
			return
		default:
		}
	}
	chk(stream.Stop())

	fmt.Println("end")
}
