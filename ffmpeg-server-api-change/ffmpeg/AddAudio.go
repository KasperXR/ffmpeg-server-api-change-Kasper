package ffmpeg

import (
	"fmt"
	"main/fileutil"
	"os/exec"
	"strconv"
	"strings"
)

func (f *FFMPEGCommand) AddAudio(audio *FFMPEGAudio, isLast bool) {
}

func (f *FFMPEGCommand) StitchAudio(fileListPath string, outpath string, ext string) {
	// read lines of fileListPath usinc wc -l
	// for each line, add to string [x:a] where x is the index of the line
	var cmd = exec.Command("wc", "-l", fileListPath)
	var out, err = cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	var inputs, readLinesErr = fileutil.ReadLines(fileListPath)

	if readLinesErr != nil {
		fmt.Println(readLinesErr)
	}

	// take the array and add -i to each line besides the last one
	var inputsString = ""
	for i, line := range inputs {
		if i != len(inputs)-1 {
			inputsString += "-i "
		}
		inputsString += line
		inputsString += " "
	}

	// get the number from out, discard everything after the number
	var fileNum, err2 = strconv.Atoi(strings.Split(string(out), " ")[0])
	if err2 != nil {
		fmt.Println(err2)
	}

	var filter = "-filter_complex \""
	// Add add the [x:a] to the filter with adelay=0.250[ax]
	for i := 0; i < fileNum; i++ {
		filter += "[" + strconv.Itoa(i) + ":a]"
		filter += "adelay=250[a" + strconv.Itoa(i) + "]"

		filter += ";"

	}

	// Add the [ax] to the filter with amix=inputs=x
	for i := 0; i < fileNum; i++ {
		filter += "[a" + strconv.Itoa(i) + "]"
	}

	filter += "concat=n=" + strconv.Itoa(fileNum) + ":v=0:a=1[outa]\" -map \"[outa]\"" +
		" -c:a libfdk_aac -b:a 320k -ar 48k -movflags +faststart -preset slow " +
		outpath + "." + ext

	var command = "ffmpeg " + inputsString + " " + filter

	f.Command = command

}

func (f *FFMPEGCommand) CombineVideoAudio(i1, i2, o string) {

	var command = "ffmpeg -i " + i1 + " -i " + i2 + " -c copy -b:a 320k -ar 48k -preset:v superfast -preset:a superfast -shortest -movflags +faststart " + o

	f.Command = command
}
