#ifndef STATE_H_INCLUDED
#define STATE_H_INCLUDED

#include "rh.h"
#include "server.h"
#include "config.h"

/* All state of WOE in a handy struct */
struct woe_state {
  struct woe_server * server;
  struct woe_config * config;
  mrb_state         * mrb;
};


#define MRB_WOE_STATE(MRB) ((struct woe_state *)(MRB->ud))
#define MRB_WOE_SERVER(MRB) (MRB_WOE_STATE(MRB)->server)
#define MRB_WOE_CONFIG(MRB) (MRB_WOE_STATE(MRB)->config)


#endif

