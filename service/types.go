package service

import (
	"encoding/json"
	"errors"
	"time"
)

type EveryType int

const (
	EveryDay    EveryType = 1
	EveryHour   EveryType = 2
	EveryMinute EveryType = 3
	EverySecond EveryType = 4
)

type SimpleConfig struct {
	Et EveryType `json:"et"`
	At int       `json:"at"`
}

func ConvertSimpleToConfig(s *SimpleConfig) string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func FromConfigToSimple(c string) *SimpleConfig {
	s := new(SimpleConfig)
	if err := json.Unmarshal([]byte(c), s); err != nil {
		panic(err)
	}
	return s
}

func CalculateNextTimeSimple(s *SimpleConfig) int {
	now := time.Now()
	switch s.Et {
	case EveryDay:
		adj := 1
		if now.Hour() < s.At {
			adj = 0
		}
		next := time.Date(now.Year(), now.Month(), now.Day(), s.At, 0, 0, 0, now.Location()).AddDate(0, 0, adj)
		return int(next.Unix())
	case EveryHour:
		adj := time.Hour
		if now.Minute() < s.At {
			adj = 0
		}
		next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), s.At, 0, 0, now.Location()).Add(adj)
		return int(next.Unix())
	case EveryMinute:
		adj := time.Minute
		if now.Second() < s.At {
			adj = 0
		}
		next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), s.At, 0, now.Location()).Add(adj)
		return int(next.Unix())
	case EverySecond:
		return int(now.Unix() + 1)
	default:
		panic(errors.New("unknown type"))
	}
}
