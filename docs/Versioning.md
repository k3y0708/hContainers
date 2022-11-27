# Versioning

This document describes the versioning of each component of the project which gets versioned.

## Versioning of the Containers

The containers are versioned in their name. The format is the following:

**v1**: `<container-name>-<portprefix>-<instance-id>-<version>`

The `<version>` is the version of the container. The `<instance-id>` is the instance id of the container. The combination of `<portprefix>` and `<instance-id>` is used to identify the port on which the container is running (Port `<portprefix><instance-id>`). The `<container-name>` is the name of the container.
