package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	var retry = 5
	var delay = 2

	if retry < 0 {
		retry = 1
	}

	for i := 0; i < retry; i++ {
		fmt.Println("hahaha")

		if i+1 < retry {
			func() {
				fmt.Println("delay")
				time.Sleep(time.Duration(delay) * time.Second)
			}()
		}
	}

	fmt.Println("done")
}
