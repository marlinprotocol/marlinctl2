# marlinctl

Marlinctl 2.X provides an improved process management command line interface for setting up the different components of the Marlin network.

# Stable releases

Current stable release: 2.1.0

If you wish to run stable releases compiled by marlin team, please use our public release artifacts. Following is marlinctl 2.0.0 for you. (marlinctl 2.0.0 will automatically upgrade to any new versions of marlinctl upon running any command later).
```
wget http://public.artifacts.marlin.pro/projects/marlinctl/2.0.0/marlinctl-2.0.0-linux-amd64 --output-document=/usr/local/bin/marlinctl
sudo chmod +x /usr/local/bin/marlinctl
```
If you run `marlinctl -v && md5sum /usr/local/bin/marlinctl`, it should return the following valid results for latest release:
```
marlinctl version 2.1.0 build master@99152b066bdb648c415a3d7de1310221268d2117
Compiled on: 18-03-2021_09-02-50@UTC
1631778fa78a11485ed4953302024717  marlinctl-2.1.0-linux-amd64
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
