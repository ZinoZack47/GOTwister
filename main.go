package main

//GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go build -buildmode=c-shared -o gotwister.dll ./

import (
	//"C"
	features "gotwister/Features"
	"time"
)

func main() {
	for {
		features.Radar()
		time.Sleep(100 * time.Millisecond)
	}
}
