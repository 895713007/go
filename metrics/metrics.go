package metrics

func Count(id string, delta int64) {
	if id == "" {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := countMap[id]; !ok {
		countMap[id] = delta
	} else {
		countMap[id] += delta
	}
}

func Gauge(id string, value interface{}) {
	if id == "" {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	switch value.(type) {
	case int:
		gaugeIntMap[id] = int64(value.(int))
	case int64:
		gaugeIntMap[id] = int64(value.(int64))
	case string:
		gaugeStrMap[id] = value.(string)
	default:
		return
	}
}

func Close() {
	mutex.Lock()
	defer mutex.Unlock()

	// report cache data
	reportStateFactory()

	// resource recovery
	close(exitChan)
	close(globalKafka.chanStateProducerValue)
	close(globalKafka.chanAlarmProducerValue)
	globalKafka.producer.Close()

	for key, _ := range countMap {
		delete(countMap, key)
	}

	for key, _ := range gaugeIntMap {
		delete(gaugeIntMap, key)
	}

	for key, _ := range gaugeStrMap {
		delete(gaugeStrMap, key)
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func StatusOK() {
	Gauge("status", STATUS_OK)
}

func StatusError() {
	Gauge("status", STATUS_ERROR)
}

func ExitWithOK() {
	Gauge("status", STATUS_OK)
	Gauge("exit_code", EXIT_CODE_OK)
}

func ExitWithErr(alarmMsg ...string) {
	if len(alarmMsg) > 0 {
		alarm(alarmMsg[0])
	}
	Gauge("status", STATUS_ERROR)
	Gauge("exit_code", EXIT_CODE_ERROR)
}

func ExitWithKill(alarmMsg ...string) {
	if len(alarmMsg) > 0 {
		alarm(alarmMsg[0])
	}
	Gauge("status", STATUS_ERROR)
	Gauge("exit_code", EXIT_CODE_KILL)
}

func Panic(err error) {
	alarm(err.Error())
	Gauge("status", STATUS_ERROR)
	Gauge("exit_code", EXIT_CODE_ERROR)
	panic(err)
}

func Alarm(alarmMsg string) {
	if len(alarmMsg) > 0 {
		alarm(alarmMsg)
	}
}
