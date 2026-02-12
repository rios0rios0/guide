# WSL Setup

> **TL;DR:** Install WSL 2 on Windows 10 (build 19041+), set it as the default version, install Ubuntu from the Microsoft Store, configure Git with SSH via 1Password, and install Docker Desktop for container support.

## Prerequisites

Windows 10 version 2004 (build 19041.264) or higher is required. Verify your version:

```powershell
winver
```

## Installation Steps

### 1. Enable Required Windows Features

Navigate to **Control Panel > Windows Features** and enable:

- Virtual Machine Platform
- Windows Subsystem for Linux

### 2. Set WSL 2 as Default

```powershell
wsl --set-default-version 2
```

### 3. Install a Linux Distribution

1. Open the **Microsoft Store** and search for **Ubuntu**.
2. Select your preferred version, click **Get**, and wait for the download.
3. Launch the distribution and configure your credentials.

Update the system:

```bash
sudo apt-get update && sudo apt-get upgrade
```

Access the Windows filesystem:

```bash
cd /mnt
```

**Note:** If you encounter a registry error during installation, download the Linux kernel update package from [Manual WSL Installation Steps](https://docs.microsoft.com/en-us/windows/wsl/install-manual).

### 4. Configure Git with SSH and 1Password

Most distributions include Git by default. If not:

```bash
sudo apt install git
```

Follow these guides to set up SSH key management via 1Password:

1. [Get Started with 1Password for SSH](https://developer.1password.com/docs/ssh/get-started/)
2. [Use the 1Password SSH Agent with WSL](https://developer.1password.com/docs/ssh/integrations/wsl/)
3. [Sign Git Commits with SSH](https://developer.1password.com/docs/ssh/git-commit-signing/)

## How to Set Up Docker in WSL

1. Complete the WSL setup steps above.
2. Download and install [Docker Desktop](https://www.docker.com/products/docker-desktop/).
3. Open your Linux distribution and start using Docker commands.

## References

- [WSL Documentation](https://docs.microsoft.com/en-us/windows/wsl/)
- [Docker Desktop for Windows](https://docs.docker.com/desktop/windows/install/)
- [1Password SSH Integration](https://developer.1password.com/docs/ssh/)
