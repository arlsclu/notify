package notify_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/arlsclu/notify"
)

func TestSend(t *testing.T) {
	msg := time.Now().GoString()
	n := notify.NewWeNotifier()
	n.Send(fmt.Sprintf("this is content %s ", msg))
}
