package types

import "testing"

func TestTimeMinString(t *testing.T) {
	timeval := TimeMin{1440}
	strtime := timeval.String()
	eqtime  := "24:00:00"

	if strtime != eqtime {
		t.Error("Not equal time", strtime, eqtime)
	}
}

func TestTimeMinScanString(t *testing.T) {
	timeval := TimeMin{}
	eqtime  := 732

	timeval.Scan("12:12:00")

	mintime := timeval.Val()

	if mintime != eqtime {
		t.Error("Not equal time", mintime, eqtime)
	}
}