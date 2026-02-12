# Azure Functions Setup

> **TL;DR:** Install Azure Functions Core Tools via npm, initialize a project with `func init`, create a custom handler with `func new --language custom`, implement the Go handler, and deploy with `func azure functionapp publish <app_name>`.

## Overview

This guide covers deploying an Azure Functions custom handler written in Go to Azure using the Azure Functions Core Tools via a Linux terminal.

## Steps

### 1. Install Azure Functions Core Tools

Requires Node.js and npm:

```bash
npm i -g azure-functions-core-tools@4 --unsafe-perm true
```

### 2. Create a New Functions Project

```bash
func init
```

This creates the necessary project structure in the current directory.

### 3. Create a New Function

```bash
func new --language custom
```

Select the appropriate template for your function trigger.

### 4. Implement the Custom Handler

Write the Go code that implements the function's business logic. Refer to the [Azure Functions Custom Handlers documentation](https://learn.microsoft.com/en-us/azure/azure-functions/functions-custom-handlers) for implementation details.

### 5. Deploy to Azure

Ensure the [Azure CLI](Install-Azure-CLI.md) is installed and authenticated:

```bash
az login
```

Deploy from the project root directory:

```bash
func azure functionapp publish <app_name>
```

This command:

1. Packages the code and dependencies into a `.zip` file.
2. Deploys the archive to the Azure Functions app.
3. Configures triggers and bindings.

Once complete, the function is live. Test it by sending HTTP requests to the endpoints listed in the Azure Portal or via the Azure CLI.

## References

- [Azure Functions Core Tools](https://learn.microsoft.com/en-us/azure/azure-functions/functions-run-local)
- [Azure Functions Custom Handlers](https://learn.microsoft.com/en-us/azure/azure-functions/functions-custom-handlers)
- [Install Azure CLI](Install-Azure-CLI.md)
