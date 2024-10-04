1. Step
  ```bash
  sudo apt update
  sudo apt install ca-certificates curl apt-transport-https lsb-release gnupg
  ```
2. Step
  ```bash
  sudo mkdir -p /etc/apt/keyrings
  curl -sLS https://packages.microsoft.com/keys/microsoft.asc |
      gpg --dearmor |
      sudo tee /etc/apt/keyrings/microsoft.gpg > /dev/null
  sudo chmod go+r /etc/apt/keyrings/microsoft.gpg
  ```
3. Step
  ```bash
  #AZ_REPO="$(lsb_release -cs)"
  AZ_REPO="bullseye"
  echo "deb [arch=`dpkg --print-architecture` signed-by=/etc/apt/keyrings/microsoft.gpg] https://packages.microsoft.com/repos/azure-cli/ $AZ_REPO main" |
    sudo tee /etc/apt/sources.list.d/azure-cli.list
  ```
4. Step
  ```bash
  sudo apt update
  sudo apt install azure-cli
  ```

## References
* https://learn.microsoft.com/pt-br/cli/azure/install-azure-cli-linux?pivots=apt
* https://learn.microsoft.com/en-us/cli/azure/run-azure-cli-docker
