package iso8601

import (
	"time"
)

type ISO8601 struct {
	time.Time
}

func (self *ISO8601) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	// Get rid of the quotes "" around the value. A second option would be to
	// include them in the date format string instead, like so below:
	//   time.Parse(`"`+time.RFC3339Nano+`"`, s)
	s = s[1 : len(s)-1]

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05.999999999", s)
	}
	self.Time = t
	return
}
