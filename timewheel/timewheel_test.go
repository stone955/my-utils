package timewheel

import (
	"testing"
	"time"
)

func TestTimeWheel(t *testing.T) {
	tw, err := New(10*time.Second, 60, func(itf interface{}) {

	})
	if err != nil {
		t.Fatal(err)
	}

	tw.Start()
	defer tw.Stop()
}
