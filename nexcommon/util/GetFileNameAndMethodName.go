package util

import (
	"runtime"
	"strings"
)

func getFrame(skipFrames int) runtime.Frame {
	targetFrameIndex := skipFrames + 2

	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

func GetFileAndMethod() (fileName string, functionName string, line int) {
	frame := getFrame(1)

	functionPath := strings.Split(frame.Function, "/")
	functionName = functionPath[len(functionPath)-1]

	filepath := strings.Split(frame.File, "/")
	fileName = filepath[len(filepath)-1]

	line = frame.Line

	return
}
