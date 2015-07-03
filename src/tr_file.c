
#include <errno.h>
#include <stdio.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <dirent.h>


#include "esh.h"
#include "toruby.h"
#include "tr_macro.h"
#include "monolog.h"
#include "rh.h"
#include "state.h"
#include "server.h"
#include "tr_file.h"


struct tr_file {
  FILE              * file;
  struct woe_config * cfg;
};

typedef struct tr_file tr_file;

static void file_close(tr_file * file) {
  if (!file) return;
  if (!file->file) return;
  fclose(file->file);
  file->file = NULL;
}

static void tr_file_free(mrb_state *mrb, void *ptr) {
  struct tr_file * file = ptr;
  file_close(file);
  mrb_free(mrb, file);
}

static FILE * file_fopen(mrb_state * mrb, char * filename, char * mode) 
{
  struct woe_config * cfg;
  struct woesb buf = { 0 };
  FILE * me = NULL;
  if (!mrb) return NULL;
  cfg = MRB_WOE_CONFIG(mrb);
  if (!cfg) return NULL;
  
  if (!woesb_new_join(&buf, cfg->data_dir, "/var/", filename, NULL)) {
    LOG_ERROR("Cannot allocate space for file name.\n");
    return NULL;
  }
  
  if (strstr(buf.text, "..")) {
    mrb_free(mrb, me);
    woesb_free(&buf);
    LOG_ERROR("Path may not contain '..' \n");
    return NULL;
  }

  me = fopen(buf.text, mode);
  if (!me) {
    LOG_ERROR("Cannot open file %s.\n", filename);
    woesb_free(&buf);
    return NULL;
  }
  
  woesb_free(&buf);
  return me;
}


static tr_file * file_open(mrb_state * mrb, char * filename, char * mode) 
{
  struct woe_config * cfg;
  tr_file * me = NULL;
  if (!mrb) return NULL;
  cfg = MRB_WOE_CONFIG(mrb);
  if (!cfg) return NULL;
  
  me  = mrb_malloc(mrb, sizeof(struct tr_file));
  if (!me) return NULL;
  me->file = file_fopen(mrb, filename, mode);
  if (!me->file) {
    mrb_free(mrb, me);
    return NULL;
  }
  me->cfg = cfg; 
  return me;
}

int woe_mkdir(struct woe_config * cfg, char * filename) {
  int res;
  DIR * dir;
  
  struct woesb buf = { 0 };
  if (!cfg) return -3;
 
  if (!woesb_new_join(&buf, cfg->data_dir, "/var/", filename, NULL)) {
    LOG_ERROR("Cannot allocate space for file name.\n");
    return -1;
  }
  
  if (strstr(buf.text, "..")) {
    woesb_free(&buf);
    LOG_ERROR("Path may not contain '..' \n");
    return -2;
  }
  
  dir = opendir(buf.text);
  if (dir) {
    LOG_DEBUG("Dir %s already exists\n");
    /* Directory already exists */
    closedir(dir);
    woesb_free(&buf);
    return 0; 
  }
  
  errno = 0;
  res = mkdir(buf.text , 0770);

  if (res < 0) {
    LOG_ERROR("Cannot make dir %s.\n", filename);
    woesb_free(&buf);
    return errno;
  }
  
  woesb_free(&buf);
  return 0;
}

int woe_link(struct woe_config * cfg, char * old_filename, char * new_filename) {
  int res;
  
  struct woesb old_buf = { 0 };
  struct woesb new_buf = { 0 };

  if (!cfg) return -3;
 
  if (!woesb_new_join(&old_buf, cfg->data_dir, "/var/", old_filename, NULL)) {
    LOG_ERROR("Cannot allocate space for file name.\n");
    return -1;
  }

  if (!woesb_new_join(&new_buf, cfg->data_dir, "/var/", new_filename, NULL)) {
    woesb_free(&old_buf);
    LOG_ERROR("Cannot allocate space for file name.\n");
    return -1;
  }

  
  if (strstr(old_buf.text, "..") || strstr(new_buf.text, "..")) {
    woesb_free(&new_buf);
    woesb_free(&old_buf);    
    LOG_ERROR("Path may not contain '..' \n");
    return -2;
  }
  
  errno = 0;
  res = link(old_buf.text, new_buf.text);

  if (res < 0) {
    LOG_ERROR("Cannot link %s to %s.\n", old_filename, new_filename);
    woesb_free(&new_buf);
    woesb_free(&old_buf);
    return errno;
  }
  
  woesb_free(&new_buf);
  woesb_free(&old_buf);
  return 0;
}


struct mrb_data_type tr_file_type = { "File", tr_file_free };

static mrb_value tr_file_wrap(mrb_state *mrb, struct RClass *tc, tr_file * file) {
  return mrb_obj_value(Data_Wrap_Struct(mrb, tc, &tr_file_type, file));
}

tr_file * tr_file_unwrap(mrb_state *mrb, mrb_value val) {
  return DATA_CHECK_GET_PTR(mrb, val, &tr_file_type, tr_file);
}

static mrb_value tr_file_open(mrb_state * mrb, mrb_value self) {
  int res; 
  char * name, * mode;
  tr_file * file; 
  
  mrb_get_args(mrb, "zz", &name, &mode);
  file = file_open(mrb, name, mode);
  
  if (!file) return mrb_nil_value();
  
  return tr_file_wrap(mrb, mrb_class_get(mrb, "File"), file);
     
}    
  
static mrb_value tr_file_close(mrb_state * mrb, mrb_value self) {
  file_close(tr_file_unwrap(mrb, self));
  return mrb_nil_value();
}

static mrb_value tr_file_write(mrb_state * mrb, mrb_value self) { 
  mrb_int res, size;
  tr_file * file; 
  char * buf;
  file = tr_file_unwrap(mrb, self);
  mrb_get_args(mrb, "s", &buf, &size);
  res = fwrite(buf, size, 1, file->file);
  return mrb_fixnum_value(res);
}

static mrb_value tr_file_read(mrb_state * mrb, mrb_value self) {
  mrb_int res, size;
  tr_file * file;
  char * mem; 
  mrb_value buf;
  
  
  file = tr_file_unwrap(mrb, self);
  mrb_get_args(mrb, "i", &size);
  mem = calloc(size, 1);
  res = fread(mem, 1, size, file->file);
  if (res > 0) {
    buf = mrb_str_new(mrb, mem, size);
    free(mem);
    return buf;
  } 
  
  if (res == 0) {
    free(mem);
    return mrb_nil_value();
  }
  
  // if (res < 0)
  free(mem);
  LOG_ERROR("Failed to read from file.\n");
  return mrb_nil_value();
}   


static mrb_value tr_file_readall(mrb_state * mrb, mrb_value self) {
  #define READALL_BUFSIZE 1024
  mrb_int res = 0, size = 0;
  FILE * file;
  char * mem = NULL, * aid = NULL; 
  char * filename = NULL;
  mrb_value buf;
  (void) self;
  
  
  mrb_get_args(mrb, "z", &filename);
  if (!filename) {
      return mrb_nil_value();
  }
  
  file = file_fopen(mrb, filename, "rb");
  
  if (!file) { 
    return mrb_nil_value();
  }
  
  while (!feof(file)) {
    size += READALL_BUFSIZE;
    aid   = realloc(mem, size);
    if (!aid) { 
      buf = mrb_nil_value();
      goto done;
    }
    mem = aid;
    res   = fread(mem + size - READALL_BUFSIZE, 1, READALL_BUFSIZE, file);
    
    if (res < READALL_BUFSIZE) {
      size = size - READALL_BUFSIZE + res;
      break;
    }
  }
  buf = mrb_str_new(mrb, mem, size);
  
  done:  
    free(mem);
    fclose(file);
    return buf;
}   
 

static mrb_value tr_file_puts(mrb_state * mrb, mrb_value self) { 
  mrb_int res, size;
  tr_file * file; 
  char * buf;
  file = tr_file_unwrap(mrb, self);
  mrb_get_args(mrb, "z", &buf);
  res = fputs(buf, file->file);
  return mrb_fixnum_value(res);
}


static mrb_value tr_file_gets(mrb_state * mrb, mrb_value self) {
  mrb_int size;
  tr_file * file; 
  char * buf;
  char * check;
  mrb_value res;
  
  file = tr_file_unwrap(mrb, self);
  mrb_get_args(mrb, "i", &size);
  buf = calloc(size, 1);
  if (!buf) {
    return mrb_nil_value();
  }
  
  check = fgets(buf, size, file->file);
  if (check) {
    res = mrb_str_new_cstr(mrb, buf);
    free(buf); 
    return res;
  } 
  
  free(buf);
  return mrb_nil_value();
}


static mrb_value tr_file_putc(mrb_state * mrb, mrb_value self) { 
  mrb_int res, chara = 0;
  tr_file * file; 
  file = tr_file_unwrap(mrb, self);
  mrb_get_args(mrb, "i", &chara);
  res = fputc(chara, file->file);
  return mrb_fixnum_value(res);
}


static mrb_value tr_file_getc(mrb_state * mrb, mrb_value self) {
  mrb_int res;
  tr_file * file; 
  file = tr_file_unwrap(mrb, self);
  res = fgetc(file->file);
  return mrb_fixnum_value(res);
} 
 
static mrb_value tr_file_eof(mrb_state * mrb, mrb_value self) {
  tr_file * file; 
  file = tr_file_unwrap(mrb, self);  
  return rh_bool_value(feof(file->file));
}

static mrb_value tr_dir_mkdir(mrb_state * mrb, mrb_value self) {
  tr_file * file; 
  char * name;
  struct woe_config * cfg;
  cfg = MRB_WOE_CONFIG(mrb);
  file = tr_file_unwrap(mrb, self);
  mrb_get_args(mrb, "z", &name);
  return mrb_fixnum_value(woe_mkdir(cfg, name));
}

static mrb_value tr_file_link(mrb_state * mrb, mrb_value self) {
  tr_file * file; 
  char * old_name, * new_name;
  struct woe_config * cfg;
  cfg = MRB_WOE_CONFIG(mrb);
  file = tr_file_unwrap(mrb, self);
  mrb_get_args(mrb, "zz", &old_name, &new_name);
  return mrb_fixnum_value(woe_link(cfg, old_name, new_name));
}


int tr_init_file(mrb_state * mrb) {
  struct RClass *krn;
  struct RClass *dir;
  struct RClass *fil;
  
  fil = mrb_define_class(mrb, "File"    , mrb_class_get(mrb, "Object"));
  dir = mrb_define_class(mrb, "Dir"     , mrb_class_get(mrb, "Object"));

  TR_CLASS_METHOD_ARGC(mrb, fil, "read" , tr_file_readall, 1);
  TR_CLASS_METHOD_ARGC(mrb, fil, "open" , tr_file_open, 2);
  TR_CLASS_METHOD_ARGC(mrb, fil, "link" , tr_file_link, 2);
  TR_CLASS_METHOD_ARGC(mrb, dir, "mkdir", tr_dir_mkdir, 1);


  TR_METHOD_NOARG(mrb, fil, "close"     , tr_file_close);
  TR_METHOD_NOARG(mrb, fil, "eof?"       , tr_file_eof);
  TR_METHOD_ARGC(mrb, fil, "read"       , tr_file_read, 1);
  TR_METHOD_ARGC(mrb, fil, "write"      , tr_file_read, 2);
  TR_METHOD_ARGC(mrb, fil, "gets"       , tr_file_gets, 1);
  TR_METHOD_ARGC(mrb, fil, "puts"       , tr_file_puts, 1);
  TR_METHOD_ARGC(mrb, fil, "putc"       , tr_file_putc, 1);
  TR_METHOD_ARGC(mrb, fil, "getc"       , tr_file_getc, 1);
  
  return 0;
}
