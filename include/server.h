#ifndef WOE_SERVER_H
#define WOE_SERVER_H

#ifndef WOE_TIMERS_MAX
#define WOE_TIMERS_MAX 32
#endif

#include <stdarg.h>
#include <mruby.h>

struct woe_server;

struct woe_server * woe_server_new(int port);
struct woe_server * woe_server_free(struct woe_server * srv);

int woe_server_listen(struct woe_server * srv);
int woe_server_update(struct woe_server * srv, int timeout);

int woe_server_busy(struct woe_server * srv);
void woe_server_request_shutdown(struct woe_server * srv);

enum woe_timer_type {
  WOE_TIMER_TYPE_REPEAT,
  WOE_TIMER_TYPE_ONCE
};



int woe_server_send_to_client(struct woe_server * srv, int client, char * data, size_t size);                           

struct woe_server * woe_server_set_mrb(struct woe_server * srv, mrb_state * mrb);
mrb_state * woe_server_get_mrb(struct woe_server * srv);

int woe_server_disconnect_id(struct woe_server * srv, int id);


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

int woe_server_make_new_timer_id(struct woe_server * srv);
struct woe_client * woe_server_remove_client(struct woe_server * srv, int index);
struct woe_timer * woe_server_remove_timer(struct woe_server * srv, int index);
int woe_server_set_timer_value(struct woe_server * srv, int index, double value, double interval);
int woe_server_get_timer_value(struct woe_server * srv, int index, double * value, double * interval);
                           
 

#endif
