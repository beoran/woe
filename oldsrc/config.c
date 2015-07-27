#include "config.h"

#define _POSIX_C_SOURCE 200801L
#define _POSIX_SOURCE 200801L

#include <stdlib.h>
#include <unistd.h>




int woe_config_init_args(struct woe_config * config, int argc, char * argv[]) {
  char opt;
  
  config->port      = 7777;
  config->data_dir  = getenv("WOE_DATA");
  if (!config->data_dir) {
    config->data_dir= "data";
  }
  
  config->log_file  = getenv("WOE_LOG");
  if (!config->log_file) {
    config->log_file = "woe.log";
  }
    
  while ((opt = getopt(argc, argv, "p:d:l:")) != -1) {
    switch (opt) {
    case 'p':
      config->port = atoi(optarg);
      break;
    
    case 'd': 
      config->data_dir = optarg;
      break;
      
    case 'l': 
      config->log_file = optarg;
      break;  
            
    default: /* '?' */
      break;
    }
  }
  return 0;
}


static struct woe_config global_woe_config;

struct woe_config woe_config_get() {
  return global_woe_config;
}

struct woe_config woe_config_put(struct woe_config * config) {
  if (config) {
    global_woe_config = *config;
  }  
  return woe_config_get();
}


