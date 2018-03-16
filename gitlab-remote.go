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
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	gitlab "github.com/xanzy/go-gitlab"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

// A Gitlab Repository object
type GitlabRemote struct {
	Config       *Config
	GitlabClient *gitlab.Client
	Repo         *gitlab.Project
}

// Function CreateRepo creates a remote and a local clone of an repository
// What is being created depends on  the config
func (g *GitlabRemote) CreateRepo() error {

	var nsid int
	var err error

	projectVisibility := gitlab.Visibility(gitlab.InternalVisibility)
	if g.Config.private {
		projectVisibility = gitlab.Visibility(gitlab.PrivateVisibility)
	}

	// We need to fetch the namespace id from our group name
	nopts := new(gitlab.ListNamespacesOptions)
	namepspaces, _, err := g.GitlabClient.Namespaces.ListNamespaces(nopts)
	if err != nil {
		return err
	}
	for _, n := range namepspaces {
		if n.Name == g.Config.Provider["gitlab"].GroupName {
			nsid = n.ID
		}
	}
	if nsid == 0 {
		return errors.New(fmt.Sprintf("Could not find namespace id for group %s", g.Config.Provider["gitlab"].GroupName))
	}

	// We create a new repository
	popts := new(gitlab.CreateProjectOptions)
	popts.Name = &g.Config.repoName
	popts.Visibility = projectVisibility
	popts.NamespaceID = &nsid
	g.Repo, _, err = g.GitlabClient.Projects.CreateProject(popts)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// We wait 1 second, just to be sure the repo was created
	time.Sleep(time.Second * 1)

	// Create a basic README.md
	readmecontent := fmt.Sprintf("# %s\n", g.Config.repoName)
	commitmsg := "Adding a README\n"
	readmepath := fmt.Sprintf("%s/%s", g.Config.Provider["gitlab"].GroupName, g.Config.repoName)
	cfopts := new(gitlab.CreateFileOptions)
	cfopts.Branch = gitlab.String("master")
	cfopts.Content = &readmecontent
	cfopts.CommitMessage = &commitmsg
	_, _, err = g.GitlabClient.RepositoryFiles.CreateFile(readmepath, "README.md", cfopts)
	if err != nil {
		return err
	}

	fmt.Printf("Repository created at %s: %s\n", g.Repo.CreatedAt.Format(time.RFC3339), g.Repo.HTTPURLToRepo)

	return nil
}

// Function CloneRepo clones the remote repository
func (g *GitlabRemote) CloneRepo() error {

	fmt.Printf("Cloning %s\n", g.Repo.WebURL)

	var err error
	var endpoint *transport.Endpoint

	// Define a git endpoint
	switch g.Config.Provider["gitlab"].CloneProtocol {
	case "ssh", "":
		endpoint, err = transport.NewEndpoint(g.Repo.SSHURLToRepo)
	case "http":
		endpoint, err = transport.NewEndpoint(g.Repo.HTTPURLToRepo)
		endpoint.User = g.Config.Provider["gitlab"].User
		endpoint.Password = g.Config.Provider["gitlab"].Password
	default:
		err = errors.New(fmt.Sprintf("Unknown clone protocol %s", g.Config.Provider["gitlab"].CloneProtocol))
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

// Function DeleteRepo deletes a (remote) repository
func (g *GitlabRemote) DeleteRepo() error {

	var pid int // Project id

	// We need to fetch the project ID for deletion
	truep := true
	plopts := new(gitlab.ListProjectsOptions)
	// Set search options (need to be pointers)
	plopts.Search = gitlab.String(g.Config.Provider["gitlab"].GroupName) // We only want projects of our group
	plopts.PerPage = 1000                                                // We set this to 1000 to get all projects, should suffice
	plopts.Owned = &truep                                                // We only want projects we are owner of

	projects, _, err := g.GitlabClient.Projects.ListProjects(plopts)
	if err != nil {
		return err
	}
	// Check for the right repo and get the id
	for _, p := range projects {
		if p.PathWithNamespace == fmt.Sprintf("%s/%s", g.Config.Provider["gitlab"].GroupName, g.Config.repoName) {
			pid = p.ID
		}

	}
	if pid == 0 {
		return errors.New(fmt.Sprintf("Could not find repository %s/%s", g.Config.Provider["gitlab"].GroupName, g.Config.repoName))
	}

	// Delete the repo
	_, err = g.GitlabClient.Projects.DeleteProject(pid)
	if err != nil {
		return err
	}

	return nil
}

// Function ListRepos lists all repos for a given GitlabClient
func (g *GitlabRemote) ListRepos() error {

	truep := true
	plopts := new(gitlab.ListProjectsOptions)

	// Set search options (need to be pointers)
	plopts.Search = gitlab.String(g.Config.Provider["gitlab"].GroupName) // We only want projects of our group
	plopts.PerPage = 1000                                                // We set this to 1000 to get all projects, should suffice
	plopts.Owned = &truep                                                // We only want projects we are owner of
	plopts.OrderBy = gitlab.String("last_activity_at")

	projects, _, err := g.GitlabClient.Projects.ListProjects(plopts)
	if err != nil {
		return err
	}
	// Loop over projects
	for _, p := range projects {
		fmt.Printf("%s - %s\n", p.LastActivityAt.Format(time.RFC3339), p.Name)
	}

	return nil
}

// Function NewGitlabRemote creates a new Remote object and returns it
func NewGitlabRemote(c *Config) (r *GitlabRemote) {

	remote := new(GitlabRemote)

	remote.Config = c
	remote.GitlabClient = gitlab.NewClient(nil, c.Provider["gitlab"].Token)
	remote.GitlabClient.SetBaseURL(c.Provider["gitlab"].HostBaseURL)
	remote.Repo = new(gitlab.Project)

	return remote
}
