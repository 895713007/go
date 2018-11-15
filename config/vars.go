package config

var (
	currentConfig *Config
)

const (
	ENV_DEV  = "dev"  // 开发环境（默认）
	ENV_BETA = "beta" // 测试环境
	ENV_PRO  = "pro"  // 生产环境
)

const (
	CFG_FROM_MONITOR_CENTER = 0 // 配置中心
	CFG_FROM_LOCAL_FILE     = 1 // 本地配置文件

	DEV_HOST_MONITOR_CENTER  = "http://dev.mars-job-platform.mytoken-local.com"  // 开发环境配置中心域名
	BETA_HOST_MONITOR_CENTER = "http://beta.mars-job-platform.mytoken-local.com" // 测试环境配置中心域名
	PRO_HOST_MONITOR_CENTER  = "http://mars-job-platform.mytoken-local.com"      // 正式环境配置中心域名
)
