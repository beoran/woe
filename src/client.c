/**
 * This file client.c, handles clients of the WOE server.
 */

#if !defined(_WIN32)
# if !defined(_POSIX_SOURCE)
#   define _POSIX_SOURCE
# endif
# if !defined(_BSD_SOURCE)
#   define _BSD_SOURCE
# endif

# include <sys/socket.h>
# include <netinet/in.h>
# include <arpa/inet.h>
# include <netdb.h>
# include <poll.h>
# include <unistd.h>
#else
# include <winsock2.h>
# include <ws2tcpip.h>

# define snprintf _snprintf
# define poll WSAPoll
# define close closesocket
# define strdup _strdup
# define ECONNRESET WSAECONNRESET
#endif

#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <ctype.h>


#include "libtelnet.h"
#include "monolog.h"
#include "rh.h"
#include "client.h"
#include "server.h"



struct woe_client * woe_client_alloc() {
  return calloc(sizeof(struct woe_client), 1);
}

struct woe_client * woe_client_init(struct woe_client * client,
  struct woe_server * server, int index, int socket, 
  struct sockaddr_in *addr, socklen_t addrlen) {
  if (!client) return NULL;
  client->server  = server;
  client->sock    = socket;
  client->addr    = *addr;
  client->addrlen = addrlen;
  client->busy    = !0;
  return client;
}

struct woe_client * woe_client_new(struct woe_server * server, int index, int socket, 
  struct sockaddr_in *addr, socklen_t addrlen) {
  struct woe_client * client = woe_client_alloc();  
  return woe_client_init(client, server, index, socket, addr, addrlen);
}

struct woe_client * woe_client_done(struct woe_client * client) {
  /* Do nothing yet, refactor later. */
  if (!client) return NULL;
  if (client->telnet) telnet_free(client->telnet);
  client->sock   = -1;
  client->telnet = NULL;
  LOG_NOTE("Connection to client %d closed.\n", client->index);
  client->index = -1;
  return client;
}


struct woe_client * woe_client_free(struct woe_client * client) {
  woe_client_done(client);
  free(client);  
  return NULL;
}


int woe_client_send(struct woe_client * client, const char *buffer, unsigned int size) {
  int res;
  if (!client) return 1;
  /* ignore on invalid socket */
  if (client->sock < 0)  return 2;
  
  /* send data */
  while (size > 0) {
    if ((res = send(client->sock, buffer, size, 0)) == -1) {
      if (errno != EINTR && errno != ECONNRESET) {
        LOG_ERROR("send() failed: %s\n", strerror(errno));
        return 3;
      } else {
        return 0;
      }
    } else if (res == 0) {
      LOG_ERROR("send() unexpectedly returned 0\n");
      return 4;
    }

    /* Update pointer and size to see if we've got more to send */
    buffer += res;
    size   -= res;
  }
  
  return 0;
}

int woe_client_input(struct woe_client * cli, const char *buffer, size_t size) {
  unsigned int i;
  mrb_state * mrb;
  LOG_NOTE("Received input for client %d");
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_on_input", "is", cli->index, buffer, size);
  }
  return 0;
}

int woe_client_zmp(struct woe_client * cli, int argc, char *argv[]) {
  unsigned int i;
  mrb_state * mrb;
  LOG_NOTE("Received ZMP reply for client %d");
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_begin_zmp", "ii", cli->index, argc);
    for (i=0; i < argc; i++) {
      rh_run_toplevel(mrb, "woe_zmp_arg", "is", cli->index, argv[i]);
    }  
    rh_run_toplevel(mrb, "woe_finish_zmp", "ii", cli->index, argc);
  }
  return 0;
}




