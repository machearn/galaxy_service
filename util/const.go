package util

import "time"

var Plan map[int32]time.Duration = map[int32]time.Duration{
	1: time.Hour * 24 * 7,
	2: time.Hour * 24 * 30,
	3: time.Hour * 24 * 365,
}
