## Database Sync: Copy PostgreSQL Databases Easily

You can copy one or more PostgreSQL databases, whether it's from your development environment to local, production to development, or other scenarios, using the `Database-Sync` tool.

Follow these steps to get started:

### 1. Install Prerequisites

Ensure the following tools are installed before using the `Database-Sync` tool:

- **Docker**: Version 18.09 or higher. [How to Install Docker](https://docs.docker.com/install)
- **Python**: Version 3.9 or higher. [Download Python](https://www.python.org/downloads/)
- **PDM**: Version 2.11.2 or higher. [Install PDM](https://pdm-project.org/latest/#installation)
- **pg_dump**: Version 15. [Install pg_dump](./INSTALL_PG_DUMP.md)
- **psycopg2**: Version 2.9.9. [Install psycopg2](https://www.psycopg.org/docs/install.html)

### 2. Clone the Repository

Clone the `database-sync` repository from GitHub:

```bash
git clone git@github.com:fnk0c/database-sync.git
```

### 3. Install the Application

Navigate to the cloned repository folder and install the required dependencies:

```bash
pdm install
```

### 4. Set Up the .env File

- Copy the provided `.env.example` file and rename it to `.env`.
- In this file, set up the required environment variables:
  - Variables starting with **`PRODUCTION_`** are for the source database (the one you're copying data from).
  - Variables starting with **`DEVELOPMENT_`** are for the destination database (the one receiving the data copy).
- The **`SYNC_TARGETS`** variable is required to specify the name(s) of the database(s) you're copying.
- The **`CONFIG_FILE_PATH=./config.yaml`** variable should remain as is, as the `config.yaml` file is already created in the repository.

### 5. Configure the `config.yaml` File

In `config.yaml`, you need to specify the following variables:

- **database**: Is Required. Is the name of the database that you're copying, the same value that you're using in the `.env` for `SYNC_TARGETS` and 
- **tables**: Is Required. If there are no tables to ignore, you need to insert (`""`).

For example, if your database is named `store` and you want to exclude the `users` table (which contains sensitive data), your configuration should look like this:

```yaml
ignore:
  - database: store
    tables:
      - "users"
```

### 6. Run the Application
Once you've installed all dependencies and configured the .env and config.yaml files, you can run the application using:

```bash
pdm start
```
