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

package cmd

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"gorepo/config"
	"gorepo/repo"
)

var (
	c = config.Config{}
	r = repo.Repo{}
)

func Run() {
	app := kingpin.New("gorepo", "Go Repo").Author(Author).Version(Version)

	repoInit := app.Command("init", "Initialize repo in the current directory").Action(initAction)
	repoInit.Flag("manifest-branch", "manifest branch or revision").Short('b').Default("master").
		StringVar(&c.Init.ManifestBranch)
	repoInit.Flag("manifest-name", "initial manifest file").Short('m').Default("default.xml").
		StringVar(&c.Init.ManifestName)
	repoInit.Flag("manifest-url", "manifest repository location").Short('u').Required().
		StringVar(&c.Init.ManifestUrl)
	repoInit.Flag("depth", "create a shallow clone with a history in the specific depth").
		IntVar(&c.Init.Depth)
	repoInit.Flag("repo-url", "repo repository location").Default("https://gerrit.googlesource.com/git-repo.git").
		StringVar(&c.Init.RepoUrl)
	repoInit.Flag("tag-since", "create a shallow clone with a history after the specific tag").
		StringVar(&c.Init.TagSince)
	repoInit.Flag("time-since", "create a shallow clone with a historoy after the specific time (format: yyyy-MM-ddTHH:mm:ss)").
		StringVar(&c.Init.TimeSince)
	repoInit.Flag("gitiles-pass", "gitiles password").Default("pass").
		StringVar(&c.Gitiles.Pass)
	repoInit.Flag("gitiles-url", "gitiles location").Default("localhost:80").
		StringVar(&c.Gitiles.Url)
	repoInit.Flag("gitiles-user", "gitiles user").Default("user").
		StringVar(&c.Gitiles.User)

	repoSync := app.Command("sync", "Update working tree to the latest revision").Action(syncAction)
	repoSync.Flag("jobs", "projects to fetch simultaneously").Short('j').Default("1").
		IntVar(&c.Sync.Jobs)
	repoSync.Flag("verbose", "show all sync output").Short('v').Default("false").
		BoolVar(&c.Sync.Verbose)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func initAction(_ *kingpin.ParseContext) error {
	if err := r.Check(); err != nil {
		return err
	}

	return r.Init(&c.Init, &c.Gitiles)
}

func syncAction(_ *kingpin.ParseContext) error {
	if err := r.Check(); err != nil {
		return err
	}

	return r.Sync(&c.Sync)
}
