package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	stdlog "log"
	"net"
	"sync"
	"time"
)

const (
	_agentTimeout = time.Duration(20 * time.Millisecond)
	_retryDelay   = 5 * time.Second
	_maxBuffer    = 10 * 1024 * 1024 // 10mb
)

var (
	_logSeparator = []byte("\u0001")
)

// AgentHandler agent struct.
type AgentHandler struct {
	c      *AgentConfig
	msgs   chan D
	waiter sync.WaitGroup
}

// AgentConfig agent config.
type AgentConfig struct {
	TaskID  string
	Proto   string
	Addr    string
	Buffer  int
	Chan    int
	Timeout time.Duration
}

// NewAgent a Agent.
func NewAgent(c *AgentConfig) (a *AgentHandler) {
	a = &AgentHandler{
		c:    c,
		msgs: make(chan D, c.Chan),
	}
	if c.Timeout == 0 {
		c.Timeout = _agentTimeout
	}
	if c.Buffer == 0 {
		c.Buffer = 1
	}
	a.waiter.Add(1)
	go a.writeProc()
	return
}

// Log log to udp statsd daemon.
func (h *AgentHandler) Log(lv Level, d D) {
	if d == nil {
		return
	}
	select {
	case h.msgs <- d:
	default:
	}
}

// writeProc write data into connection.
func (h *AgentHandler) writeProc() {
	var (
		buf   bytes.Buffer
		conn  net.Conn
		err   error
		count int
		quit  bool
	)
	defer h.waiter.Done()
	taskID := []byte(h.c.TaskID)
	tick := time.NewTicker(_retryDelay)
	enc := json.NewEncoder(&buf)
	for {
		select {
		case f := <-h.msgs:
			if f == nil {
				quit = true
				goto DUMP
			}
			if buf.Len() >= _maxBuffer {
				buf.Reset() // avoid oom
			}
			now := time.Now()
			buf.Write(taskID)
			buf.Write([]byte(fmt.Sprintf("%d", now.UnixNano()/1e6)))
			enc.Encode(f)
			if count++; count < h.c.Buffer {
				buf.Write(_logSeparator)
				continue
			}
		case <-tick.C:
		}
		if conn == nil || err != nil {
			if conn, err = net.DialTimeout(h.c.Proto, h.c.Addr, time.Duration(h.c.Timeout)); err != nil {
				stdlog.Printf("net.DialTimeout(%s:%s) error(%v)\n", h.c.Proto, h.c.Addr, err)
				continue
			}
		}
	DUMP:
		if conn != nil && buf.Len() > 0 {
			count = 0
			if _, err = conn.Write(buf.Bytes()); err != nil {
				stdlog.Printf("conn.Write(%d bytes) error(%v)\n", buf.Len(), err)
				conn.Close()
			} else {
				// only succeed reset buffer, let conn reconnect.
				buf.Reset()
			}
		}
		if quit {
			if conn != nil && err == nil {
				conn.Close()
			}
			return
		}
	}
}

// Close close the connection.
func (h *AgentHandler) Close() (err error) {
	h.msgs <- nil
	h.waiter.Wait()
	return nil
}
