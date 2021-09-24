package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/gordonklaus/portaudio"
)

// stream sound to two soundcards
func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing required argument:  input file name")
		return
	}
	fmt.Println("Playing.  Press Ctrl-C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	fileName := os.Args[1]
	f, err := os.Open(fileName)
	chk(err)
	defer f.Close()

	id, data, err := readChunk(f)
	chk(err)
	if id.String() != "FORM" {
		fmt.Println("bad file format")
		return
	}
	_, err = data.Read(id[:])
	chk(err)
	if id.String() != "AIFF" {
		fmt.Println("bad file format")
		return
	}
	var c commonChunk
	var audio io.Reader
	for {
		id, chunk, err := readChunk(data)
		if err == io.EOF {
			break
		}
		chk(err)
		switch id.String() {
		case "COMM":
			chk(binary.Read(chunk, binary.BigEndian, &c))
		case "SSND":
			chunk.Seek(8, 1) //ignore offset and block
			audio = chunk
		default:
			fmt.Printf("ignoring unknown chunk '%s'\n", id)
		}
	}

	//assume 48000 sample rate, mono, 32 bit, AIFF formatted file
	portaudio.Initialize()
	defer portaudio.Terminate()

	out := make([]int32, 8192)
	outA := make([]int32, len(out))
	outB := make([]int32, len(out))

	streamA, err := openStream("CTAG", 0, 1, 48000, len(outA), &outA)
	chk(err)
	defer streamA.Close()
	chk(streamA.Start())
	defer streamA.Stop()

	streamB, err := openStream("RAVENNA", 0, 1, 48000, len(outB), &outB)
	chk(err)
	defer streamB.Close()
	chk(streamB.Start())
	defer streamB.Stop()

	for remaining := int(c.NumSamples); remaining > 0; remaining -= len(out) {
		if len(out) > remaining {
			out = out[:remaining]
		}
		err := binary.Read(audio, binary.BigEndian, out)
		if err == io.EOF {
			break
		}
		chk(err)

		// copy file buffer to stream buffers
		copy(outA, out)
		copy(outB, out)
		chk(streamA.Write())
		chk(streamB.Write())
		copy(out, outA)

		select {
		case <-sig:
			return
		default:
		}
	}
}

func readChunk(r readerAtSeeker) (id ID, data *io.SectionReader, err error) {
	_, err = r.Read(id[:])
	if err != nil {
		return
	}
	var n int32
	err = binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return
	}
	off, _ := r.Seek(0, 1)
	data = io.NewSectionReader(r, off, int64(n))
	_, err = r.Seek(int64(n), 1)
	return
}

type readerAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

type ID [4]byte

func (id ID) String() string {
	return string(id[:])
}

type commonChunk struct {
	NumChans      int16
	NumSamples    int32
	BitsPerSample int16
	SampleRate    [10]byte
}
