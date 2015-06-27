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
#include "libtelnet.h"
#include "tr_file.h"

#include <signal.h>
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
  

#define WRAP_SERVER_TIMER_BEGIN(NAME) \
static mrb_value NAME(mrb_state * mrb, mrb_value self) {          \
  int res;                                                        \
  mrb_int timer = -1;                                             \
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


WRAP_SERVER_TIMER_BEGIN(tr_server_new_timer) { 
  (void) timer; 
  res = woe_server_make_new_timer_id(srv);
} WRAP_SERVER_END()


WRAP_SERVER_TIMER_BEGIN(tr_server_set_timer) { 
  mrb_float value = 0.0, interval = 0.0;
  mrb_get_args(mrb, "iff", &timer, &value, &interval);
  res = woe_server_set_timer_value(srv, timer, value, interval);
} WRAP_SERVER_END()

static mrb_value tr_server_get_timer(mrb_state * mrb, mrb_value self) {
  int res;                                                        
  mrb_int timer = -1;                                             
  struct woe_server * srv = MRB_WOE_SERVER(mrb);                  
  (void) self;
  double value = 0.0, interval = 0.0;
  mrb_get_args(mrb, "i", &timer);
  res = woe_server_get_timer_value(srv, timer, &value, &interval);
  { 
    mrb_value vals[3] = { mrb_fixnum_value(res), mrb_float_value(mrb, value), mrb_float_value(mrb, interval) };
    return mrb_ary_new_from_values(mrb, 3, vals);
  }
}



/* Initializes the functionality that Eruta exposes to Ruby. */
int tr_init(mrb_state * mrb) {
  // luaL_dostring(lua, "print 'Hello!' ");
  struct RClass *woe;
  struct RClass *srv;
  struct RClass *krn;
  struct RClass *tel;
  struct RClass *sig;
  struct RClass *fil;
  
  tr_init_file(mrb);
 
  woe = mrb_define_module(mrb, "Woe");
  srv = mrb_define_module_under(mrb, woe, "Server"); 
  tel = mrb_define_module(mrb, "Telnet");
  sig = mrb_define_module(mrb, "Signal");
  krn = mrb_module_get(mrb, "Kernel");
  if(!krn) return -1;
  
  TR_CLASS_METHOD_NOARG(mrb, woe, "quit"  , tr_server_done);
  TR_CLASS_METHOD_NOARG(mrb, srv, "quit"  , tr_server_done);
  TR_CLASS_METHOD_ARGC(mrb, srv, "send_to_client"  , tr_send_to_client, 2);
  TR_CLASS_METHOD_NOARG(mrb, srv, "disconnect"  , tr_disconnect_client);
  
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

  TR_CLASS_METHOD_NOARG(mrb, srv, "new_timer"  , tr_server_new_timer);
  TR_CLASS_METHOD_ARGC(mrb, srv, "set_timer"  , tr_server_set_timer, 3);
  TR_CLASS_METHOD_ARGC(mrb, srv, "get_timer"  , tr_server_get_timer, 1);
  
  /* Telnet constants, commands, etc. */
  
  TR_CONST_INT_VALUE(mrb, tel, TELNET_IAC);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_DONT);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_DO);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_WILL);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_WONT);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_SB);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_GA);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_EL);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_EC);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_AYT);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_AO);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_IP);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_BREAK);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_DM);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_NOP);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_SE);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_EOR);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_ABORT);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_SUSP);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_EOF);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_BINARY);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_ECHO);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_RCP);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_SGA);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_NAMS);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_STATUS);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_TM);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_RCTE);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_NAOL);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_NAOP);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_NAOCRD);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_NAOHTS);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_NAOHTD);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_NAOFFD);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_NAOVTS);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_NAOVTD);
  
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_NAOLFD);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_XASCII);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_LOGOUT);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_BM);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_DET);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_SUPDUP);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_SUPDUPOUTPUT);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_SNDLOC);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_TTYPE);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_EOR);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_3270REGIME);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_X3PAD);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_NAWS);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_TSPEED);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_LFLOW);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_LINEMODE);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_XDISPLOC);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_ENVIRON);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_AUTHENTICATION);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_ENCRYPT);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_NEW_ENVIRON);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_MSSP);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_COMPRESS);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_COMPRESS2);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_ZMP);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_EXOPL);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TELOPT_MCCP2);

  TR_CONST_INT_VALUE(mrb, tel, TELNET_TTYPE_IS);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_TTYPE_SEND);

  TR_CONST_INT_VALUE(mrb, tel, TELNET_ENVIRON_IS);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_ENVIRON_SEND);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_ENVIRON_INFO);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_ENVIRON_VAR);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_ENVIRON_VALUE);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_ENVIRON_ESC);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_ENVIRON_USERVAR);

  TR_CONST_INT_VALUE(mrb, tel, TELNET_MSSP_VAL);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_MSSP_VAR);

  TR_CONST_INT_VALUE(mrb, tel, TELNET_EOK);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_EBADVAL);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_ENOMEM);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_EOVERFLOW);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_EPROTOCOL);
  TR_CONST_INT_VALUE(mrb, tel, TELNET_ECOMPRESS);

  TR_CONST_INT_VALUE(mrb, sig, SIGHUP);
  TR_CONST_INT_VALUE(mrb, sig, SIGINT);
  TR_CONST_INT_VALUE(mrb, sig, SIGQUIT);
  TR_CONST_INT_VALUE(mrb, sig, SIGILL);
  TR_CONST_INT_VALUE(mrb, sig, SIGABRT);
  TR_CONST_INT_VALUE(mrb, sig, SIGFPE);
  TR_CONST_INT_VALUE(mrb, sig, SIGKILL);
  TR_CONST_INT_VALUE(mrb, sig, SIGSEGV);
  TR_CONST_INT_VALUE(mrb, sig, SIGPIPE);
  TR_CONST_INT_VALUE(mrb, sig, SIGALRM);
  TR_CONST_INT_VALUE(mrb, sig, SIGTERM);
  TR_CONST_INT_VALUE(mrb, sig, SIGUSR1);
  TR_CONST_INT_VALUE(mrb, sig, SIGUSR2);
  TR_CONST_INT_VALUE(mrb, sig, SIGCHLD);
  TR_CONST_INT_VALUE(mrb, sig, SIGCONT);
  TR_CONST_INT_VALUE(mrb, sig, SIGSTOP);
  TR_CONST_INT_VALUE(mrb, sig, SIGTSTP);
  TR_CONST_INT_VALUE(mrb, sig, SIGTTIN);
  TR_CONST_INT_VALUE(mrb, sig, SIGTTOU);
  TR_CONST_INT_VALUE(mrb, sig, SIGBUS);
  TR_CONST_INT_VALUE(mrb, sig, SIGPOLL);
  TR_CONST_INT_VALUE(mrb, sig, SIGPROF);
  TR_CONST_INT_VALUE(mrb, sig, SIGSYS);
  TR_CONST_INT_VALUE(mrb, sig, SIGTRAP);
  TR_CONST_INT_VALUE(mrb, sig, SIGURG);
  TR_CONST_INT_VALUE(mrb, sig, SIGVTALRM);
  TR_CONST_INT_VALUE(mrb, sig, SIGXCPU);  
  TR_CONST_INT_VALUE(mrb, sig, SIGXFSZ);
  TR_CONST_INT_VALUE(mrb, sig, SIGIOT);
  TR_CONST_INT_VALUE(mrb, sig, SIGSTKFLT);
  TR_CONST_INT_VALUE(mrb, sig, SIGIO);
  TR_CONST_INT_VALUE(mrb, sig, SIGCLD);
  TR_CONST_INT_VALUE(mrb, sig, SIGPWR);
  TR_CONST_INT_VALUE(mrb, sig, SIGWINCH);
  TR_CONST_INT_VALUE(mrb, sig, SIGUNUSED);

  
  
  TR_METHOD_ARGC(mrb, krn, "woe_warn"     , tr_warn   , 1);
  TR_METHOD_ARGC(mrb, krn, "woe_warning"  , tr_warn   , 1);
  TR_METHOD_ARGC(mrb, krn, "woe_log"      , tr_log    , 1);
  TR_METHOD_ARGC(mrb, krn, "woe_log_to"   , tr_log_to , 2);
  TR_METHOD_ARGC(mrb, krn, "log_enable"   , tr_log_disable , 1);
  TR_METHOD_ARGC(mrb, krn, "log_disable"  , tr_log_enable  , 1);
  TR_METHOD_ARGC(mrb, krn, "script"       , tr_script , 1);

   
  // must restore gc area here ????
  mrb_gc_arena_restore(mrb, 0);
  
  return 0;
}








