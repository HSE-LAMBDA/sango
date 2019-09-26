# Welcome to SANgo!

GoSAN is the discrete-event based simulator with reinforcement learning support. 

# What it simulates?

 - Resource usage
	 - CPU
	 - Links
	 - Storage
		 - SSD
		 - HDD
 - Load generation
	 - Read mode
	 - Write mode

# Installation
Execute the following command: 

    go get -u gitlab.com/lambda-hse/tatlin-hse/gosan

# Example

To demonstrate how GoSAN can be used we implemented a toy version of actual storage array network. It consists of:

 - Load generator
 - 2 controllers 
 - 8 storage units:
	 - 4 SSD
	 - 4 HDD

Client sends data blocks to controllers. They, in own turn, write data to storage units in the round-robin algorithm.  

 
![enter image description here](https://sun9-64.userapi.com/c854216/v854216025/101217/U2KF-i4FyVc.jpg)

# Available flags

The following flags are available at the current moment:
1. Input configuration files:

```
-platform — Tatlin storage array controller (hard-drives, controllers, etc.)
-packet — how much time packets (data blocks) consume on the different Tatlin components depending on their size with self-explainable field tags
-atm_dep - atmosphere (temperature, humidity, pressure, vibration) dependencies
-atm_control - file with atmospheric time-series data in JSON format
-output — file where measured data will be written to
```

2. Load generation:

```
-file_amount — the amount of W/R files
-file_size — file size interval [lower, upper)
-load_range — file read/write frequency [minTime, maxTime)
```

3. Failures:

```
-anomaly_type — failure type; VESNIN1ANOMALY, VESNIN2ANOMALY, VESNIN1CLIENTLINK, VESNIN2CLIENTLINK are available
-anomaly_amount — the number of failures
-anomaly_time_range — time interval [minTime, maxTime). A random number from this interval is chosen, which will be the time of the next failure 
-anomaly_duration —  time interval [minTime, maxTime). Repair time is sampled from this interval. 
```

4. Hard-drive amount:

```
-disk_amount — the number of disks in the JBOD. Overwrite disk amount, given in the -platform flag
```

DeepController synchronization:

```
Turning on/off DC connection.
--controling_mode=true/false(or 1/0), by default = false
Protocol:
--protocol={tcp,udp}, by default tcp
Connection hostname:
--host, by default localhost
Connection port:
--port, by default 1337
DC-GoTatlin communication delay:
--json_delay, by default 1(sec)
```


# Additional

The below command makes the output file well-formed:
```
sed -i -e '1i \[' -e 's/}{/},{/g' -e "\$a]" $(jsonfilename)
```

Description of each field in `json` format can be found [here](https://docviewer.yandex.ru/view/399819079/?*=tRg8KHkaftG96OxyVNSWpjO5zT17InVybCI6InlhLWRpc2s6Ly8vZGlzay9ZYWRyby9BU2Fwcm9ub3Yvc3RvcmFnZS1kYXRhLWpzb24tZGVzY3JpcHRpb24ueGxzeCIsInRpdGxlIjoic3RvcmFnZS1kYXRhLWpzb24tZGVzY3JpcHRpb24ueGxzeCIsInVpZCI6IjM5OTgxOTA3OSIsInl1IjoiMTA1MDIzODgxNTIwMzM0MTM1Iiwibm9pZnJhbWUiOmZhbHNlLCJ0cyI6MTUyMDkzODM3MzI1M30%3D).

For profiling (benchmarking) run:
```
go tool pprof http://localhost:6060/debug/pprof/profile
```
