package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"text_editor/teletype"
)

func main() {
	tty, err := teletype.Open()
	if err != nil {
		fmt.Printf("Error opening TTY: %v\n", err)
		return
	}
	defer func() {
		if err := tty.Close(); err != nil {
			fmt.Printf("Error closing TTY: %v\n", err)
		}
	}()

	signal.Notify(tty.Signal(), syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("TTY test program. Press Ctrl+C to exit.")

	go func() {
		<-tty.Signal()
		fmt.Println("Exiting...")
		os.Exit(0)
	}()

	for {
		r, err := tty.Read()
		if err != nil {
			fmt.Printf("Error reading from TTY: %v\n", err)
			continue
		}
		_, err = tty.Output().WriteString("You pressed: " + string(r) + "\n")
		if err != nil {
			return
		}
	}
}
