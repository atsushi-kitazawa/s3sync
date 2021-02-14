package util

import (
	"fmt"
	"testing"
)

func TestListDir(t *testing.T) {
    ret := ListDir("../../testdata")
    for _, s := range ret {
	fmt.Println(s)
    }
}
