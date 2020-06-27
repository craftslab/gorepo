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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	OP_BRANCH = "branch:"
	OP_COMMIT = "commit:"
	OP_TAG    = "tag:"

	OP_DELIMITER = " "
	OP_GROUPS    = 2
)

type Gitiles struct {
	pass string
	url  string
	user string
}

func (g *Gitiles) Init(url, user, pass string) error {
	g.url = url
	g.user = user
	g.pass = pass

	return nil
}

// Example:
// branch:BRANCH: https://android.googlesource.com/platform/build/soong/+/refs/heads/master?format=JSON
// commit:COMMIT: https://android.googlesource.com/platform/build/soong/+/42ada5cff3fca011b5a0d017955f14dc63898807?format=JSON
//       tag:TAG: https://android.googlesource.com/platform/build/soong/+/refs/tags/android-vts-10.0_r4
func (g Gitiles) Get(project, operator string) (map[string]interface{}, error) {
	var buf map[string]interface{}
	var err error

	if project == "" || operator == "" || len(strings.Split(operator, OP_DELIMITER)) >= OP_GROUPS {
		return nil, errors.New("parameter invalid")
	}

	if strings.HasPrefix(operator, OP_BRANCH) {
		branch := strings.TrimPrefix(operator, OP_BRANCH)
		buf, err = g.request(g.url+"/"+project+"/+/"+"refs/heads/"+branch+"?format=JSON", g.user, g.pass)
	} else if strings.HasPrefix(operator, OP_COMMIT) {
		commit := strings.TrimPrefix(operator, OP_COMMIT)
		buf, err = g.request(g.url+"/"+project+"/+/"+commit+"?format=JSON", g.user, g.pass)
	} else {
		err = errors.New("operator invalid")
	}

	if err != nil {
		return nil, err
	}

	return buf, nil
}

// Example:
//               branch:BRANCH: https://android.googlesource.com/platform/build/soong/+log/refs/heads/master?format=JSON
// branch:BRANCH commit:COMMIT: https://android.googlesource.com/platform/build/soong/+log/refs/heads/master/?s=42ada5cff3fca011b5a0d017955f14dc63898807&format=JSON
//                     tag:TAG: https://android.googlesource.com/platform/build/soong/+log/refs/tags/android-vts-10.0_r4?format=JSON
//       tag:TAG commit:COMMIT: https://android.googlesource.com/platform/build/soong/+log/refs/tags/android-vts-10.0_r4/?s=9863d53618714a36c3f254d949497a7eb2d11863&format=JSON
func (g Gitiles) Query(project, operator string) (map[string]interface{}, error) {
	parser := func(op string) (string, string, string, error) {
		var branch, commit, tag string

		buf := strings.Split(op, OP_DELIMITER)
		if len(buf) > OP_GROUPS {
			return "", "", "", errors.New("operator invalid")
		}

		for _, val := range buf {
			if strings.HasPrefix(val, OP_BRANCH) {
				branch = strings.TrimPrefix(val, OP_BRANCH)
			} else if strings.HasPrefix(val, OP_COMMIT) {
				commit = strings.TrimPrefix(val, OP_COMMIT)
			} else if strings.HasPrefix(val, OP_TAG) {
				tag = strings.TrimPrefix(val, OP_TAG)
			} else {
				continue
			}
		}

		if branch != "" && tag != "" {
			return "", "", "", errors.New("operator invalid")
		}

		if len(buf) == 1 && commit != "" {
			return "", "", "", errors.New("operator invalid")
		}

		return branch, commit, tag, nil
	}

	var buf map[string]interface{}
	var err error

	if project == "" || operator == "" {
		return nil, errors.New("parameter invalid")
	}

	branch, commit, tag, err := parser(operator)
	if err != nil {
		return nil, err
	}

	if branch != "" {
		if commit != "" {
			buf, err = g.request(g.url+"/"+project+"/+log/"+"refs/heads/"+branch+"/?s="+commit+"&format=JSON", g.user, g.pass)
		} else {
			buf, err = g.request(g.url+"/"+project+"/+log/"+"refs/heads/"+branch+"?format=JSON", g.user, g.pass)
		}
	} else if tag != "" {
		if commit != "" {
			buf, err = g.request(g.url+"/"+project+"/+log/"+"refs/tags/"+tag+"/?s="+commit+"&format=JSON", g.user, g.pass)
		} else {
			buf, err = g.request(g.url+"/"+project+"/+log/"+"refs/tags/"+tag+"?format=JSON", g.user, g.pass)
		}
	} else {
		err = errors.New("operator invalid")
	}

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (g Gitiles) request(url, user, pass string) (map[string]interface{}, error) {
	var buf map[string]interface{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client failed")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("client failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read failed")
	}

	body = []byte(strings.ReplaceAll(string(body), ")]}'", ""))

	if err := json.Unmarshal(body, &buf); err != nil {
		return nil, errors.New("unmarshal failed")
	}

	return buf, nil
}
