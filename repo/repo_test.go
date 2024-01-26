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

package repo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"gorepo/config"
)

func TestInit(t *testing.T) {
	i := config.Init{
		ManifestBranch: "master",
		ManifestName:   "default.xml",
		ManifestUrl:    "https://android.googlesource.com/a/platform/manifest",
		RepoUrl:        "https://gerrit.googlesource.com/git-repo.git",
	}

	g := config.Gitiles{
		Pass: "",
		Url:  "https://android.googlesource.com/",
		User: "",
	}

	r := Repo{}

	i.Depth = 1
	i.TagSince = "android10-release"
	i.TimeSince = "2020-06-25T00:00:00"
	err := r.Init(&i, &g)
	assert.NotEqual(t, nil, err)

	i.Depth = 1
	i.TagSince = ""
	i.TimeSince = ""
	err = r.Init(&i, &g)
	assert.Equal(t, nil, err)

	_ = os.RemoveAll(".repo")
}

func TestCheck(t *testing.T) {
	r := Repo{}

	err := r.Check()
	assert.Equal(t, nil, err)
}

func TestDepthAfterTag(t *testing.T) {
	// TODO: FIXME
}

func TestShallowAfterTag(t *testing.T) {
	// TODO: FIXME
}

func TestDepthAfterTime(t *testing.T) {
	c := config.Gitiles{
		Pass: "",
		Url:  "https://android.googlesource.com",
		User: "",
	}

	r := Repo{}

	_, err := r.DepthAfterTime("platform/build/soong", "master", "2020-06-25T00:00:00", &c)
	assert.Equal(t, nil, err)
}

func TestShallowAfterTime(t *testing.T) {
	c := config.Gitiles{
		Pass: "",
		Url:  "https://android.googlesource.com",
		User: "",
	}

	r := Repo{}

	err := r.ShallowAfterTime("../test/manifest-1.xml", "2020-06-25T00:00:00", &c)
	assert.Equal(t, nil, err)
}
