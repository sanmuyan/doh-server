package cmd

import (
	"context"
	"doh-server/pkg/config"
	"doh-server/pkg/configpost"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
	"runtime"
)

var rootCtx context.Context

var rootCmd = &cobra.Command{
	Use:   "doh-server",
	Short: "DoH Server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := initConfig()
		if err != nil {
			logrus.Fatalf("init config error: %v", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		configpost.PostInit(rootCtx)
	},
	Example: "doh-server -c config.yaml",
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
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.Flags().IntP("log-level", "l", logLevel, "log level")
	rootCmd.Flags().String("server-bind", serverBind, "server bind addr")
	rootCmd.Flags().String("udp-bind", udpBind, "udp bind addr")
	rootCmd.Flags().String("tcp-bind", tcpBind, "tcp bind addr")
	rootCmd.Flags().BoolP("cache", "C", cache, "enable cache")
	rootCmd.Flags().IntP("cache-ttl", "T", cacheTTL, "cache ttl (ms)")
	rootCmd.Flags().StringP("upstream-server", "s", upstreamServer, "upstream dns server")
	rootCmd.Flags().StringP("upstream-net", "n", upstreamNet, "upstream dns net (udp|tcp|tcp-tls|doh)")
	rootCmd.Flags().IntP("upstream-timeout", "t", upstreamTimeout, "upstream dns timeout (ms)")
}

func initConfig() error {
	// 设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return frame.Function, fileName
		},
	})

	viper.SetConfigName("config")
	// 配置文件和命令行参数都不指定时的默认配置
	// viper.SetDefault("conn_timeout", 10)

	// 设置默认配置文件，如果配置文件是可选的则可以不设置
	//if len(configFile) == 0 {
	//	dir, err := os.Getwd()
	//	if err != nil {
	//		return err
	//	}
	//	configFile = dir + "/config.yaml"
	//	configFile = filepath.Clean(configFile)
	//}

	// 读取配置文件
	if len(configFile) > 0 {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
	}

	// 绑定命令行参数到配置项
	// 配置项优先级：命令行参数 > 配置文件 > 默认命令行参数
	_ = viper.BindPFlag("log_level", rootCmd.Flags().Lookup("log-level"))
	_ = viper.BindPFlag("server_bind", rootCmd.Flags().Lookup("server-bind"))
	_ = viper.BindPFlag("udp_bind", rootCmd.Flags().Lookup("udp-bind"))
	_ = viper.BindPFlag("tcp_bind", rootCmd.Flags().Lookup("tcp-bind"))
	_ = viper.BindPFlag("cache", rootCmd.Flags().Lookup("cache"))
	_ = viper.BindPFlag("cache_ttl", rootCmd.Flags().Lookup("cache-ttl"))
	_ = viper.BindPFlag("upstream_server", rootCmd.Flags().Lookup("upstream-server"))
	_ = viper.BindPFlag("upstream_net", rootCmd.Flags().Lookup("upstream-net"))
	_ = viper.BindPFlag("upstream_timeout", rootCmd.Flags().Lookup("upstream-timeout"))

	err := viper.Unmarshal(&config.Conf)
	if err != nil {
		return err
	}
	logrus.SetLevel(logrus.Level(config.Conf.LogLevel))
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	if logrus.Level(config.Conf.LogLevel) >= logrus.DebugLevel {
		gin.SetMode(gin.DebugMode)
		gin.DefaultWriter = os.Stdout
		logrus.SetReportCaller(true)
	}
	logrus.Debugf("config init completed: %+v", config.Conf)
	return nil
}

func Execute(ctx context.Context) {
	rootCtx = ctx
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("cmd execute error: %v", err)
	}
}
