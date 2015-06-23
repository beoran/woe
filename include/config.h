
#ifndef WOE_CONFIG_H
#define WOE_CONFIG_H


struct woe_config {
  int     port;
  char *  data_dir;
  char *  log_file;
};

int woe_config_init_args(struct woe_config * config, int argc, char * argv[]);

struct woe_config woe_config_get();
struct woe_config woe_config_put(struct woe_config * config);


  
#endif  
