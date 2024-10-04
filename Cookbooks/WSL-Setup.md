## Verify Windows Version
In order to use WSL you must have Windows 10 version 2004 (build 19041.264) or higher.
To verify the version you can go to Windows PowerShell and run the following command:
```
winver
```

## Turn on Virtual Machine and WSL on your Local Machine
Go to Control Panel → Windows Features Turn on or off → Check the following boxes there:

* Virtual Machine Platform
* Windows Subsystem for Linux

## Set the WSL2 as default
Make the WSL2 as default so all the apps will use version 2 by default. You can do that as follow:
```
Windows Powershell → wsl --set-default-version 2
```

## Install Linux Distribution (Ubuntu in this case)
Go to Microsoft Store → Search Ubuntu → Open your desired version → Click on "get" → Let the download complete → Install → Launch the distribution Setup your credentials.
Upgrade the Linux CLI as follows:
```bash
sudo apt-get update
sudo apt-get upgrade
```
You can go to your windows filesystem by:
```bash
cd /mnt
```

**Note:** you get the registry error when you are installing the distribution from the Microsoft store, you can download the Linux kernel update package from [Manual installation steps for older versions of WSL](https://docs.microsoft.com/en-us/windows/wsl/install-manual).

## How to Setup Docker in WSL?
1. Follow the steps above until you have the WSL working.
2. Download Install Docker Desktop from [Download Docker Desktop | Docker](https://www.docker.com/products/docker-desktop/).
3. Now you can open your Linux distribution and start using docker commands.
