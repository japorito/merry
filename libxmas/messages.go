package xmas

import (
	"fmt"
	"time"
)

func PrintHolidayMessage(start time.Time) {
	fmt.Printf("Merry Christmas! Solutions calculated in %v", time.Since(start))
}
