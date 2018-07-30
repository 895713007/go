package log

import (
	stdlog "log"
	"net"
	"time"

	"github.com/gin-gonic/gin/json"
)

const (
	_agentTimeout = time.Duration(20 * time.Millisecond)
	_retryDelay   = 5 * time.Second
)

// AgentHandler agent struct.
type AgentHandler struct {
	conf *AgentConfig
	conn net.Conn
	msgs chan D
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
		conf: c,
		msgs: make(chan D, c.Chan),
	}
	if c.Timeout == 0 {
		c.Timeout = _agentTimeout
	}
	if c.Buffer == 0 {
		c.Buffer = 1
	}

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
		b    []byte
		conn net.Conn
		err  error
	)
	tick := time.NewTicker(_retryDelay)
	for {
		if conn == nil || err != nil {
			if conn, err = net.DialTimeout(h.conf.Proto, h.conf.Addr, time.Duration(h.conf.Timeout)); err != nil {
				stdlog.Printf("net.DialTimeout(%s:%s) error(%v)\n", h.conf.Proto, h.conf.Addr, err)
				continue
			}
		}
		select {
		case f := <-h.msgs:
			b, err = json.Marshal(f)
			if err != nil {
				stdlog.Printf("json.Marshal(%v) error(%v)\n", f, err)
				continue
			}
			if _, err = conn.Write(b); err != nil {
				stdlog.Printf("conn.Write(%s bytes) error(%v)\n", string(b), err)
				conn.Close()
			}
		case <-tick.C: // time out retry
		}
	}
}

// Close close the connection.
func (h *AgentHandler) Close() (err error) {
	h.msgs <- nil
	return nil
}
