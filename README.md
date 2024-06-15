# vnh1

( EN | [DE](../main/README_DE.md) )

## Description

**vnh1** offers a robust solution for securely and isolatedly running scripts by combining containerization and sandboxing techniques. With vnh1, V8go instances can be executed in separate processes and containers, ensuring strong isolation and security.

## Features

- **Containerized Sandboxing Environment**: Each V8go instance runs in its own container, ensuring complete isolation from other instances and the host system.
- **Process Isolation**: V8go interpreters are executed in separate processes, ensuring that each process manages its own resources.
- **Restricted Access**: The executed scripts have no access to the host filesystem or network, minimizing the risk of damage or data exfiltration.

## Containerization Details

- **Linux**: 
  On Linux, vnh1 utilizes namespaces to create an isolated container environment for each V8go instance. This ensures that each container has its own set of users, processes, and filesystem. By leveraging Linux namespaces, vnh1 can provide a high level of security and isolation, effectively preventing any interference between running scripts and the host system.

- **MacOS**:
  For macOS, vnh1 employs launchd to manage the containerization of VM processes. Launchd is a service management framework for macOS that can start, stop, and manage daemons, applications, processes, and scripts. By using launchd, vnh1 ensures that each VM process is properly isolated and managed, providing a secure environment for script execution.

- **Windows**:
  On Windows, vnh1 does not use traditional containerization techniques. Instead, each V8go instance runs as a separate process. While these processes are isolated from each other, they do not benefit from the additional layer of security provided by containerization. This means that although resources are managed independently for each VM, they are not as securely isolated as on Linux or macOS.

- **BSD**:
  For BSD systems, vnh1 plans to implement containerization similar to that on Linux. BSD offers several containerization technologies such as Jails, which provide a way to partition the operating system into several independent mini-systems. Each Jail has its own filesystem, users, and processes, ensuring strong isolation and security. By utilizing BSD Jails, vnh1 aims to create a secure and isolated environment for each V8go instance, similar to the approach used on Linux with namespaces.

## Functionality

- **RPC Function Calls** ✅
- **Network Functions, Sockets, etc.** ⚠️
- **HTTP Client / Server Functions** ⚠️
- **Electrum Support** ⚠️
- **Lightning Support** ⚠️
- **Nostr Support** ⚠️
- **Database Support (MySQL/SQLite, MongoDB)** ⚠️
- **Crypto Functions (SSL, ECC, RSA, PGP)** ⚠️
- **NodeJS (JS) Console Functions** ⚠️
- **Filesystem Access** ⚠️
- **Wireguard Management** ⚠️

## Roadmap

- **Linux**: Currently in development ⚠️
- **Windows**: Planned with limited support ⚠️
- **MacOS**: Planned support ⚠️
- **BSD**: Planned support ⚠️