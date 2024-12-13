#include <stdio.h>
#include <string.h>
#include <winsock2.h>

#define PORT 8080
#define SERVER "127.0.0.1"  

int main() {
    WSADATA wsaData;
    SOCKET sock;
    struct sockaddr_in server;
    char message[1024];

    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0) {
        printf("WSAStartup failed\n");
        return 1;
    }

    sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock == INVALID_SOCKET) {
        printf("Could not create a socket!\n");
        WSACleanup();
        return 1;
    }

    server.sin_family = AF_INET;
    server.sin_port = htons(PORT);
    server.sin_addr.s_addr = inet_addr(SERVER);

    if (connect(sock, (struct sockaddr *)&server, sizeof(server)) == SOCKET_ERROR) {
        printf("Could not connect!\n");
        closesocket(sock);
        WSACleanup();
        return 1;
    }

    printf("Connected to the server!\n");

    int recv_size;
    if ((recv_size = recv(sock, message, sizeof(message), 0)) == SOCKET_ERROR) {
        printf("Could not get the message!\n");
        closesocket(sock);
        WSACleanup();
        return 1;
    }

    message[recv_size] = '\0'; 
    printf("Server's message: %s\n", message);

    printf("Enter your guess: ");
    fgets(message, sizeof(message), stdin);

    message[strcspn(message, "\n")] = '\0';

    printf("Sending: %s\n", message);

    int send_size = send(sock, message, strlen(message), 0);
    if (send_size == SOCKET_ERROR) {
        printf("Error sending data!\n");
        closesocket(sock);
        WSACleanup();
        return 1;
    }

    closesocket(sock);
    WSACleanup();

    return 0;
}
