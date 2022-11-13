# hContainers

![pipeline status](https://github.com/hContainers/hContainers/actions/workflows/CI.yml/badge.svg?branch=main)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)
![Terminal](https://badgen.net/badge/icon/terminal?icon=terminal&label)
[![Open Source? Yes!](https://badgen.net/badge/Open%20Source/Yes%21/blue?icon=github)](https://github.com/Naereen/badges/)

hContainers is a tool to help you manage your Docker containers at the hetzner cloud in a simple way without the need of managing the vm's yourself.

# Setup and Installation

## Install hContainers

### Use a package manager

_Follows_

### Download the binary

You can download the binary from the [releases page](https://github.com/hContainers/hContainers/releases). There are binaries for Linux (32 and 64 bit), Windows (32 and 64 bit) and MacOS (64 bit).

### Build from source

These instructions are the same as in the [development docs](docs/Development.md) under the section "Build".

## Perpare Hetzner Cloud

First you need to setup a hetzner cloud account and create a project. This project will be used to create the vm's and all the other stuff. Because of this the project needs to be empty. If you already have a project with servers in it you need to create a new one.

After you have created the project you need to create a new api token. This token will be used to authenticate against the hetzner cloud api. You can create the token in the hetzner cloud console in the project under "Security" -> "API Tokens". You need to give the token the permissions "Read" and "Write" for the project.

## Configure hContainers

The api token from Hetzner needs to be stored in the environment variable `HCLOUD_TOKEN`.\
The path of the ssh key that will be used to connect to the vm's needs to be stored in the environment variable `HCONTAINERS_SSH_KEY_PATH`. The path should NOT includ the `.pub` extension. This will be added automatically. _Please take note that this program also uses the private key to connect to the vm's. So make sure that the private key is not readable by other users._

You can set tje environment variables with the following commands:

```bash
export HCLOUD_TOKEN="your-token"
export HCONTAINERS_SSH_KEY_PATH="~/.ssh/id_ed25519"
```

# Concept

The base of hContainers are runners. This are vms which get managed by hContainers. The only thing you need to think about is how many containers should run on a runner and how much memory and cpu the runner should get. The rest is managed by hContainers.

_More follows_

# Tools used by hContainers

- [Hetzner Cloud API](https://docs.hetzner.cloud/)
- [Containerd](https://containerd.io/)
- [NerdCtl](https://github.com/containerd/nerdctl)
