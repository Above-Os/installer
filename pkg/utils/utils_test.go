package utils

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	var a, _ = GeneratePassword(15)
	fmt.Println(a)
}
