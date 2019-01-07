package protocol

import "encoding/json"

func (m *LogEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Index   uint64 `json:"index"`
		Term    uint64 `json:"term"`
		Content []byte `json:"content"`
	}{
		Index:   m.Index,
		Term:    m.Term,
		Content: m.Content,
	})
}

func (m *LogEntry) UnmarshalJSON(data []byte) error {
	s := &struct {
		Index   uint64 `json:"index"`
		Term    uint64 `json:"term"`
		Content []byte `json:"content"`
	}{}
	err := json.Unmarshal(data, s)
	if err != nil {
		return err
	}
	m.Index = s.Index
	m.Term = s.Term
	m.Content = s.Content
	return nil
}
