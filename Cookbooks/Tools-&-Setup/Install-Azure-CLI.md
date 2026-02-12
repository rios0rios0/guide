# Install Azure CLI

> **TL;DR:** Install the Azure CLI on Debian/Ubuntu by adding the Microsoft GPG key and APT repository, then running `apt install azure-cli`.

## Installation Steps

### 1. Install Prerequisites

```bash
sudo apt update
sudo apt install ca-certificates curl apt-transport-https lsb-release gnupg
```

### 2. Add the Microsoft GPG Key

```bash
sudo mkdir -p /etc/apt/keyrings
curl -sLS https://packages.microsoft.com/keys/microsoft.asc |
    gpg --dearmor |
    sudo tee /etc/apt/keyrings/microsoft.gpg > /dev/null
sudo chmod go+r /etc/apt/keyrings/microsoft.gpg
```

### 3. Add the Azure CLI APT Repository

```bash
AZ_REPO=$(lsb_release -cs)
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/microsoft.gpg] https://packages.microsoft.com/repos/azure-cli/ $AZ_REPO main" |
    sudo tee /etc/apt/sources.list.d/azure-cli.list
```

### 4. Install the Azure CLI

```bash
sudo apt update
sudo apt install azure-cli
```

## References

- [Install Azure CLI on Linux](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli-linux?pivots=apt)
- [Run Azure CLI in a Docker Container](https://learn.microsoft.com/en-us/cli/azure/run-azure-cli-docker)
