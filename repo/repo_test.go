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
	if err := r.Init(&i, &g); err == nil {
		t.Error("FAIL")
	}

	i.Depth = 1
	i.TagSince = ""
	i.TimeSince = ""
	if err := r.Init(&i, &g); err != nil {
		t.Error("FAIL")
	}

	_ = os.RemoveAll(".repo")
}

func TestCheck(t *testing.T) {
	r := Repo{}

	if err := r.Check(); err != nil {
		t.Error("FAIL")
	}
}

func TestDepthAfterTag(t *testing.T) {
	// TODO
}

func TestShallowAfterTag(t *testing.T) {
	// TODO
}

func TestDepthAfterTime(t *testing.T) {
	c := config.Gitiles{
		Pass: "",
		Url:  "https://android.googlesource.com",
		User: "",
	}

	r := Repo{}

	if _, err := r.DepthAfterTime("platform/build/soong", "master", "2020-06-25T00:00:00", &c); err != nil {
		t.Error("FAIL")
	}
}

func TestShallowAfterTime(t *testing.T) {
	c := config.Gitiles{
		Pass: "",
		Url:  "https://android.googlesource.com",
		User: "",
	}

	r := Repo{}

	if err := r.ShallowAfterTime("../test/manifest-1.xml", "2020-06-25T00:00:00", &c); err != nil {
		t.Error("FAIL")
	}
}
