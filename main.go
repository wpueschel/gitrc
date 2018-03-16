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

func main() {

	// Vars we need here
	var remote Remote
	var config *Config
	var err error
	var conftypes = []string{"gitea", "github", "gitlab"}

	// Conftype is allways the last command line parameter
	conftype := os.Args[len(os.Args)-1]

	switch conftype {

	case "gitea":
		// Read/fetch a config
		config, err = NewConfig()
		if err != nil {
			log.Fatalf("Could not read config: %s", err)
		}
		// Set remote
		remote = NewGiteaRemote(config)

	case "gitlab":
		// Read/fetch a config
		config, err = NewConfig()
		if err != nil {
			log.Fatalf("Could not read config: %s", err)
		}
		// Set remote
		remote = NewGitlabRemote(config)

	case "github":
		// Read/fetch a config
		config, err = NewConfig()
		if err != nil {
			log.Fatalf("Could not read config: %s", err)
		}
		// Set remote
		remote = NewGithubRemote(config)

	default:
		log.Fatalf("Unknown config type: %s\nTry -h [configtype] where configtype can be one of: %s", conftype, conftypes)
	}

	// List repos
	if config.list {
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
