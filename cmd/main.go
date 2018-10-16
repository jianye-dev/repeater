package main

import (
	"os"
	"os/signal"
	"vnc_repeater/repeater"
)

func main() {

	sig := make(chan os.Signal)

	signal.Notify(sig, os.Interrupt)

	rep := repeater.New()

	go func() {
		rep.ListenAndServe(":5500", ":5902", "")
	}()

	<-sig

	rep.Shutdown()

	<-sig

	close(sig)
}
