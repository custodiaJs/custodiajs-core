# vnh1

( EN | [DE](../main/README_DE.de) )

## Description

**vnh1** offers a robust solution for securely and isolatedly running scripts by combining containerization and sandboxing techniques. With vnh1, V8go instances can be executed in separate processes and containers, ensuring strong isolation and security.

## Features

- **Containerized Sandboxing Environment**: Each V8go instance runs in its own container, ensuring complete isolation from other instances and the host system.
- **Process Isolation**: V8go interpreters are executed in separate processes, ensuring that each process manages its own resources.
- **Restricted Access**: The executed scripts have no access to the host filesystem or network, minimizing the risk of damage or data exfiltration.

## Benefits

- **Security**: Combining containerization and sandboxing creates a secure execution environment.
- **Flexibility**: Multiple scripts can be executed in parallel and independently without affecting each other.
- **Isolation**: Strict separation between the scripts and the host system prevents unwanted interactions and enhances security.

## Use Cases

- Secure execution of custom JavaScript code.
- Provision of an isolated environment for script-based automations.
- Development and testing of JavaScript code in a controlled environment.
