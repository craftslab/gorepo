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
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gorepo/config"
	"gorepo/gitiles"
	"gorepo/manifest"
)

const (
	Major  = 2
	Minor  = 4
	Prefix = "repo launcher version"
	Total  = 2

	SHA1 = "^[0-9a-f]{40}$"

	Time1 = "2006-01-02T15:04:05"
	Time2 = "Mon Jan 2 15:04:05 2006 -0700"
)

type Repo struct {
}

func (r Repo) Init(i *config.Init, g *config.Gitiles) error {
	if (i.Depth > 0 && i.TagSince != "") ||
		(i.Depth > 0 && i.TimeSince != "") ||
		(i.TagSince != "" && i.TimeSince != "") {
		return errors.New("config invalid")
	}

	cmd := exec.Command("repo", "init",
		"--manifest-branch="+i.ManifestBranch,
		"--manifest-name="+i.ManifestName,
		"--manifest-url="+i.ManifestUrl,
		"--repo-url="+i.RepoUrl,
		"--depth="+strconv.Itoa(i.Depth))

	if err := r.Run(cmd); err != nil {
		return errors.Wrap(err, "init failed")
	}

	cmd = exec.Command("repo", "manifest",
		"--manifest-name="+i.ManifestName,
		"--output-file=manifest.xml")

	if err := r.Run(cmd); err != nil {
		return errors.Wrap(err, "manifest failed")
	}

	cmd = exec.Command("mv", "manifest.xml", ".repo/manifest.xml")

	if err := r.Run(cmd); err != nil {
		return errors.Wrap(err, "mv failed")
	}

	if i.TagSince != "" {
		if err := r.ShallowAfterTag(".repo/manifest.xml", i.TagSince, g); err != nil {
			return errors.Wrap(err, "shallow failed")
		}
	}

	if i.TimeSince != "" {
		if err := r.ShallowAfterTime(".repo/manifest.xml", i.TimeSince, g); err != nil {
			return errors.Wrap(err, "shallow failed")
		}
	}

	return nil
}

func (r Repo) Sync(s *config.Sync) error {
	var verbose string

	if s.Verbose {
		verbose = "--verbose"
	}

	cmd := exec.Command("repo", "sync", "--jobs="+strconv.Itoa(s.Jobs), verbose)

	if err := r.Run(cmd); err != nil {
		return errors.Wrap(err, "run failed")
	}

	return nil
}

func (r Repo) Check() error {
	var version []string

	_, err := exec.LookPath("repo")
	if err != nil {
		return errors.Wrap(err, "repo not found")
	}

	cmd := exec.Command("repo", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "version not found")
	}

	buf := strings.Split(string(out), "\n")
	for _, item := range buf {
		if strings.HasPrefix(item, Prefix) {
			b := strings.Trim(item, Prefix)
			version = strings.Split(b, ".")
			break
		}
	}

	if version == nil || len(version) < Total {
		return errors.New("version invalid")
	}

	major, err := strconv.Atoi(version[0])
	if err != nil {
		return errors.Wrap(err, "major version invalid")
	}

	minor, err := strconv.Atoi(version[1])
	if err != nil {
		return errors.Wrap(err, "minor version invalid")
	}

	if major < Major || minor < Minor {
		return errors.New(strconv.Itoa(Major) + "." + strconv.Itoa(Minor) + "+ required")
	}

	return nil
}

func (r Repo) Run(c *exec.Cmd) error {
	var outBuf, errBuf bytes.Buffer
	var errOut, errErr error

	outPipe, _ := c.StdoutPipe()
	errPipe, _ := c.StderrPipe()

	outWriter := io.MultiWriter(os.Stdout, &outBuf)
	errWriter := io.MultiWriter(os.Stderr, &errBuf)

	err := c.Start()
	if err != nil {
		return errors.Wrap(err, "start failed")
	}

	go func() {
		_, errOut = io.Copy(outWriter, outPipe)
	}()

	go func() {
		_, errErr = io.Copy(errWriter, errPipe)
	}()

	err = c.Wait()
	if err != nil {
		return errors.Wrap(err, "wait failed")
	}

	if errOut != nil || errErr != nil {
		return errors.Wrap(err, "copy failed")
	}

	log.Println(outBuf.String())
	log.Println(errBuf.String())

	return nil
}

func (r Repo) DepthAfterTag(project, tag string, c *config.Gitiles) (int, error) {
	// TODO
	return 0, nil
}

func (r Repo) ShallowAfterTag(name, tag string, c *config.Gitiles) error {
	// TODO
	return nil
}

func (r Repo) DepthAfterTime(project, branch, _time string, c *config.Gitiles) (int, error) {
	g := gitiles.Gitiles{}

	if err := g.Init(c.Url, c.User, c.Pass); err != nil {
		return 0, errors.Wrap(err, "init failed")
	}

	buf, err := g.Query(project, "branch:"+branch)
	if err != nil {
		return 0, errors.Wrap(err, "query failed")
	}

	depth := 0
	t, _ := time.Parse(Time1, _time)

	// TODO
	for _, val := range buf["log"].([]interface{}) {
		c := val.(map[string]interface{})["committer"].(map[string]interface{})
		b, _ := time.Parse(Time2, c["time"].(string))
		if b.UTC().Before(t.UTC()) {
			break
		}
		depth++
	}

	return depth, nil
}

func (r Repo) ShallowAfterTime(name, _time string, c *config.Gitiles) error {
	m := manifest.Manifest{}

	if err := m.Load(name); err != nil {
		return errors.Wrap(err, "load failed")
	}

	projects, err := m.Projects()
	if err != nil {
		return errors.Wrap(err, "projects failed")
	}

	for index, val := range projects {
		d, n, _, rev, err := m.Project(val.(map[string]interface{}))
		if err != nil {
			return errors.Wrap(err, "project failed")
		}
		if match, _ := regexp.MatchString(SHA1, rev); match {
			continue
		}
		if _, err := strconv.Atoi(d); err != nil {
			if depth, err := r.DepthAfterTime(n, rev, _time, c); err == nil {
				if depth > 0 {
					projects[index].(map[string]interface{})["-clone-depth"] = strconv.Itoa(depth)
				}
			}
		}
	}

	if err := m.Update(projects); err != nil {
		return errors.Wrap(err, "update failed")
	}

	if err := m.Write(name); err != nil {
		return errors.Wrap(err, "write failed")
	}

	return nil
}
