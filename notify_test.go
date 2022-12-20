package notify

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	msg := time.Now().GoString()
	n := NewWeNotifier()
	n.Send(fmt.Sprintf("this is content %s ", msg))
}

func TestConfig(t *testing.T) {
	if corpID == "" || corpSecret == "" {
		log.Fatal("empty notify config")
	}
}
