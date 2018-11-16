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

func GetCount(id string) (int64, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	if v, ok := countMap[id]; !ok {
		return 0, false
	} else {
		return v, true
	}
}

func GetCountMap() map[string]int64 {
	mutex.Lock()
	defer mutex.Unlock()

	return countMap
}

func GetGaugeInt64(id string) (int64, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	if v, ok := gaugeIntMap[id]; !ok {
		return 0, false
	} else {
		return v, true
	}
}

func GetGaugeInt64Map() map[string]int64 {
	mutex.Lock()
	defer mutex.Unlock()

	return gaugeIntMap
}

func GetGaugeStr(id string) (string, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	if v, ok := gaugeStrMap[id]; !ok {
		return "", false
	} else {
		return v, true
	}
}

func GetGaugeStrMap() map[string]string {
	mutex.Lock()
	defer mutex.Unlock()

	return gaugeStrMap
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
