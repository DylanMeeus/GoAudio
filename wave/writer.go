package wave

import (
	"encoding/binary"
	"io"
	"math"
	"os"
)

// Consts that appear in the .WAVE file format
var (
	ChunkID          = []byte{0x52, 0x49, 0x46, 0x46} // RIFF
	BigEndianChunkID = []byte{0x52, 0x49, 0x46, 0x58} // RIFX
	WaveID           = []byte{0x57, 0x41, 0x56, 0x45} // WAVE
	Format           = []byte{0x66, 0x6d, 0x74, 0x20} // FMT
	Subchunk2ID      = []byte{0x64, 0x61, 0x74, 0x61} // DATA
)

type intsToBytesFunc func(i int) []byte

var (
	// intsToBytesFm to map X-bit int to byte functions
	intsToBytesFm = map[int]intsToBytesFunc{
		16: int16ToBytes,
		32: int32ToBytes,
	}
)

// WriteFrames writes the slice to disk as a .wav file
// the WaveFmt metadata needs to be correct
// WaveData and WaveHeader are inferred from the samples however..
func WriteFrames(samples []Frame, wfmt WaveFmt, file string) error {
	return WriteWaveFile(samples, wfmt, file)
}

func WriteWaveFile(samples []Frame, wfmt WaveFmt, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return WriteWaveToWriter(samples, wfmt, f)
}

func WriteWaveToWriter(samples []Frame, wfmt WaveFmt, writer io.Writer) error {
	wfb := fmtToBytes(wfmt)
	data, databits := framesToData(samples, wfmt)
	hdr := createHeader(data)

	_, err := writer.Write(hdr)
	if err != nil {
		return err
	}
	_, err = writer.Write(wfb)
	if err != nil {
		return err
	}
	_, err = writer.Write(databits)
	if err != nil {
		return err
	}

	return nil
}

func int16ToBytes(i int) []byte {
	b := make([]byte, 2)
	in := uint16(i)
	binary.LittleEndian.PutUint16(b, in)
	return b
}

func int32ToBytes(i int) []byte {
	b := make([]byte, 4)
	in := uint32(i)
	binary.LittleEndian.PutUint32(b, in)
	return b
}

func framesToData(frames []Frame, wfmt WaveFmt) (WaveData, []byte) {
	b := []byte{}
	raw := samplesToRawData(frames, wfmt)

	// We receive frames but have to store the size of the samples
	// The size of the samples is frames / channels..
	subchunksize := (len(frames) * wfmt.NumChannels * wfmt.BitsPerSample) / 8
	subBytes := int32ToBytes(subchunksize)

	// construct the data part..
	b = append(b, Subchunk2ID...)
	b = append(b, subBytes...)
	b = append(b, raw...)

	wd := WaveData{
		Subchunk2ID:   Subchunk2ID,
		Subchunk2Size: subchunksize,
		RawData:       raw,
		Frames:        frames,
	}
	return wd, b
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
func samplesToRawData(samples []Frame, props WaveFmt) []byte {
	raw := []byte{}
	for _, s := range samples {
		// the samples are scaled - rescale them?
		rescaled := rescaleFrame(s, props.BitsPerSample)
		bits := intsToBytesFm[props.BitsPerSample](rescaled)
		raw = append(raw, bits...)
	}
	return raw
}

// rescale frames back to the original values..
func rescaleFrame(s Frame, bits int) int {
	rescaled := float64(s) * float64(maxValues[bits])
	return int(rescaled)
}

func fmtToBytes(wfmt WaveFmt) []byte {
	b := []byte{}

	subchunksize := int32ToBytes(wfmt.Subchunk1Size)
	audioformat := int16ToBytes(wfmt.AudioFormat)
	numchans := int16ToBytes(wfmt.NumChannels)
	sr := int32ToBytes(wfmt.SampleRate)
	br := int32ToBytes(wfmt.ByteRate)
	blockalign := int16ToBytes(wfmt.BlockAlign)
	bitsPerSample := int16ToBytes(wfmt.BitsPerSample)

	b = append(b, wfmt.Subchunk1ID...)
	b = append(b, subchunksize...)
	b = append(b, audioformat...)
	b = append(b, numchans...)
	b = append(b, sr...)
	b = append(b, br...)
	b = append(b, blockalign...)
	b = append(b, bitsPerSample...)

	return b
}

// turn the sample to a valid header
func createHeader(wd WaveData) []byte {
	// write chunkID
	bits := []byte{}

	chunksize := 36 + wd.Subchunk2Size
	cb := int32ToBytes(chunksize)

	bits = append(bits, ChunkID...) // in theory switch on endianness..
	bits = append(bits, cb...)
	bits = append(bits, WaveID...)

	return bits
}
