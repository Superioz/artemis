package protocol

import (
	"github.com/golang/protobuf/proto"
	"reflect"
)

var registry = make(map[uint16]proto.Message)

func init() {
	registry[0x01] = &RequestVoteCall{}
	registry[0x02] = &RequestVoteRespond{}
	registry[0x03] = &AppendEntriesCall{}
	registry[0x04] = &AppendEntriesRespond{}
}

// returns the respective proto message of
// given `id`
func FromId(id uint16) proto.Message {
	reg := registry[id]
	if reg == nil {
		return nil
	}

	val := reflect.ValueOf(registry[id])
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	v := reflect.New(val.Type()).Interface().(proto.Message)
	return v
}

// gets the respective packet id of given message
func ToId(message proto.Message) uint16 {
	var packetId uint16
	for id, p := range registry {
		if reflect.TypeOf(p) == reflect.TypeOf(message) {
			packetId = id
		}
	}
	return packetId
}
