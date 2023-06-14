package ffmpeg

type FFMPEG interface {
	AddAudio()
	AddText()
	Configure()
	CombineVideoAudio()
	StitchVideos()
	StitchAudio()

	MakeCommand() string
}

type FFMPEGCommand struct {
	Flags           []string
	Command         string
	InputFile       bool
	Input           string
	Out             string
	FileType        string
	ShouldCopy      bool
	HasComplexAudio bool
	VideoCodec      string
	AudioCodec      string
}

type FFMPEGText struct {
	TextFile    bool
	Data        string
	FontFile    string
	FontSize    int
	FontColor   string
	LineHeight  int
	X           int
	Y           int
	HasDuration bool
	Delay       float64
	TimeFrom    float64
	TimeTo      float64
	FadeIn      float64
	FadeOut     float64
}

type FFMPEGAudio struct {
	Input    string
	FileType string
	Duration float64
	Delay    float64
}

type FFMPEGVideo struct {
	FileName string
	Duration float64
}
