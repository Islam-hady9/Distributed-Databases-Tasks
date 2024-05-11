# Distributed Databases Tasks

## [Task 1 : Shutdown](https://github.com/Islam-hady9/Distributed-Databases-Tasks/tree/main/Shutdown)

This project involves a network-based application using Go (Golang) that allows communication between a master server and multiple client devices (slaves). The application is structured into two main parts: the server code and the client code, each handling specific roles in the communication process.

### Server Code
The server operates as the master device. It is responsible for:
1. **Listening for Incoming Connections:** It starts by listening on a TCP port for connections from client devices.
2. **Accepting Connections:** Each client connection is accepted and handled in its own goroutine, allowing concurrent processing of multiple clients.
3. **Handling Client Data:** For each connected client, it reads data sent from the client, processes it, and can send responses back. It specifically listens for a "shutdown" command from the master to terminate the client.
4. **Managing Client Connections:** It maintains a map of active connections which can be dynamically updated as clients connect or disconnect.
5. **Shutting Down:** It listens for a "shutdown" command on the command line interface to terminate all active client connections and shut down the server gracefully.

### Client Code
The client operates as a slave device, handling:
1. **Connecting to the Server:** It connects to the specified server's IP address and port.
2. **Reading Server Messages:** In a dedicated goroutine, it continuously reads messages from the server. If a "shutdown" command is received, it triggers the shutdown sequence on the client machine.
3. **Sending User Messages:** It reads user input from the standard input and sends it to the server. This allows interactive communication between the user at the client-side and the server.
4. **Handling Shutdown Command:** Upon receiving a shutdown command from the server, the client shuts down its operation based on the operating system it's running on (Windows or Linux).

### Operating System Specific Operations
The application incorporates functionality to handle different operating systems by checking the file path separators to determine if the client is running on Windows or Linux, adjusting the shutdown commands accordingly.

### Overall Architecture
The project effectively demonstrates network programming, concurrent processing using goroutines, basic input/output operations, and conditional logic based on operating system specifics. It can be used as a foundational framework for more complex client-server applications where real-time data exchange and remote command execution are required.

## [Task 2 : FileShare](https://github.com/Islam-hady9/Distributed-Databases-Tasks/tree/main/FileShare)

This project entails a distributed file processing application using Go (Golang), designed to demonstrate networking, file handling, and concurrent operations across multiple networked machines. The system architecture includes a master server, multiple slave clients, and an optional client interface for initiating operations. Here’s a breakdown of each component’s roles and interactions:

### Master Server
The master server is the central coordinator in this system. Its main responsibilities include:
1. **Receiving Files from Clients:** It accepts connections from clients who may send files to the server.
2. **File Division:** Upon receiving a file, it divides the file into two parts and saves them as separate files.
3. **Distributing File Parts:** It connects to two slave servers (slaves 1 and 2) and sends each one a part of the file for potential further processing.
4. **Receiving Processed Parts:** (implied but not explicitly shown in the given code) It may receive processed parts back from the slaves.
5. **File Reassembly:** It reassembles the parts back into a single file after processing.
6. **Sending Combined File to Clients:** Finally, it sends the reassembled file back to the client.

### Client
The client has a simple user interface that allows it to either send a file to the master server for processing or receive a processed file back from the server. The client functionality is split into two commands:
1. **Send:** The client sends a file to the master for processing.
2. **Receive:** The client receives the processed (combined) file back from the master.

### Slaves
There are two slave servers in this system, both having similar functionalities but handling different parts of the file:
1. **Receiving File Parts:** Each slave server listens for a connection from the master, receives a part of the file, and potentially processes it.
2. **Acknowledgment to Master:** After receiving the file part, each slave sends an acknowledgment back to the master. This could be used to signal readiness or successful receipt and processing of the data.

### Operational Flow
- **Start:** The master listens for incoming connections. The client, upon user command, connects and either sends a file or prepares to receive one.
- **File Processing:** If sending, the file is transmitted to the master, split, distributed to slaves, potentially processed, reassembled, and sent back to the client.
- **End:** The slaves and the master close connections once the operations are completed.

### Use Case
This setup could be useful for scenarios where files need to be processed in parallel, reducing the time for processing large files by utilizing multiple machines. Examples include distributed video processing, data analysis, or any task that benefits from parallel processing to speed up handling large datasets.

This explanation outlines a robust example of how distributed computing and networking concepts are implemented using Go, demonstrating both file handling and network communications effectively. The system architecture facilitates practical learning and exploration of concurrent programming and network interactions in distributed systems.
