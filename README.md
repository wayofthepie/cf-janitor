[![Build Status]
    (https://travis-ci.org/wayofthepie/cf-janitor.svg?branch=master)
] (https://travis-ci.org/wayofthepie/cf-janitor) 
[![Coverage Status]
    (https://coveralls.io/repos/wayofthepie/cf-janitor/badge.svg?branch=master&service=github)
](https://coveralls.io/github/wayofthepie/cf-janitor?branch=master)

# Cloud Foundry Janitor
A Cloud Foundry cli plugin for analyzing usage and cleaning. 

## Installation
Install the Cloud Foundry cli, see https://github.com/cloudfoundry/cli. Clone this repo and from the root run:
```bash
go get -d
```
You will likely run into this issue (see https://github.com/cloudfoundry/cli/issues/529):

```bash
package github.com/cloudfoundry/cli/cf/resources: cannot find package "github.com/cloudfoundry/cli/cf/resources" 
```
To fix:

```bash
$ cd $GOPATH/src/github.com/cloudfoundry/cli
$ bin/generate-language-resources
$ cd - # Back to the root of this project
$ go install
```
Then you can install the plugin:
```bash
$ cf install-plugin $GOPATH/bin/cf-janitor    
$ cf janitor -h
NAME:
   janitor - Delete applications last uploaded before a certain date, filtering out certain applications by name using a regular expression.

USAGE:
   janitor --before ("now"|RFC3339 Timestamp) [--ignore regex]

```

## Usage
```bash
# Find all apps uploaded before now, ignoring apps starting with spring.
$ cf janitor --before "2015-11-24T21:22:53.108Z" --ignore "^spring.*"
Ignoring spring-music

# The --before flag takes an RFC3339 encoded timestamp or the string "now" for the current time
$ cf janitor --before "2015-11-24T21:22:53.108Z" 
spring-music last uploaded 2015-11-24 16:53:06 +0000 UTC
```
