/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vfs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/huminghe/infini-framework/core/util"
	"io/ioutil"
	"testing"
)

func TestFiles(t *testing.T) {

	util.FilePutContent("/tmp/test_gopa.txt", "hello")

	RegisterFS(StaticFS{StaticFolder: "/", CheckLocalFirst: true})

	fs := VFS()

	f, e := fs.Open("/tmp/test_gopa.txt")

	b, e := ioutil.ReadAll(f)

	fmt.Println(e)

	fmt.Println(string(b))

	assert.Equal(t, "hello", string(b))
}
