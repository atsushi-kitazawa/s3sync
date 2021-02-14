package util

import (
	_ "fmt"
	"io/ioutil"
	"strings"
)

func ListDir(dir string) []string {
    list := traverseDir(dir)

    // trim argument dir name
    ret := make([]string, 0)
    for _, f := range list {
	ret = append(ret, strings.TrimPrefix(f, dir))
    }

    return ret
}

func traverseDir(dir string) []string {
    ret := make([]string, 0)
    files, err := ioutil.ReadDir(dir)
    if err != nil {
	panic(err)
    }

    for _, f := range files {
	if !f.IsDir() {
	    ret = append(ret, (dir + "/" + f.Name()))
	} else {
	    ret = append(ret, traverseDir(dir + "/" + f.Name())...)
	}
    }

    return ret
}
