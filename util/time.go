package babelcoin

import (
	"time"
	"encoding/json"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var unixtime int64

    if error := json.Unmarshal(b, &unixtime); error != nil {
    	return error
    }

	t.Time = time.Unix(unixtime, 0)
    return nil
}