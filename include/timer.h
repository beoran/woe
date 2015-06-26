#ifndef WOE_TIMER_H
#define WOE_TIMER_H

#if !defined(_POSIX_C_SOURCE)
#define _POSIX_C_SOURCE 200801L
#endif

#if !defined(_POSIX_SOURCE)
#define _POSIX_SOURCE
#endif
#if !defined(_BSD_SOURCE)
#define _BSD_SOURCE
#endif

#include <signal.h>
#include <time.h>

#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <ctype.h>


struct woe_server;

struct woe_timer {  
  struct woe_server * server;
  int                 index;
  timer_t             timer;
  int                 set;
  struct sigaction    sa;
};



struct woe_timer * woe_timer_new(struct woe_server * server, int index);
struct woe_timer * woe_timer_free(struct woe_timer * me);
int woe_timer_get(struct woe_timer * me, double * value, double * interval);
int woe_timer_set(struct woe_timer * me, double value, double interval);
int woe_timer_passed(struct woe_timer * me);
int woe_timer_callback(struct woe_timer * me);

#endif
