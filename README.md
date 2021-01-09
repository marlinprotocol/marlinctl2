# Marlinctl 2.0.0

Marlinctl 2.0 provides an improved process management command line interface for setting up the different components of the Marlin network.

# Cloning

 ```sh
$ git clone https://github.com/marlinprotocol/marlinctl2.git
```

# Building

Prerequisites: go >= 1.15.1, make, supervisord, supervisorctl

NOTE: master may be in dev and may not reflect public release states. Build at you own risk.
If you wish to run stable editions, please use the following:
```
wget http://public.artifacts.marlin.pro/projects/marlinctl/2.0.0/marlinctl-2.0.0-linux-amd64
chmod +x marlinctl-2.0.0-linux-amd64
mv marlinctl-2.0.0-linux-amd64 marlinctl 
cp marlinctl /usr/local/bin/
```
To build, run
```
$ make
$ make install
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
