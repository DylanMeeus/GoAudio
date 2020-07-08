package internal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func Test() {
	fmt.Println("hello world")
}

type WaveHeader struct {
	ChunkID   []byte // should be RIFF on little-endian or RIFX on big-endian systems..
	ChunkSize int
	Format    string // sanity-check, should be WAVE
}

// ParseFloatFrames for audio
func ReadFloatFrames(f string) ([]float32, error) {
	// open as read-only file
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// determine size?
	info, _ := file.Stat()
	fmt.Printf("size: %v\n", info.Size)

	data := make([]byte, info.Size())
	bytesread, err := file.Read(data)
	fmt.Printf("Bytes read: %v\n", bytesread)
	readHeader(data)
	return nil, nil
}

func readHeader(b []byte) WaveHeader {
	// the start of the bte slice..
	hdr := WaveHeader{}
	chunkID := b[0:4]
	hdr.ChunkID = b[0:4]
	fmt.Printf("chunkID: %v (%v) \n", chunkID, string(chunkID))
	if string(hdr.ChunkID) != "RIFF" {
		panic("Invalid file")
	}

	chunkSize := b[4:8]
	fmt.Printf("chunkSize: %v\n", chunkSize)
	var size uint32
	buf := bytes.NewReader(chunkSize)
	err := binary.Read(buf, binary.LittleEndian, &size)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", size)
	hdr.ChunkSize = int(size) // easier to work with ints

	format := b[8:12]
	if string(format) != "WAVE" {
		panic("Format should be WAVE")
	}
	hdr.Format = string(format)
	fmt.Printf("format: %v (%v)\n", format, hdr.Format)
	return hdr
}
