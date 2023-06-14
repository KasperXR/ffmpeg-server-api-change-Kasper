package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/uuid"

	ffmpeg "nrt/ffmpeg"
	omniglyph "nrt/omniglyph"
)

// FFMPEG Video stiching command
var FFMPEG_COMMAND = "ffmpeg -f concat -safe 0 -i vFile.txt -c copy videos/output/"

// FFMPEG Top-right text logo
var LOGO = `drawtext="fontfile=./fonts/TitilliumWeb-SemiBold.ttf:text='MIT HJERTE':x=w-tw-15:y=15:fontsize=48:fontcolor=#B40031"`

// FFMPEG Intro text in the center of the screen for the first 5 seconds
var INTRO_TEXT = `drawtext="fontfile=./fonts/TitilliumWeb-SemiBold.ttf:text='MIT HJERTE':x=(w-text_w)/2:y=(h-text_h)/2:fontsize=96:fontcolor=#B40031:enable='between(t,0,5)'"`

// FFMPEG Outro disclaimer in the center of the screen for the last 5 seconds
var OUTRO_DISCLAIMER_TEXT = `drawtext="fontfile=./fonts/TitilliumWeb-SemiBold.ttf:text='De medicinske/sundhedsmæssige oplysninger gives kun til generelle informations- og uddannelsesformål og er ikke en erstatning for professionel rådgivning. Derfor opfordrer vi dig til at rådføre dig med de relevante fagfolk, før du tager nogen handlinger baseret på sådanne oplysninger. Vi yder ingen form for medicinsk eller sundhedsmæssig rådgivning. Brugen af eller tilliden til enhver information i denne video er på eget ansvar.':x=(w-text_w)/2:y=(h-text_h)/2:fontsize=32:fontcolor=#B40031`

// Generate outro text between, it takes the full dureatoin
func GenerateOutroFromDur(duration float64) string {
	return `:enable='between(t,` + strconv.FormatFloat(duration-5, 'f', 2, 64) + `,` + strconv.FormatFloat(duration, 'f', 2, 64) + `)'` + `"`
}

// Textfile generation function for each Parent option in VideoStruct, and each its sub Options.
// It creates two files, one for the title and one for the text.
// Takes choice.name as title, and subOption.name as text.
// It returns the filename
func GenerateTextFile(parentOpt ParentOption) string {
	// Create a new UUID
	UUID := uuid.New()

	// Create a file for title and text
	titleFile, _ := os.Create("text/" + UUID.String() + "-title.txt")
	textFile, _ := os.Create("text/" + UUID.String() + "-text.txt")

	var titleWrapper = omniglyph.WordWrapper{
		Joiner:    " ",
		NewLine:   "\n",
		Separator: " ",
		Text:      parentOpt.Name,
		Width:     120,
	}
	var textWrapper = omniglyph.WordWrapper{
		Joiner:    " ",
		NewLine:   "\n",
		Separator: " ",
		Width:     60,

		IndentAmount: 5,
		IndentAll:    true,
		IndentStart:  false,
		IndentGlyph:  " ",

		Prefix:      "  • ",
		PrefixStart: true,
	}

	var replace = []string{"%"}
	var replaceWith = []string{"\\%"}

	titleWrapper.ReplaceGlyphs(replace, replaceWith)
	titleFile.WriteString(titleWrapper.Text)

	for _, option := range parentOpt.Options {
		if option.Active {
			// fmt.Println("Checked:", option.Name) // Debugging
			textWrapper.Text = option.Name
			textWrapper.ReplaceGlyphs(replace, replaceWith)
			textWrapper.Wrap()
			textFile.WriteString(textWrapper.Text + textWrapper.NewLine)
		}
	}

	titleFile.Close()
	textFile.Close()

	return UUID.String()
}

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

// Main video generation function
func GenerateVideo(fileName string, videoChoiceArr []VideoObj) {
	var optArrText []SanitizedOption
	var optArrAudio []ffmpeg.FFMPEGAudio

	// Cleanup
	removeVFileErr := os.Remove("vFile.txt")
	removeAFileErr := os.Remove("aFile.txt")

	if removeVFileErr != nil {
		fmt.Println("vFile.txt does not exist, proceeding to create it")
	}
	if removeAFileErr != nil {
		fmt.Println("aFile.txt does not exist, proceeding to create it")
	}

	vFile, createVFileErr := os.Create("vFile.txt")
	aFile, createAFileErr := os.Create("aFile.txt")

	if createVFileErr != nil {
		panic("Error creating vFile.txt")
	}
	if createAFileErr != nil {
		panic("Error creating aFile.txt")
	}

	var videoName string
	var totalDuration float64 = 10 // 10 seconds extra for intro and outro

	// Loop through each video choice and make data array for both text and audio
	for _, v := range videoChoiceArr {
		fmt.Println("Video:", v.Id) // Debugging
		videoName += v.Id + "_Long" // Set video name

		for idx, parentOpt := range v.ParentOptions {
			var parentOptDur float64 = 0

			// If first parent option, consider intro text duration

			var audioFileName = "audio/" + parentOpt.AudioName + ".aac"

			// Use an introduction if defined
			if parentOpt.Introduction != "" {

				var audioDur, err = getDurationInSeconds("audio/" + parentOpt.Introduction + ".aac")
				if err != nil {
					fmt.Println("Error getting duration of audio file:", "audio/"+parentOpt.Introduction+".aac")
				}

				parentOptDur += audioDur

				fmt.Println("Parent option duration:", parentOptDur) // Debugging

				var introTextWrapper = omniglyph.WordWrapper{
					Joiner:    " ",
					NewLine:   "\n",
					Separator: " ",
					Text:      parentOpt.Description,
					Width:     60,

					IndentAmount: 5,
					IndentAll:    true,
					IndentStart:  false,
					IndentGlyph:  " ",
				}

				var replace = []string{"%"}
				var replaceWith = []string{"\\%"}

				introTextWrapper.ReplaceGlyphs(replace, replaceWith)
				introTextWrapper.Wrap()

				// Create a file for the intro text
				introTextFile, _ := os.Create("text/" + parentOpt.AudioName + "-intro.txt")
				introTextFile.WriteString(introTextWrapper.Text)
				introTextFile.Close()

				aFile.WriteString("file 'audio/" + parentOpt.Introduction + ".aac'" + "\n")

				optArrAudio = append(optArrAudio, ffmpeg.FFMPEGAudio{
					Duration: audioDur,
					Delay:    0.250,
					FileType: "aac",
					Input:    "audio/" + parentOpt.Introduction,
				})
			} else if parentOpt.AudioName != "" {
				// get duration of audio file
				var audioDur, err = getDurationInSeconds(audioFileName)
				if err != nil {
					fmt.Println("Error getting duration of audio file:", audioFileName)
				}

				parentOptDur += audioDur + 0.250 // Add 250ms delay

				fmt.Println("Parent option duration:", parentOptDur) // Debugging

				aFile.WriteString("file 'audio/" + parentOpt.AudioName + ".aac'" + "\n")

				optArrAudio = append(optArrAudio, ffmpeg.FFMPEGAudio{
					Duration: audioDur,
					Delay:    0.250,
					FileType: "aac",
					Input:    "audio/" + parentOpt.AudioName,
				})
			}

			// Add each sub option's audio if it's defined
			for _, option := range parentOpt.Options {
				if option.Active {

					var audioDur, err = getDurationInSeconds("audio/" + option.AudioName + ".aac")

					if err != nil {
						fmt.Println("Error getting duration of audio file:", "audio/"+option.AudioName+".aac")
					}

					//	totalDurationMs += audioDur * 1000
					parentOptDur += audioDur + option.Delay

					fmt.Println("Option duration:", audioDur) // Debugging

					aFile.WriteString("file 'audio/" + option.AudioName + ".aac'" + "\n")

					optArrAudio = append(optArrAudio, ffmpeg.FFMPEGAudio{
						Duration: audioDur,
						Delay:    option.Delay,
						FileType: "aac",
						Input:    "audio/" + option.AudioName,
					})
				}
			}

			// Audioname depends on whether or not there is an introduction
			var audioName string
			if parentOpt.Introduction != "" {
				audioName = parentOpt.Introduction
			} else {
				audioName = parentOpt.AudioName
			}

			optArrText = append(optArrText, SanitizedOption{
				Id:        idx,
				Text:      GenerateTextFile(parentOpt), // Generate a textfile for each option's title & text
				AudioName: audioName,
				Duration:  parentOptDur, // In seconds
				Delay:     0.250,
			})

			totalDuration += parentOptDur
			fmt.Println("Total duration:", totalDuration) // Debugging
		}

	}

	vFile.Close()
	aFile.Close()

	var workingDir, err = os.Getwd()
	var audioFileName = StitchAudio(workingDir + "/" + aFile.Name())

	if err != nil {
		fmt.Println(err)
	}

	// Get duration of video
	cmd := exec.Command("ffprobe", "-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		"videos/"+videoName+".mp4")
	output, durErr := cmd.Output()
	if durErr != nil {
		fmt.Println(durErr)
	}

	// Convert the output to a float
	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Video duration:", duration)
	// Cut the video from the duration to the end
	stitchOut, stitchErr := exec.Command(
		"ffmpeg", "-i", "videos/"+videoName+".mp4",
		"-ss", strconv.FormatFloat(duration+10, 'f', -1, 64),
		"-c", "copy",
		"videos/output/"+fileName+".mp4").CombinedOutput()

	fmt.Println("Stitching video...")

	fmt.Printf("\n\n")

	fmt.Printf("\n\n")
	fmt.Println("", string(stitchOut))

	if stitchErr != nil {
		fmt.Println("Error stitching videos:", stitchErr)
	}

	// Make A/V combine function
	var combiner = ffmpeg.FFMPEGCommand{}

	combiner.CombineVideoAudio(
		"videos/output/"+fileName+".mp4",
		audioFileName,
		"videos/output/"+fileName+"-av.mp4",
	)

	combinedAVCmd := combiner.Command
	combinedAVOut, combinedAVErr := exec.Command("sh", "-c", combinedAVCmd).CombinedOutput()

	fmt.Printf("COMBINED OUT\n\n")
	fmt.Println("", string(combinedAVOut))
	fmt.Printf("\n")

	if combinedAVErr != nil {
		fmt.Println("Error combining audio and video:", combinedAVErr)
	}

	// Make the final video using the base video
	var finalVideoCmd = ffmpeg.FFMPEGCommand{
		Input:      "videos/output/" + fileName + "-av.mp4",
		Out:        "videos/output/" + fileName + "-final",
		FileType:   "mp4",
		ShouldCopy: false,
	}

	finalVideoCmd.Configure()

	// Add text to the video
	finalVideoCmd.Command += " -vf "
	// Add logo
	finalVideoCmd.Command += LOGO + "," + INTRO_TEXT + "," + OUTRO_DISCLAIMER_TEXT + GenerateOutroFromDur(totalDuration) + ","
	addText(&finalVideoCmd, optArrText)

	//	fmt.Println("Making final video...")
	//	fmt.Printf("\n\n")

	command := finalVideoCmd.MakeCommand("ultrafast", "ultrafast", true)

	fmt.Println(command)
	fmt.Printf("\n")

	finalOut, finalErr := exec.Command("sh", "-c", command).CombinedOutput()

	fmt.Printf("FINAL OUT:\n\n")
	fmt.Println("", string(finalOut))
	fmt.Printf("\n\n")

	if finalErr != nil {
		fmt.Println("Error running command:", finalErr)
	}

	// Call the unsued funtions
	UNUSED(finalOut)
	UNUSED(stitchOut)
	UNUSED(combinedAVOut)
}

// Stitch audio files together
func StitchAudio(fileList string) string {
	// Remove mediator file if it exists
	removeMediatorErr := os.Remove("audio/output/audioMediator.aac")

	if removeMediatorErr != nil {
		fmt.Println("audioMediator.aac does not exist, proceeding to create it")
	}

	var audioCmd = ffmpeg.FFMPEGCommand{
		Input:    fileList,
		Out:      "audio/output",
		FileType: "aac",
	}

	audioCmd.StitchAudio(fileList, "audio/output/audioMediator", "aac")

	// debug output
	fmt.Println("Audio command:", audioCmd.Command)

	fmt.Printf("\n")
	audioOut, audioErr := exec.Command("sh", "-c", audioCmd.Command).CombinedOutput()

	fmt.Println("\n\nAUDIO OUT ", string(audioOut))
	fmt.Printf("\n\n")

	if audioErr != nil {
		fmt.Println("Error running command:", audioErr)
	}
	UNUSED(audioOut)
	return "audio/output/audioMediator.aac"
}

func getDurationInSeconds(filename string) (float64, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer

	// Run probe command and capture the duration metadata
	command := exec.Command("sh", "-c", "ffprobe -show_entries format=duration -v error -of csv='p=0' -i "+filename)

	command.Stdout = &out
	command.Stderr = &stderr

	err := command.Run()

	if err != nil {
		println("Error getting duration of audio: ", err.Error())
		return 0, err
	}

	// Convert the output to a string
	duration := out.String()

	// Remove the newline character
	duration = strings.TrimSuffix(duration, "\n")

	// Round the duration to two decimal places
	durationFloat, err := strconv.ParseFloat(duration, 64)

	roundedDuration := fmt.Sprintf("%.2f", durationFloat)

	roundedDurationFloat, err := strconv.ParseFloat(roundedDuration, 64)

	// Debug
	// println("DURATION: ", durationFloat)
	// println("ROUNDED DUR: ", roundedDuration)

	if err != nil {
		println("Error converting duration to float: ", err.Error())
		return 0, err
	}

	return roundedDurationFloat, nil
}

// Function for generating the text with ffmpeg
// It takes an array of textfiles to be used
// It makes a command that concatenates all the textfiles onto a single video
func addText(f *ffmpeg.FFMPEGCommand, options []SanitizedOption) {
	var xPosTitle = 52
	var xPosText = 52

	var yPosTitle = 64
	var yPosText = 124

	var delay = float64(0.25)

	var durSoFar float64 = 0
	var durFrom = durSoFar
	var durTo = float64(options[0].Duration)

	for idx, opt := range options {
		// add 5 seconds to initial text
		if idx == 0 {
			durSoFar += 5
			durFrom += 5
			durTo += 5
		}
		// Increment duration variables
		// Also add delay to from and to variables
		if idx != 0 {
			durFrom += durSoFar
			durTo += float64(opt.Duration) + delay
		}

		var titleText = ffmpeg.FFMPEGText{
			TextFile:    true,
			Data:        "text/" + opt.Text + "-title.txt",
			HasDuration: true,
			TimeFrom:    durSoFar,
			TimeTo:      durTo,
			FadeIn:      1.3,
			FadeOut:     2,
			FontFile:    "fonts/TitilliumWeb-SemiBold.ttf",
			LineHeight:  2,
			FontSize:    52,
			FontColor:   "black",
			X:           xPosTitle,
			Y:           yPosTitle,
		}

		var textText = ffmpeg.FFMPEGText{
			TextFile:    true,
			Data:        "text/" + opt.Text + "-text.txt",
			HasDuration: true,
			TimeFrom:    durSoFar,
			TimeTo:      durTo,
			FadeIn:      2,
			FadeOut:     2,
			FontFile:    "fonts/TitilliumWeb-SemiBold.ttf",
			LineHeight:  2,
			FontSize:    40,
			FontColor:   "black",
			X:           xPosText,
			Y:           yPosText,
		}

		f.AddText(&titleText, false)

		// Last title is not last, the text that comes after is
		if idx == len(options)-1 {
			f.AddText(&textText, true)
		} else {
			f.AddText(&textText, false)
		}

		durSoFar += float64(opt.Duration) + delay
	}
}

func main() {
	StartServer(false)

	// A story written by Github copilot directed by Mathias Wøbbe
	// It starts with a guy in a hat
	// And he's a funny guy
	// He has a funny hat
	// Sometimes he likes to play with his funny hat
	// But that doesn't mean he's a funny guy
	// He's just a funny guy
	// One day he's going to go to the store
	// He's going to buy a new hat
	// But he realized he doesn't have enough money
	// He's going to go to the bank to get more money
	// But when he get to the bank he realizes he doesn't have a funny hat
	// So he tried to be funny and it didn't work out for him :(
	// Security kicked him out of the bank and he's homeless now :(
	// But two months later he has a great idea to change his life for a better one :)
	// He's going to buy a new house and he's going to buy a new funny hat
	// He's going to buy a new car and he's going to buy a new funny hat
	// His idea was amazing, it was the best idea he ever had. It is about to happen!
	// The idea revovled around the world and it's spreading fast
	// It's spreading fast and it's spreading fast and it's spreading fast, hold on to your funny hats!
	// Everybody liked the idea so much that they bought the funny hats
	// His idea was to make a new funny hat that can be used by everyone
	// The moral of the story is that everyone should have a funny hat, even if they don't have money. :)

}
