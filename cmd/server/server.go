package main

import (
	"context"
	"doh-server/cmd/server/cmd"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigs
		logrus.Warn("shutting down server...")
		cancel()
	}()
	cmd.Execute(ctx)
}
