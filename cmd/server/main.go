package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"vnet/internal/server"
)

func handleShutdown(stop func()) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	stop()
}

func main() {
	srv := server.NewLayerTwoSwitch(context.Background(), &server.LayerTwoSwitchParams{
		ListenAddr: []byte{0x7F, 0x00, 0x00, 0x01},
		ListenPort: 10001,
		DebugMode:  true,
		Logger:     log.Default(),
		SaveFrames: true,
	})

	go handleShutdown(srv.Close)
	fmt.Println("listen finished", srv.Listen())
}
