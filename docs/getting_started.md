# marlinctl 2.0

## Getting it on your system
```
curl http://public.artifacts.marlin.pro/projects/marlinctl/2.0.0/marlinctl-2.0.0-linux-amd64 -o marlinctl
chmod +x marlinctl
md5sum -c - <<< "6bdc336c87bbedb172a0c7a10c0881c9 marlinctl"
sudo cp marlinctl /usr/local/bin/
```

Verify you have the marlinctl setup using
`marlinctl --version`

## Runtimes
Currently marlinctl only supports platform `linux-amd64` and `linux-amd64.supervisor` runtime. Support of other runtimes will come soon

## Setting up beacon
To set up a beacon use
```
sudo marlinctl beacon create
````
To read status of beacon use
```
sudo marlinctl beacon status
````
To read the logs do
```
sudo marlinctl beacon logs
```
To destroy a beacon use
```
sudo marlinctl beacon destroy
```
You can spin up a beacon using custom runtime arguments as well. For example:
```
sudo marlinctl beacon create --discovery-addr "0.0.0.0:9002" --heartbeat-addr "0.0.0.0:9003" --bootstrap-addr "127.0.0.1:8003" --keystore-path ~/Downloads/keystore --keystore-pass-path ~/Downloads/pass 
```
Use `marlinctl beacon -h` for more flags

## Setting up relay

To set up a relay use
```
sudo marlinctl relay eth create
````
To read status of relay use
```
sudo marlinctl relay eth status
````
To read the logs do
```
sudo marlinctl relay eth logs
```
To destroy a relay use
```
sudo marlinctl relay eth destroy
```
you can spin up a relay using custom runtime arguments as well. For example:
```
sudo marlinctl relay eth create --discovery-addr "0.0.0.0:9002" --heartbeat-addrs "0.0.0.0:9003" --sync=mode "light"
```
Use `marlinctl relay eth -h` for more flags