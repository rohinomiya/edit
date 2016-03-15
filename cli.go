package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

// file match pattern and Command
type UserCommand struct {
	FilePattern string
	Command     string
	Option      string
	LineOption  string
}

// Sets of User command
type UserCommandSlice struct {
	UserCommands []UserCommand
}

var commandSet UserCommandSlice

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

func exchangeFileExt(filePath string, newExt string) string {
	dir, filename := path.Split(filePath)
	ext := path.Ext(filePath)

	newFilename := filename[0:len(filename)-len(ext)] + newExt
	return path.Join(dir, newFilename)
}

func editFile(line string) {
	filePath := ""
	lineNum := ""

	//GREP形式の場合には、行番号を指定して開く
	//GREP形式は2種類あるので、どちらにも対応するように
	//TODO: もっとすっきり書いて
	rGrepFormat1 := regexp.MustCompile(`^(.+?):([0-9]+):`)
	rGrepFormat2 := regexp.MustCompile(`^(.+?)\s*\(([0-9]+)\)`)
	if rGrepFormat1.MatchString(line) {
		result := rGrepFormat1.FindStringSubmatch(line)
		filePath = result[1]
		lineNum = result[2]
	} else if rGrepFormat2.MatchString(line) {
		result := rGrepFormat2.FindStringSubmatch(line)
		filePath = result[1]
		lineNum = result[2]
	} else {
		filePath = line
		lineNum = ""
	}

	//fmt.Println(filePath, lineNum)

	//ext := path.Ext(filePath)
	_, filename := path.Split(filePath)

	appPath := ``
	option := ""
	lineOption := ""

	for _, c := range commandSet.UserCommands {
		re := regexp.MustCompile(c.FilePattern)

		//todo:remove
		//fmt.Println("file:" + filename + "]")
		//fmt.Println("regex:" + c.FilePattern + "]")
		if re.MatchString(filename) {
			appPath = c.Command
			option = c.Option
			lineOption = strings.Replace(c.LineOption, "[num]", lineNum, 1)
			break
		}
	}

	//実行する（終わりを待たない）
	err := exec.Command(appPath, option, lineOption, filePath).Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		line    int
		noStdin bool
		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.IntVar(&line, "line", 0, "line number")
	flags.IntVar(&line, "l", 0, "line number(Short)")

	flags.BoolVar(&noStdin, "noStdin", false, "do not input from stdin.")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	_ = line

	//todo:ファイル名部分もEXE名から生成できるようにしてpath.Baseだと切り出せないので工夫が必要
	exePath := strings.Replace(ExecuteFilePath(), `\`, "/", -1)
	commandFile := exchangeFileExt(exePath, ".json")

	f, err := os.Open(commandFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "File:"+commandFile)
		os.Exit(1)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&commandSet)
	if err != nil {
		fmt.Fprintln(os.Stderr, "json decode error: "+commandFile)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	//fmt.Printf("%+v\n", commandSet)

	//TODO: もうちょっと綺麗にかけない？
	files := []string{}

	//引数が空の場合には、標準入力からパスを取得して実行する
	if len(flags.Args()) != 0 {
		files = flags.Args()
	} else if !noStdin {

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			files = append(files, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}

	for _, arg := range files {
		editFile(arg)
	}

	return ExitCodeOK
}
