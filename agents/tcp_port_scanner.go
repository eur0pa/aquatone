package agents

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/eur0pa/aquatone/core"
)

type TCPPortScanner struct {
	session *core.Session
}

func NewTCPPortScanner() *TCPPortScanner {
	return &TCPPortScanner{}
}

func (d *TCPPortScanner) ID() string {
	return "agent:tcp_port_scanner"
}

func (a *TCPPortScanner) Register(s *core.Session) error {
	s.EventBus.SubscribeAsync(core.Host, a.OnHost, false)
	a.session = s
	return nil
}

func (a *TCPPortScanner) OnHost(host string) {
	a.session.Out.Debug("[%s] Received new host: %s\n", a.ID(), host)
	if strings.Contains(host, ":") {
		x := strings.Split(host, ":")
		host = x[0]
		port, _ := strconv.Atoi(x[1])
		a.session.WaitGroup.Add()
		go func(port int, host string) {
			defer a.session.WaitGroup.Done()
			a.session.EventBus.Publish(core.TCPPort, port, host)
		}(port, host)
	}
}

func (a *TCPPortScanner) scanPort(port int, host string) bool {
	conn, _ := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Duration(*a.session.Options.ScanTimeout)*time.Millisecond)
	if conn != nil {
		conn.Close()
		return true
	}
	return false
}
