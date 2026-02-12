# Tools & Setup

> **TL;DR:** The primary development platform is Linux. Install Docker and Docker Compose for container-based workflows. Windows users should configure WSL 2 first.

## Overview

Following a project's `README.md` may get it running, but productive development requires additional tooling and configuration. This page provides guidance on setting up a complete development environment.

## Operating System

The primary development platform is **Linux**. All documentation and scripts are validated for Linux environments. Platform-specific notes:

### Windows

- Set up [Windows Subsystem for Linux (WSL)](Tools-&-Setup/WSL-Setup.md) first.
- For SSH key management on WSL 2, see [Use an SSH agent in WSL](https://pscheit.medium.com/use-an-ssh-agent-in-wsl-with-your-ssh-setup-in-windows-10-41756755993e).

## Containers

Most projects use Docker and Docker Compose for local development and testing:

- [Docker Installation](https://docs.docker.com/get-docker/)
- [Docker Compose Installation](https://docs.docker.com/compose/install/)
- [Docker with WSL Setup](Tools-&-Setup/WSL-Setup.md#how-to-set-up-docker-in-wsl)

## Additional Guides

- [WSL Setup](Tools-&-Setup/WSL-Setup.md)
- [Install Azure CLI](Tools-&-Setup/Install-Azure-CLI.md)
- [Azure Functions Setup](Tools-&-Setup/Azure-Functions-Setup.md)
- [Database Sync](Tools-&-Setup/Database-Sync.md)
