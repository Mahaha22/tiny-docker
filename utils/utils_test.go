package utils

import (
	"fmt"
	"testing"
)

func TestHash(T *testing.T) {
	ret := GetUniqueId()
	fmt.Println(ret)
}

func TestPort(t *testing.T) {
	err := IsPortAvilable("8888")
	fmt.Println("err = ", err)
}
