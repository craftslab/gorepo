// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manifest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	m := Manifest{}

	err := m.Load("../test/manifest-1.xml")
	assert.Equal(t, nil, err)
}

func TestProjects(t *testing.T) {
	m := Manifest{}

	err := m.Load("../test/manifest-1.xml")
	assert.Equal(t, nil, err)

	_, err = m.Projects()
	assert.Equal(t, nil, err)
}

func TestProject(t *testing.T) {
	m := Manifest{}

	err := m.Load("../test/manifest-1.xml")
	assert.Equal(t, nil, err)

	projects, err := m.Projects()
	assert.Equal(t, nil, err)

	for _, val := range projects {
		_, _, _, _, err := m.Project(val.(map[string]interface{}))
		assert.Equal(t, nil, err)
	}
}

func TestUpdate(t *testing.T) {
	m := Manifest{}

	err := m.Load("../test/manifest-1.xml")
	assert.Equal(t, nil, err)

	buf := make([]interface{}, 1)

	buf[0] = map[string]string{
		"groups": "pdk",
		"name":   "platform/build",
		"path":   "build/make",
	}

	err = m.Update(buf)
	assert.Equal(t, nil, err)
}

func TestWrite(t *testing.T) {
	m := Manifest{}

	err := m.Load("../test/manifest-1.xml")
	assert.Equal(t, nil, err)

	err = m.Write("../test/manifest-1-new.xml")
	assert.Equal(t, nil, err)

	_ = os.RemoveAll("../test/manifest-1-new.xml")
}
