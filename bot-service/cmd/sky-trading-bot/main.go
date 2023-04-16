package main

import (
	"fmt"
	"github.com/serhiq/skye-trading-bot/internal/config"
	"github.com/serhiq/skye-trading-bot/internal/logger"
	"github.com/serhiq/skye-trading-bot/internal/server"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {

	cfg, err := config.New()

	if err != nil {
		log.Panic(err)
	}

	err = logger.InitLogger(*cfg)
	defer logger.Sync()

	s, err := server.Serve(*cfg)
	if err != nil {
		log.Panic(err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	go func() {
		sig := <-exit
		logger.SugaredLogger.Info(fmt.Sprintf("Exit. %s  received.", sig.String()))
		fmt.Printf("\n\n Goroutines: %d \n", runtime.NumGoroutine())

		s.Stop()
		fmt.Println("Shutdown server")
	}()

	err = s.Start()
	if err != nil {
		logger.SugaredLogger.Fatalf("Server start error: %s", err)
	}
}
