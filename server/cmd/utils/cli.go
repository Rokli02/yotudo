package utils

import "fmt"

func WaitForExit() chan struct{} {
	exitSignal := make(chan struct{})

	go func() {
		var line string

		for {
			fmt.Scanln(&line)
			if line == "exit" {
				close(exitSignal)
				return
			}
		}
	}()

	return exitSignal
}
