package main

import (
	features "gotwister/Features"
	"time"
)

func main() {
	for {
		features.Radar()
		time.Sleep(100 * time.Millisecond)
	}
}
