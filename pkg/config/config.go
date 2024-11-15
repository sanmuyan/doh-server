package config

// Config 全局配置
type Config struct {
	// 日志等级
	LogLevel int `mapstructure:"log_level"`
	// 服务绑定地址
	ServerBind string `mapstructure:"server_bind"`
	// UDP 绑定地址
	UDPBind string `mapstructure:"udp_bind"`
	// TCP 绑定地址
	TCPBind string `mapstructure:"tcp_bind"`
	// 是否启用缓存
	Cache bool `mapstructure:"cache"`
	// 缓存过期时间(s)
	CacheTTL int `mapstructure:"cache_ttl"`
	// 上游 DNS 服务器
	UpstreamServer string `mapstructure:"upstream_server"`
	// 上游 DNS 服务器网络类型
	UpstreamNet string `mapstructure:"upstream_net"`
	// 上游 DNS 超时时间(s)
	UpstreamTimeout int `mapstructure:"upstream_timeout"`
}

var Conf Config
