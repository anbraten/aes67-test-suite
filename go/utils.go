package main

import (
	"fmt"
	"strings"

	bbhw "github.com/btittelbach/go-bbhw"
	"github.com/gordonklaus/portaudio"
)

func getDevice(name string) (*portaudio.DeviceInfo, error) {
	devices, err := portaudio.Devices()
	chk(err)

	for _, d := range devices {
		// fmt.Println(d.Name)
		if strings.Contains(d.Name, name) {
			return d, nil
		}
	}

	return nil, fmt.Errorf("Can't find device")
}

func openStream(deviceName string, inChannels int, outChannels int, sampleRate float64, framesPerBuffer int, args ...interface{}) (*portaudio.Stream, error) {
	device, err := getDevice(deviceName)
	chk(err)
	fmt.Println("Found device", device.Name, "#", device.MaxOutputChannels)

	var p portaudio.StreamParameters
	if inChannels == 0 && outChannels == 0 {
		return nil, fmt.Errorf("Input or output has to be used")
	} else if inChannels != 0 {
		p = portaudio.LowLatencyParameters(device, nil) // input
	} else {
		p = portaudio.LowLatencyParameters(nil, device) // output
	}

	p.Input.Channels = inChannels
	p.Output.Channels = outChannels
	p.SampleRate = sampleRate
	p.FramesPerBuffer = framesPerBuffer
	return portaudio.OpenStream(p, args...)
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func enablePin(pin *bbhw.MMappedGPIO, enable bool) {
	err := pin.SetState(enable)
	chk(err)
}
