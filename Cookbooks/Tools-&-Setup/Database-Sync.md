# Database Sync

> **TL;DR:** Use the `database-sync` tool to copy PostgreSQL databases between environments. Configure source/destination in `.env`, specify tables to exclude in `config.yaml`, and run with `pdm start`.

## Overview

The Database-Sync tool copies one or more PostgreSQL databases between environments (e.g., production to development, development to local).

## Prerequisites

| Tool     | Minimum Version | Installation                                                  |
|----------|-----------------|---------------------------------------------------------------|
| Docker   | 18.09+          | [Install Docker](https://docs.docker.com/install)             |
| Python   | 3.9+            | [Download Python](https://www.python.org/downloads/)          |
| PDM      | 2.11.2+         | [Install PDM](https://pdm-project.org/latest/#installation)   |
| pg_dump  | 15              | See `INSTALL_PG_DUMP.md` in the repository                    |
| psycopg2 | 2.9.9           | [Install psycopg2](https://www.psycopg.org/docs/install.html) |

## Setup

### 1. Clone the Repository

```bash
git clone git@github.com:fnk0c/database-sync.git
```

### 2. Install Dependencies

```bash
cd database-sync
pdm install
```

### 3. Configure Environment Variables

Copy `.env.example` to `.env` and configure:

| Prefix             | Description                                      |
|--------------------|--------------------------------------------------|
| `PRODUCTION_*`     | Source database (the one being copied)           |
| `DEVELOPMENT_*`    | Destination database (receiving the copy)        |
| `SYNC_TARGETS`     | Database name(s) to synchronize                  |
| `CONFIG_FILE_PATH` | Path to `config.yaml` (default: `./config.yaml`) |

### 4. Configure Table Exclusions

In `config.yaml`, specify tables to exclude from the sync (e.g., tables containing sensitive data):

```yaml
ignore:
  - database: store
    tables:
      - "users"
```

If no tables should be excluded, use an empty string: `""`.

### 5. Run the Sync

```bash
pdm start
```
