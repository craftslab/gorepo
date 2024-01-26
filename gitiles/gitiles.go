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
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	urlConcat = "/+/"
	urlFormat = "format=JSON"
	urlHeads  = "refs/heads/"
	urlLog    = "/+log/"
	urlSearch = "/?s="
	urlTags   = "refs/tags/"
)

const (
	opBranch    = "branch:"
	opCommit    = "commit:"
	opDelimiter = " "
	opGroups    = 2
	opTag       = "tag:"
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

// Get
//
// Example:
//
// branch:BRANCH: https://android.googlesource.com/platform/build/soong/+/refs/heads/main?format=JSON
//
// commit:COMMIT: https://android.googlesource.com/platform/build/soong/+/42ada5cff3fca011b5a0d017955f14dc63898807?format=JSON
//
// tag:TAG: https://android.googlesource.com/platform/build/soong/+/refs/tags/android-vts-10.0_r4?format=JSON
//
// nolint: lll
func (g Gitiles) Get(project, operator string) (map[string]interface{}, error) {
	var buf map[string]interface{}
	var err error

	if project == "" || operator == "" || len(strings.Split(operator, opDelimiter)) >= opGroups {
		return nil, errors.New("parameter invalid")
	}

	if strings.HasPrefix(operator, opBranch) {
		branch := strings.TrimPrefix(operator, opBranch)
		buf, err = g.request(g.url+"/"+project+urlConcat+urlHeads+branch+"?"+urlFormat, g.user, g.pass)
	} else if strings.HasPrefix(operator, opCommit) {
		commit := strings.TrimPrefix(operator, opCommit)
		buf, err = g.request(g.url+"/"+project+urlConcat+commit+"?"+urlFormat, g.user, g.pass)
	} else if strings.HasPrefix(operator, opTag) {
		tag := strings.TrimPrefix(operator, opTag)
		buf, err = g.request(g.url+"/"+project+urlConcat+urlTags+tag+"?"+urlFormat, g.user, g.pass)
	} else {
		err = errors.New("operator invalid")
	}

	if err != nil {
		return nil, err
	}

	return buf, nil
}

// Query
//
// Example:
//
// branch:BRANCH: https://android.googlesource.com/platform/build/soong/+log/refs/heads/main?format=JSON
//
// branch:BRANCH commit:COMMIT: https://android.googlesource.com/platform/build/soong/+log/refs/heads/main/?s=42ada5cff3fca011b5a0d017955f14dc63898807&format=JSON
//
// tag:TAG: https://android.googlesource.com/platform/build/soong/+log/refs/tags/android-vts-10.0_r4?format=JSON
//
// tag:TAG commit:COMMIT: https://android.googlesource.com/platform/build/soong/+log/refs/tags/android-vts-10.0_r4/?s=9863d53618714a36c3f254d949497a7eb2d11863&format=JSON
//
// nolint: gocyclo,lll
func (g Gitiles) Query(project, operator string) (map[string]interface{}, error) {
	parser := func(op string) (string, string, string, error) {
		var branch, commit, tag string

		buf := strings.Split(op, opDelimiter)
		if len(buf) > opGroups {
			return "", "", "", errors.New("operator invalid")
		}

		for _, val := range buf {
			if strings.HasPrefix(val, opBranch) {
				branch = strings.TrimPrefix(val, opBranch)
			} else if strings.HasPrefix(val, opCommit) {
				commit = strings.TrimPrefix(val, opCommit)
			} else if strings.HasPrefix(val, opTag) {
				tag = strings.TrimPrefix(val, opTag)
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
			buf, err = g.request(g.url+"/"+project+urlLog+urlHeads+branch+urlSearch+commit+"&"+urlFormat, g.user, g.pass)
		} else {
			buf, err = g.request(g.url+"/"+project+urlLog+urlHeads+branch+"?"+urlFormat, g.user, g.pass)
		}
	} else if tag != "" {
		if commit != "" {
			buf, err = g.request(g.url+"/"+project+urlLog+urlTags+tag+urlSearch+commit+"&"+urlFormat, g.user, g.pass)
		} else {
			buf, err = g.request(g.url+"/"+project+urlLog+urlTags+tag+"?"+urlFormat, g.user, g.pass)
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

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read failed")
	}

	body = []byte(strings.ReplaceAll(string(body), ")]}'", ""))

	if err := json.Unmarshal(body, &buf); err != nil {
		return nil, errors.New("unmarshal failed")
	}

	return buf, nil
}
