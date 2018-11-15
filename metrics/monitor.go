package metrics

import (
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cihub/seelog"
	"github.com/mytokenio/go/log"
)

func cronMonitor() {
	t := time.NewTicker(cronInterval)

	for {
		select {
		case <-exitChan:
			return
		case <-t.C:
			reportStateFactory()
		}
	}
}

func reportStateFactory() {
	now := time.Now().Unix()
	rs := ReportStatePkg{
		JobID:       globalServiceInfo.jobID,
		ServiceName: globalServiceInfo.serviceName,
		EnvType:     globalServiceInfo.envType,
		Host:        globalServiceInfo.host,
		ProcessID:   globalServiceInfo.processID,
		Extend:      make(map[string]interface{}),
	}

	if v, ok := gaugeIntMap["status"]; !ok {
		rs.Status = STATUS_OK
	} else {
		rs.Status = int(v)
	}

	if v, ok := gaugeIntMap["start_time"]; ok {
		rs.StartTime = v
	}

	if v, ok := gaugeIntMap["hear_time"]; ok {
		rs.HearTime = v
	} else {
		rs.HearTime = now
	}

	if v, ok := gaugeIntMap["exit_code"]; ok {
		rs.ExitCode = int(v)
		if v > 0 {
			rs.StopTime = now
		}
	}

	// set extend
	for key, value := range countMap {
		if !isDefaultKey[key] {
			rs.Extend[key] = value
		}
	}
	for key, value := range gaugeIntMap {
		if !isDefaultKey[key] {
			rs.Extend[key] = value
		}
	}
	for key, value := range gaugeStrMap {
		if !isDefaultKey[key] {
			rs.Extend[key] = value
		}
	}

	value, err := json.Marshal(rs)
	if err != nil {
		log.Errorf("json marshal reportState: %+v err: %v", rs, err)
		return
	}

	globalKafka.chanStateProducerValue <- string(value)
}

func alarm(content string) {
	if content != "" {
		ra := ReportAlarmPkg{
			JobID:       globalServiceInfo.jobID,
			ServiceName: globalServiceInfo.serviceName,
			Content:     content,
			HearTime:    time.Now().Unix(),
		}
		pkg, _ := json.Marshal(ra)
		globalKafka.chanAlarmProducerValue <- string(pkg)
	}
}

func reportMonitorCenter() {
	var value string
	var isNotClosed bool

	for {
		select {

		// receive value from chan_alarm
		case value, isNotClosed = <-globalKafka.chanAlarmProducerValue:
			if !isNotClosed {
				return
			}

			globalKafka.producer.Input() <- &sarama.ProducerMessage{
				Topic: globalKafka.reportAlarmTopic,
				Key:   sarama.StringEncoder(alarm_producter_msg_key),
				Value: sarama.ByteEncoder(value),
			}

		// receive value from chan_state
		case value, isNotClosed = <-globalKafka.chanStateProducerValue:
			if !isNotClosed {
				return
			}

			globalKafka.producer.Input() <- &sarama.ProducerMessage{
				Topic: globalKafka.reportStateTopic,
				Key:   sarama.StringEncoder(state_producter_msg_key),
				Value: sarama.ByteEncoder(value),
			}
		}
	}
}

func callback() error {
	var sucValue, failValue []byte
	var suc *sarama.ProducerMessage
	var fail *sarama.ProducerError

	for {
		select {
		case <-exitChan:
			return nil
		case suc = <-globalKafka.producer.Successes():
			sucValue, _ = suc.Value.Encode()
			seelog.Infof("send alarm msg success. [T:%s P:%d O:%d M:%s]",
				suc.Topic, suc.Partition, suc.Offset, string(sucValue))
		case fail = <-globalKafka.producer.Errors():
			failValue, _ = fail.Msg.Value.Encode()
			seelog.Errorf("send alarm msg failed: [T:%s M:%s], err: %s",
				suc.Topic, string(failValue), fail.Err.Error())
		}
	}
}
