package cmd

import (
	"context"
	"doh-server/pkg/configpost"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCtx context.Context

var rootCmd = &cobra.Command{
	Use:   "doh-server",
	Short: "DoH Server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := initConfig(cmd)
		if err != nil {
			logrus.Fatalf("init config error: %v", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		configpost.PostInit(rootCtx)
	},
}

var configFile string

const (
	logLevel        = 4
	serverBind      = ":8053"
	udpBind         = ""
	tcpBind         = udpBind
	cache           = false
	cacheTTL        = 60
	upstreamServer  = "8.8.8.8:53"
	upstreamNet     = "udp"
	upstreamTimeout = 2
)

func init() {
	// 初始化命令行参数
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().IntP("log-level", "l", logLevel, "log level")
	rootCmd.PersistentFlags().BoolP("pprof-server", "", false, "enable pprof server")
	rootCmd.Flags().String("server-bind", serverBind, "server bind addr")
	rootCmd.Flags().String("udp-bind", udpBind, "udp bind addr")
	rootCmd.Flags().String("tcp-bind", tcpBind, "tcp bind addr")
	rootCmd.Flags().BoolP("cache", "C", cache, "enable cache")
	rootCmd.Flags().IntP("cache-ttl", "T", cacheTTL, "cache ttl (ms)")
	rootCmd.Flags().StringP("upstream-server", "s", upstreamServer, "upstream dns server")
	rootCmd.Flags().StringP("upstream-net", "n", upstreamNet, "upstream dns net (udp|tcp|tcp-tls|doh)")
	rootCmd.Flags().IntP("upstream-timeout", "t", upstreamTimeout, "upstream dns timeout (ms)")
}

func Execute(ctx context.Context) {
	rootCtx = ctx
	if err := rootCmd.Execute(); err != nil {
		logrus.Tracef("cmd execute error: %v", err)
	}
}
