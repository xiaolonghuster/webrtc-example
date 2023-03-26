package tools

import (
	"fmt"
	"testing"
)

func TestSrs(t *testing.T) {

	call := func() string {
		fmt.Println("123")
		return "222"
	}

	fmt.Printf("call is %s\n", call())
	fmt.Printf("call is %s\n", call())
}
