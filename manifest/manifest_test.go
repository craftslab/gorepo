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
)

func TestLoad(t *testing.T) {
	m := Manifest{}

	if err := m.Load("../test/manifest-1.xml"); err != nil {
		t.Error("FAIL")
	}
}

func TestProjects(t *testing.T) {
	m := Manifest{}

	if err := m.Load("../test/manifest-1.xml"); err != nil {
		t.Error("FAIL")
	}

	if _, err := m.Projects(); err != nil {
		t.Error("FAIL")
	}
}

func TestProject(t *testing.T) {
	m := Manifest{}

	if err := m.Load("../test/manifest-1.xml"); err != nil {
		t.Error("FAIL")
	}

	projects, err := m.Projects()
	if err != nil {
		t.Error("FAIL")
	}

	for _, val := range projects {
		if _, _, _, _, err := m.Project(val.(map[string]interface{})); err != nil {
			t.Error("FAIL")
		}
	}
}

func TestUpdate(t *testing.T) {
	m := Manifest{}

	if err := m.Load("../test/manifest-1.xml"); err != nil {
		t.Error("FAIL")
	}

	buf := make([]interface{}, 1)

	buf[0] = map[string]string{
		"groups": "pdk",
		"name":   "platform/build",
		"path":   "build/make",
	}

	if err := m.Update(buf); err != nil {
		t.Error("FAIL")
	}
}

func TestWrite(t *testing.T) {
	m := Manifest{}

	if err := m.Load("../test/manifest-1.xml"); err != nil {
		t.Error("FAIL")
	}

	if err := m.Write("../test/manifest-1-new.xml"); err != nil {
		t.Error("FAIL")
	}

	_ = os.RemoveAll("../test/manifest-1-new.xml")
}
