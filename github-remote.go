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
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GithubRemote implements Remote
type GithubRemote struct {
	Config       *Config
	GithubClient *github.Client
	Repo         *github.Repository
	oauthclient  *http.Client
	ctx          context.Context
}

// CreateRepo creates a remote repository
func (g *GithubRemote) CreateRepo() error {

	var err error

	// Set some options for creation
	g.Repo.Name = &g.Config.repoName
	g.Repo.Private = &g.Config.private

	// Create repo
	g.Repo, _, err = g.GithubClient.Repositories.Create(g.ctx, "", g.Repo)
	if err != nil {
		return err
	}

	// We wait 1 second, just to be sure the repo was created
	time.Sleep(time.Second * 1)

	// Create a basic README
	opt := new(github.RepositoryContentFileOptions)
	opt.Content = []byte(fmt.Sprintf("# %s", g.Repo.GetName()))
	opt.Message = func(s string) *string { return &s }("Added a README")

	_, _, err = g.GithubClient.Repositories.CreateFile(g.ctx, g.Config.Provider["github"].User, g.Config.repoName, "README.md", opt)
	if err != nil {
		return err
	}

	fmt.Printf("Repository created at %s: %s\n", g.Repo.GetCreatedAt().Format(time.RFC3339), g.Repo.GetHTMLURL())

	return nil
}

// CloneRepo clones the remote repository
func (g *GithubRemote) CloneRepo() error {

	fmt.Printf("Cloning %s\n", g.Repo.GetURL())

	var err error
	var endpoint *transport.Endpoint

	// Define a git endpoint
	switch g.Config.Provider["github"].CloneProtocol {
	case "ssh", "":
		endpoint, err = transport.NewEndpoint(g.Repo.GetSSHURL())
	case "http":
		endpoint, err = transport.NewEndpoint(g.Repo.GetHTMLURL())
		endpoint.User = g.Config.Provider["github"].User
		endpoint.Password = g.Config.Provider["github"].Password
	default:
		err = fmt.Errorf("Unknown clone protocol %s", g.Config.Provider["github"].CloneProtocol)
	}
	if err != nil {
		log.Fatalf("Error creating endpoint: %s\n", err)
		return err
	}

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

// DeleteRepo deletes a (remote) repository
func (g *GithubRemote) DeleteRepo() error {

	_, err := g.GithubClient.Repositories.Delete(g.ctx, g.Config.Provider["github"].User, g.Config.repoName)
	if err != nil {
		return err
	}

	return nil
}

// ListRepos lists all repos for a given GithubClient
func (g *GithubRemote) ListRepos() error {

	opt := new(github.RepositoryListOptions)
	opt.PerPage = 1000
	opt.Type = g.Config.Provider["github"].User
	opt.Sort = "updated"

	repositories, _, err := g.GithubClient.Repositories.List(g.ctx, g.Config.Provider["github"].User, opt)
	if err != nil {
		return err
	}

	fmt.Println("Repolist:")

	if g.Config.listLong {

		switch g.Config.Provider["github"].CloneProtocol {
		case "ssh":
			for _, r := range repositories {
				fmt.Printf("%s - %-36s %s\n", r.GetUpdatedAt().Format(time.RFC3339), r.GetName(), r.GetSSHURL())
			}
		case "http":
			for _, r := range repositories {
				fmt.Printf("%s - %-36s %s\n", r.GetUpdatedAt().Format(time.RFC3339), r.GetName(), r.GetHTMLURL())
			}
		default:
			return fmt.Errorf("Unknown cloning protocol: %s", g.Config.Provider["github"].CloneProtocol)
		}

	} else {
		for _, r := range repositories {
			fmt.Printf("%s - %s\n", r.GetUpdatedAt().Format(time.RFC3339), r.GetName())
		}
	}

	return nil
}

// NewGithubRemote creates a new Remote object and returns it
func NewGithubRemote(c *Config) (r *GithubRemote) {

	remote := new(GithubRemote)

	remote.Config = c
	// Create an oauth client
	remote.ctx = context.Background()
	token := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: remote.Config.Provider["github"].Token})
	remote.oauthclient = oauth2.NewClient(remote.ctx, token)
	// Create Github Client
	remote.GithubClient = github.NewClient(remote.oauthclient)
	// Create a github repo object
	remote.Repo = new(github.Repository)

	return remote
}
