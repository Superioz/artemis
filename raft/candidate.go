package raft

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/pkg/transport"
	"github.com/superioz/artemis/pkg/uid"
	"github.com/superioz/artemis/raft/protocol"
	"time"
)

func (n *Node) candidateLoop() {
	timeout := n.generateTimeout()

	timeoutTimer := time.NewTimer(timeout)
	packetChan := n.transport.Receive()
	hardBeetTimer := time.NewTicker(n.heartbeatInterval * time.Millisecond)

	// increment term and vote for himself
	n.currentTerm++
	n.votedFor = n.id
	n.currentVotes[n.id] = true

	// send request vote packet function
	logrus.WithFields(logrus.Fields{
		"prefix": "out",
	}).Debugln("candidate.o.voterequest", n.id)
	n.sendRequestVote()

candidateLoop:
	for n.state == Candidate {
		select {
		case p := <-packetChan:
			m, err := transport.Decode(p.Packet.Data)
			if err != nil {
				logrus.Errorln("couldn't decode packet", err)
				break
			}

			switch m.(type) {
			case *protocol.RequestVoteCall:
				// * sent back false, as we already voted for ourself

				reqVote := *m.(*protocol.RequestVoteCall)
				logrus.WithFields(logrus.Fields{
					"prefix": "in",
				}).Debugln("candidate.i.voterequest")

				// if term of packet is higher than current term
				// update and step back
				if reqVote.Term > n.currentTerm {
					n.currentTerm = reqVote.Term
					n.SetState(Follower)
					break candidateLoop
				}

				n.processRequestVote(reqVote, p.Source.String())
				break
			case *protocol.RequestVoteRespond:
				// * count positive votes until now and decide whether to be
				// leader or not

				vote := *m.(*protocol.RequestVoteRespond)

				n.currentVotes[p.Source] = vote.VoteGranted

				count := 0
				for _, v := range n.currentVotes {
					if v {
						count++
					}
				}
				quorum := n.QuorumSize()

				logrus.WithFields(logrus.Fields{
					"id": n.id,
					"quorum": quorum,
				}).Infoln(fmt.Sprintf("Received votes: %d/%d", count, n.config.ClusterSize))
				if count >= quorum {
					n.currentVotes = make(map[uid.UID]bool)

					logrus.WithFields(logrus.Fields{
						"prefix": "in",
					}).Debugln("candidate.i.voterespond", n.id)
					n.SetState(Leader)
				}
				break candidateLoop
			case *protocol.AppendEntriesCall:
				// * step back from being candidate, as there is already a leader
				// sending append entries.
				// * also reset timeout

				appendEntr := *m.(*protocol.AppendEntriesCall)
				logrus.WithFields(logrus.Fields{
					"prefix": "in",
				}).Debugln("candidate.i.appendEntries", n.id, appendEntr)

				// reset timeout
				timeout = n.generateTimeout()
				timeoutTimer.Reset(timeout)

				// set leader
				n.leader, _ = uuid.FromString(appendEntr.LeaderId)

				n.SetState(Follower)
				break candidateLoop
			case *protocol.AppendEntriesRespond:
				// * ignore packet, as we are a candidate and can't
				// process these packets. Only the leader can.
				break
			}
			break
		case <-timeoutTimer.C:
			// * begin a new term and try again receiving votes.

			logrus.Debugln("candidate.timeout", n.id, timeout)
			break candidateLoop
		case <-hardBeetTimer.C:
			// * try again to receive votes from followers, cause maybe not every follower
			// received the packet or responded yet.

			logrus.WithFields(logrus.Fields{
				"prefix": "out",
			}).Debugln("candidate.o.heartbeat", n.id)
			n.sendRequestVote()
			break
		}
	}

	timeoutTimer.Stop()
	hardBeetTimer.Stop()
}
