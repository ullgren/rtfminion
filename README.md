# Anypoint Runtime Fabric Minion

This minion helps you automate tasks in Anypoint Runtime Fabric by leveraring the Anypoint platform API. Primarly the API related to Runtime Manager and Exchange. In this way it differs from rtfctl which mainly interact with the Kubernetes cluster directly and is meant to be run on one of the Kubernetes cluster nodes.

Using RTF Minion you can currently execute tasks that otherwise normaly is done through the Anypoint Runtime Manager UI or in somecases using the Mule Maven Plugin.

## Build 

This application is writen in go and requires the go tool chain to build. Please see the go.dev webpage for more information.

### Geting and building

Start by cloning this repository to your local machine.

```
$ git clone https://github.com/ullgren/rtfminion.git
$ cd rtfminion
$ go build -o rtfminion
$ ./rtfminion
```

## Run

To run rtfminion simply execute the resulting binary. It will give you help on how to use each command and sub-command.


## Configure

Instead of providing options such as username and password using commandline flags these can be picked up from a yaml based configuration file.

By default RTF minion will look for a file called `.rtfminion.yaml` located in the users home folder. You can also use the `--config file` option to specify another config file.

The structure of the config file is described in docs/dot-rtfminion.yaml
