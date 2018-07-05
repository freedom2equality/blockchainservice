package p2p

type Reactor interface {
	// GetChannels returns the list of channel descriptors.
	//GetChannels() []*conn.ChannelDescriptor
	Receive(chID byte, conn PeerConn, msgBytes []byte)
}
