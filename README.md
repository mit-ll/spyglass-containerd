containerd
==========

An abstraction layer between the webapp system and the container
creation logic.

This needs to run in the background for Spyglass to work correctly. Make
sure it starts using your system's init system of choice (systemd, upstart, 
sysvinit)

## Use

Simply: `containerd /path/to/containerd.conf`

What does `containerd.conf` look like? It looks like:
```
{
	"DataHost": "10.10.10.10",
        "DataPort": 3306,
	"DataUser": "spyglass",
	"DataPass": "yams",
	"DataBase": "spyglass"
}
```
## Future Work
This code could use:

* More Resiliency. Currently it crashes if it can't find the container, and
  this is not a good thing

