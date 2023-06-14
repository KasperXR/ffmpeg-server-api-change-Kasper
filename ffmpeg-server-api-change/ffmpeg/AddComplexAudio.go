package ffmpeg

import "strconv"

func (f *FFMPEGCommand) AddComplexAudio(audio *[]FFMPEGAudio) {
	var inputs = ""
	var command = ""
	var amixCount = ""

	// Initialise helper variables for durSoFar
	var durSoFar float64 = 0
	var timeTo float64 = 0
	var timeFrom float64 = 0

	// Make the input graph, add delay if applicable
	for idx, a := range *audio {

		if idx == 0 {
			timeFrom = 0
			inputs += ` `
			timeTo = a.Duration
		} else {
			timeFrom = durSoFar
			timeTo = durSoFar + a.Duration
		}

		inputs += `-i ` + a.Input + `.` + a.FileType + ` `

		command += `[` + strconv.Itoa(idx+1) + `]`
		command += `adelay=`

		if a.Delay != 0 {
			if idx == 0 {
				command += strconv.FormatFloat(timeFrom, 'f', 2, 64)
			} else {
				command += strconv.FormatFloat(timeFrom+a.Delay, 'f', 2, 64)
			}

			command += `|`
			command += strconv.FormatFloat(timeTo+a.Delay, 'f', 2, 64)
		} else {
			// No delay for this audio
			command += strconv.FormatFloat(timeFrom, 'f', 2, 64)
			command += `|`
			command += strconv.FormatFloat(timeTo, 'f', 2, 64)
		}

		command += `[a` + strconv.Itoa(idx+1) + `]`

		// if not last
		if idx != len(*audio)-1 {
			command += `;`
		}

		// Increment duration, and delay if applicable
		durSoFar += a.Duration

		if a.Delay != 0 && idx != 0 {
			durSoFar += a.Delay
		}

		// Define all audio inputs for ffmpeg graph
		amixCount += `[a` + strconv.Itoa(idx+1) + `]`
	}

	command = inputs + `-filter_complex` + ` "` + command + ";" + amixCount

	// Define final audio mixing
	command += ` amix=` + strconv.Itoa(len(*audio)) + `[a]`

	command += `"`

	f.Command += command
	f.HasComplexAudio = true
}
