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
	hdr := readHeader(data)
	fmt.Printf("%v\n", hdr)

	wfmt := readFmt(data)
	fmt.Printf("%v\n", wfmt)

	wavdata := readData(data, wfmt)
	fmt.Printf("%v\n", wavdata)
	return nil, nil
}

func bits16ToInt(b []byte) int {
	if len(b) != 2 {
		panic("Expected size 4!")
	}
	var payload uint16
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &payload)
	if err != nil {
		// TODO: make safe
		panic(err)
	}
	return int(payload) // easier to work with ints
}

// turn a 32-bit byte array into an int
func bits32ToInt(b []byte) int {
	if len(b) != 4 {
		panic("Expected size 4!")
	}
	var payload uint32
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &payload)
	if err != nil {
		// TODO: make safe
		panic(err)
	}
	return int(payload) // easier to work with ints
}

func readData(b []byte, wfmt WaveFmt) WaveData {
	wd := WaveData{}

	start := 36 + wfmt.ExtraParamSize
	subchunk2ID := b[start : start+4]
	wd.Subchunk2ID = subchunk2ID

	subsize := bits32ToInt(b[start+8 : start+12])
	wd.Subchunk2Size = subsize

	wd.Data = b[start+12:]

	return wd
}

// readFmt parses the FMT portion of the WAVE file
// assumes the entire binary representation is passed!
func readFmt(b []byte) WaveFmt {
	wfmt := WaveFmt{}
	subchunk1ID := b[12:16]
	wfmt.Subchunk1ID = subchunk1ID

	subchunksize := bits32ToInt(b[16:20])
	wfmt.Subchunk1Size = subchunksize

	format := bits16ToInt(b[20:22])
	wfmt.AudioFormat = format

	numChannels := bits16ToInt(b[22:24])
	wfmt.NumChannels = numChannels

	sr := bits32ToInt(b[24:28])
	wfmt.SampleRate = sr

	br := bits32ToInt(b[28:32])
	wfmt.ByteRate = br

	ba := bits16ToInt(b[32:34])
	wfmt.BlockAlign = ba

	bps := bits16ToInt(b[34:36])
	wfmt.BitsPerSample = bps

	// parse extra (optional) elements..

	if subchunksize != 16 {
		// only for compressed files (non-PCM)
		extraSize := bits16ToInt(b[36:38])
		wfmt.ExtraParamSize = extraSize
		wfmt.ExtraParams = b[38 : 38+extraSize]
	}

	return wfmt
}

// TODO: make safe.
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
