# marlinctl2

Marlinctl 2.0 provides an improved process management command line interface for setting up the different components of the Marlin network.

# Stable releases

If you wish to run stable releases compiled by marlin team, please use our public release artifacts. Following is our latest release artifact for you.
```
wget http://public.artifacts.marlin.pro/projects/marlinctl/2.0.0/marlinctl-2.0.0-linux-amd64 --output-document=/usr/local/bin/marlinctl
sudo chmod +x /usr/local/bin/marlinctl
```
If you run `marlinctl -v && md5sum /usr/local/bin/marlinctl`, it should return the following valid results for latest release:
```
marlinctl version 2.0.0 build master@ea6097ddcacf62be617fbc605e64c19db426fd1f
Compiled on: 03-01-2021_09-18-09@UTC
bec6caaa7c6336485964165b9e29edfd  marlinctl-2.0.0-linux-amd64
```

Always try running the latest version of marlinctl. Marlinctl will auto-update by default if new versions are found upstream.


# Cloning

 ```sh
$ git clone https://github.com/marlinprotocol/marlinctl2.git
```

# Building

Only for development purposes, not for release builds unless by marlin team.

Prerequisites: go >= 1.15.1, make, supervisord, supervisorctl

To build marlinctl2 tagged with version 2.0.0 from repository, run
```
$ sh mk.sh 2.0.0
```
A `marlinctl` executable should be built inside the `build` directory

# Usage

Root access is needed to run commands, be sure to run it with sudo if you are not the root user.

To get list of available commands, run

```
$ sudo marlinctl --help
```

The cli is fully explorable, so every subcommand at all depths has a `--help` option. For example, running
```
$ sudo marlinctl beacon --help
```
will list the subcommands available w.r.t the beacon and running
```
$ sudo marlinctl beacon create --help
```
will print the usage and the cli options available.
