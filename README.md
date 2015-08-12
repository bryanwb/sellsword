# sellsword

Sellsword is a generic command-line tool for switching between application configurations

Technology consultants, such as this project's original author, have to manage different configurations
for a given application for each customer. Doing this separately for each application just sucks. Sellsword
is a generic command-line tool for switching between arbitrary application configurations. Sellsword
currently supports two mechanisms for switching between applications:

#. loading environment variables
#. changing symlinks to either directories or individual files

There are two components to sellsword, `ssw` and `sellsword`. The `sellsword` binary does all the
work but cannot source any environment variables into the parent shell. `source /path/to/ssw [ARG ...]`
executes sellsword with t the supplied arguments and loads any changed environment variables into
the parent shell. This is a huge *hack* but it is the only way I know how to load environment variables into the parent shell.

Sellsword is only supported for the BASH shell and on the OS X and linux operating systems. Sellsword is implemented primarily in Go because writing complex logic in BASH dramatically shortens one's life expectancy.


## Installation

#. Download the tarball
#. `tar xvzf sellsword.tgz --strip-dir=1 -C /usr/local/bin`
#. Add the following to your `.bashrc` file

```
source $(which ssw)      # this loads environment variables for current configurations
alias ssw='source $(which ssw)'   
```



## Configuration

Sellsword knows about a few applications by default but these can be overriden

.ssw/
     app_name/
             app_name-env.ssw  # this is just a yaml file, .ssw extension is used to avoid conflicts
                           # with the application's own configuration files


Example configuration for AWS

.ssw/
     aws/
        acme-dev-env.ssw
        acme-qa-env.ssw
        acme-prod-env.ssw
        megacorp-dev-env.ssw
        megacorp-qa-env.ssw
        megacorp-prod-env.ssw
        current-env.ssw  # symlink to current configuration
        
Any file ending in '-env.ssw' will be sourced to the parent shell and should only include key/value pairs.
Note that this actually a yaml file. Also note that we are using lower-case values here and not the typical
uppercase used by environment variables. This is because most self-respecting developers swap caps lock
with control, making long uppercase names inconvenient to type.

```
# file ~/.ssw/aws/acme-dev-env.ssw
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

Example Configuration for Chef Server

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
# file chef.ssw
type: directory
target: ~/.chef
```

## Usage

ssw list chef   # list possible chef configurations
ssw show chef   # show current configuration in use
ssw use chef acme-qa 
ssw rm chef acme-qa   # remove acme-1 configuration
ssw list chef   # show available chef configurations



# License and Author

Apache 2.0, Copyright 2015 [Bryan W. Berry](mailto:bryan.berry@gmail.com)  