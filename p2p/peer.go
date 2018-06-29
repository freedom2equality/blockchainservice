package p2p

import (
	"fmt"
	"net"
	"time"
)

// PeerConn contains the raw connection
type PeerConn struct {
	outbound   bool
	persistent bool
	conn       net.Conn
}

func newPeerConn(rawConn net.Conn, outbound, persistent bool) (PeerConn, error) {
	return PeerConn{
		outbound:   outbound,
		persistent: persistent,
		conn:       rawConn,
	}, nil
}

// CloseConn should be called if the peer was created but never started.
func (pc *PeerConn) CloseConn() {
	pc.conn.Close()
}

// HandshakeTimeout performs the P2P handshake between a given node and the peer by exchanging their NodeInfo.
func (pc *PeerConn) HandshakeTimeout(timeout time.Duration, stage chan<- *PeerConn) error {
	// Set deadline for handshake so we don't block forever on conn.ReadFull
	if err := pc.conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		log.Error("Error setting deadline")
		pc.CloseConn()
		return fmt.Errorf("Error setting deadline")
	}
	if pc.outbound {
		// 发送版本号
		if err := pc.writeLocalVersionMsg(); err != nil {
			log.Error("writeLocalVersionMsg")
			pc.CloseConn()
			return fmt.Errorf("writeLocalVersionMsg")
		}
		// 读取版本号
		if err := pc.readRemoteVersionMsg(); err != nil {
			log.Error("readRemoteVersionMsg")
			pc.CloseConn()
			return fmt.Errorf("readRemoteVersionMsg")
		}
	} else {
		// 读取版本号
		if err := pc.readRemoteVersionMsg(); err != nil {
			log.Error("readRemoteVersionMsg")
			pc.CloseConn()
			return fmt.Errorf("readRemoteVersionMsg")
		}
		// 发送版本号
		if err := pc.writeLocalVersionMsg(); err != nil {
			log.Error("writeLocalVersionMsg")
			pc.CloseConn()
			return fmt.Errorf("writeLocalVersionMsg")
		}
	}

	// Remove deadline
	if err := pc.conn.SetDeadline(time.Time{}); err != nil {
		log.Error("Error removing deadline")
		return fmt.Errorf("Error removing deadline")
	}
	stage <- pc
	return nil
}

func (pc *PeerConn) readRemoteVersionMsg() error {
	return nil
}

func (pc *PeerConn) writeLocalVersionMsg() error {
	return nil
}
