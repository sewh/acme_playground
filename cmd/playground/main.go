package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"9fans.net/go/acme"
)

func main() {
	// Make sure we are run from acme
	winIDStr := os.Getenv("winid")
	if winIDStr == "" {
		fmt.Println("Must be run from an acme window!")
		os.Exit(0)
	}

	winID, err := strconv.Atoi(winIDStr)
	if err != nil {
		fmt.Printf("Invalid winid: %s\n", winIDStr)
		os.Exit(1)
	}

	// Make sure we have a program to execute
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("usage: playground [executable to run]")
		os.Exit(0)
	}

	// Open the acme windows
	win, err := acme.Open(winID, nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	pwd, _ := os.Getwd()
	outputWin, err := acme.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	outputWin.Name(pwd + "/+playground")

	event := make(chan bool)

	go acmeEventHandler(win, event)
	go acmeEventHandler(outputWin, event)

	for {
		func() {
			// Clean the buffer
			outputWin.Addr(",")
			outputWin.Write("data", nil)

			// Set up command
			inPipeR, inPipeW, err := os.Pipe()
			if err != nil {
				fmt.Println("pipe error: " + err.Error())
				os.Exit(1)
			}
			defer inPipeR.Close()
			defer inPipeW.Close()

			var cmd *exec.Cmd
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if len(args) == 1 {
				cmd = exec.CommandContext(ctx, args[0])
			} else {
				cmd = exec.CommandContext(ctx, args[0], args[1:]...)
			}

			if err != nil {
				fmt.Println("exec error: " + err.Error())
				os.Exit(1)
			}

			cmd.Stdin = inPipeR

			buffContents, err := win.ReadAll("body")
			if err != nil {
				fmt.Println("read window error: " + err.Error())
				os.Exit(1)
			}

			go func() {
				inPipeW.Write(buffContents)
				inPipeW.Close()
			}()

			res, err := cmd.CombinedOutput()
			if err != nil {
				outputWin.Write("data", []byte(err.Error()))
				outputWin.Write("data", []byte("\n"))
			}

			outputWin.Write("data", res)
			outputWin.Write("ctl", []byte("clean"))

			<-event
		}()
	}

}

func acmeEventHandler(win *acme.Win, event chan bool) {
	for e := range win.EventChan() {
		if e.C1 == 'K' && (e.C2 == 'D' || e.C2 == 'I') {
			event <- true
		}
		win.WriteEvent(e)
	}

	os.Exit(0)
}
