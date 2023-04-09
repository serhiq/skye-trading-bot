package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOrderNumber() string {
	now := time.Now()
	dayOfYear := now.YearDay()
	rand.Seed(now.UnixNano())
	randNum := rand.Intn(10000)
	orderNum := fmt.Sprintf("%03d%04d", dayOfYear, randNum)
	return orderNum
}
