containerd
==========

An abstraction layer between the webapp system and the container
creation logic.

This needs to run in the background for Spyglass to work correctly. Make
sure it starts using your system's init system of choice (systemd, upstart, 
sysvinit)

## Use

Simply: `containerd /path/to/containerd.conf`

`containerd.conf` is a simple JSON-based config file:
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
## Disclaimer
---

DISTRIBUTION STATEMENT A. Approved for public release: distribution unlimited.

This material is based upon work supported by the Under Secretary of Defense for Research and Engineering under Air Force Contract No. FA8721-05-C-0002 and/or FA8702-15-D-0001. Any opinions, findings, conclusions or recommendations expressed in this material are those of the author(s) and do not necessarily reflect the views of the Under Secretary of Defense for Research and Engineering.

The software/firmware is provided to you on an As-Is basis
