package utils

import (
	"fmt"
	"testing"
)

func TestHash(T *testing.T) {
	ret := GetUniqueId()
	fmt.Println(ret)
}
