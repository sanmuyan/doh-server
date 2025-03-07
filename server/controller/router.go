package controller

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RunServer(ctx context.Context, addr string) {
	r := gin.Default()
	router(r)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			if err != nil {
				logrus.Fatalf("server run error: %s", err)
			}
		}
	}()
	logrus.Infof("server listening on %s", addr)
	<-ctx.Done()
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("server shutdown error: %s", err)
	}
	logrus.Info("server has been shutdown")
}

func router(r *gin.Engine) {
	r.GET("/dns-query", DoHQuery)
	r.POST("/dns-query", DoHQuery)
	r.GET("/resolve", DJAQuery)
}

func RunDNServer(ctx context.Context, addr string, net string) {
	dns.HandleFunc(".", DNSQuery)
	server := &dns.Server{Addr: addr, Net: net}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Fatalf("%s server run error: %s", net, err)
		}
	}()
	<-ctx.Done()
	if err := server.Shutdown(); err != nil {
		logrus.Errorf("%s server shutdown error: %s", net, err)
	}
	logrus.Warnf("%s server has been shutdown", net)
}
