package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {

	return drawTree(out, "", path, printFiles)
}

func drawTree(out io.Writer, drawPath string, path string, printFiles bool) error {

	dir, err := os.Open(path)

	if err != nil {
		return fmt.Errorf("Error")
	}

	defer dir.Close()

	outputDirFiles, err := dir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("Error")
	}

	if !printFiles {
		var newOutputDirFiles []os.FileInfo
		for i := range outputDirFiles {
			outputFileHere := outputDirFiles[i]
			if outputFileHere.IsDir() {
				newOutputDirFiles = append(newOutputDirFiles, outputDirFiles[i])
			}
		}
		outputDirFiles = newOutputDirFiles
	}

	for i := range outputDirFiles {
		outputFileHere := outputDirFiles[i]
		fileName := outputFileHere.Name()
		lastChild := i+1 == len(outputDirFiles)

		if outputFileHere.IsDir() {

			drawElement(out, drawPath, outputFileHere, lastChild)
			var newDrawPath string
			if lastChild {
				newDrawPath = drawPath + "	"
			} else {
				newDrawPath = drawPath + "│	"
			}
			drawTree(out, newDrawPath, path+string(os.PathSeparator)+fileName, printFiles)
		} else if printFiles {
			drawElement(out, drawPath, outputFileHere, lastChild)
		}

	}

	return nil

}

func drawElement(out io.Writer, drawPath string, outputFileHere os.FileInfo, lastChild bool) {

	out.Write([]byte(drawPath))
	if lastChild {
		out.Write([]byte("└───"))
	} else {
		out.Write([]byte("├───"))
	}
	out.Write([]byte(outputFileHere.Name()))
	if !outputFileHere.IsDir() {
		out.Write([]byte(" (" + getSizeFormat(outputFileHere.Size()) + ")"))
	}
	out.Write([]byte("\n"))
}

func getSizeFormat(size int64) string {
	if size > 0 {
		return strconv.Itoa(int(size)) + "b"
	}
	return "empty"
}
