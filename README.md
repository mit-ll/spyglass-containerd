containerd
==========

An abstraction layer between the webapp system and the container
creation logic.

This needs to run in the background for Spyglass to work correctly. Make
sure it starts using your system's init system of choice (systemd, upstart, 
sysvinit)
