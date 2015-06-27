#ifndef FILE_H_INCLUDED
#define FILE_H_INCLUDED

#include <stdlib.h>
#include <stdarg.h>

#include "config.h"

struct woe_file;
typedef struct woe_file woe_file;

woe_file * woe_file_open(struct woe_config * cfg, char * filename, char * mode);
void woe_file_close(woe_file * file);

size_t woe_file_write(woe_file * file, void * buf, size_t size);
size_t woe_file_read(woe_file * file, void * buf, size_t size);

int woe_file_puts(woe_file * file, char * buf);
int woe_file_gets(woe_file * file, char * buf, int size);

int woe_file_putc(woe_file * file, char c);
int woe_file_getc(woe_file * file, char c);

int woe_mkdir(struct woe_config * cfg, char * filename);

int woe_file_eof(woe_file * file);

int tr_init_file(mrb_state * mrb);

#endif
