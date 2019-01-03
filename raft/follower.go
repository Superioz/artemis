package raft

import (
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/pkg/transport"
	"github.com/superioz/artemis/pkg/uid"
	"github.com/superioz/artemis/raft/protocol"
	"time"
)

func (n *Node) followerLoop() {
	timeout := n.generateTimeout()
	tc := time.NewTimer(timeout)
	pc := n.transport.Receive()

	// reset some values
	n.votedFor = uid.UID{}

followerLoop:
	for n.state == Follower {
		select {
		case p := <-pc:
			m, err := transport.Decode(p.Packet.Data)
			if err != nil {
				logrus.Errorln("couldn't decode packet", err)
				break
			}

			switch m.(type) {
			case *protocol.RequestVoteCall:
				// * check if he can grant the vote and respond with
				// his answer.

				reqVote := *m.(*protocol.RequestVoteCall)

				logrus.WithFields(logrus.Fields{
					"prefix": "in",
				}).Debugln("follower.i.voterequest", n.id, m)

				n.processRequestVote(reqVote, p.Source.String())
				break
			case *protocol.RequestVoteRespond:
				// * ignore as followers can't process this packet. only
				// candidates can.
				break
			case *protocol.AppendEntriesCall:
				// * the leader sends an append entries call, so
				// update the leader and reset the timeout.

				appendEntr := *m.(*protocol.AppendEntriesCall)
				logrus.WithFields(logrus.Fields{
					"prefix": "in",
				}).Debugln("follower.i.appendEntries", n.id, appendEntr)

				// reset timeout
				timeout = n.generateTimeout()
				tc.Reset(timeout)

				// append entries
				err := n.processAppendEntries(appendEntr, p.Source.String())
				n.sendAppendEntriesResponse(err == nil, p.Source.String())
				break
			case *protocol.AppendEntriesRespond:
				// * ignore as followers can't process this packet. only
				// the leader can.
				break
			}
			break
		case <-tc.C:
			// * timeout and try to become leader by sending request votes
			// also step up to being candidate
			logrus.Debugln("follower.timeout", n.id, timeout)

			// if the node is only passive, don't try to ever get leader
			if n.Passive {
				break followerLoop
			}

			n.SetState(Candidate)
			break
		}
	}
}
