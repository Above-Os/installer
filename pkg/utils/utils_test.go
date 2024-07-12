package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	var a = time.Now().Format("2006-01-02T15:04:05MST")
	fmt.Println(a)
}
