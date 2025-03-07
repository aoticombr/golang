package mail

import "time"

func Now() time.Time {
	return time.Date(2014, 06, 25, 17, 46, 0, 0, time.UTC)
}
