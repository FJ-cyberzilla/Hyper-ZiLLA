package network

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

type ServiceDetector struct {
	timeout time.Duration
}

func NewServiceDetector() *ServiceDetector {
	return &ServiceDetector{
		timeout: 2 * time.Second,
	}
}

func (sd *ServiceDetector) Detect(ctx context.Context, target string, port int) string {
	address := net.JoinHostPort(target, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", address, sd.timeout)
	if err != nil {
		return "Unknown"
	}
	defer conn.Close()

	// Set deadline for banner grab
	conn.SetDeadline(time.Now().Add(sd.timeout))

	// Send generic probe
	fmt.Fprintf(conn, "HEAD / HTTP/1.0\r\n\r\n")

	reader := bufio.NewReader(conn)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	if line != "" {
		if strings.Contains(line, "HTTP") {
			return "Web Server (HTTP)"
		}
		if strings.Contains(strings.ToLower(line), "ssh") {
			return "SSH"
		}
		return line // Return raw banner if matched
	}

	return sd.commonPortMap(port)
}

func (sd *ServiceDetector) commonPortMap(port int) string {
	ports := map[int]string{
		21:   "FTP",
		22:   "SSH",
		23:   "Telnet",
		25:   "SMTP",
		53:   "DNS",
		80:   "HTTP",
		443:  "HTTPS",
		3306: "MySQL",
		5432: "PostgreSQL",
		8080: "HTTP-Proxy/Alt",
	}
	if name, ok := ports[port]; ok {
		return name
	}
	return "Unknown Service"
}
