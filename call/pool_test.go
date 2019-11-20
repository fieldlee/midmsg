package call

import (
	"testing"
	"time"
)

func TestTimerCallPool(t *testing.T) {
	go TimerCallPool()

	time.Sleep(time.Second * 1000)
}
