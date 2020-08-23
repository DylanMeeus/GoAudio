package wave

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

// type aliases for conversion functions
type (
	bytesToIntF   func([]byte) int
	bytesToFloatF func([]byte) float64
)

var (
	// figure out which 'to int' function to use..
	byteSizeToIntFunc = map[int]bytesToIntF{
		16: bits16ToInt,
		32: bits32ToInt,
	}

	byteSizeToFloatFunc = map[int]bytesToFloatF{
		16: bitsToFloat,
		32: bitsToFloat,
		64: bitsToFloat,
	}

	// max value depending on the bit size
	maxValues = map[int]int{
		8:  math.MaxInt8,
		16: math.MaxInt16,
		32: math.MaxInt32,
		64: math.MaxInt64,
	}
)

// ReadWaveFile parses a .wave file into a Wave struct
func ReadWaveFile(f string) (Wave, error) {
	// open as read-only file
	file, err := os.Open(f)
	if err != nil {
		return Wave{}, err
	}
	defer file.Close()

	// determine size?
	info, _ := file.Stat()
	fmt.Printf("size: %v\n", info.Size())

	data, err := ioutil.ReadFile(f)
	if err != nil {
		return Wave{}, err
	}
	fmt.Printf("Bytes read: %v\n", len(data))
	hdr := readHeader(data)

	wfmt := readFmt(data)

	wavdata := readData(data, wfmt)

	frames := parseRawData(wfmt, wavdata.RawData)
	wavdata.Frames = frames

	return Wave{
		WaveHeader: hdr,
		WaveFmt:    wfmt,
		WaveData:   wavdata,
	}, nil
}

// for our wave format we expect double precision floats
func bitsToFloat(b []byte) float64 {
	var bits uint64
	switch len(b) {
	case 2:
		bits = uint64(binary.LittleEndian.Uint16(b))
	case 4:
		bits = uint64(binary.LittleEndian.Uint32(b))
	case 8:
		bits = binary.LittleEndian.Uint64(b)
	default:
		panic("Can't parse to float..")
	}
	float := math.Float64frombits(bits)
	return float
}

func bits16ToInt(b []byte) int {
	if len(b) != 2 {
		panic("Expected size 4!")
	}
	var payload int16
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
	var payload int32
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

	subsize := bits32ToInt(b[start+4 : start+8])
	wd.Subchunk2Size = subsize
	fmt.Printf("subchunk size: %v\n", wd.Subchunk2Size)

	wd.RawData = b[start+8:]

	return wd
}

// Should we do n-channel separation at this point?
func parseRawData(wfmt WaveFmt, rawdata []byte) []Frame {
	bytesSampleSize := wfmt.BitsPerSample / 8
	// TODO: sanity-check that this is a power of 2? I think only those sample sizes are
	// possible

	// TODO: READ AS INT??????

	frames := []Frame{}
	// read the chunks
	for i := 0; i < len(rawdata); i += bytesSampleSize {
		rawFrame := rawdata[i : i+bytesSampleSize]
		unscaledFrame := byteSizeToIntFunc[wfmt.BitsPerSample](rawFrame)
		scaled := scaleFrame(unscaledFrame, wfmt.BitsPerSample)
		frames = append(frames, scaled)
	}

	return frames
}

func scaleFrame(unscaled, bits int) Frame {
	maxV := maxValues[bits]
	return Frame(float64(unscaled) / float64(maxV))

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
