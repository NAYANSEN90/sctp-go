package sctp_go

import (
	"net"
	"sync/atomic"
	"syscall"
	"time"
)

type SCTPConn struct {
	sock int64
}

func NewSCTPConn(sock int) *SCTPConn {
	return &SCTPConn{
		sock: int64(sock),
	}
}

func (conn *SCTPConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (conn *SCTPConn) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (conn *SCTPConn) SendMsg(b []byte, info *SCTPSndRcvInfo) (int, error) {
	var buffer []byte
	if nil != info {
		hdr := &syscall.Cmsghdr{
			Level: syscall.IPPROTO_SCTP,
			Type:  SCTP_SNDRCV,
			Len:   uint64(syscall.CmsgSpace(SCTPSndRcvInfoSize)),
		}
		buffer = append(buffer, Pack(hdr)...)
		buffer = append(buffer, Pack(info)...)
	}
	return syscall.SendmsgN(int(conn.sock), b, buffer, nil, 0)
}

func (conn *SCTPConn) Close() error {
	if !conn.ok() {
		return syscall.EINVAL
	}
	sock := atomic.SwapInt64(&conn.sock, -1)
	if sock > 0 {
		msg := &SCTPSndRcvInfo{
			Flags: SCTP_EOF,
		}
		_, _ = conn.SendMsg(nil, msg)
		_ = syscall.Shutdown(int(sock), syscall.SHUT_RDWR)
		return syscall.Close(int(sock))
	}
	return syscall.EBADFD
}

func (conn *SCTPConn) LocalAddr() net.Addr {
	return nil
}

func (conn *SCTPConn) RemoteAddr() net.Addr {
	return nil
}

func (conn *SCTPConn) SetDeadline(t time.Time) error {
	return nil
}

func (conn *SCTPConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (conn *SCTPConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (conn *SCTPConn) ok() bool {
	if nil != conn && conn.sock > 0 {
		return true
	}
	return false
}

func DialSCTP(network string, local, remote *SCTPAddr, init *SCTPInitMsg) (*SCTPConn, error) {
	switch network {
	case "sctp", "sctp4", "sctp6":
	default:
		return nil, &net.OpError{
			Op:     "dial",
			Net:    network,
			Source: local.Addr(),
			Addr:   remote.Addr(),
			Err:    net.UnknownNetworkError(network),
		}
	}
	conn := &SCTPConn{}
	return conn, nil
}
