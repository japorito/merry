package sleigh

import (
	"fmt"
	"os"
)

type InputResolver interface {
	ResolveInput(args ...string) (*os.File, error)
}

type DefaultInputResolver struct{}

func (resolver DefaultInputResolver) ResolveInput(args ...string) (*os.File, error) {
	if len(args) == 1 {
		inputFile, err := os.Open(args[0])
		if err == nil {
			return inputFile, nil
		}

		if wd, err := os.Getwd(); err == nil {
			inputFile, err := os.Open(wd + string(os.PathSeparator) + args[0])
			if err == nil {
				return inputFile, nil
			}
		}
	}

	if stat, err := os.Stdin.Stat(); err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		return os.Stdin, nil
	}

	return nil, fmt.Errorf("no input found")
}

var InputSource InputResolver = DefaultInputResolver{}
