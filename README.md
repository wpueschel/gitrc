# gitrc - Git Remote Control

## Installation

```
go get github.com/wpueschel/gitrc
```

If you don't have a working go setup, you can also download binaries for Linux or Windows from the release tab.

## Usage

```
gitrc [options] [repotype]
```

Repotype may be gitea, github or gitlab.

Detailed usage information will be given by issuing 

```
gitrc -h [repotype]
```

### Examples

#### Create a new repository on github and clone it

```
mkdir test-repo
cd test-repo
gitrc -N github
```

This will create a new repository named test-repo on github, put a basic README.md in it and clone it into the directory test-repo.

#### Only create a remote repository, no clone

```
gitrc -n test-repo github
```

This will create a new remote repository named test-repo with a basic README.md in it.

## Config file

Some of the options have to be put into a config file (mainly credentials). The config file has to be valid JSON.
See the examples directory for an example config file.
  
The default location where gitrc will look for the config file will be ```$HOME/.gitrc.json```

## Gitea

Right now, for gitea, only http/https will work for cloning a remote repository (-N).

## GitLab 

For Gitlab, if you don't set a group in the config file, the group will be the username.
  
If you chose ssh as cloning protocol, which is the default, you will need a running and configured ssh agent.

Password in the config will only be needed if you clone via http.

## GitHub

ssh cloning is the default. Same as with GitLab. You will need a runnign and configured ssh-agent and the host github.com should be already in your known_host file.

Password in the config will only be needed if you clone via http.
