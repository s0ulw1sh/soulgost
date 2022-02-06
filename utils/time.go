package utils

import "time"

func ParseDateTime(bdate []byte, loc *time.Location) time.Time {
	// YYYY-MM-DD HH:MM:SS.MMMMMM

	l   := len(bdate)
	z   := [7]int{0, 0, 0, 0, 0, 0, 0}
	switch l {
	case 10, 19, 21, 22, 23, 24, 25, 26:

		// YEAR
		z[0] += int(bdate[0] - '0')
		z[0] *= 10
		z[0] += int(bdate[1] - '0')
		z[0] *= 10
		z[0] += int(bdate[2] - '0')
		z[0] *= 10
		z[0] += int(bdate[3] - '0')

		// MONTH
		z[1] += int(bdate[5] - '0')
		z[1] *= 10
		z[1] += int(bdate[6] - '0')

		// DAY
		z[2] += int(bdate[8] - '0')
		z[2] *= 10
		z[2] += int(bdate[9] - '0')

		if l == 10 {
			goto ParseDateTime_out
		}

		// HOUR
		z[3] += int(bdate[11] - '0')
		z[3] *= 10
		z[3] += int(bdate[12] - '0')

		// MIN
		z[4] += int(bdate[14] - '0')
		z[4] *= 10
		z[4] += int(bdate[15] - '0')

		// SECOND
		z[5] += int(bdate[17] - '0')
		z[5] *= 10
		z[5] += int(bdate[18] - '0')

		if l == 19 {
			goto ParseDateTime_out
		}

		// NANO-SECONDS
		for i := 20; i < l; i++ {
			if i > 0 {
				z[6] *= 10
			}

			z[6] += int(bdate[i] - '0')
		}

		z[6] *= 1000
	}

ParseDateTime_out:
	return time.Date(z[0], time.Month(z[1]), z[2], z[3], z[4], z[5], z[6], loc)
}