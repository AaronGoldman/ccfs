package fuse

import (
	"io/ioutil"
	"testing"
)

const mountpiont = "../../mountpiont"

func TestPwd(t *testing.T) {
	fileInfos, _ := ioutil.ReadDir(mountpiont)
	t.Logf("pwd: %s", fileInfos)
}
