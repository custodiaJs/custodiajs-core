### CustodiaJS

**Description:**

CustodiaJS offers a robust solution for securely and isolatedly running JavaScript scripts by combining containerization and sandboxing techniques. With CustodiaJS, V8go instances can be executed in separate processes and containers, ensuring strong isolation and security.

**Features:**

- **Containerized Sandboxing Environment:** Each V8go instance runs in its own container, ensuring complete isolation from other instances and the host system.
- **Process Isolation:** V8go interpreters are executed in separate processes, ensuring that each process manages its own resources.
- **Restricted Access:** The executed scripts have no access to the host filesystem or network, minimizing the risk of damage or data exfiltration.

**Containerization Details:**

- **Linux:** On Linux, CustodiaJS uses namespaces to create an isolated container environment for each V8go instance. This ensures that each container has its own set of users, processes, and filesystem.
- **MacOS:** For macOS, CustodiaJS employs launchd to manage the containerization of VM processes.
- **Windows:** On Windows, traditional containerization techniques are not used. Instead, each V8go instance runs as a separate process.
- **BSD:** For BSD systems, CustodiaJS plans to implement containerization similar to Linux, utilizing technologies like Jails.

**Functionality:**

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

**Roadmap:**

- **Linux:** Currently in development ⚠️
- **Windows:** Planned with limited support ⚠️
- **MacOS:** Planned support ⚠️
- **BSD:** Planned support ⚠️