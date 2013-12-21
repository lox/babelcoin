package babelcoin

import (
	"encoding/json"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var unixtime int64

	if err := json.Unmarshal(b, &unixtime); err != nil {
		return err
	}

	t.Time = time.Unix(unixtime, 0)
	return nil
}
