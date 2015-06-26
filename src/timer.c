/**
 * This file client.c, handles clients of the WOE server.
 */

#define _POSIX_C_SOURCE 200801L
#define _POSIX_SOURCE 200801L

#include <signal.h>
#include <time.h>

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

struct woe_timer * woe_timer_init(struct woe_timer * me,
  struct woe_server * server, int index) {
  sigevent_t sev;
  if (!me) return NULL;
  me->server  = server;
  me->index   = index;
  me->set     = 0;
  sev.sigev_notify= SIGEV_NONE;
  
  sev.sigev_notify = SIGEV_SIGNAL;
  sev.sigev_signo  = SIGRTMIN + (index % (SIGRTMAX - SIGRTMIN - 1));
  sev.sigev_value.sival_ptr = me;
  
  
  if (timer_create(CLOCK_MONOTONIC, &sev, &me->timer) < 0) { 
    LOG_ERROR("Could not create timer %d.", index);
    return NULL;
  }
  LOG_NOTE("Timer %d initialized: %d.\n", me->index, me->timer);

  return me;
}

struct woe_timer * woe_timer_new(struct woe_server * server, int index) {
  struct woe_timer * me = woe_timer_alloc();
  if (!me) return NULL;  
  if (!woe_timer_init(me, server, index)) { 
    free(me);
    return NULL;
  }
  return me;
}

struct woe_timer * woe_timer_done(struct woe_timer * me) {
  if (!me) return NULL;
  if (me->timer)  timer_delete(me->timer);
  me->timer     = NULL;
  LOG_NOTE("Timer %d destroyed.\n", me->index);
  me->index     = -1;
  return me;
}

struct woe_timer * woe_timer_free(struct woe_timer * me) {
  woe_timer_done(me);
  free(me);  
  return NULL;
}

struct timespec * timespec_init(struct timespec * tv, double sec) {
  if (!tv) return NULL;
  if (sec <= 0.0) {
    tv->tv_sec  = 0;
    tv->tv_nsec = 0;
  } else {
    tv->tv_sec  = floor(sec);
    tv->tv_nsec = (sec - floor(sec)) * 1000000000;
  }
  return tv;
} 

double timespec_to_s(struct timespec * tv) {
  return tv->tv_sec + ( tv->tv_nsec / 1000000000.0 );
}

int woe_timer_get(struct woe_timer * me, double * value, double * interval) {
  int res;
  if (!me) return -1;
  struct itimerspec get;
  /* timespec_init(&nv.it_value, value); */ 
  res = timer_gettime(me->timer, &get);
  if (res == 0) {
    if (value)    (*value)    = timespec_to_s(&get.it_value);
    if (interval) (*interval) = timespec_to_s(&get.it_interval);
  }
  return res;  
}

int woe_timer_set(struct woe_timer * me, double value, double interval) {
   if (!me) return -1;
   me->set = 1;
   struct itimerspec nv, ov;
   timespec_init(&nv.it_value, value);
   timespec_init(&nv.it_interval, interval);
   return timer_settime(me->timer, 0, &nv, &ov);
}

int woe_timer_passed(struct woe_timer * me) {
  double value, interval;
  if (!me) return 0;
  if (!me->set) return 0;
  if (woe_timer_get(me, &value, &interval) != 0) { 
    LOG_ERROR("Trouble getting timer %d\n", me->index);
    return 0;
  }
  return value <= 0;
}

int woe_timer_callback(struct woe_timer * me) {
  mrb_state * mrb;
  LOG_DEBUG("Timer passed: %d\n", me->index);
  mrb = woe_server_get_mrb(me->server);  
  if (mrb) {  
    double value = 0.0, interval = 0.0;
    woe_timer_get(me, &value, &interval); 
    rh_run_toplevel(mrb, "woe_on_timer", "iff", me->index, value, interval);
  }
  return 0;
}


