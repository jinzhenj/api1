package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jinzhenj/api1/pkg/all"
	"github.com/jinzhenj/api1/pkg/utils"
)

func main() {
	files, err := utils.ListFiles(".", isApiFile)
	if err != nil {
		fatal(err)
	}
	if len(files) == 0 {
		info("No API files found")
		return
	}

	info("Found API files:")
	for _, file := range files {
		info("    %s", file)
	}

	render := all.NewRender()
	codeFiles, err := render.RenderFiles(files)
	if err != nil {
		fatal(err)
	}
	if len(codeFiles) == 0 {
		info("No file generated")
		return
	}

	info("Will write files:")
	for _, codeFile := range codeFiles {
		exists, err := utils.FileExists(codeFile.Name)
		if err != nil {
			fatal(err)
		}
		state := "create"
		if exists {
			state = "update"
		}
		info("    (%s) %s", state, codeFile.Name)
	}

	fmt.Fprintf(os.Stderr, "Please confirm [y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	if strings.ToLower(input) != "y" {
		info("Give up")
		return
	}

	for _, codeFile := range codeFiles {
		info("Writing %s ...", codeFile.Name)
		if err := codeFile.WriteFile(); err != nil {
			fatal(err)
		}
	}

	info("Done")
}

func isApiFile(file string) bool {
	return strings.HasSuffix(file, ".api")
}

func info(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func fatal(err error, s ...string) {
	message := "Error"
	if len(s) > 0 {
		var args []interface{}
		for _, arg := range s[1:] {
			args = append(args, arg)
		}
		message = fmt.Sprintf(s[0], args...)
	}
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}
