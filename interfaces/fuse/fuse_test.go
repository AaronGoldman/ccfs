//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
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
