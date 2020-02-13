# sshmon-check-heketi
Nagios/Checkmk-compatible SSHMon-check for Heketi-Clusters

## Installation
* Download [latest Release](https://github.com/indece-official/sshmon-check-heketi/releases/latest)

## Usage
```
$> sshmon_check_heketi -host 10.2.0.2 -user <youruser> -key <yourkey>
```

```
Usage of ./dist/bin/sshmon_check_heketi:
  -dns string
        Use other dns server
  -host string
        Host
  -key string
        Key
  -port int
        Port (default 5080)
  -service string
        Service name (defaults to Heketi_<host>)
  -user string
        Username
  -v    Print the version info and exit
```

Output:
```
0 Heketi_10.2.0.2 - OK - Heketi controller on 10.2.0.2 is up and running
0 Heketi_10.2.0.2_f71522478ae95e57343704b0971ae85d - OK - All 3 nodes of heketi cluster 'f71522478ae95e57343704b0971ae85d' are healthy
```

## Development
### Build the binary

```
$> make --always-make
```