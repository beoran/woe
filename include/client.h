#ifndef WOE_CLIENT_H
#define WOE_CLIENT_H

#if !defined(_POSIX_SOURCE)
#define _POSIX_SOURCE
#endif
#if !defined(_BSD_SOURCE)
#define _BSD_SOURCE
#endif

#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <poll.h>
#include <unistd.h>

#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <ctype.h>

#include "libtelnet.h"


#ifndef WOE_CLIENTS_MAX
/* Must beless than ulimit -nor */
#define WOE_CLIENTS_MAX 1000
#endif

#ifndef WOE_CLIENT_BUFFER_SIZE 
#define WOE_CLIENT_BUFFER_SIZE 256
#endif

struct woe_server;

struct woe_client {  
  struct woe_server * server;
  int                 index;
  int                 sock;
  telnet_t *          telnet;
  char                linebuf[256];
  int                 linepos;
  int                 busy;    
  struct sockaddr_in  addr;
  socklen_t           addrlen;
};


int woe_clients_max(void);

int woe_client_input(struct woe_client * cli, const char * buffer, size_t size);
int woe_client_zmp(struct woe_client * cli, int argc, char *argv[]);


struct woe_client * 
woe_client_new(struct woe_server * srv,  int index, int socket, 
               struct sockaddr_in *addr, socklen_t addrlen);

struct woe_client * woe_client_free(struct woe_client * cli);



#endif
