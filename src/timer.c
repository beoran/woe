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
#include <math.h>

#include "libtelnet.h"
#include "monolog.h"
#include "rh.h"
#include "timer.h"
#include "server.h"



struct woe_timer * woe_timer_alloc() {
  return calloc(sizeof(struct woe_timer), 1);
}

struct woe_client * woe_timer_init(struct woe_timer * me,
  struct woe_server * server, int index) {
  sigevent_t sev;
  if (!me) return NULL;
  me->server  = server;
  me->index   = index;
  sev.sigev_notify= SIGEV_NONE;
  if (timer_create(CLOCK_MONOTONIC, &sev, &me->timer) < 0) { 
    LOG_ERROR("Could not create timer %d.", index)
    return NULL;
  }
  LOG_NOTE("Timer %d initialized: %d.\n", client->index, client->timer);

  return me;
}

struct woe_timer * woe_timer_new(struct woe_server * server, int index) {
  struct woe_timer * me = woe_timer_alloc();
  if (!me) return NULL;  
  if (!woe_timer_init(me, index)) { 
    free(me);
    return NULL;
  }
  return me;
}

struct woe_timer * woe_timer_done(struct woe_timer * me) {
  if (!me) return NULL;
  if (me->timer > -1)  timer_delete(me->timer);
  me->timer     = -1;
  LOG_NOTE("Timer %d destroyed.\n", client->index);
  me->index     = -1;
  return me;
}

struct woe_timer * woe_timer_free(struct woe_timer * me) {
  woe_client_done(client);
  free(client);  
  return NULL;
}

struct timespec * timespec_init(struct timespec * tv, double sec) {
  if !tv) return NULL;
  if (sec <= 0.0) {
    tv->tv_sec  = 0;
    tv->tv_nsec = 0;
  } else {
    tv->tv_sec  = floor(sec);
    tv->tv_nsec = (sec - floor(sec)) * 1000000000;
  }
  return tv;
} 

int woe_timer_set(struct woe_timer * me, double value, double interval) {
  if (!me) return -1;
  itimerspec nv, ov;
  timespec_init(&nv.it_value, value);
  timespec_init(&nv.it_interval, interval);
  return timer_settime(me->timer, 0, &nv, &ov);
}




