package ffmpeg

type Choices struct {
	Id        int
	Text      string
	AudioName string
	Duration  int
}

func (f *FFMPEGCommand) Configure() {
	command := "ffmpeg"
	command += ` `

	// Set input differently if input is a list of files
	if f.InputFile {
		command += "-f concat -i " + f.Input + ".txt" + ` `
	} else {
		command += "-i " + f.Input + " "
	}

	if f.ShouldCopy {
		command += "-c copy "
	} // Implement Audio and Video codecs

	f.Command = command

}

func (f *FFMPEGCommand) MakeCommand(Preset string, APreset string, final bool) string {
	command := f.Command

	if f.HasComplexAudio {
		command += ` -map 0:v -map "[a]"`
	}

	if final {
		command += ` -c:a copy -movflags +faststart -tune fastdecode -crf 31 -pix_fmt yuv420p -level 4.2 `
	}

	if Preset != "" {
		command += ` -preset:v ` + Preset + ` `
	}

	if APreset != "" {
		command += ` -preset:a ` + APreset + ` `
	}

	command += f.Out + "." + f.FileType

	return command
}
