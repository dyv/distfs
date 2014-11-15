package paxos

type MsgType uint

const (
	Empty MsgType = iota
	Prepare
	Promise
	Nack
	AcceptRequest
	Accepted
	ClientRequest
	ClientResponse
	ClientRedirect
	ClientConn
	Heartbeat
	Error
)

func (m MsgType) String() string {
	switch m {
	case Empty:
		return "Empty"
	case Prepare:
		return "Prepare"
	case Promise:
		return "Promise"
	case Nack:
		return "Nack"
	case AcceptRequest:
		return "AcceptRequest"
	case Accepted:
		return "Accepted"
	case ClientRequest:
		return "ClientRequest"
	case ClientResponse:
		return "ClientResponse"
	case ClientRedirect:
		return "Redirect"
	case Heartbeat:
		return "Heartbeat"
	case Error:
		return "Error"
	}
	return "INVALID"
}

type RequestInfo struct {
	Id  int   `json:"id"`
	No  int   `json:"no"`
	Val Value `json:"val"`
}

type Msg struct {
	Type          MsgType `json:"type"`
	FromAddress   string  `json:"fromaddress"`
	FromPort      string  `json:"fromport"`
	LeaderAddress string  `json:"leaderAddress"`
	LeaderPort    string  `json:"leaderPort"`
	Request       RequestInfo
	Entry         int        `json:"entry"` // what entry in the log this is meant for
	Round         int        `json:"round"`
	Value         Value      `json:"value"`
	RoundValue    RoundValue `json:"roundvalue"` // For Previous Value (is this necessary)
	Error         string     `json:"error"`
}
