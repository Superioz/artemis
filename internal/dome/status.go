package dome

type Status struct {
	BrokerConnected bool `json:"broker"`
}

func CurrentStatus() Status {
	// TODO
	return Status{}
}
