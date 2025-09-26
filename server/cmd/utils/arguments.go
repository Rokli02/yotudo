package utils

import (
	"fmt"
	"strconv"
)

type State uint64

const (
	Help State = 1 << iota
	Host
	HostFrontend
)

type Arguments struct {
	Port  int
	state State
}

func (a *Arguments) setState(state State, value bool) {
	if value {
		a.state |= state
	} else {
		a.state &= ^state
	}
}

func (a *Arguments) GetState(state State) bool {
	return a.state&state != 0
}

func ParseArguments(args []string) *Arguments {
	arguments := Arguments{}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--port":
			fallthrough
		case "-p":
			if i++; i >= len(args) {
				break
			}

			if _port, err := strconv.Atoi(args[i]); err != nil {
				fmt.Println("WARN: Passed port is invalid")
			} else if _port < 1000 || _port > 65535 {
				fmt.Printf("WARN: Passed port is out of bound (0 < %d < 65535)\n", _port)
			} else {
				arguments.Port = _port
			}
		case "--host":
			arguments.setState(Host, true)
		case "--localhost":
			arguments.setState(Host, false)
		case "--host-frontend":
			arguments.setState(HostFrontend, true)
		case "--help":
			fallthrough
		case "-h":
			arguments.setState(Help, true)
		default:
			fmt.Printf("WARN: Unknown argument was passed: '%s'\n", args[i])
		}
	}

	return &arguments
}
