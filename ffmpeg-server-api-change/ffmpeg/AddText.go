package ffmpeg

import "strconv"

func (f *FFMPEGCommand) AddText(txt *FFMPEGText, isLast bool) {
	var command = `drawtext="`

	// Add text
	if txt.TextFile {
		command += `textfile=` + txt.Data + `:`
	} else {
		command += `text=` + txt.Data + `:`
	}

	// Add font
	command += `fontfile=` + txt.FontFile + `:`

	// Add font size
	command += `fontsize=` + strconv.Itoa(txt.FontSize) + `:`

	// Add font color
	command += `fontcolor=` + txt.FontColor + `:`

	// Add lineheight
	command += `line_spacing=` + strconv.Itoa(txt.LineHeight) + `:`

	// Add x and y
	command += `x=` + strconv.Itoa(txt.X) + `:`
	command += `y=` + strconv.Itoa(txt.Y) + `:`

	// Add time from and time to
	if txt.HasDuration {
		command += `enable='between(t,` + strconv.FormatFloat(txt.TimeFrom, 'f', 2, 64) + `,` + strconv.FormatFloat(txt.TimeTo+txt.Delay, 'f', 2, 64) + `)'`
	} else {
		// Remove trailing colon
		command = command[:len(command)-1]
	}

	// Add fade in and out using an alpha channel
	// Then continue with the drawtext filter by adding the `,` separator
	// If FadeIn is undefined or 0, then don't add it
	// This is an example of how it should look:
	// alpha='if(lt(t,52.56),0,if(lt(t,53.86),(t-52.56)/2,if(lt(t,59.47),1,if(lt(t,61.47),1-(t-59.47)/2.00,0))))'
	// We want to fade the text in from transparent to opaque, and then fade it out from opaque to transparent

	// Fade in
	if txt.FadeIn > 0 && txt.FadeOut > 0 {

		// Fade in

		// Add alpha channel
		command += `:alpha='`

		// Add if statement for time less than time from
		command += `if(lt(t,` +
			strconv.FormatFloat(txt.TimeFrom+txt.Delay, 'f', 2, 64) + `),0,`

		// Add if statement for time less than time from + fade in
		command += `if(lt(t,` +
			strconv.FormatFloat(txt.TimeFrom+txt.Delay+txt.FadeIn, 'f', 2, 64) +
			`),(t-` + strconv.FormatFloat(txt.TimeFrom+txt.Delay, 'f', 2, 64) +
			`)/` + strconv.FormatFloat(txt.FadeIn, 'f', 2, 64) + `,`

		// Add if statement for time less than time to - fade out
		command += `if(lt(t,` +
			strconv.FormatFloat(txt.TimeTo+txt.Delay-txt.FadeOut, 'f', 2, 64) +
			`),1,`

		// Add if statement for time less than time to
		command += `if(lt(t,` +
			strconv.FormatFloat(txt.TimeTo+txt.Delay, 'f', 2, 64) +
			`),1-(t-` +
			strconv.FormatFloat(txt.TimeTo+txt.Delay-txt.FadeOut, 'f', 2, 64) +
			`)/` +
			strconv.FormatFloat(txt.FadeOut, 'f', 2, 64) +
			`,0))))'`

	} else {
	}

	//	if txt.FadeIn != 0.0 && txt.FadeOut != 0.0 {
	//		command += `:alpha='if(lt(t,` +
	//			strconv.FormatFloat(txt.TimeFrom, 'f', 2, 64) +
	//			`),0,if(lt(t,` +
	//			strconv.FormatFloat(txt.TimeFrom+txt.FadeIn, 'f', 2, 64) +
	//			`),(t-` + strconv.FormatFloat(txt.TimeFrom, 'f', 2, 64) +
	//			`)/` + strconv.FormatFloat(txt.FadeIn, 'f', 2, 64) +
	//			`,if(lte(t,` + strconv.FormatFloat(txt.TimeTo, 'f', 2, 64) +
	//			`),1,if(gte(t,` + strconv.FormatFloat(txt.TimeTo-2, 'f', 2, 64) +
	//			`),1-(t-` + strconv.FormatFloat(txt.TimeTo-2, 'f', 2, 64) +
	//			`)/` + strconv.FormatFloat(txt.FadeOut, 'f', 2, 64) + `,0))))'`
	//	}

	command += `"`

	// Add comma to end of command if not last
	if !isLast {
		command += `,`
	}

	f.Command += command

}
