/*
Copyright 2018 Wilhelm Peter Püschel

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"
	"os"
)

// version will be set during build (-ldflags)
var version string

func main() {

	// Vars we need here
	var remote Remote
	var config *Config
	var err error
	var conftypes = []string{"gitea", "github", "gitlab"}

	// Conftype is allways the last command line parameter
	conftype := os.Args[len(os.Args)-1]

	// Get a config
	config, err = NewConfig()
	if err != nil {
		log.Fatalf("Could not read config: %s", err)
	}

	switch conftype {

	case "gitea":
		// Set remote
		remote = NewGiteaRemote(config)

	case "gitlab":
		// Set remote
		remote = NewGitlabRemote(config)

	case "github":
		// Set remote
		remote = NewGithubRemote(config)

	case "version":
		log.Printf("Version: %s\n", version)
		os.Exit(0)

	default:
		log.Fatalf("Unknown config type: %s\nTry -h [configtype] where configtype can be one of: %s", conftype, conftypes)
	}

	// List repos
	if config.list || config.listLong {
		err = remote.ListRepos()
		if err != nil {
			log.Fatalf("Could not list repos: %s\n", err)
		}
	}

	// Create a remote repo
	if config.repoName != "" && !config.del {
		err := remote.CreateRepo()
		if err != nil {
			log.Fatalf("Could not create repository %s: %s\n", config.repoName, err)
		}
		// Clone the remote repo
		if config.newrepo {
			err := remote.CloneRepo()
			if err != nil {
				log.Fatalf("Could not clone the remote repository %s: %s\n", config.repoName, err)
			}
		}
	}

	// Delete a remote repo
	if config.repoName != "" && config.del {
		err := remote.DeleteRepo()
		if err != nil {
			log.Fatalf("Could not delete repository %s: %s\n", config.repoName, err)
		}
	}

}
