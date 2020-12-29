package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
)

type options struct {
	// Command string `short:"c" long:"command" description:"The name of the input Markdown file" value-name:"INPUT"`
	Input  string `short:"i" long:"input" description:"The name of the input file" value-name:"INPUT"`
	Output string `short:"o" long:"output" description:"The name of the output file" value-name:"OUTPUT"`
}

func main() {
	opts := &options{}
	flags.Parse(opts)

	input, err := getInput(opts)
	if err != nil {
		log.Fatalf("unable to open input file: %v", err)
	}
	defer input.Close()

	output, err := getOutput(opts)
	if err != nil {
		log.Fatalf("unable to open output file: %v", err)
	}
	defer output.Close()

	E := color.New(color.FgHiRed).PrintfFunc()
	W := color.New(color.FgHiYellow).PrintfFunc()
	I := color.New(color.FgHiGreen).PrintfFunc()
	D := color.New(color.FgHiWhite).PrintfFunc()
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {

		// fmt.Fprintf(output, "%s\n", scanner.Text())
		m := map[string]interface{}{}
		err := json.Unmarshal(scanner.Bytes(), &m)
		if err != nil {
			panic(err)
		}
		switch m["level"] {
		case "error", "fatal":
			dumpMap(E, "", m)
		case "warn", "warning":
			dumpMap(W, "", m)
		case "info":
			dumpMap(I, "", m)
		case "debug":
			dumpMap(D, "", m)
		}

	}
}

func dumpMap(print func(format string, a ...interface{}), space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			print("{ \"%v\": \n", k)
			dumpMap(print, space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			print("%v %v : %v\n", space, k, v)
		}
	}
}

// getInput returns the input Reader to use; if a filename argument is provided,
// open the file to read from it, otherwise return STDIN; the Reader must be
// closed by the method's caller.
func getInput(opts *options) (*os.File, error) {
	if opts.Input != "" {
		return os.Open(opts.Input)
	}
	return os.Stdin, nil
}

// getOutput returns the output Writer to use; if a filename argument is provided,
// open the file to write to it, otherwise return STDOUT; the Writer must be
// closed by the method's caller.
func getOutput(opts *options) (*os.File, error) {
	if opts.Output != "" {
		return os.Create(opts.Output)
	}
	return os.Stdout, nil
}
