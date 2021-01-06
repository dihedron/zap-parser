package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/mattn/go-isatty"
	"gopkg.in/yaml.v2"
)

type options struct {
	Input  string `short:"i" long:"input" description:"The name of the input file" value-name:"INPUT"`
	Output string `short:"o" long:"output" description:"The name of the output file" value-name:"OUTPUT"`
}

type print func(format string, a ...interface{})

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

	var E, W, I, D print
	if isatty.IsTerminal(output.Fd()) {
		E = color.New(color.FgHiRed).PrintfFunc()
		W = color.New(color.FgHiYellow).PrintfFunc()
		I = color.New(color.FgHiGreen).PrintfFunc()
		D = color.New(color.FgWhite).PrintfFunc()
	} else {
		write := func(format string, a ...interface{}) {
			fmt.Fprintf(output, format, a...)
		}
		E = write
		W = write
		I = write
		D = write
	}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		m := map[string]interface{}{}
		err := json.Unmarshal(scanner.Bytes(), &m)
		if err != nil {
			panic(err)
		}
		switch m["level"] {
		case "error", "fatal":
			// dumpMap(E, "", m)
			write(E, m)
		case "warn", "warning":
			// dumpMap(W, "", m)
			write(W, m)
		case "info":
			// dumpMap(I, "", m)
			write(I, m)
		case "debug":
			// dumpMap(D, "", m)
			write(D, m)
		}

	}
}

// Message is a structured log message.
type Message struct {
	Application string                 `yaml:"application,omitempty"`
	Level       string                 `yaml:"level,omitempty"`
	Message     string                 `yaml:"message,omitempty"`
	Data        map[string]interface{} `yaml:"data,omitempty"`
}

func write(print func(format string, a ...interface{}), m map[string]interface{}) {
	msg := Message{
		Data: map[string]interface{}{},
	}
	for k, v := range m {
		switch k {
		case "application":
			if v, ok := v.(string); ok {
				msg.Application = v
			}
		case "level":
			if v, ok := v.(string); ok {
				msg.Level = v
			}
		case "message":
			if v, ok := v.(string); ok {
				msg.Message = v
			}
		default:
			msg.Data[k] = v
		}
	}
	data, err := yaml.Marshal(msg)
	if err == nil {
		print(string(data))
		print("\n")
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
