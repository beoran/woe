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


#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <ctype.h>

#include "client.h"
#include "server.h"
#include "monolog.h"

#include "libtelnet.h"

#include "rh.h"


static const telnet_telopt_t woe_telnet_opts[] = {

  { TELNET_TELOPT_COMPRESS2,  TELNET_WILL, TELNET_DONT },
  { TELNET_TELOPT_ECHO     ,  TELNET_WONT, TELNET_DONT },

  { -1, 0, 0 }
};


#define WOE_SERVER_BUFFER_SIZE      10000
#define WOE_SERVER_LISTEN_BACKLOG   8

struct woe_server {
  char                buffer[WOE_SERVER_BUFFER_SIZE];
  short               listen_port;
  int                 listen_sock;
  int                 busy;
  struct sockaddr_in  addr;
  socklen_t           addrlen;
  mrb_state         * mrb;
  
  struct pollfd       pfd[WOE_CLIENTS_MAX + 1];
  struct woe_client * clients[WOE_CLIENTS_MAX];  
  void (*event_handler)      (telnet_t *telnet, telnet_event_t *ev, void *user_data);
  void (*disconnect_handler) (struct woe_server * srv, struct woe_client * cli, void *user_data);
};



int woe_server_disconnect(struct woe_server * srv, struct woe_client * client);


int woe_send(int sock, const char *buffer, unsigned int size) {
  int res;

  if (sock == -1)
    return -1;

  /* send data */
  while (size > 0) {
    if ((res = send(sock, buffer, size, 0)) == -1) {
      if (errno != EINTR && errno != ECONNRESET) {
        LOG_ERROR("send() failed: %s\n", strerror(errno));
        return -2;
      } else {
        return 0;
      }
    } else if (res == 0) {
      LOG_ERROR("send() unexpectedly returned 0\n");
      return -3;
    }

    /* update pointer and size to see if we've got more to send */
    buffer += res;
    size -= res;
  }
  return 0;
}



static void woe_event_handler(telnet_t *telnet, telnet_event_t *ev, void *user_data) {
  struct woe_client * client = (struct woe_client *) user_data;

  switch (ev->type) {
  /* data received */
  case TELNET_EV_DATA:
    woe_client_input(client, ev->data.buffer, ev->data.size);
    /* telnet_negotiate(telnet, TELNET_WONT, TELNET_TELOPT_ECHO);
      telnet_negotiate(telnet, TELNET_WILL, TELNET_TELOPT_ECHO); */
    break;
    
  /* data must be sent */
  case TELNET_EV_SEND:
    woe_send(client->sock, ev->data.buffer, ev->data.size);
    break;
        
  /* enable compress2 if accepted by client */
  case TELNET_EV_DO:
    if (ev->neg.telopt == TELNET_TELOPT_COMPRESS2) telnet_begin_compress2(telnet);
    woe_client_negotiate(client, TELNET_DO, ev->neg.telopt);
    break;


  case TELNET_EV_DONT:
    woe_client_negotiate(client, TELNET_DONT, ev->neg.telopt);
    break;
    
  case TELNET_EV_WILL:  
    woe_client_negotiate(client, TELNET_WILL, ev->neg.telopt);
    break;
    
  case TELNET_EV_WONT:  
    woe_client_negotiate(client, TELNET_WONT, ev->neg.telopt);
    break;
    
  case TELNET_EV_IAC:
    woe_client_iac(client, ev->iac.cmd);
    break;
    
  case TELNET_EV_SUBNEGOTIATION:
    woe_client_subnegotiate(client, ev->sub.buffer, ev->sub.size, ev->sub.telopt);
    break; 

  case TELNET_EV_TTYPE:
    woe_client_ttype(client, ev->ttype.cmd, ev->ttype.name);
    break;
    
  case TELNET_EV_COMPRESS: 
    woe_client_compress(client, ev->compress.state);
    break;
  
  case TELNET_EV_ENVIRON:
    woe_client_environ(client, ev->environ.cmd, ev->environ.values, ev->environ.size);
    break;
  
  case TELNET_EV_MSSP:
    woe_client_mssp(client, ev->mssp.values, ev->mssp.size);
    break; 

  /* warning */
  case TELNET_EV_WARNING:  
    LOG_WARNING("Telnet warning for client %d %s.\n", client->index, ev->error.msg);    
    woe_client_warning(client, ev->error.errcode, ev->error.msg);
    break;
  
  /* error */
  case TELNET_EV_ERROR:  
    LOG_ERROR("Telnet error for client %d %s.\n", client->index, ev->error.msg);    
    woe_client_error(client, ev->error.errcode, ev->error.msg);
    woe_server_disconnect(client->server, client);
    break;
  
  case TELNET_EV_ZMP:
    woe_client_zmp(client, ev->zmp.argc, ev->zmp.argv);
    break;
     
  default:
     LOG_NOTE("Ignored telnet event %d.\n", ev->type);
    /* ignore */
    break;
  }
}


void woe_server_request_shutdown(struct woe_server * srv) {  
  if (srv) srv->busy = 0;
}

int woe_server_busy(struct woe_server * srv) {
  if (!srv) return 0;
  return srv->busy;
}


struct woe_server * woe_server_free(struct woe_server * srv) {
  close(srv->listen_sock);
  free(srv);
  return NULL;
}



struct woe_server * woe_server_init(struct woe_server * srv, int port) {
  int index;
  if (!srv) return NULL;
  srv->listen_sock = -1;
  srv->listen_port = port;
  srv->busy        = !0;
  srv->mrb         = NULL;
  memset(srv->buffer  , '\0'  , sizeof(srv->buffer));
  memset(srv->pfd     , 0     , sizeof(srv->pfd));
  memset(&srv->addr   , 0     , sizeof(srv->addr));
  
  for (index = 0; index < WOE_CLIENTS_MAX; ++index) {
     srv->clients[index] = NULL;
  }
  
   srv->event_handler = woe_event_handler;
  return srv;
}

struct woe_server * woe_server_new(int port) {
  struct woe_server * srv = calloc(1, sizeof(struct woe_server));
  if (!srv) return NULL;  
  return woe_server_init(srv, port);
}


static const telnet_telopt_t telopts[] = {
  { TELNET_TELOPT_COMPRESS2,  TELNET_WILL, TELNET_DONT },
  { -1, 0, 0 }
};


struct woe_server * woe_server_set_mrb(struct woe_server * srv, mrb_state * mrb) {
  if (!srv) return NULL;  
  srv->mrb = mrb;
  return srv;
}

mrb_state * woe_server_get_mrb(struct woe_server * srv) {
  if (!srv) return NULL;
  return srv->mrb;
}


/** Sets up the server to listen at its configured port. */
int woe_server_listen(struct woe_server * srv) {
  int         res;
  /* create listening socket */
  if ((srv->listen_sock = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
    LOG_ERROR("socket() failed: %s\n", strerror(errno));
    return 1;
  }

  /* Reuse address option */
  res = 1;
  if (setsockopt(srv->listen_sock, SOL_SOCKET, SO_REUSEADDR, (void*)&res, sizeof(res))
     == -1) {
    LOG_ERROR("setsockopt() failed: %s\n", strerror(errno));
    return 2;
  }
     
  /* Bind to listening addr/port */  
  srv->addr.sin_family       = AF_INET;
  srv->addr.sin_addr.s_addr  = INADDR_ANY;
  srv->addr.sin_port         = htons(srv->listen_port);
  if (bind(srv->listen_sock, (struct sockaddr *)&srv->addr, sizeof(srv->addr)) == -1) {
    LOG_ERROR("setsockopt() failed: %s\n", strerror(errno));
    return 3;
  }

  /* listen for clients */
  if (listen(srv->listen_sock, WOE_SERVER_LISTEN_BACKLOG) == -1) {
    fprintf(stderr, "listen() failed: %s\n", strerror(errno));
    return 4;
  }

  LOG_NOTE("Listening on port %d\n", srv->listen_port);

  /* Initialize listening descriptors */
  srv->pfd[WOE_CLIENTS_MAX].fd = srv->listen_sock;
  srv->pfd[WOE_CLIENTS_MAX].events = POLLIN;
  
  return 0;
}

/** Returns one of the clients of the server or NULL if not in use or out of 
 * range */
struct woe_client * woe_server_get_client(struct woe_server * srv, int index) {
  if (!srv)                           return NULL;
  if (index < 0)                      return NULL;
  if (index >= WOE_CLIENTS_MAX)       return NULL;
  return srv->clients[index];
}

/** Stores a client of the server at the given index*/
struct woe_client * woe_server_put_client(struct woe_server * srv, int index, struct woe_client * cli) {
  if (!srv)                           return NULL;
  if (!cli)                           return NULL;
  if (index < 0)                      return NULL;
  if (index >= WOE_CLIENTS_MAX)       return NULL;
  if (srv->clients[index]) {
    woe_client_free(srv->clients[index]);
  }  
  srv->clients[index] = cli;
  return cli; 
}

/** Removes a client of the server at the given index*/
struct woe_client * woe_server_remove_client(struct woe_server * srv, int index) {
  if (!srv)                           return NULL;
  if (index < 0)                      return NULL;
  if (index >= WOE_CLIENTS_MAX)       return NULL;
  if (srv->clients[index]) {
    woe_client_free(srv->clients[index]);
  }  
  srv->clients[index] = NULL;
  return NULL; 
}


/** Find an index to put a new user and returns a pointer to it. 
 *  Returns -1 if no free space is available. 
 **/
int woe_server_get_available_client_index(struct woe_server * srv) {
   int i;
   for (i = 0; i < WOE_CLIENTS_MAX; ++i) {
     struct woe_client * client = woe_server_get_client(srv, i);
     if (!client) {
       return i;
     } 
   }
   return -1;
}

/** Creates a new client for this server. Return null if no memory or no space 
 * for a new client. */
struct woe_client * woe_server_make_new_client(struct woe_server * srv, 
  int socket, struct sockaddr_in *addr, socklen_t addrlen) {
   struct woe_client * new;
   int index = woe_server_get_available_client_index(srv);
   if (index < 0) return NULL;
   new = woe_client_new(srv, index, socket, addr, addrlen);
   return woe_server_put_client(srv, index, new);
} 



/* Handles a new connection to this server. */
int woe_server_handle_connect(struct woe_server * srv) {
  struct woe_client * client;
  int res;
  struct sockaddr_in addr;
  socklen_t addrlen;
 
  /* accept the connection */
  addrlen = sizeof(addr);
  if ((res = accept(srv->listen_sock, (struct sockaddr *)&addr, &addrlen)) == -1) {
    LOG_ERROR("accept() failed: %s\n", strerror(errno));
    return 1;
  }
  
  LOG_NOTE("Connection received.\n");
  
  client = woe_server_make_new_client(srv, res, &addr, addrlen);
  
  /* No space for a new client */
  if (!client) {
    LOG_WARNING("Connection rejected (too many users or OOM)\n");
    write(res, "Too many users.\r\n", 14);
    close(res);  
    return 2;
  }

  /* init, welcome */
  client->telnet  = telnet_init(woe_telnet_opts, srv->event_handler, 0, client);
  if (!client->telnet) {
    LOG_ERROR("Could not initialize telnet connection for user.");
    woe_server_disconnect(srv, client);
    return 3;    
  }
  /*telnet_negotiate(client->telnet, TELNET_DO, TELNET_TELOPT_ECHO);*/
  telnet_negotiate(client->telnet, TELNET_WILL, TELNET_TELOPT_COMPRESS2);
  telnet_printf(client->telnet, "Welcome to WOE!\r\n");
  /* telnet_negotiate(client->telnet, TELNET_WILL, TELNET_TELOPT_ECHO); */
  
  if (srv->mrb) {
    rh_run_toplevel(srv->mrb, "woe_on_connect", "i", client->index);
  }
  
  return 0;
}


/** Sends a telnet command to the numbered client. */
int woe_server_iac(struct woe_server * srv, int client, int command) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_iac(pclient->telnet, command);
  return 0;
}


/** Send a telnet negotiation to the numbered client. */
int woe_server_negotiate(struct woe_server * srv, int client, int how, int option) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_negotiate(pclient->telnet, how, option);
  return 0;
}

/** Sends a telnet start of subnegotiation to the numbered client. */
int woe_server_begin_sb(struct woe_server * srv, int client, int telopt) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_begin_sb(pclient->telnet, telopt);
  return 0;
}

/** Sends a telnet end of subnegotiation to the numbered client. */
int woe_server_finish_sb(struct woe_server * srv, int client) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_finish_sb(pclient->telnet);
  return 0;
}


/** Sends a complete telnet subnegotiation  buffer to the numbered client. */
int woe_server_subnegotiation(struct woe_server * srv, int client, int telopt, char * buffer, int size) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_subnegotiation(pclient->telnet, telopt, buffer, size);
  return 0;
}

/** Begin sending compressed data to the to the numbered client. */
int woe_server_begin_compress2(struct woe_server * srv, int client) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_begin_compress2(pclient->telnet);
  return 0;
}


/** Send formated output with newline escaping to the to the numbered client. */
int woe_server_vprintf(struct woe_server * srv, int client, const char *fmt, va_list va) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_vprintf(pclient->telnet, fmt, va);  
  return 0;
}

/** Send formated output with newline escaping to the to the numbered client. */
int woe_server_printf(struct woe_server * srv, int client, const char *fmt, ...) {
  va_list va;
  int res;
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  va_start(va, fmt);
  telnet_vprintf(pclient->telnet, fmt, va);
  va_end(va);
  return 0;
}

/** Send formated output without newline escaping to the to the numbered client. */
int woe_server_raw_vprintf(struct woe_server * srv, int client, const char *fmt, va_list va) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_raw_vprintf(pclient->telnet, fmt, va);
  return 0;  
}

/** Send formated output without newline escaping to the to the numbered client. */
int woe_server_raw_printf(struct woe_server * srv, int client, const char *fmt, ...) {
  va_list va;
  int res;
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  va_start(va, fmt);
  res = telnet_raw_vprintf(pclient->telnet, fmt, va);
  va_end(va);
  return 0;
}


/** Begin a NEW-ENVIRON subnegotiation with the numbered client. */
int woe_server_begin_newenviron(struct woe_server * srv, int client, int type) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_begin_newenviron(pclient->telnet, type);
  return 0;
}

/** Send a NEW-ENVIRON variable name or value to the numbered client. */
int woe_server_newenviron_value(struct woe_server * srv, int client, int type, char * value) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_newenviron_value(pclient->telnet, type, value);
  return 0;
}

/** Finish a NEW-ENVIRON subnegotiation with the numbered client. */
int woe_server_finish_newenviron(struct woe_server * srv, int client) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_finish_newenviron(pclient->telnet);
  return 0;
}

/** Send a TERMINAL-TYPE SEND command to the numbered client. */
int woe_server_ttype_send(struct woe_server * srv, int client) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_ttype_send(pclient->telnet);
  return 0;
}

/** Send a TERMINAL-TYPE IS command to the numbered client. */
int woe_server_ttype_is(struct woe_server * srv, int client, char * ttype) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_ttype_is(pclient->telnet, ttype);
  return 0;
}


/** Send a ZMP command to the numbered client. */
int woe_server_send_zmp(struct woe_server * srv, int client, int argc, const char ** argv) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_send_zmp(pclient->telnet, argc, argv);
  return 0;
}

/** Send a ZMP command to the numbered client. */
int woe_server_send_vzmpv(struct woe_server * srv, int client, va_list va) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_send_vzmpv(pclient->telnet, va);
  return 0;
}

/** Send a ZMP command to the numbered client. */
int woe_server_send_zmpv(struct woe_server * srv, int client, ...) {
  va_list va;
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  va_start(va, client);
  telnet_send_vzmpv(pclient->telnet, va);
  va_end(va);
  return 0;
}

/** Begin sending a ZMP command to the numbered client. */
int woe_server_begin_zmp(struct woe_server * srv, int client, const char * cmd) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_begin_zmp(pclient->telnet, cmd);
  return 0;
}


/** Send a ZMP command argument to the numbered client. */
int woe_server_zmp_arg(struct woe_server * srv, int client, const char * arg) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_zmp_arg(pclient->telnet, arg);
  return 0;
}

/** Finish sending a ZMP command to the numbered client. */
int woe_server_finish_zmp(struct woe_server * srv, int client, const char * cmd) {
  struct woe_client * pclient;
  if (!srv)     return -1;
  pclient = woe_server_get_client(srv, client);
  if (!pclient) return -2;
  telnet_finish_zmp(pclient->telnet);
  return 0;
}


/** Disconnect a client from the server. */
int woe_server_disconnect(struct woe_server * srv, struct woe_client * client) {
  int index;
  
  if (!srv)     return 1;
  if (!client)  return 2; 
  close(client->sock);
  if (srv->disconnect_handler) {
    srv->disconnect_handler(srv, client, NULL);
  }
  index = client->index;
  
  if (srv->mrb) {
    rh_run_toplevel(srv->mrb, "woe_on_disconnect", "i", index);
  }
    
  /* Get rid of client, will also free memory asociated. */
  woe_server_remove_client(srv, index);  
  return 0;
}


/** Forcfullly disconnect a client from the server by id. 
 * Set a quit flag that woe_server_update will check. */
int woe_server_disconnect_id(struct woe_server * srv, int id) {
  struct woe_client * client = woe_server_get_client(srv, id);
  if (!client) return -1;
  client->busy = 0;
  return 0;
}

/** Polls the server once and updates any of the clients if needed. */
int woe_server_update(struct woe_server * srv, int timeout) {
  int i, res;

  /* prepare for poll */
  memset(srv->pfd     , 0     , sizeof(srv->pfd));
  
  
  for (i = 0; i != WOE_CLIENTS_MAX; ++i) {
   struct woe_client * client = woe_server_get_client(srv, i);
   if (client) {
     srv->pfd[i].fd     = client->sock;
     srv->pfd[i].events = POLLIN;
   } else {
     srv->pfd[i].fd     = -1;
     srv->pfd[i].events = 0;
   }
  }
  
  /* Also listen for connnect events. */
  srv->pfd[WOE_CLIENTS_MAX].fd = srv->listen_sock;
  srv->pfd[WOE_CLIENTS_MAX].events = POLLIN;
  
  

  /* Poll for activity */
  res = poll(srv->pfd, WOE_CLIENTS_MAX + 1, timeout);

  /* Check for time out */
  if (res == 0) {
   /* Time out but that's OK. */
   return 0;
  }

  /* Log errors. */
  if (res == -1 && errno != EINTR) {
    LOG_ERROR("poll() failed: %s\n", strerror(errno));
    return 1;
  }

  /* Handle new connection connection */
  if (srv->pfd[WOE_CLIENTS_MAX].revents & POLLIN) {
    woe_server_handle_connect(srv);
  }

   /* Read from clients */
  for (i = 0; i < WOE_CLIENTS_MAX; ++i) {
    struct woe_client * client = woe_server_get_client(srv, i);
    if (!client) continue;
    
    /* Input from clients. */
    if (srv->pfd[i].revents & POLLIN) {
      res = recv(client->sock, srv->buffer, sizeof(srv->buffer), 0);
      if (res < 0) {
        LOG_ERROR("recv(client) failed: %s\n", strerror(errno));
      } else if (res == 0) {
        /* Disconnect the client. */
        woe_server_disconnect(srv, client);
      } else {
        /* Let telnet lib process incoming data. */
        telnet_recv(client->telnet, srv->buffer, res);
        // telnet_send(client->telnet, srv->buffer, res);
        // telnet_send(telnet, ev->data.buffer, ev->data.size);
      }
    }
  }
  
  /* Disconnect clients that should quit */
  for (i = 0; i < WOE_CLIENTS_MAX; ++i) {
    struct woe_client * client = woe_server_get_client(srv, i);
    if (!client) continue;
    if (!client->busy) {
      woe_server_disconnect(srv, client);
    }  
  }
  
  return 0;
}


int woe_server_send_to_client(struct woe_server * srv, int client, char * data, size_t size) {
  struct woe_client * pclient = woe_server_get_client(srv, client);
  if (!pclient) return -1;
  telnet_send(pclient->telnet, data, size);
  return size;
}



