  
#include "rh.h"
#include "mem.h"
#include "monolog.h"
#include "state.h"

#include <string.h>

#include <stdarg.h>
#include <mruby.h>
#include <mruby/error.h>
#include <mruby/variable.h>

#include "esh.h"
#include "config.h"

/* Debugging level for console output. */ 
#define LOG_CONSOLE(FORMAT, ...) LOG_LEVEL("CONSOLE", FORMAT, __VA_ARGS__)
#define LOG_ENABLE_CONSOLE()     LOG_ENABLE(CONSOLE)
#define LOG_DISABLE_CONSOLE()    LOG_DISABLE(CONSOLE)

/* Amount of parameters a mruby function can be called with using _va functions */
#define RH_ARGS_MAX 64

/*
* RH contains helper functions for the mruby ruby interpreter.
*/

/*
 * Converts C like arguments to an array of mruby mrb_values.
 * Returns amount of arguments parsed, negative on failure.
 * 
 * Format specifier characters for the format argument: 
 *
 * z: String from null terminated C string      [char*]
 * s: String                                    [char*, int] 
 * Y: sYmbol from null terminated C string      [char*]
 * y: symbol                                    [char*, int] 
 * 
 * f: Float                                     [double]
 * i: Integer                                   [int]
 * b: Boolean                                   [int] 
 * 0: nil
 * Other characters: Error. Will return -(index or erroneous character) .
 */ 
int rh_args_va(Ruby * ruby, mrb_value * values,  int size,  char * format, va_list list) {
  int ch; int index;
  int i; double d; const char * str;
  index = 0; 
  
  for (ch = (*format) ; (ch) && (index < size) ; format++ , ch = (*format)) {
    mrb_value val;
    int error = FALSE;
    switch(ch) {
      case 's':
        str = va_arg(list, const char*);
        i   = va_arg(list, int);
        val = mrb_str_new(ruby, str, i);
        break;
        
      case 'z':
        str = va_arg(list, const char*);
        val = mrb_str_new_cstr(ruby, str);
        break;
        
      case 'y':
        str = va_arg(list, const char*);
        i   = va_arg(list, int);
        val = mrb_symbol_value(mrb_intern(ruby, str, (size_t)i));
        break;
        
      case 'Y':
        str = va_arg(list, const char*);
        val = mrb_symbol_value(mrb_intern_cstr(ruby, str));
        break; 
        
      case 'i':
        i   = va_arg(list, int);
        val = mrb_fixnum_value(i);
        break;
        
      case 'f':
        d   = va_arg(list, double);
        val = mrb_float_value(ruby, d);
        break;
        
      case 'b': 
        i   = va_arg(list, int);
        val = ( i ? mrb_true_value() : mrb_false_value());
        break;
      
      case '0':
        val = mrb_nil_value();
        break;
        
      default:
        error = TRUE;
        break;
    }
    if (error) { 
      return -index;
    }
    values[index] = val;
    index ++;
  }
  return index;
}


/* strdup isn't ANSI C, just posix... :p so we need our own version.*/
char *rh_strdup(const char *str) {
    char * res = malloc(strlen(str) + 1);
    if(res) { strcpy(res, str); }
    return res;
}

/* Helps convert C values to mruby values in an array. */
int rh_args(Ruby * ruby, mrb_value * values,  int size,  char * format, ...) {
  int res;
  va_list list;
  va_start(list, format);
  res = rh_args_va(ruby, values, size, format, list);
  va_end(list);
  return res;
}


/** Calulates the execption string. Result only tempoarily available..
XXX: check if this doesn't leak memory... you must free the results manually.
*/
char * rh_exceptionstring(Ruby * self) {
  char      * result;
  mrb_value   value;
  mrb_value   backtrace;
  mrb_value   backtrace_str;
  
  if (!self->exc) return NULL; 
  //
  /* XXX: Too bad, the backtrace doesn't seem to be filled in for some reason...
   * Should figure out how to fix this.
   */
  mrb_print_backtrace(self);
  backtrace = // mrb_get_backtrace(self);
  mrb_funcall(self, mrb_obj_value(self->exc), "backtrace", 0);

  backtrace_str = mrb_funcall(self, backtrace, "join", 1, mrb_str_new_lit(self, "\n"));
  LOG_ERROR("backtrace: %s\n", mrb_string_value_cstr(self, &backtrace_str));
  value  = mrb_funcall(self, mrb_obj_value(self->exc), "inspect", 0);
  
  
  
  // reset exception since we have it's string value.
  // Does this leak memory or not???
  self->exc = NULL;
  return rh_strdup(mrb_string_value_cstr(self, &value));
}


/** Allocates and initialzes a new ruby state. */
Ruby * rh_new() {
   Ruby * self = mrb_open();
   /*mrb_define_method(self, self->kernel_module, 
                     "path", tr_Path, ARGS_REQ(1));*/
   return self;
}

/** Frees a ruby state. */
Ruby * rh_free(Ruby * self) {
  mrb_close(self);
  return NULL;
}

/** Returns an mrb_value that contains the value of object.inspect. 
*/
mrb_value rh_inspect(mrb_state *mrb, mrb_value obj) {
  return mrb_inspect(mrb, obj);
}

char * rh_inspect_cstr(mrb_state *mrb, mrb_value value) {
  mrb_value res = rh_inspect(mrb, value);
  /* XXX: check that it's not modified anywere or export the const? */
  return (char *) mrb_string_value_cstr(mrb, &res);
}


/* Does the actual reporting depending on the current state of 
 ruby and the returned value. */
int rh_make_report(Ruby * self, mrb_value v) { 
  int res = 0;
  char * str;
 
  /* Report exceptions */
  str = rh_exceptionstring(self);
  if(str) {
    LOG_WARNING("mruby exception: %s\n", str);
    free(str);
    return 0;
  }
  
  /* Report result value if it's not nil on debug and console levels.
   */
  if (!mrb_nil_p(v)) {
    str = rh_inspect_cstr(self, v);
    LOG_DEBUG("mruby result: %s\n", str);
    LOG_CONSOLE("-> %s\n", str);
    return 0;
  }
  return 0;
}


/* Runs a file and reports any errors over the monolog logging system. */
int rh_run_file(Ruby * self, const char * filename, FILE * file) {
  int res;
  char * str;
  mrbc_context * c ; 
  mrb_value v;
  int ai;
  ai = mrb_gc_arena_save(self);
  c = mrbc_context_new(self);  
  mrbc_filename(self, c, filename);
  v = mrb_load_file_cxt(self, file, c);
  mrbc_context_free(self, c);  
  /* Report exceptions */
  res = rh_make_report(self, v);  
  mrb_gc_arena_restore(self, ai);
  return res;
}

int rh_run_filename(Ruby * self, const char * filename) {
  FILE * file = fopen(filename, "rt");
  int res;
  if (!file) {
    LOG_ERROR("No such ruby file: %s\n", filename);
    return -1;
  }
  res = rh_run_file(self, filename, file);
  fclose(file);
  return 0;
}


/**
* Executes a ruby file in Eruta's data/script directory with reporting.
* Returns -2 if the file was not found.
* Returns -3 if the path wasn't found.
*/
int rh_run_script(Ruby * self, const char * filename) {
  struct woesb path = {0};
  int runres;
  
  if (!woesb_new_join(&path, MRB_WOE_CONFIG(self)->data_dir, "/script/", filename, NULL)) {
    woesb_free(&path);
    LOG_ERROR("Out of memory when joining path.\n");
    return -3;
  }
  
  if (strstr(path.text, "..")) {
    woesb_free(&path);
    LOG_ERROR("Path may not contain '..' \n");
    return -4;
  }
  
  runres = rh_run_filename(self, path.text);
  woesb_free(&path);
  return runres;
}


/* Executes a ruby command string. 
Errors are reported to the reporter callback if it isn't NULL. */
int rh_dostring(Ruby * self, const char * command) {
  int res = 0;
  char * str;
  mrb_value v;
  int ai;
  ai = mrb_gc_arena_save(self);  
  #ifdef RH_USE_CONTEXT
  mrbc_context      * c  = mrbc_context_new(self);
  mrbc_filename(self, c, "command");
  v = mrb_load_string_cxt(self, command, c);
  mrbc_context_free(self, c); 
  #else
  v = mrb_load_string(self, command);
  #endif
  /* Report exceptions */
  res = rh_make_report(self, v);
  /* mrb GC area seems to be an area for 1024 "new" objects for the generational 
   * GC. It can overflow if a lot of new objects are generated 
   * (say exceptions, etc) on the C side. To prevent this the area must be saved 
   * and restored anywhere many ruby objects may have been generated.
   * It seems that here too this is needed.
   */
  mrb_gc_arena_restore(self, ai);
  return res;
}


/* Executes a ruby function with parameters. 
Errors are reported to the reporter callback if it isn't NULL. */
mrb_value rh_run_function_args(Ruby * self, mrb_value rubyself, char * funcname,    
                        int argc, mrb_value * argv) {
  int res = 0;
  char * str;
  mrb_value v;
  mrb_sym symname = mrb_intern_cstr(self, funcname);
  int ai;
  if(!mrb_respond_to(self, rubyself, symname)) {
    return mrb_nil_value();
  }  
  
  ai = mrb_gc_arena_save(self);
  v = mrb_funcall_argv(self, rubyself, symname, argc, argv);
  res = rh_make_report(self, v);
  mrb_gc_arena_restore(self, ai);
  return v;
}


/** Runs a function in the ruby interpreter, with C arguments according to the 
 * given format string, logging results and errors back to
 * the reporter. The limit is RH_ARGS_MAX arguments.
 */
mrb_value 
rh_run_function_va(Ruby * self, mrb_value rubyself, char * funcname,
                        char * format, va_list list) {
  mrb_value argv[RH_ARGS_MAX];
  int argc;
  argc = rh_args_va(self, argv, RH_ARGS_MAX, format, list);
  if (argc < 0) return mrb_nil_value();
  return rh_run_function_args(self, rubyself, funcname, argc, argv);
}


/** Runs a function in the ruby interpreter, under the toplevel self. 
 * This logs results and errors using monolog.h interfaces.
 */
mrb_value rh_run_toplevel_args(Ruby * ruby, char * name, int argc, mrb_value * argv) {
  return rh_run_function_args(ruby, mrb_top_self(ruby), name, argc, argv);
} 


/** Runs a function in the ruby interpreter, under the toplevel self. 
 * This logs results and errors using monolog.h interfaces.
 */
mrb_value rh_run_toplevel_va(Ruby * ruby, char * name, char * format, va_list list) {
  return rh_run_function_va(ruby, mrb_top_self(ruby), name, format, list);
} 


/** Runs a function in the ruby interpreter, logging results and errors
 * using monolog.h interfaces.
 */
mrb_value rh_run_function(Ruby * ruby, mrb_value rubyself, 
                          char * name, char * format, ...) {
  mrb_value res; 
  va_list list;
  va_start(list, format);
  res = rh_run_function_va(ruby, rubyself, name, format, list);
  va_end(list);
  return res;
}

/** Runs a function in the ruby interpreter, under the toplevel self. 
 * This logs results and errors using monolog.h interfaces.
 */
mrb_value rh_run_toplevel(Ruby * ruby, char * name, char * format, ...) {
  mrb_value res; 
  va_list list;
  va_start(list, format);
  res = rh_run_function_va(ruby, mrb_top_self(ruby), name, format, list);
  va_end(list);
  return res;
}


/* Calls a function, doesn't log anything. */
mrb_value rh_simple_funcall(Ruby * ruby, char * name) {
  int ai;  
  Ruby * mrb = (Ruby *) ruby;  
  mrb_value args[16];
  ai = mrb_gc_arena_save(mrb);
  LOG("GC Area: %d\n", ai);
  // if(ai> 99) exit(0);
  mrb_value v = mrb_funcall_argv(mrb, mrb_top_self(mrb), mrb_intern_cstr(mrb, name), 
                    0, args);
  
  if (mrb->exc) {
    if (!mrb_undef_p(v)) {        
      mrb_p(mrb, mrb_obj_value(mrb->exc));
    }
    return mrb_nil_value();
  }
  mrb_gc_arena_restore(mrb, 0);
  return v;
}

/* Maybe wrap this too? 
mrb_value mrb_funcall_argv(mrb_state*, mrb_value, mrb_sym, int, mrb_value*);
*/

/* Converts and mrb_value to a C boolean as per the ruby interpretation. */
int rh_tobool(mrb_value v) {
  if (mrb_nil_p(v))             return FALSE;
  return mrb_bool(v);
}

/* Documentation for mrb_get_args
  retrieve arguments from mrb_state.

  mrb_get_args(mrb, format, ...)
  
  returns number of arguments parsed.

  format specifiers:

   o: Object [mrb_value]
   S: String [mrb_value]
   A: Array [mrb_value]
   H: Hash [mrb_value]
   s: String [char*,int]
   z: String [char*]
   a: Array [mrb_value*,int]
   f: Float [mrb_float]
   i: Integer [mrb_int]
   n: Symbol [mrb_sym]
   &: Block [mrb_value]
   *: rest argument [mrb_value*,int]
   |: optional
*/



/* Stores a pointer in the mrb state. Handy to avoid having to use global 
 * variables o,n the C side. It's done by adding a ruby global variable
 * with a name that cannot be used from the mruby side. */
void * rh_store_pointer(mrb_state * mrb, const char * name, void * ptr) {
  struct woesb hname;
  mrb_value      val;
  mrb_sym      rname;

  if (!woesb_new_join(&hname, "--", name, "--", NULL)) {
    woesb_free(&hname);
    LOG_ERROR("Out of memory when making hidden name.\n");
    return NULL;
  }
  
  rname = mrb_intern_cstr(mrb, hname.text);
  val   = mrb_cptr_value(mrb, ptr);  
  mrb_gv_set(mrb, rname, val);
  woesb_free(&hname);
  return ptr;
}

/** Fetches a previously stored pointer. */
void * rh_fetch_pointer(mrb_state * mrb, const char * name) {
  struct woesb hname;
  mrb_value      val;
  mrb_sym      rname;

  if (!woesb_new_join(&hname, "--", name, "--", NULL)) {
    woesb_free(&hname);
    LOG_ERROR("Out of memory when making hidden name.\n");
    return NULL;
  }
  
  rname = mrb_intern_cstr(mrb, hname.text);
  val   = mrb_gv_get(mrb, rname);
  woesb_free(&hname);
  return mrb_cptr(val);
}


/* Tries to (re-)load the main ruby file, output to console. */
int rh_load_main() { 
  return -1; /*rh_run_script(woe_global_ruby, "main.rb"); */
}




