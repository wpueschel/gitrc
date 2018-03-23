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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Provider contains all the necessary config settings of a git provider
type Provider struct {
	Token         string `json:"token"`
	TokenName     string `json:"token_name"`
	HostBaseURL   string `json:"host_base_url"`
	User          string `json:"user"`
	Password      string `json:"password"`
	GroupName     string `json:"group_name"`
	CloneProtocol string `json:"clone_protocol"`
}

// Config contains all necessary config settings
type Config struct {
	Provider   map[string]Provider `json:"provider"`
	repoName   string
	localdir   string
	configfile string
	newrepo    bool
	list       bool
	listLong   bool
	private    bool
	del        bool
}

func (c *Config) readFlags() error {

	//var mandatoryFlags map[string]struct{}

	// cmd line flags
	flag.StringVar(&c.configfile, "c", fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".gitrc.json"), "Config file")
	flag.StringVar(&c.repoName, "n", "", "Repository name")

	flag.BoolVar(&c.list, "l", false, "List remote repository names and last commit timestamp")
	flag.BoolVar(&c.listLong, "L", false, "List remote repository names, cloning urls and last commit timestamp")
	flag.BoolVar(&c.private, "P", false, "Create a private repository")
	flag.BoolVar(&c.del, "D", false, "Delete a remote repository, has to be used with -n")
	flag.BoolVar(&c.newrepo, "N", false, "Create a local and remote repo based on the current directory name")
	flag.Parse()

	return nil
}

func (c *Config) readFile(fname string) error {

	// Read raw config file -> []byte
	rawConfFile, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("Could not read config file %s: %s\n", fname, err)

		return err
	}

	// Unmarshall json from rawConfFile into giteaconf
	err = json.Unmarshal(rawConfFile, &c.Provider)
	if err != nil {
		log.Fatalf("Could not unmarshall json from %s: %s\n", fname, err)

		return err
	}

	return nil
}

// NewConfig creates and returns a new gitea config object
func NewConfig() (*Config, error) {

	c := new(Config)

	err := c.readFlags()
	if err != nil {

		return &Config{}, err
	}

	// Check if we have a config file
	if _, err = os.Stat(c.configfile); err == nil {
		err = c.readFile(c.configfile)
		if err != nil {

			return &Config{}, err
		}
	} else {
		log.Fatalf("Could not read config file %s: %s", c.configfile, err)
		return &Config{}, err
	}

	// if -N is set, we dont need a repo name and make the current directory name the reponame
	if c.newrepo {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {

			return &Config{}, err
		}
		_, c.repoName = filepath.Split(dir)
		c.localdir = dir

	}

	return c, nil
}
