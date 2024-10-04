You can deploy an Azure Functions custom handler written in Go to Azure using the Azure Functions Core Tools (previously known as the Azure Functions CLI) via a Linux terminal.
The following steps outline the basic process:

1. Install the Azure Functions Core Tools: To install the Azure Functions Core Tools, you'll need to have Node.js and npm installed.
  You can install the tools by running the following command in your terminal:
  ```bash
  npm i -g azure-functions-core-tools@4 --unsafe-perm true
  ```

2. Create a new Functions project: To create a new Functions project, run the following command in your terminal:
  ```bash
  func init
  ```
  This command will create a new Functions project in the current directory, complete with the necessary files and directories.

3. Create a new Function: To create a new Function, run the following command in your terminal:
  ```bash
  func new --language custom
  ```
  This command will create a new Go function in the current project. You can choose the specific template for the function you want to create.

4. Implement your custom handler: To implement your custom handler, you'll need to write Go code that implements the necessary logic for your function.
You can find the implementation details in the Go documentation for Azure Functions.

5. Deploy your Function to Azure: To deploy your Function to Azure, you'll need to create an Azure Functions app and then deploy your code to it. To create an Azure Functions app, you can use the Azure CLI or the Azure portal. Once you have created the app, you can deploy your code to it by running the following command in your terminal:
  ```bash
  func azure functionapp publish <app_name>
  ```
  Where `<app_name>` is the name of your Azure Functions app. This command will deploy your code to the Azure Functions app, making it available for execution.

  Before running this command, make sure that you are in the root directory of your Azure Functions project.
  Additionally, you will need to have the Azure CLI installed (see [here](Install-Azure-CLI.md)) and be logged in to your Azure account via the CLI. You can log in using the following command:
  ```bash
  az login
  ```
  After logging in, you can publish your code to the Azure Functions app by running the `func azure functionapp publish` command. This command will perform the following steps:

  * Package your code and dependencies into a .zip file
  * Deploy the .zip file to your Azure Functions app
  * Configure the necessary resources, such as function triggers and bindings, in your Azure Functions app

  Once the deployment is complete, your Azure Functions custom handler will be live and ready for use.
  You can test it by sending HTTP requests to the appropriate endpoints, which you can find in the Azure portal or by using the Azure CLI to list the functions in your app.

These are the basic steps for deploying an Azure Functions custom handler written in Go to Azure via a Linux terminal.
You may need to make additional modifications to your code or your Azure Functions app configuration to meet your specific requirements.
