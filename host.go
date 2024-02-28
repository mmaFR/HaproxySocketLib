package HaproxySocketLib

import (
	"crypto/tls"
	"io"
	"net"
)

type connectionType uint8

type Host struct {
	connType connectionType
	tcpAddr  *net.TCPAddr
	unixAddr *net.UnixAddr
	useTls   bool
	tcpConn  *net.TCPConn
	tls      *tls.Conn
	unixConn *net.UnixConn
}

func (h *Host) SetTcpAddress(t *net.TCPAddr) *Host {
	h.connType = connTcp
	h.tcpAddr = t
	return h
}
func (h *Host) SetUnixAddress(u *net.UnixAddr) *Host {
	h.connType = connUnix
	h.unixAddr = u
	return h
}
func (h *Host) UseTLS(b bool) *Host {
	h.useTls = b
	return h
}

func (h *Host) establishConnection() (err error) {
	switch h.connType {
	case connTcp:
		err = h.establishTcpConnection()
	case connUnix:
		err = h.establishUnixConnection()
	default:
		err = unknownConnectionType("host->establishConnection")
	}
	return
}
func (h *Host) establishTcpConnection() (err error) {
	if h.tcpConn, err = net.DialTCP(h.tcpAddr.Network(), nil, h.tcpAddr); err != nil {
		return
	}
	if _, err = h.tcpConn.Write([]byte("prompt\n")); err != nil {
		return
	}
	_, err = h.tcpConn.Read(make([]byte, 3))
	return
}
func (h *Host) establishUnixConnection() (err error) {
	if h.unixConn, err = net.DialUnix(h.unixAddr.Network(), nil, h.unixAddr); err != nil {
		return
	}
	_, err = h.unixConn.Write([]byte("prompt\n"))
	_, _ = io.ReadAll(h.unixConn)
	return
}
func (h *Host) connectionIsAlive(c net.Conn) (b bool) {
	if _, e := c.Write([]byte{0x0a}); e != nil {
		return false
	} else {
		if _, ee := h.tcpConn.Read(make([]byte, 3)); ee != nil {
			return false
		}
		return true
	}
}
func (h *Host) getConnection() (c net.Conn, e error) {
	switch h.connType {
	case connTcp:
		if h.tcpConn == nil { // connection not yet established
			if e = h.establishConnection(); e != nil { // failed to establish the new connection
				return
			} else { // new connection established with success
				c = h.tcpConn
				return
			}
		} else { // connection already established
			if h.connectionIsAlive(h.tcpConn) { // the connection is still alive
				c = h.tcpConn
				return
			} else { // the connection is dead
				if e = h.establishConnection(); e != nil { // fail to establish a new connection
					return
				} else { // new connection established with success
					c = h.tcpConn
					return
				}
			}
		}
	case connUnix:
		if h.unixConn == nil { // connection not yet established
			if e = h.establishConnection(); e != nil { // failed to establish the new connection
				return
			} else { // new connection established with success
				c = h.unixConn
				return
			}
		} else { // connection already established
			if h.connectionIsAlive(h.unixConn) { // the connection is still alive
				c = h.unixConn
				return
			} else { // the connection is dead
				if e = h.establishConnection(); e != nil { // fail to establish a new connection
					return
				} else { // new connection established with success
					c = h.unixConn
					return
				}
			}
		}
	default:
		e = unknownConnectionType("host->getConnection")
	}
	return
}
func (h *Host) sendCommand(cmd string) (r []byte, e error) {
	var c net.Conn
	if c, e = h.getConnection(); e != nil {
		return
	} else {
		if _, e = c.Write([]byte(cmd + "\n")); e != nil {
			return
		} else {
			var buffer1, buffer2 []byte = make([]byte, 0, 16384), make([]byte, 16384)
			var n int
			for loop := true; loop; {
				if n, e = c.Read(buffer2); e != nil {
					loop = false
				} else {
					buffer1 = append(buffer1, buffer2[:n]...)
					loop = n == 16384
				}
			}
			return buffer1, e
		}
	}
}

func NewHost() *Host {
	return new(Host)
}
