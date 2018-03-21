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
	"fmt"
	"log"
	"os"
	"time"

	"code.gitea.io/sdk/gitea"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	git "gopkg.in/src-d/go-git.v4"
)

// A Gitea Remote object
type GiteaRemote struct {
	Config      *Config
	GiteaClient *gitea.Client
	Repo        *gitea.Repository
}

// Function CreateRepo creates a remote repository
func (g *GiteaRemote) CreateRepo() error {

	var err error

	// Create an empty new Repo
	opts := gitea.CreateRepoOption{
		Name:     g.Config.repoName,
		Private:  g.Config.private,
		Readme:   "Default",
		AutoInit: true,
	}
	g.Repo, err = g.GiteaClient.CreateRepo(opts)
	if err != nil {
		return err
	}
	fmt.Printf("Repository created at %s: %s\n", g.Repo.Created.Format(time.RFC3339), g.Repo.CloneURL)

	// We wait 1 second, just to be sure the repo was created
	time.Sleep(time.Second * 1)

	return nil
}

// Function CloneRepo clones the remote repository
func (g *GiteaRemote) CloneRepo() error {

	fmt.Printf("Cloning %s\n", g.Repo.CloneURL)

	// Define a git endpoint
	endpoint, err := transport.NewEndpoint(g.Repo.CloneURL)
	if err != nil {
		log.Fatalf("Error creating endpoint: %s\n", err)
		return err
	}
	endpoint.User = g.Config.Provider["gitea"].User
	endpoint.Password = g.Config.Provider["gitea"].Password

	// Clone the repository
	_, err = git.PlainClone(g.Config.localdir, false, &git.CloneOptions{
		URL:               endpoint.String(),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	})

	if err != nil {
		return err
	}

	return nil
}

// Function DeleteRepo deletes a (remote) repository
func (g *GiteaRemote) DeleteRepo() error {

	err := g.GiteaClient.DeleteRepo(g.Config.Provider["gitea"].User, g.Config.repoName)
	if err != nil {
		return err
	}

	return nil
}

// Function ListRepos lists all repos for a given GiteaCLient
func (g *GiteaRemote) ListRepos() error {

	repos, err := g.GiteaClient.ListMyRepos()
	if err != nil {
		return err
	}

	fmt.Println("Repolist:")

	if g.Config.listLong {
		for _, r := range repos {
			fmt.Printf("%s - %-36s %s\n", r.Updated.Format(time.RFC3339), r.Name, r.HTMLURL)
		}

	} else {
		for _, r := range repos {
			fmt.Printf("%s - %s\n", r.Updated.Format(time.RFC3339), r.Name)
		}

	}

	return nil
}

// Function NewGiteaRemote creates a new Remote object and returns it
func NewGiteaRemote(c *Config) (r *GiteaRemote) {

	remote := new(GiteaRemote)

	remote.Config = c
	remote.GiteaClient = gitea.NewClient(remote.Config.Provider["gitea"].HostBaseURL, remote.Config.Provider["gitea"].Token)
	remote.Repo = new(gitea.Repository)

	return remote
}
