package raft

import (
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/pkg/transport"
	"github.com/superioz/artemis/raft/protocol"
	"time"
)

func (n *Node) leaderLoop() {
	pc := n.transport.Receive()
	hb := time.NewTicker(n.heartbeatInterval * time.Millisecond)

	// send packet
	logrus.WithFields(logrus.Fields{
		"prefix": "out",
	}).Debugln("leader.o.heartbeat", n.id)
	n.sendHeartbeat()

	for n.state == Leader {
		select {
		case p := <-pc:
			m, err := transport.Decode(p.Packet.Data)
			if err != nil {
				logrus.Errorln("couldn't decode packet", err)
				break
			}

			switch m.(type) {
			case *protocol.RequestVoteCall:
				// * that can't possibly happen. But if so just ignore until
				// another heartbeat gets sent
				break
			case *protocol.RequestVoteRespond:
				// * we didn't call for a vote so just ignore this
				// packet.
				break
			case *protocol.AppendEntriesCall:
				// * a client redirected an append entries call.

				break
			case *protocol.AppendEntriesRespond:
				// * a client responded with the result of the append entries
				// call.

				appendEntrResp := *m.(*protocol.AppendEntriesRespond)

				logrus.WithFields(logrus.Fields{
					"prefix": "in",
				}).Debugln("leader.i.appendEntries.respond", appendEntrResp)

				break
			}
			break
		case <-hb.C:
			// * send normal heartbeat with no information
			// just to keep his authority.
			logrus.WithFields(logrus.Fields{
				"prefix": "out",
			}).Debugln("leader.o.heartbeat")
			n.sendHeartbeat()
			break
		}
	}
}
