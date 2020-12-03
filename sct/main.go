package main

import (
	"os"
	"strconv"

	"github.com/d4l3k/go-sct"
)

func main() {
	args := os.Args
	temp := 6500
	if len(args) == 2 {
		if ntemp, err := strconv.Atoi(args[1]); err == nil {
			temp = ntemp
		}
	}
	sct.SetColorTemp(temp)
}
