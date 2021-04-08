package webhooks

import "encoding/json"

type OpType string

const (
	OpTest    OpType = "test"
	OpRemove  OpType = "remove"
	OpAdd     OpType = "add"
	OpReplace OpType = "replace"
	OpMove    OpType = "move"
	OpCopy    OpType = "copy"
)

type PatchStruct struct {
	Op    OpType      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func PatchData(data []PatchStruct) []byte {
	playLoadBytes, _ := json.Marshal(data)
	return playLoadBytes
}
