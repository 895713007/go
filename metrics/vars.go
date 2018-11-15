package metrics

import (
	"sync"
	"time"
)

var (
	mutex             sync.Mutex
	countMap          map[string]int64
	gaugeIntMap       map[string]int64
	gaugeStrMap       map[string]string
	globalServiceInfo serviceInfo
	globalKafka       kafkaInfo
	exitChan          chan struct{}
)

var (
	cronInterval = 120 * time.Second // 定时任务的间隔时间
)

var isDefaultKey = map[string]bool{
	"job_id":       true,
	"service_name": true,
	"status":       true,
	"env_type":     true,
	"start_time":   true,
	"stop_time":    true,
	"hear_time":    true,
	"exit_code":    true,
	"host":         true,
	"process_id":   true,
	"memory":       true,
	"load":         true,
	"net_in":       true,
	"net_out":      true,
}

// ---------------------------------------------------------------------------------------------------------------------

const (
	STATUS_ERROR     = 0 // 状态异常
	STATUS_OK        = 1 // 状态正常
	EXIT_CODE_UNEXIT = 0 // 服务未退出
	EXIT_CODE_OK     = 1 // 服务正常退出
	EXIT_CODE_ERROR  = 2 // 服务异常退出
	EXIT_CODE_KILL   = 3 // 服务被杀
)

const (
	ENV_DEV  = "dev"
	ENV_BETA = "beta"
	ENV_PRO  = "pro"

	ENV_TYPE_DEV  = 0
	ENV_TYPE_BETA = 1
	ENV_TYPE_PRO  = 4

	ENV_JOB_ID       = "JOB_ID"
	ENV_SERVICE_NAME = "SERVICE_NAME"
	ENV_ENV_TYPE     = "ENV_TYPE"

	DEF_SERVICE_NAME = "undefined"
)

const (
	default_map_caps           = 20
	default_producer_msg_caps  = 20
	default_roport_state_topic = "roport_state_topic"
	default_roport_alarm_topic = "roport_alarm_topic"
	dev_default_kafka_broker   = "todo"
	beta_default_kafka_broker  = "todo"
	pro_default_kafka_broker   = "todo"
)

const (
	state_producter_msg_key = "state_monitor_center"
	alarm_producter_msg_key = "alarm_monitor_center"
)
