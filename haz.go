package msg

import (
	"encoding/json"
)

// Haz is a useful wire format.  Clients will typically expect only one
// of the members to be non nil.
// Implements MessageTx
type Haz struct {
	Quake         *Quake
	HeartBeat     *HeartBeat
	err           error
	receiptHandle string
}

// HazDecode decodes JSON and returns MessageTx with a concrete type Haz.
func HazDecode(b []byte, receiptHandle string) (MessageTx, error) {
	n := Haz{receiptHandle: receiptHandle}
	err := json.Unmarshal(b, &n)
	return n, err
}

func (h Haz) ReceiptHandle() string {
	return h.receiptHandle
}

// Err returns the first non nil error of h, h.Quake, h.HeartBeat otherwise nil.
func (h Haz) Err() error {
	if h.err != nil {
		return h.err
	}

	if h.Quake != nil && h.Quake.err != nil {
		return h.Quake.err
	}

	if h.HeartBeat != nil && h.HeartBeat.err != nil {
		return h.HeartBeat.err
	}

	return nil
}

func (h *Haz) SetErr(err error) {
	h.err = err
}

// Encode encodes Haz as JSON.
func (h Haz) Encode() ([]byte, error) {
	if h.Err() != nil {
		return nil, h.Err()
	}

	return json.Marshal(h)
}
