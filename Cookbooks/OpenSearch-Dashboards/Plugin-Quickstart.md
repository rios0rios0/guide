# Plugin Quickstart

> **TL;DR:** Clone OpenSearch Dashboards, place your plugin in the `plugins/` directory, run `yarn osd bootstrap` to install dependencies, and start the dev server with `yarn start`. OSD plugins are full-stack (frontend in `/public`, backend in `/server`).

## Overview

OpenSearch Dashboards (OSD) is an open-source fork of Kibana, primarily built with Node.js and React. This guide covers setting up a local development environment for plugin development. Familiarity with HTML, SASS, TypeScript, JavaScript, and React is expected (see [JavaScript Conventions](../../Code-Style/JavaScript.md)).

## Repository Setup

Clone the required repositories:

```bash
git clone git@github.com:opensearch-project/OpenSearch-Dashboards.git
git clone <your-plugin-repo> OpenSearch-Dashboards/plugins/<plugin-name>
```

### Key Repositories

| Repository                                                                                     | Purpose                          |
|------------------------------------------------------------------------------------------------|----------------------------------|
| [OpenSearch Dashboards](https://github.com/opensearch-project/OpenSearch-Dashboards)           | Base platform                    |
| [Security Dashboards Plugin](https://github.com/opensearch-project/security-dashboards-plugin) | Authentication and authorization |

## Architecture

![](.assets/app_relationship.png)

### Base Stack

- **OpenSearch** (Elasticsearch fork) runs in containers via Docker Compose.
- **OpenSearch Dashboards** (Kibana fork) depends on custom plugins and is started via the development server.

### Plugin Structure

OSD plugins are **full-stack**: the `/public` directory contains frontend code and `/server` contains backend code.

The Security Dashboards plugin provides authentication and authorization. It is used by other plugins to retrieve user information before forwarding requests to backend APIs.

## Getting Started

### 1. Bootstrap Dependencies

```bash
cd OpenSearch-Dashboards
yarn osd bootstrap
```

This installs all dependencies and builds internal packages. Run this command whenever you update packages or sync with upstream changes.

**Network timeout?** Add a timeout flag:

```bash
yarn osd bootstrap --network-timeout 1000000
```

Or configure it in `.yarnrc`:

```
network-timeout 1000000
```

To start fresh:

```bash
yarn osd clean
```

### 2. Start the Development Server

```bash
yarn start
```

Wait for the server to be ready:

```
[info][listening] Server running at http://localhost:5603/pgt
```

Bundle compilation may take a few minutes. Wait for:

```
[success][@osd/optimizer] 28 bundles compiled successfully after 145.9 sec, watching for changes
```

### Troubleshooting

**Docker memory error:**

```
Error: max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]
```

Fix on Linux:

```bash
sudo sysctl -w vm.max_map_count=262144
```

Fix on Windows (WSL):

```bash
wsl -d docker-desktop
sysctl -w vm.max_map_count=262144
```

## References

- [OpenSearch Dashboards Developer Guide](https://github.com/opensearch-project/OpenSearch-Dashboards/blob/main/DEVELOPER_GUIDE.md)
- [OpenSearch Downloads](https://opensearch.org/downloads.html)
