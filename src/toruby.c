/*
* toruby.c helps expose functionality from C to Mruby for Eruta.
* All functions are prefixed with tr_.
* Look at the tr_*.c files.
* */

#include "toruby.h"
#include "tr_macro.h"
#include "monolog.h"
#include "rh.h"
#include "state.h"
#include "server.h"

#include <mruby/hash.h>
#include <mruby/class.h>
#include <mruby/data.h>
#include <mruby/array.h>
/*
#include "tr_macro.h"
#include "tr_audio.h"
#include "tr_graph.h"
#include "tr_store.h"
#include "tr_sprite.h"
*/

/* Documentation of mrb_get_args: 
 
  retrieve arguments from mrb_state.

  mrb_get_args(mrb, format, ...)

  returns number of arguments parsed.

  format specifiers:

   o: Object [mrb_value]
   S: String [mrb_value]
   A: Array [mrb_value]
   H: Hash [mrb_value]
   s: String [char*,int]
   z: String [char*] nul terminated
   a: Array [mrb_value*,mrb_int]
   f: Float [mrb_float]
   i: Integer [mrb_int]
   b: Binary [int]
   n: Symbol [mrb_sym]
   &: Block [mrb_value]
   *: rest argument [mrb_value*,int]
   |: optional
 */




/** Writes a NOTE message to the log. */
static mrb_value tr_log(mrb_state * mrb, mrb_value self) {
  (void) self; (void) mrb;
  
  mrb_value text    = mrb_nil_value();
  mrb_get_args(mrb, "S", &text);
  
  LOG_NOTE("%s\n", RSTRING_PTR(text));
  return self;
}

/** Writes a messageto a certain log level log. */
static mrb_value tr_log_to(mrb_state * mrb, mrb_value self) {
  (void) self; (void) mrb;
  
  mrb_value level   = mrb_nil_value();
  mrb_value text    = mrb_nil_value();

  mrb_get_args(mrb, "SS", &level, &text);
  
  LOG_LEVEL(RSTRING_PTR(level), "%s\n", RSTRING_PTR(text));
  return self;
}


/** Cause a warning to be logged */
static mrb_value tr_warn(mrb_state * mrb, mrb_value self) {
  (void) self; (void) mrb;
  
  mrb_value text    = mrb_nil_value();
  mrb_get_args(mrb, "S", &text);
  LOG_WARNING("%s\n", RSTRING_PTR(text));
  return self;
}


/** Enables a certain log level */
static mrb_value tr_log_enable(mrb_state * mrb, mrb_value self) {
  (void) self; (void) mrb;
  
  mrb_value text    = mrb_nil_value();
  mrb_get_args(mrb, "S", &text);
  monolog_enable_level(RSTRING_PTR(text));
  return self;
}

/** Disables a certain log level */
static mrb_value tr_log_disable(mrb_state * mrb, mrb_value self) {
  (void) self; (void) mrb;

  mrb_value text    = mrb_nil_value();
  mrb_get_args(mrb, "S", &text);
  monolog_disable_level(RSTRING_PTR(text));
  return self;
}


/* Loads another script from the script directory. */
static mrb_value tr_script(mrb_state * mrb, mrb_value self) {
  int res; 
  char * command;
  
  (void) self;
  
  mrb_value text        = mrb_nil_value();
  mrb_get_args(mrb, "S", &text);
  command               = mrb_str_to_cstr(mrb, text);
  
  res = rh_run_script(mrb, command);
  return mrb_fixnum_value(res);
}


/* Sends data to a given client */
static mrb_value tr_send_to_client(mrb_state * mrb, mrb_value self) {
  int res; 
  mrb_int client = -1;
  char * data    = NULL;
  int    size    = 0;
  struct woe_server * srv = MRB_WOE_SERVER(mrb);
  
  (void) self;
  mrb_get_args(mrb, "is", &client, &data, &size);
  res = woe_server_send_to_client(srv, client, data, size);
  
  return mrb_fixnum_value(res);
}

/* Shuts down a given client */
static mrb_value tr_server_done(mrb_state * mrb, mrb_value self) {
  struct woe_server * srv = MRB_WOE_SERVER(mrb);
  
  (void) self;
  woe_server_request_shutdown(srv);
  return mrb_nil_value(); 
}


/* Disconnects a given client by id */
static mrb_value tr_disconnect_client(mrb_state * mrb, mrb_value self) {
  int res; 
  mrb_int client;
  struct woe_server * srv = MRB_WOE_SERVER(mrb);
  
  (void) self;
  mrb_get_args(mrb, "i", &client);
  res = woe_server_disconnect_id(srv, client);
  
  return mrb_fixnum_value(res);
}


/* Yeah, I know, but this reduces boilerplate. */
#define WRAP_SERVER_BEGIN(NAME) \
static mrb_value NAME(mrb_state * mrb, mrb_value self) {          \
  int res;                                                        \
  mrb_int client = -1;                                            \
  struct woe_server * srv = MRB_WOE_SERVER(mrb);                  \
  (void) self;                                                    \
  
  
#define WRAP_SERVER_END() \
  return mrb_fixnum_value(res); \
}
    

WRAP_SERVER_BEGIN(tr_server_iac) { 
  mrb_int command = 0;  
  mrb_get_args(mrb, "ii", &client, &command);
  res = woe_server_iac(srv, client, command);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_negotiate) { 
  mrb_int how = 0, option = 0;  
  mrb_get_args(mrb, "iii", &client, &how, &option);
  res = woe_server_negotiate(srv, client, how, option);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_begin_sb) { 
  mrb_int telopt = 0;  
  mrb_get_args(mrb, "ii", &client, &telopt);
  res = woe_server_begin_sb(srv, client, telopt);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_finish_sb) { 
  mrb_get_args(mrb, "i", &client);
  res = woe_server_finish_sb(srv, client);
} WRAP_SERVER_END()


WRAP_SERVER_BEGIN(tr_server_subnegotiation) { 
  mrb_int telopt = 0, size = 0;
  char * data = NULL;  
  mrb_get_args(mrb, "iis", &client, &telopt, &data, &size);
  res = woe_server_subnegotiation(srv, client, telopt, data, size);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_begin_compress2) { 
  mrb_get_args(mrb, "i", &client);
  res = woe_server_begin_compress2(srv, client);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_puts) { 
  char * fmt = NULL;  
  mrb_get_args(mrb, "iz", &client, &fmt);
  res = woe_server_printf(srv, client, fmt);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_raw_puts) { 
  char * fmt = NULL;  
  mrb_get_args(mrb, "iz", &client, &fmt);
  res = woe_server_raw_printf(srv, client, fmt);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_begin_newenviron) { 
  mrb_int type;
  mrb_get_args(mrb, "ii", &client, &type);
  res = woe_server_begin_newenviron(srv, client, type);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_newenviron_value) { 
  mrb_int type;
  char * value = NULL;  
  mrb_get_args(mrb, "iiz", &client, &type, &value);
  res = woe_server_newenviron_value(srv, client, type, value);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_finish_newenviron) { 
  mrb_get_args(mrb, "i", &client);
  res = woe_server_finish_newenviron(srv, client);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_ttype_send) { 
  mrb_get_args(mrb, "i", &client);
  res = woe_server_ttype_send(srv, client);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_ttype_is) { 
  mrb_int type;
  char * ttype = NULL;  
  mrb_get_args(mrb, "iz", &client, &ttype);
  res = woe_server_ttype_is(srv, client, ttype);
} WRAP_SERVER_END()

/*
int woe_server_send_zmp(struct woe_server * srv, int client, int argc, const char ** argv);
int woe_server_send_vzmpv(struct woe_server * srv, int client, va_list va);
int woe_server_send_zmpv(struct woe_server * srv, int client, ...);
*/
WRAP_SERVER_BEGIN(tr_server_begin_zmp) { 
  char * command = NULL;
  mrb_get_args(mrb, "iz", &client, &command);
  res = woe_server_begin_zmp(srv, client, command);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_zmp_arg) { 
  char * arg = NULL;  
  mrb_get_args(mrb, "iz", &client, &arg);
  res = woe_server_zmp_arg(srv, client, arg);
} WRAP_SERVER_END()

WRAP_SERVER_BEGIN(tr_server_finish_zmp) { 
  char * command;
  mrb_get_args(mrb, "iz", &client, &command);
  res = woe_server_finish_zmp(srv, client, command);
} WRAP_SERVER_END()

/* Initializes the functionality that Eruta exposes to Ruby. */
int tr_init(mrb_state * mrb) {
  // luaL_dostring(lua, "print 'Hello!' ");
  struct RClass *woe;
  struct RClass *srv;
  struct RClass *krn;
 
  woe = mrb_define_module(mrb, "Woe");
  srv = mrb_define_module_under(mrb, woe, "Server"); 
  TR_CLASS_METHOD_NOARG(mrb, woe, "quit"  , tr_server_done);
  TR_CLASS_METHOD_NOARG(mrb, srv, "quit"  , tr_server_done);
  TR_CLASS_METHOD_ARGC(mrb, srv, "send_to_client"  , tr_send_to_client, 2);
  TR_CLASS_METHOD_NOARG(mrb, srv, "disconnect"  , tr_disconnect_client);

int woe_server_iac(struct woe_server * srv, int client, int command);
int woe_server_negotiate(struct woe_server * srv, int client, int how, int option);
int woe_server_begin_sb(struct woe_server * srv, int client, int telopt);
int woe_server_finish_sb(struct woe_server * srv, int client);
int woe_server_subnegotiation(struct woe_server * srv, int client, int telopt, char * buffer, int size);
int woe_server_begin_compress2(struct woe_server * srv, int client);
int woe_server_vprintf(struct woe_server * srv, int client, const char *fmt, va_list va);
int woe_server_printf(struct woe_server * srv, int client, const char *fmt, ...);
int woe_server_raw_vprintf(struct woe_server * srv, int client, const char *fmt, va_list va);
int woe_server_raw_printf(struct woe_server * srv, int client, const char *fmt, ...);
int woe_server_begin_newenviron(struct woe_server * srv, int client, int type);
int woe_server_newenviron_value(struct woe_server * srv, int client, int type, char * value);
int woe_server_finish_newenviron(struct woe_server * srv, int client);
int woe_server_ttype_send(struct woe_server * srv, int client);
int woe_server_ttype_is(struct woe_server * srv, int client, char * ttype);
int woe_server_send_zmp(struct woe_server * srv, int client, int argc, const char ** argv);
int woe_server_send_vzmpv(struct woe_server * srv, int client, va_list va);
int woe_server_send_zmpv(struct woe_server * srv, int client, ...);
int woe_server_begin_zmp(struct woe_server * srv, int client, const char * cmd);
int woe_server_zmp_arg(struct woe_server * srv, int client, const char * arg);
int woe_server_finish_zmp(struct woe_server * srv, int client, const char * cmd);


  TR_CLASS_METHOD_ARGC(mrb, srv, "iac"  , tr_server_iac, 2);

  TR_CLASS_METHOD_ARGC(mrb, srv, "negotiate"      , tr_server_negotiate     , 3);
  TR_CLASS_METHOD_ARGC(mrb, srv, "begin_sb"       , tr_server_begin_sb      , 2);
  TR_CLASS_METHOD_ARGC(mrb, srv, "finish_sb"      , tr_server_finish_sb     , 1);
  TR_CLASS_METHOD_ARGC(mrb, srv, "subnegotiation" , tr_server_subnegotiation, 3);
  TR_CLASS_METHOD_ARGC(mrb, srv, "begin_compress2", tr_server_begin_compress2, 2);
  TR_CLASS_METHOD_ARGC(mrb, srv, "puts"           , tr_server_puts, 2);
  TR_CLASS_METHOD_ARGC(mrb, srv, "raw_puts"       , tr_server_raw_puts, 2);
  TR_CLASS_METHOD_ARGC(mrb, srv, "begin_newenviron" , tr_server_begin_newenviron, 2);
  TR_CLASS_METHOD_ARGC(mrb, srv, "newenviron_value" , tr_server_newenviron_value, 3);
  TR_CLASS_METHOD_ARGC(mrb, srv, "finish_newenviron", tr_server_finish_newenviron, 1);
  TR_CLASS_METHOD_ARGC(mrb, srv, "ttype_send"  , tr_server_ttype_send, 1);
  TR_CLASS_METHOD_ARGC(mrb, srv, "ttype_is"    , tr_server_ttype_is, 2);
  TR_CLASS_METHOD_ARGC(mrb, srv, "begin_zmp"  , tr_server_finish_zmp, 2);
  TR_CLASS_METHOD_ARGC(mrb, srv, "zmp_arg"  , tr_server_finish_zmp, 2);
  TR_CLASS_METHOD_ARGC(mrb, srv, "finish_zmp"  , tr_server_finish_zmp, 2);

  
  krn = mrb_module_get(mrb, "Kernel");
  if(!krn) return -1;
  
  TR_METHOD_ARGC(mrb, krn, "warn"         , tr_warn   , 1);
  TR_METHOD_ARGC(mrb, krn, "warning"      , tr_warn   , 1);
  TR_METHOD_ARGC(mrb, krn, "log"          , tr_log    , 1);
  TR_METHOD_ARGC(mrb, krn, "log_to"       , tr_log_to , 2);
  TR_METHOD_ARGC(mrb, krn, "log_enable"   , tr_log_disable , 1);
  TR_METHOD_ARGC(mrb, krn, "log_disable"  , tr_log_enable  , 1);
  TR_METHOD_ARGC(mrb, krn, "script"       , tr_script , 1);

   
  // must restore gc area here ????
  mrb_gc_arena_restore(mrb, 0);
  
  return 0;
}








