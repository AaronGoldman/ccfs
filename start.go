package main

import ()

func start() {
	go func() { //defining, calling and throwing to a different thread
		ch := make(chan os.Signal, 1) //ch is the name of the channel.
		signal.Notify(ch, os.Interrupt, os.Kill)
		sig := <-ch
		log.Printf("Got signal: %s", sig)
		log.Printf("Stoping...")
		stop()
	}()
}
