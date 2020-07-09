package internal

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
)

var (
	ChunkID          = []byte{0x52, 0x49, 0x46, 0x46}
	BigEndianChunkID = []byte{0x52, 0x49, 0x46, 0x58}
	Subchunk2ID      = []byte{0x64, 0x61, 0x74, 0x61}
)

// WriteSamples writes the slice to disk as a .wav file
// the WaveFmt metadata needs to be correct
// WaveData and WaveHeader are inferred from the samples however..
func WriteSamples(samples []Sample, wfmt WaveFmt, file string) error {

	// construct this in reverse (Data -> Fmt -> Header0
	// as Fmt needs info of Data, and Hdr needs to know entire length of file

	// write chunkSize
	bits := []byte{}
	//hdr := sampleToHeader(samples, wfmt)
	//	bits = append(bits, hdr...)

	data := samplesToData(samples, wfmt)
	bits = append(bits, data...)

	// 1. Write header
	// 2. Write FMT
	// 3. Write Data
	return ioutil.WriteFile("out.wav", bits, 0644)
}

func int32ToBytes(i int) []byte {
	b := make([]byte, 4)
	in := uint32(i)
	binary.LittleEndian.PutUint32(b, in)
	return b
}

func samplesToData(samples []Sample, wfmt WaveFmt) []byte {
	b := []byte{}
	raw := samplesToRawData(samples, wfmt)

	fmt.Printf("raw length: %v\n", len(raw))
	bytesPerSample := wfmt.BitsPerSample / 8
	subchunksize := len(samples) * wfmt.NumChannels * bytesPerSample
	subBytes := int32ToBytes(subchunksize)

	// construct the data part..
	b = append(b, Subchunk2ID...)
	b = append(b, subBytes...)
	b = append(b, raw...)
	return b
}

func floatToBytes(f float64, nBytes int) []byte {
	bits := math.Float64bits(f)
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, bits)
	// trim padding
	switch nBytes {
	case 2:
		return bs[:2]
	case 4:
		return bs[:4]
	}
	return bs
}

// Turn the samples into raw data...
func samplesToRawData(samples []Sample, props WaveFmt) []byte {
	raw := []byte{}
	for _, s := range samples {
		bits := floatToBytes(float64(s), props.BitsPerSample/8)
		raw = append(raw, bits...)
	}
	return raw
}

// turn the sample to a valid header
func sampleToHeader(samples []Sample, props WaveFmt) []byte {
	// write chunkID
	bits := []byte{}
	bits = append(bits, ChunkID...) // in theory switch on endianness..

	return bits
}
