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
  mrb_state * mrb;
  LOG_NOTE("Received input for client %d\n", cli->index);
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_on_input", "is", cli->index, buffer, size);
  }
  return 0;
}

int woe_client_zmp(struct woe_client * cli, int argc, const char *argv[]) {
  unsigned int i;
  mrb_state * mrb;
  LOG_NOTE("Received ZMP reply for client %d\n", cli->index);
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_begin_zmp", "ii", cli->index, argc);
    for (i=0; i < argc; i++) {
      rh_run_toplevel(mrb, "woe_zmp_arg", "iiz", cli->index, i, argv[i]);
    }  
    rh_run_toplevel(mrb, "woe_finish_zmp", "ii", cli->index, argc);
  }
  return 0;
}

int woe_client_iac(struct woe_client * cli, int cmd) {
  mrb_state * mrb;
  LOG_NOTE("Received iac for client %d %d\n", cli->index, cmd);
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_on_iac", "i", cli->index, cmd);
  }
  return 0;
}


int woe_client_negotiate(struct woe_client * cli, int how, int option) {
  mrb_state * mrb;
  LOG_NOTE("Received negotiate for client %d %d %d\n", cli->index, how, option);
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_on_negotiate", "iii", cli->index, how, option);
  }
  return 0;
}

int woe_client_subnegotiate(struct woe_client * cli, const char * buf, int len, int telopt) {
  mrb_state * mrb;
  LOG_NOTE("Received subnegotiate for client %d\n", cli->index);
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_on_subnegotiate", "iis", cli->index, telopt, buf, len);
  }
  return 0;
}

int woe_client_ttype(struct woe_client * cli, int cmd, const char * name) {
  mrb_state * mrb;
  LOG_NOTE("Received ttype for client %d %d %s\n", cli->index, cmd, name);
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_on_ttype", "iz", cli->index, cmd, name);
  }
  return 0;
}


int woe_client_error(struct woe_client * cli, int code, const char * msg) {
  mrb_state * mrb;
  LOG_NOTE("Received error for client %d %d %s\n", cli->index, code, msg);
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_on_error", "iz\n", cli->index, code, msg);
  }
  return 0;
}


int woe_client_warning(struct woe_client * cli, int code, const char * msg) {
  mrb_state * mrb;
  LOG_NOTE("Received warning for client %d %d %s\n", cli->index, code, msg);
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_on_warning", "iz", cli->index, code, msg);
  }
  return 0;
}


int woe_client_compress(struct woe_client * cli, int state) {
  mrb_state * mrb;
  LOG_NOTE("Received compress for client %d %d\n", cli->index, state);
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_on_compress", "ii", cli->index, state);
  }
  return 0;
}

int woe_client_environ(struct woe_client * cli, int cmd, const struct telnet_environ_t *values, size_t size) {
  int i;
  mrb_state * mrb;
  LOG_NOTE("Received environ for client %d %d\n", cli->index, cmd);
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_begin_environ", "iii", cli->index, cmd, size);
    for (i=0; i < size; i++) {
      rh_run_toplevel(mrb, "woe_environ_arg", "iiizz", cli->index, i, values[i].type, values[i].var, values[i].value);
    }  
   rh_run_toplevel(mrb, "woe_finish_environ", "iii", cli->index, cmd, size);
  }
  return 0;
}

int woe_client_mssp(struct woe_client * cli, const struct telnet_environ_t *values, size_t size) {
  int i;
  mrb_state * mrb;
  LOG_NOTE("Received mssp for client %d\n", cli->index);
  mrb = woe_server_get_mrb(cli->server);  
  if (mrb) {  
    rh_run_toplevel(mrb, "woe_begin_mssp", "ii", cli->index, size);
    for (i=0; i < size; i++) {
      rh_run_toplevel(mrb, "woe_mssp_arg", "iiizz", cli->index, i, values[i].type, values[i].var, values[i].value);
    }  
   rh_run_toplevel(mrb, "woe_finish_mssp", "iii", cli->index, size);
  }
  return 0;
}

