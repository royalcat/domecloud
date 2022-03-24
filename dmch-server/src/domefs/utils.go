package domefs

import "time"

// MAYBE тут может быть умная функция
func getTimestamps(duration time.Duration) []time.Duration {
	stamps := make([]time.Duration, 0, 1)
	stampDuration := duration / 10
	for i := 0; i < int(duration); i += int(stampDuration) {
		stamps = append(stamps, time.Duration(i))
	}
	return stamps
}
