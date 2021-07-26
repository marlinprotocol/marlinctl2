# marlinctl

Marlinctl 2.X provides an improved process management command line interface for setting up the different components of the Marlin network.

# Stable releases

If you wish to run the latest stable releases compiled by marlin team, please use our public release artifacts. Following is marlinctl 2.5.1 for you (which automatically upgrades to latest publically released marlinctl upon running the following).
```sh
sudo wget http://public.artifacts.marlin.pro/projects/marlinctl/2.5.1/marlinctl-2.5.1-linux-amd64 --output-document=/usr/local/bin/marlinctl
if [[ `md5sum /usr/local/bin/marlinctl | cut -d' ' -f1` == "2acbdb08c09ffadf2ce4fe57bbbd9f96" ]]; then  echo "verified md5sum" ; else echo "wrong md5sum, deleting marlinctl" && sudo rm /usr/local/bin/marlinctl;  fi
sudo chmod +x /usr/local/bin/marlinctl
sudo marlinctl --registry-sync
```
If you run `marlinctl -v`, it should return the latest release of marlinctl. For example (for illustration purposes only):
```
marlinctl version 2.5.1 build master@76eacbd3b31dc8955caffb0313d133ed1e44c0ea
Compiled on: 08-06-2021_04-40-35@UTC
```

Always try running the latest version of marlinctl. Marlinctl will auto-update by default or on calling `marlinctl --registry-sync` if new versions are found upstream.

# Cloning

 ```sh
git clone https://github.com/marlinprotocol/marlinctl2.git
```

# Building

Only for development purposes, not for release builds unless by marlin team.

Prerequisites: go >= 1.15.1, make, supervisord, supervisorctl

To build marlinctl2 tagged with version 2.0.0 from repository, run
```sh
sh mk.sh 2.0.0
```
A `marlinctl` executable should be built inside the `build` directory

# Usage

Root access is needed to run commands, be sure to run it with sudo if you are not the root user.

To get list of available commands, run

```sh
sudo marlinctl --help
```

The cli is fully explorable, so every subcommand at all depths has a `--help` option. For example, running
```sh
sudo marlinctl beacon --help
```
will list the subcommands available w.r.t the beacon and running
```sh
sudo marlinctl beacon create --help
```
will print the usage and the cli options available.
