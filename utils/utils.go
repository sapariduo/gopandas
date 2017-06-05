package utils

import (
	"gopandas/types"
	"log"
	"time"
)

func DateRangeUTC(start, end string, d types.Duration) []types.Time {
	ret := []types.Time{}
	duration, err := time.ParseDuration(string(d))
	if err != nil {
		log.Println("Error interval syntax is incorrect")
		return nil
	}
	tstart, err := time.Parse("2006-01-02", start)
	if err != nil {
		log.Println("Time format is 2006-01-02")
		return nil
	}
	tend, err := time.Parse("2006-01-02", end)
	if err != nil {
		log.Println("Time format is 2006-01-02")
		return nil
	}
	for tstart.Unix() <= tend.Unix() {
		ret = append(ret, types.Time(tstart))
		tstart = tstart.Add(duration)
	}
	return ret
}
