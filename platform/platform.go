package platform

import (
	"os"
	"runtime"
)

type OS string

const (
	MacOS OS = "macos"
	Linux OS = "linux"
)

var OSs []OS = []OS{MacOS, Linux}

func Target() OS {
	if env := os.Getenv("BLOS"); env != "" {
		switch OS(env) {
		case MacOS, Linux:
			return OS(env)
		default:
			panic("unsupported target OS: " + env)
		}
	}
	return Host()
}

func Host() OS {
	switch runtime.GOOS {
	case "darwin":
		return MacOS
	case "linux":
		return Linux
	default:
		panic("unsupported target OS: " + runtime.GOOS)
	}
}
