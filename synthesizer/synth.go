package synthesizer

// ADSR creates an attack -> decay -> sustain -> release envelope
func ADSR(maxamp, duration, attacktime, decaytime, sus, releasetime, controlrate float64, currenttime int) float64 {
	dur := duration * controlrate
	at := attacktime * controlrate
	dt := decaytime * controlrate
	rt := releasetime * controlrate
	cnt := float64(currenttime)

	amp := 0.0
	if cnt < dur {
		if cnt <= at {
			// attack
			amp = cnt * (maxamp / at)
		} else if cnt <= (at + dt) {
			// decay
			amp = ((sus-maxamp)/dt)*(cnt-at) + maxamp
		} else if cnt <= dur-rt {
			// sustain
			amp = sus
		} else if cnt > (dur - rt) {
			// release
			amp = -(sus/rt)*(cnt-(dur-rt)) + sus
		}
	}

	return amp
}
