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
)

func TestInit(t *testing.T) {
	g := Gitiles{}

	if err := g.Init("https://android.googlesource.com", "", ""); err != nil {
		t.Error("FAIL")
	}
}

func TestGet(t *testing.T) {
	g := Gitiles{}

	if err := g.Init("https://android.googlesource.com", "", ""); err != nil {
		t.Error("FAIL")
	}

	if _, err := g.Get("platform/build/soong", "branch:master"); err != nil {
		t.Error("FAIL")
	}

	if _, err := g.Get("platform/build/soong", "commit:42ada5cff3fca011b5a0d017955f14dc63898807"); err != nil {
		t.Error("FAIL")
	}

	if _, err := g.Get("platform/build/soong", "tag:android-vts-10.0_r4"); err == nil {
		t.Error("FAIL")
	}

	if _, err := g.Get("platform/build/soong", "branch:master tag:android-vts-10.0_r4"); err == nil {
		t.Error("FAIL")
	}
}

func TestQuery(t *testing.T) {
	g := Gitiles{}

	if err := g.Init("https://android.googlesource.com", "", ""); err != nil {
		t.Error("FAIL")
	}

	if _, err := g.Query("platform/build/soong", "branch:master"); err != nil {
		t.Error("FAIL")
	}

	if _, err := g.Query("platform/build/soong", "branch:master commit:42ada5cff3fca011b5a0d017955f14dc63898807"); err != nil {
		t.Error("FAIL")
	}

	if _, err := g.Query("platform/build/soong", "tag:android-vts-10.0_r4"); err != nil {
		t.Error("FAIL")
	}

	if _, err := g.Query("platform/build/soong", "tag:android-vts-10.0_r4 commit:9863d53618714a36c3f254d949497a7eb2d11863"); err != nil {
		t.Error("FAIL")
	}

	if _, err := g.Query("platform/build/soong", "branch:master tag:android-vts-10.0_r4"); err == nil {
		t.Error("FAIL")
	}

	if _, err := g.Query("platform/build/soong",
		"branch:master commit:42ada5cff3fca011b5a0d017955f14dc63898807 tag:android-vts-10.0_r4"); err == nil {
		t.Error("FAIL")
	}
}

func TestRequest(t *testing.T) {
	g := Gitiles{}

	if _, err := g.request("https://android.googlesource.com/platform/build/soong/+/refs/heads/master?format=JSON", "", ""); err != nil {
		t.Error("FAIL")
	}
}
