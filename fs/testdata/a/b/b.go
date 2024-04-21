package b

import (
	"time"

	"github.com/oklog/ulid/v2"
)

func Print() string {
	return ulid.MustNew(uint64(time.Now().UnixNano()), nil).String()
}
