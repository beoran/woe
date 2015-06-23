#ifndef ESH_H_INCLUDED
#define EHS_H_INCLUDED

/* Explicit String Handling. These functions help handling strings directly 
 * using struct wrappers, you need to pass in the parameters explicitly and 
 * most often by reference.
 * 
 * Nevertheless a struct wrapper is also provided , named woesb, 
 * or WOE String Buffer.
 */
#include <stdlib.h> 
#include <stdio.h>
#include <stdarg.h> 

 
int    esh_bad_p(char **me, size_t * size, size_t * space);
char * esh_make_empty(char ** me, size_t * size, size_t * space);
char * esh_nul_terminate(char ** me, size_t * size, size_t * space);
char * esh_alloc(char ** me, size_t * size, size_t * space, size_t init_space);
char * esh_grow(char ** me, size_t  * size, size_t * space, size_t new_space);
char * esh_ensure_space(char ** me, size_t * size, size_t * space, size_t grow_by);
char * esh_append_char(char ** me, size_t * size, size_t * space, char ch);
char * esh_append_buf(char  ** me, size_t * size, size_t * space, char * buf, size_t bufsize);
char * esh_append_cstr(char ** me, size_t * size, size_t * space, char * str);
char * esh_init_buf(char ** me, size_t * size, size_t * space, char * buf, size_t bufsize);
char * esh_init_cstr(char ** me, size_t * size, size_t * space, char * buf);
char * esh_new_buf(char **me, size_t * size, size_t * space, char * buf, size_t bufsize);
char * esh_new(char **me, size_t * size, size_t * space, char * init);
char * esh_new_empty(char **me, size_t * size, size_t * space);
char * esh_free(char ** me, size_t * size, size_t * space);
char * esh_read_file(char ** me, size_t * size, size_t * space, FILE * file);
char * esh_read_filename(char ** me, size_t * size, size_t * space, char * fn);


char * esh_join(char ** me, size_t * size, size_t * space,  ...);
char * esh_join_va(char ** me, size_t * size, size_t * space, va_list strings);
char * esh_init_join_cstr(char ** me, size_t * size, size_t * space, char * first, ...);
char * esh_init_join_cstr_va(char ** me, size_t * size, size_t * space, char * first, va_list strings);



struct woesb {
  char      * text;
  size_t      size;
  size_t      space;  
};
 
int woesb_bad_p(struct woesb * me);
struct woesb * woesb_make_empty(struct woesb * me);
struct woesb * woesb_nul_terminate(struct woesb * me);
struct woesb * woesb_alloc(struct woesb * me, size_t init_space);
struct woesb * woesb_grow(struct woesb * me, size_t new_space);
struct woesb * woesb_ensure_space(struct woesb * me, size_t grow_by);
struct woesb * woesb_append_char(struct woesb * me, char ch);
struct woesb * woesb_append_buf(struct woesb * me, char * buf, size_t bufsize);
struct woesb * woesb_append_cstr(struct woesb * me, char * str);
struct woesb * woesb_init_buf(struct woesb * me, char * buf, size_t bufsize);
struct woesb * woesb_init_cstr(struct woesb * me, char * buf);
struct woesb * woesb_new_buf(struct woesb * me, char * buf, size_t bufsize);
struct woesb * woesb_new(struct woesb * me, char * init);
struct woesb * woesb_new_empty(struct woesb * me);
struct woesb * woesb_free(struct woesb * me);
struct woesb * woesb_read_file(struct woesb * me, FILE * file);
struct woesb * woesb_read_filename(struct woesb * me, char * fn); 

struct woesb * woesb_join(struct woesb * me, ...);
struct woesb * woesb_join_va(struct woesb * me, va_list strings);
struct woesb * woesb_init_join(struct woesb * me, char * first, ...);
struct woesb * woesb_init_join_va(struct woesb * me, char * first, va_list strings);
struct woesb * woesb_new_join_va(struct woesb * me, char * first, va_list strings);
struct woesb * woesb_new_join(struct woesb * me, char * first, ...);


#endif




