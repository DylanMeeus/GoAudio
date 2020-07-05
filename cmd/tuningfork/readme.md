# Tuning fork

This small sample programs allows you to generate a sine-wave signal of variable duration,
frequency, sample rate and amplitude.

A 32bit floating-point (raw) binary file will be created, this can be listened to with ffplay /  Audacity or anything else supporting these files.


**Only works on LittleEndian systems, change the binary encoding to BigEndian for BigEndian
systems!**

An exponential decay is applied to the signal.

```
go run main.go [duration] [frequency] [samplerate] [amplitude] [outfile]
go run main.go 2 440 44100 1 out.bin
```

## ffplay

Play with ffplay:

```
ffplay -f f32le -ar 44100 out.bin
```

Or visualize with:

```
ffplay -f f32le -ar 44100 -showmode 1 out.bin
```

