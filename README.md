[![Build Status](https://travis-ci.org/wpueschel/gitrc.svg?branch=master)](https://travis-ci.org/wpueschel/gitrc)

# gitrc - Git Remote Control

## Installation

### Downloading binaries

You can download binaries for Linux, MacOS or Windows from the release tab and put them in your path.

Linux example:

```sh
cd /tmp
wget https://github.com/wpueschel/gitrc/releases/download/v0.2.1/gitrc-linux-amd64
sudo cp gitrc-linux-amd64 /usr/local/bin/gitrc
chmod 755 /usr/local/bin/gitrc
rm gitrc-linux-amd64
```

### Building from source with make

Building from source requires a setup [Go environment](https://golang.org/doc/install) with GOPATH etc. correctly set. Also, [git](https://git-scm.com/) and [make](https://www.gnu.org/software/make/) are mandatory. 

Linux example:

```sh
mkdir -p $GOPATH/src/github.com/wpueschel
cd $GOPATH/src/github.com/wpueschel
git clone https://github.com/wpueschel/gitrc.git
cd gitrc
# git checkout [tag] if you want to build a stable version, master is not guaranteed to work at all times.
make dep linux
sudo cp gitrc-linux-amd64 /usr/local/bin/gitrc
make clean
```

Building for other systems: 

- ```make``` or ```make all``` Builds gitrc for Linux, MacOS and Windows
- ```make dep darwin``` Builds gitrc for MacOS
- ```make dep windows``` Builds gitrc for Windows 

All build targets build an amd64 binary.

### Building and installing with just "go get" (not recommended)

Requires a setup go environment.

Linux example:

```
go get github.com/wpueschel/gitrc
sudo cp $GOPATH/bin/gitrc /usr/local/bin
```

This will work, as long a the master branch builds/works. However, ```gitrc version``` will not work properly.

## Usage

```sh
gitrc [options] [repotype]
```

Repotype may be gitea, github or gitlab.

Detailed usage information will be given by issuing 

```sh
gitrc -h [repotype]
```

### Examples

#### Create a new repository on github and clone it

```sh
mkdir test-repo
cd test-repo
gitrc -N github
```

This will create a new repository named test-repo on github, put a basic README.md in it and clone it into the directory test-repo. If you use this with an additional -P, it will create a private repository.

#### Create a remote repository, no clone

```sh
gitrc -n test-repo github
```

This will create a new remote repository named test-repo with a basic README.md in it. With an added -P it will create a private repository.

#### List remote repositories

```sh
gitrc -l gitea
```

This will list all repositories of your configured gitea, sorted by last commit.

```sh
gitrc -L gitlab 
```

This will list all repositories on your configured gitlab with name and clone URL, sorted by last commit.

#### Delete a remote repository

```sh
gitrc -n test-repo -D github
```

This will delete an existing repository on github. Be carefull though, there's no second thought. It's just being deleted.

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

ssh cloning is the default. Same as with GitLab. You will need a running and configured ssh-agent and the host github.com should be already in your known_host file.

Password in the config will only be needed if you clone via http.

