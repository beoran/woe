
#include "esh.h"
#include <stdio.h>
#include <string.h>

/* Returns true if the string is in an unusable state, false if not. */
int esh_bad_p(char **me, size_t * size, size_t * space) {
  if (!me) return !0;
  if (!*me) return !0;
  if (!size) return !0;
  if (!space) return !0;
  return 0;
}


/* Makes the string an empty string by writing a nul character at positon 0
 * and setting size to 0. */
char * esh_make_empty(char ** me, size_t * size, size_t * space) {
  if (esh_bad_p(me, size, space)) return NULL;
  if ((*space) < 1) return NULL;
  (*me)[0] = '\0';
  (*size)  = 0;
  return (*me);
}

/* Ensures the string is NUL terminated. */
char * esh_nul_terminate(char ** me, size_t * size, size_t * space) {
  if (esh_bad_p(me, size, space)) return NULL;
  if (((*size) + 1) > (*space)) {
    (*me)[(*space)] = '\0';
    (*size) = (*space) - 1;
  }  else { 
    (*me)[(*size)+1] = '\0';
  }
  return (*me);
}


/* allocates a new string buffer with init_space. If init_space == 0 
 *  uses 1024 in stead.*/
char * esh_alloc(char ** me, size_t * size, size_t * space, size_t init_space) {
  if (!me) return NULL;
  if (!space) return NULL;
  if (!size) return NULL;
  if (init_space < 1) init_space = 1024; 
  (*me) = calloc(init_space, 1);
  (*space) = init_space;
  if (!*me) return NULL;
  return esh_make_empty(me, size, space);
}

/* Grows the given string's space. */
char * esh_grow(char ** me, size_t  * size, size_t * space, size_t new_space) {
  char * aid;
  (void) size;
  if (!me) return NULL;
  if (!*me) return NULL;
  if (new_space <= (*space)) return NULL;
  
  aid = realloc(*me, new_space);
  if (!aid) return NULL;
    
  memset(aid + (*space), '\0', new_space - (*space));
  (*space) = new_space;
  (*me)    = aid;   
  return (*me);
} 

/* Makes sure there is enough space to add amount characters. */
char * esh_ensure_space(char ** me, size_t * size, size_t * space, size_t grow_by) {
  if (!me) return NULL;
  if (!*me) return NULL;
  if (!size) return NULL;
  if (!space) return NULL;
  
  if ((*space) < ((*size) + grow_by)) { 
    if (!esh_grow(me, size, space, (*space) + grow_by + 255)) return NULL;
  }
  
  return (*me);
} 

char * esh_append_char(char ** me, size_t * size, size_t * space, char ch) {
  if (!esh_ensure_space(me, size, space, 1)) return NULL;
    
  (*me)[(*size)] = ch;
  (*size) = (*size) + 1;
  return (*me);
}


char * esh_append_buf(char  ** me, size_t * size, size_t * space, char * buf, size_t bufsize) {
  if (!esh_ensure_space(me, size, space, bufsize + 1)) return NULL;
  strncpy((*me) + (*size), buf, bufsize);
  (*size) = (*size) + bufsize;
  return esh_nul_terminate(me, size, space);
}

char * esh_append_cstr(char ** me, size_t * size, size_t * space, char * str) {
  return esh_append_buf(me, size, space, str, strlen(str));
}

char * esh_init_buf(char ** me, size_t * size, size_t * space, char * buf, size_t bufsize) {
  if (!size) return NULL;
  if (!esh_make_empty(me, size, space)) return NULL; 
  return esh_append_buf(me, size, space, buf, bufsize);
} 

char * esh_init_cstr(char ** me, size_t * size, size_t * space, char * buf) {
  return esh_init_buf(me, size, space, buf, strlen(buf));
} 

char * esh_new_cstr(char **me, size_t * size, size_t * space, char * buf) {
  size_t bufsize;
  if (!buf) return esh_new_empty(me, size, space);
  bufsize = strlen(buf);
  return esh_new_buf(me, size, space, buf, bufsize);
}


char * esh_new_buf(char **me, size_t * size, size_t * space, char * buf, size_t bufsize) {
  if (!esh_alloc(me, size, space, bufsize)) return NULL;
  return esh_init_buf(me, size, space, buf, bufsize);
}

/* Creates a new string with the given space and initialies it from init. */
char * esh_new(char **me, size_t * size, size_t * space, char * init) {
  return esh_new_buf(me, size, space, init, strlen(init));
}

/* Creates a new empty string with enough space. */
char * esh_new_empty(char **me, size_t * size, size_t * space) {
  return esh_new(me, size, space, "");
}

/* Frees the string and sets it to NULL. */
char * esh_free(char ** me, size_t * size, size_t * space) {
  if (!me) return NULL;
  if (!*me) return NULL;
  free(*me);
  *me    = NULL;
  *size  = 0;
  *space = 0;
  return NULL;
}

/** Reads in a file into the buffer. */
char * esh_read_file(char ** me, size_t * size, size_t * space, FILE * file) {
  char aid[1024];
  int read;
  if (esh_bad_p(me, size, space)) return NULL;
  if (!file) return NULL;
  
  while(!feof(file)) { 
    read = fread(aid, 1, sizeof(aid), file);
    if (read > 0) { 
      esh_append_buf(me, size, space, aid, read);
    }
  }
  return (*me);
}

/** Reads a named file into the buffer. */
char * esh_read_filename(char ** me, size_t * size, size_t * space, char * fn) {
  char * res;
  FILE * file = fopen(fn, "r");
  res = esh_read_file(me, size, space, file);
  if (file) fclose(file);
  return res;
}

/** Joins a NULL terminated list onto va*/
char * esh_join_va(char ** me, size_t * size, size_t * space, va_list strings) {
  char * next;
  while ( (next = va_arg(strings, char *)) ) {
    if (!esh_append_cstr(me, size, space, next)) return NULL;
  }
  return *me;
}

/* Joins a NULL terminated lists of strings onto me */
char * esh_join(char ** me, size_t * size, size_t * space, ...) {
  char * res;
  va_list strings;
  va_start(strings, space);
  res = esh_join_va(me, size, space, strings);
  va_end(strings);
  return res;
}

char * esh_init_join_cstr_va(char ** me, size_t * size, size_t * space, char * first, va_list strings) {
  if (!esh_init_cstr(me, size, space, first)) return NULL;
  return esh_join_va(me, size, space, strings);
}

char * esh_init_join_cstr(char ** me, size_t * size, size_t * space, char * first, ...) 
{
  char * res;
  va_list strings;
  va_start(strings, first);
  res = esh_init_join_cstr_va(me, size, space, first, strings);
  va_end(strings);
  return res;
}


char * esh_new_join_cstr_va(char ** me, size_t * size, size_t * space, char * first, va_list strings) {
  if (!esh_new_cstr(me, size, space, first)) return NULL;
  return esh_join_va(me, size, space, strings);
}

char * esh_new_join_cstr(char ** me, size_t * size, size_t * space, char * first, ...) 
{
  char * res;
  va_list strings;
  va_start(strings, first);
  res = esh_init_join_cstr_va(me, size, space, first, strings);
  va_end(strings);
  return res;
}



/** Deletes n bytes */

#define SWIS_EXPAND(SWIS) &((SWIS)->text), &((SWIS)->size), &((SWIS)->space)

int woesb_bad_p(struct woesb * me) {
  return (esh_bad_p(SWIS_EXPAND(me)));
}

struct woesb * woesb_make_empty(struct woesb * me) {
  if(!esh_make_empty(SWIS_EXPAND(me))) return NULL;
  return me;
}


struct woesb * woesb_nul_terminate(struct woesb * me) {
  if(!esh_nul_terminate(SWIS_EXPAND(me))) return NULL;
  return me;  
}

struct woesb * woesb_alloc(struct woesb * me, size_t init_space) {
  if(!esh_alloc(SWIS_EXPAND(me), init_space)) return NULL;
  return me;
}

struct woesb * woesb_grow(struct woesb * me, size_t new_space) {
  if(!esh_grow(SWIS_EXPAND(me), new_space)) return NULL;
  return me;
}

struct woesb * woesb_ensure_space(struct woesb * me, size_t grow_by) {
  if(!esh_ensure_space(SWIS_EXPAND(me), grow_by)) return NULL;
  return me;
}


struct woesb * woesb_append_char(struct woesb * me, char ch) {
  if(!esh_append_char(SWIS_EXPAND(me), ch)) return NULL;
  return me;
}


struct woesb * woesb_append_buf(struct woesb * me, char * buf, size_t bufsize) {
  if(!esh_append_buf(SWIS_EXPAND(me), buf, bufsize)) return NULL;
  return me;
}

struct woesb * woesb_append_cstr(struct woesb * me, char * str) {
  if(!esh_append_cstr(SWIS_EXPAND(me), str)) return NULL;
  return me;
}

struct woesb * woesb_init_buf(struct woesb * me, char * buf, size_t bufsize) {
  if(!esh_init_buf(SWIS_EXPAND(me), buf, bufsize)) return NULL;
  return me;
}

struct woesb * woesb_init_cstr(struct woesb * me, char * buf) {
  if(!esh_init_cstr(SWIS_EXPAND(me), buf)) return NULL;
  return me;
}


struct woesb * woesb_new_buf(struct woesb * me, char * buf, size_t bufsize) {
  if(!esh_new_buf(SWIS_EXPAND(me), buf, bufsize)) return NULL;
  return me;
}

struct woesb * woesb_new(struct woesb * me, char * init) {
  if(!esh_new(SWIS_EXPAND(me), init)) return NULL;
  return me;
}

struct woesb * woesb_new_empty(struct woesb * me) {
  if(!esh_new_empty(SWIS_EXPAND(me))) return NULL;
  return me;
}

struct woesb * woesb_free(struct woesb * me) {
  esh_free(SWIS_EXPAND(me)); 
  return NULL;
}

struct woesb * woesb_read_file(struct woesb * me, FILE * file) {
  if(!esh_read_file(SWIS_EXPAND(me), file)) return NULL;
  return me;
}
 
struct woesb * woesb_read_filename(struct woesb * me, char * fn) {
  if(!esh_read_filename(SWIS_EXPAND(me), fn)) return NULL;
  return me;
}

struct woesb * woesb_join_va(struct woesb * me, va_list strings) {
  if (!esh_join_va(SWIS_EXPAND(me), strings)) return NULL;
  return me;
}

struct woesb * woesb_join(struct woesb * me, ...) {
  struct woesb * res;
  va_list strings;
  va_start(strings, me);
  res = woesb_join_va(me, strings);
  va_end(strings);
  return res;
}

struct woesb * woesb_init_join_va(struct woesb * me, char * first, va_list strings) {
  if (!esh_init_join_cstr_va(SWIS_EXPAND(me), first, strings)) return NULL;
  return me;
}


struct woesb * woesb_init_join(struct woesb * me, char * first, ...) {
  struct woesb * res;
  va_list strings;
  va_start(strings, first);
  res = woesb_init_join_va(me, first, strings);
  va_end(strings);
  return res;  
}


struct woesb * woesb_new_join_va(struct woesb * me, char * first, va_list strings) {
  if (!esh_new_join_cstr_va(SWIS_EXPAND(me), first, strings)) return NULL;
  return me;
}

struct woesb * woesb_new_join(struct woesb * me, char * first, ...) {
  struct woesb * res;
  va_list strings;
  va_start(strings, first);
  res = woesb_new_join_va(me, first, strings);
  va_end(strings);
  return res;  
}


