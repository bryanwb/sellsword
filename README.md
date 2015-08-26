# sellsword

Sellsword is a generic command-line tool for switching between application environments

Technology consultants, such as this project's original author, have to manage different environments
for a given application for each customer. Doing this separately for each application just sucks. Sellsword
is a generic command-line tool for switching between arbitrary application environments. Sellsword
currently supports two mechanisms for switching between applications:

* loading environment variables
* changing symlinks to either directories or individual files

There are two components to sellsword, `ssw` and `sellsword`. The `sellsword` binary does all the
work but cannot source any environment variables into the parent shell. `source /path/to/ssw [ARG ...]`
executes sellsword with t the supplied arguments and loads any changed environment variables into
the parent shell. This is a huge *hack* but it is the only way I know how to load environment variables into the parent shell.

Sellsword is only supported for the BASH shell and on the OS X and linux operating systems. Sellsword is implemented primarily in Go because writing complex logic in BASH dramatically shortens one's life expectancy.

Sellsword has two core concepts, *applications* which are defined by YAML files in ~/.ssw/config/ and
*environments* per application stored in ~/.ssw/appname/. The environment can be a set of environment
variables, a directory containing arbitrary files, or a single file to be linked to a particular
location.

There are two types of environments, directory and *environment*. I realize the naming of this second type is very confusing so suggestions are most welcome.

## Installation

* [Download the tarball](https://github.com/bryanwb/sellsword/releases)
* `tar xvzf sellsword*.gz -C /usr/local/bin`
* Add the following to your `.bashrc` file:

        # this loads environment variables for current configurations
        alias ssw='source $(which ssw)'
        ssw load

## Configuration

Sellsword knows about a few applications by default but these can be overridden:

```
.ssw/
     config/
             aws.ssw  # this is just a yaml file, .ssw extension is used to avoid conflicts
                      # with the application's own configuration files
```


### Example environment for AWS

```
.ssw/
     aws/
        acme-dev
        acme-qa
        acme-prod
        megacorp-dev
        megacorp-qa
        megacorp-prod
        current  # symlink to current environment
```
        
The environment file which current points to will be sourced to the parent shell and should only include key/value pairs.
Note that this actually a yaml file. Also note that we are using lower-case values here and not the typical
uppercase used by environment variables. This is because most self-respecting developers swap caps lock
with control, making long uppercase names inconvenient to type.

```
# file ~/.ssw/aws/acme-dev
access_key: ASDFAFASDFSDAF...
secret_key: asdfasdfadsf...
```

There should be a corresponding configuration file that maps the keys to environment variables. Notice
that you can map a single key to multiple environment variables.

```
# file ~/.ssw/config/aws.ssw yaml
type: environment

variables:
  - access_key=AWS_ACCESS_KEY_ID
  - access_key=AWS_ACCESS_ID
  - region=AWS_DEFAULT_REGION
  - region=AWS_REGION
```

It is important to note that the variables section is a list of key/value pairs where the same keys
can be present multiple times. This is so that the same key can be mapped to multiple values. The
reason for this is that different applications often use different names for the same environment
variables.

Example Setup for Chef Server

```
.ssw/
     chef/
        acme-dev/
                knife.rb
                chefadmin.pem
                default-validator.pem
        acme-qa/
        acme-prod/
     config/
        chef.ssw
```

```
# file chef.ssw
type: directory
target: ~/.chef
```

Sellsword supports running arbitrary shell command when an environment is loaded and unloaded. In practice this
is only relevant to directory environments, at least in my experience.

For example, you could use sellsword to manage which ssh keys are loaded into your local SSH Agent

```
# file ssh.ssw
type: directory
target: ~/ssh
load_action: ssh-add $SSW_CURRENT/*.pem $SSW_CURRENT/*.priv
unload_action: ssh-add -D
```

Here is what the contents of ~/.ssw/ssh/acme directory look like

```
acme-prod.pem
acme-dev.priv
acme-dev.pub
```

## Usage

```
ssw list chef          # list possible chef environments
ssw show chef          # show current environment in use
ssw load               # load default environments
ssw new aws acme-prod  # wizard to create new aws environment
ssw use aws acme-qa
ssw unlink aws         # unlink default environment but do not delete the
                       # actual environment
ssw rm aws acme-qa     # remove acme-qa environment  TODO
```

For applications with the *environment* type, `sellsword new app_name env_name` will interactively prompt you
for the values needed.

# Development

This project uses [godep](https://github.com/tools/godep) for managing dependencies. Godep has some
quirks so it is worth reading a
[tutorial](https://blog.codeship.com/godep-dependency-management-in-golang/) or two. The trickiest
part is that you check out the sellsword to a location inside your GOPATH/src. This author locates
the sellsword code at $GOPATH/src/github.com/bryanwb/sellsword.

To install godep, `go get github.com/tools/godep`

I also recommend setting the GOBIN environment variable to $GOPATH/bin

# License and Author

Apache 2.0, Copyright 2015 [Bryan W. Berry](mailto:bryan.berry@gmail.com)  
