package udp

import (
	"fmt"
	"net"
	"time"
)

// The maximum "safe" UDP payload is 508 bytes.
// This is useful for buffer sizing.
const MAX_SAFE_PAYLOAD_SIZE = 508

// WiZ bulbs communicate on UDP port 38899
const PORT = 38899

type UDPSession struct {
	conn   net.PacketConn
	remote net.Addr
}

func NewSession(ip string, timeout time.Duration) (*UDPSession, error) {
	conn, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", PORT))
	if err != nil {
		return nil, err
	}
	conn.SetReadDeadline(time.Now().Add(timeout))
	remote, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", ip, PORT))
	if err != nil {
		return nil, err
	}
	session := &UDPSession{
		conn:   conn,
		remote: remote,
	}
	return session, nil
}

func (s *UDPSession) Read(buf []byte) (int, net.Addr, error) {
	return s.conn.ReadFrom(buf)
}

func (s *UDPSession) Write(buf []byte) (int, error) {
	return s.conn.WriteTo([]byte(buf), s.remote)
}

func (s *UDPSession) Close() error {
	return s.conn.Close()
}
