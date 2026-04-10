package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixMicro())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}