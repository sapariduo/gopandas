package utils

import (
	"gopandas/indices"
	"gopandas/types"
	"log"
	"time"
)

type DateRange []types.Time

func DateRangeUTC(start, end string, d types.Duration) DateRange {
	ret := DateRange{}
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

func (d DateRange) ToIndex() indices.Indices {
	ret := make(indices.Indices, len(d))

	for i, t := range d {
		ret[i] = indices.Index(t)
	}
	return ret
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
