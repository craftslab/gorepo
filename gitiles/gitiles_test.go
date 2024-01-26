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

package gitiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	g := Gitiles{}

	err := g.Init("https://android.googlesource.com", "", "")
	assert.Equal(t, nil, err)
}

func TestGet(t *testing.T) {
	g := Gitiles{}

	err := g.Init("https://android.googlesource.com", "", "")
	assert.Equal(t, nil, err)

	_, err = g.Get("platform/build/soong", "branch:master")
	assert.Equal(t, nil, err)

	_, err = g.Get("platform/build/soong", "commit:42ada5cff3fca011b5a0d017955f14dc63898807")
	assert.Equal(t, nil, err)

	_, err = g.Get("platform/build/soong", "tag:android-vts-10.0_r4")
	assert.Equal(t, nil, err)
}

func TestQuery(t *testing.T) {
	g := Gitiles{}

	err := g.Init("https://android.googlesource.com", "", "")
	assert.Equal(t, nil, err)

	_, err = g.Query("platform/build/soong", "branch:master")
	assert.Equal(t, nil, err)

	_, err = g.Query("platform/build/soong", "branch:master commit:42ada5cff3fca011b5a0d017955f14dc63898807")
	assert.Equal(t, nil, err)

	_, err = g.Query("platform/build/soong", "tag:android-vts-10.0_r4")
	assert.Equal(t, nil, err)

	_, err = g.Query("platform/build/soong", "tag:android-vts-10.0_r4 commit:9863d53618714a36c3f254d949497a7eb2d11863")
	assert.Equal(t, nil, err)
}

func TestRequest(t *testing.T) {
	g := Gitiles{}

	_, err := g.request("https://android.googlesource.com/platform/build/soong/+/refs/heads/master?format=JSON", "", "")
	assert.Equal(t, nil, err)
}
