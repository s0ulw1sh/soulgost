package db

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

func TestTimeEmptyScanString(t *testing.T) {
	test := "2038-01-19 03:14:07"
	timeval := TimeEmpty{}
	timeval.Scan(test)
	res := timeval.String()

	if res != test {
		t.Error("Not equal time", test, res)
	}
}