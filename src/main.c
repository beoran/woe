#define _POSIX_C_SOURCE 200801L
#define _POSIX_SOURCE 200801L

#include <stdlib.h>
#include <unistd.h>
#include <signal.h>
#include <errno.h>
#include <time.h>


#include "server.h"
#include "monolog.h"
#include "config.h"
#include "state.h"
#include "rh.h"
#include "toruby.h"



DEFINE_FILE_LOGGER(woe_file_logger);
DEFINE_STDOUT_LOGGER(woe_stdout_logger);
DEFINE_STDERR_LOGGER(woe_stderr_logger);


FILE * woe_start_monolog(char * logname) {
  FILE * fout;

  monolog_init();
  LOG_ENABLE_ERROR();
  monolog_enable_level("WARNING");
  monolog_enable_level("NOTE");

  fout = fopen(logname, "w");
  
  monolog_add_logger(NULL, &woe_stdout_logger);
  
  if (fout) {
    monolog_add_logger(fout, &woe_file_logger);
    fprintf(stderr, "Opened log file %s %p\n", logname, fout);
  } else {
    fprintf(stderr, "Could not open log file %s\n", logname);
  }
    
  return fout;
}



void woe_stop_monolog() {  
  monolog_done();
}


int main(int argc, char * argv[]) {
  sigset_t            signal_set  = { { 0 } };
  sigset_t            signal_old  = { { 0 } };
  siginfo_t           signal_info = { 0 };
  struct timespec     signal_time = { 0, 10 };
  
  struct woe_state    state  = { 0 };
  struct woe_config   config = { 0 };  
  state.config     = &config;
    
  FILE * logfile   = NULL;
  state.mrb        = mrb_open();
  state.mrb->ud    = &state;
  tr_init(state.mrb);
  woe_config_init_args(state.config, argc, argv);
  
  logfile          = woe_start_monolog(state.config->log_file);
  
  state.server     = woe_server_new(state.config->port);
  
  woe_server_set_mrb(state.server, state.mrb);
  
  /* Handle all signals, except a few. */
  sigfillset(&signal_set);
  sigdelset(&signal_set, SIGKILL);
  sigdelset(&signal_set, SIGTERM);
  sigdelset(&signal_set, SIGINT);
  sigdelset(&signal_set, SIGSEGV);
  sigdelset(&signal_set, SIGTSTP);
  
  /* Ignore them all, because we will wait for them in stead. */
  sigprocmask(SIG_SETMASK, &signal_set, &signal_old);

  
  rh_run_script(state.mrb, "main.rb");
  
  if (woe_server_listen(state.server) > 0) {
     LOG_ERROR("Cannot listen. Stop.");   
  } else {
    rh_run_toplevel(state.mrb, "woe_on_start", "");
    while (woe_server_busy(state.server)) {
      int caught = 0;
      
      woe_server_update(state.server, 1);
      caught = sigtimedwait(&signal_set, NULL, &signal_time);
      if (caught > 0) {
        LOG_NOTE("Received signal %d\n", caught);
        rh_run_toplevel(state.mrb, "woe_on_signal", "i", caught);
        /* woe_server_request_shutdown(state.server); */
      }
    }
  }
 
  rh_run_toplevel(state.mrb, "woe_on_stop", "");
  woe_server_free(state.server);
  mrb_close(state.mrb);
  LOG_NOTE("Shutting down WOE\n");
  
  woe_stop_monolog();
  return 0;
}
