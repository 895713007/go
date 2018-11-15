package metrics

import (
	"os"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/mytokenio/go/log"
)

func init() {
	exitChan = make(chan struct{})
	countMap = make(map[string]int64, default_map_caps)
	gaugeIntMap = make(map[string]int64, default_map_caps)
	gaugeStrMap = make(map[string]string, default_map_caps)

	var envType int
	var broker, serviceName string
	switch os.Getenv(ENV_ENV_TYPE) {
	case ENV_BETA:
		envType = ENV_TYPE_BETA
		broker = beta_default_kafka_broker
	case ENV_PRO:
		envType = ENV_TYPE_PRO
		broker = pro_default_kafka_broker
	default:
		envType = ENV_TYPE_DEV
		broker = dev_default_kafka_broker
	}

	// init serviceInfo value
	host, _ := os.Hostname()
	jobId, _ := strconv.ParseInt(os.Getenv(ENV_JOB_ID), 10, 64)
	if serviceName = os.Getenv(ENV_SERVICE_NAME); serviceName == "" {
		serviceName = DEF_SERVICE_NAME
	}
	globalServiceInfo.jobID = jobId
	globalServiceInfo.serviceName = serviceName
	globalServiceInfo.envType = envType
	globalServiceInfo.host = host
	globalServiceInfo.processID = os.Getpid()

	err := initKafka([]string{broker}, default_roport_state_topic, default_roport_alarm_topic)
	if err != nil {
		log.Errorf("default init kafka err: %v", err)
		return
	}
}

func Init(brokers []string, stateTopic, alarmTopic string) error {

	return initKafka(brokers, stateTopic, alarmTopic)
}

// ---------------------------------------------------------------------------------------------------------------------

func initKafka(brokers []string, stateTopic, alarmTopic string) error {
	if p, err := createProducer(brokers); err != nil {
		return err
	} else {
		globalKafka.producer = p
		globalKafka.brokers = brokers
		globalKafka.reportStateTopic = stateTopic
		globalKafka.reportAlarmTopic = alarmTopic
		globalKafka.chanStateProducerValue = make(chan string, default_producer_msg_caps)
		globalKafka.chanAlarmProducerValue = make(chan string, default_producer_msg_caps)

		if !globalKafka.isInitialized {
			go cronMonitor()
			go reportMonitorCenter()
			go callback()
		}

		globalKafka.isInitialized = true

		Gauge("start_time", time.Now().Unix())
	}

	return nil
}

func createProducer(brokers []string) (sarama.AsyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Version = sarama.V0_11_0_2
	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
