package udp

import (
	"fmt"
	"net"
	"time"
)

// MaxSafePayloadSize is the maximum "safe" UDP payload (508 bytes).
// Useful for sizing read buffers.
const MaxSafePayloadSize = 508

// Port is the UDP port WiZ bulbs communicate on.
const Port = 38899

type UDPSession struct {
	conn   net.PacketConn
	remote net.Addr
}

func NewSession(ip string, timeout time.Duration) (*UDPSession, error) {
	conn, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", Port))
	if err != nil {
		return nil, err
	}
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		conn.Close()
		return nil, err
	}
	remote, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", ip, Port))
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &UDPSession{conn: conn, remote: remote}, nil
}

func (s *UDPSession) Read(buf []byte) (int, net.Addr, error) {
	return s.conn.ReadFrom(buf)
}

func (s *UDPSession) Write(buf []byte) (int, error) {
	return s.conn.WriteTo(buf, s.remote)
}

func (s *UDPSession) Close() error {
	return s.conn.Close()
}
