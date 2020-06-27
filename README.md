# gorepo

[![Build Status](https://travis-ci.com/craftslab/gorepo.svg?branch=master)](https://travis-ci.com/craftslab/gorepo)
[![Coverage Status](https://coveralls.io/repos/github/craftslab/gorepo/badge.svg?branch=master)](https://coveralls.io/github/craftslab/gorepo?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/craftslab/gorepo)](https://goreportcard.com/report/github.com/craftslab/gorepo)
[![License](https://img.shields.io/github/license/craftslab/gorepo.svg?color=brightgreen)](https://github.com/craftslab/gorepo/blob/master/LICENSE)
[![Tag](https://img.shields.io/github/tag/craftslab/gorepo.svg?color=brightgreen)](https://github.com/craftslab/gorepo/tags)



## Introduction

*Go Repo* is a tool to manage Git repositories using *[Git](https://github.com/git/git)*, *[Gitiles](https://gerrit.googlesource.com/gitiles)* and *[Repo](https://gerrit.googlesource.com/git-repo)*.



## Features

- Support to fetch repositories from Manifest.

- Support to fetch repositories based on depth, tag, or time.



## Prerequisites

- Git 2.26+

- Gitiles 0.3+

- Repo 2.4+



## Usage

```bash
usage: gorepo [<flags>] <command> [<args> ...]

Go Repo

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --version  Show application version.

Commands:
  help [<command>...]
    Show help.


  init --manifest-url=MANIFEST-URL [<flags>]
    Initialize repo in the current directory

    -b, --manifest-branch="master"
                                 manifest branch or revision
    -m, --manifest-name="default.xml"
                                 initial manifest file
    -u, --manifest-url=MANIFEST-URL
                                 manifest repository location
        --depth=DEPTH            create a shallow clone with a history in the
                                 specific depth
        --repo-url="https://gerrit.googlesource.com/git-repo.git"
                                 repo repository location
        --tag-since=TAG-SINCE    create a shallow clone with a history after the
                                 specific tag
        --time-since=TIME-SINCE  create a shallow clone with a historoy after
                                 the specific time (format: yyyy-MM-ddTHH:mm:ss)
        --gitiles-pass="pass"    gitiles password
        --gitiles-url="localhost:80"
                                 gitiles location
        --gitiles-user="user"    gitiles user

  sync [<flags>]
    Update working tree to the latest revision

    -j, --jobs=1   projects to fetch simultaneously
    -v, --verbose  show all sync output
```



## Run

- **Depth mode**

```bash
gorepo init -u https://android.googlesource.com/a/platform/manifest --depth=1
gorepo sync
```



- **Tag mode**

```bash
gorepo init -u https://android.googlesource.com/a/platform/manifest --tag-since=android10-release
gorepo sync
```



- **Time mode**

```bash
gorepo init -u https://android.googlesource.com/a/platform/manifest --time-since=2020-01-01T00:00:00
gorepo sync
```



## License

Project License can be found [here](LICENSE).
