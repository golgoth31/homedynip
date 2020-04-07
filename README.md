# HOMEDYNIP

![](https://github.com/golgoth31/homedynip/workflows/Lint/badge.svg) ![](https://github.com/golgoth31/homedynip/workflows/Release/badge.svg)

Small utility to update dynamic DNS from inside your network when your ISP doesn't give you fixed public IP. If your internet box or other router doesn't integrate the dyndns client you want, deploy the homedynip server on any cloud provider and use the client part to get your own public IP.

This project is under development.

## Usage

### Binary

    ```bash
    ❯ homedynip help

    Homedynip help

    Usage:
    homedynip [command]

    Available Commands:
    client      Request my own IP address and send it a dyndns server
    help        Help about any command
    server      Simple server to send back public IP to client
    version     Show Homedynip version

    Flags:
    -c, --config string      config file
    -h, --help               help for homedynip
        --logFormat string   config file (default "console")
        --logLevel string    config file (default "error")

    Use "homedynip [command] --help" for more information about a command.
    ```

### Docker

Use the yaml files from "deployment/client" or "deployment/server" as samples to deploy Homedynip into a kubernetes cluster.
Generate the needed secret with the following command:

    ```bash
    kubectl create secret generic homedynip --from-literal=username=<your user> --from-literal=password=<your password> --from-literal=hostname=<your hostname>
    ```
