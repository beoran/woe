
/* Callbacks from the C side into the mruby side.
 * Used to signal several events such as collisions or sprite
 * animation.
*/
#include "state.h"
#include "rh.h"
#include "spritestate.h"
#include "callrb.h"

/* Sprite event handler. Calls an mruby callback. */
int callrb_sprite_event(SpriteState * spritestate, int kind, void * data) { 
  mrb_value res;
  Sprite * sprite;
  State * state;
  int spriteid, thingid, pose, direction;
  Ruby * ruby;
  void * thing;
  sprite    = spritestate_sprite(spritestate);
  spriteid  = sprite_id(sprite);
  thing     = spritestate_data(spritestate);  
  pose      = spritestate_pose(spritestate);
  direction = spritestate_direction(spritestate);
  state     = state_get();
  ruby      = state_ruby(state);
  res       = rh_run_toplevel(ruby, "eruta_on_sprite", "iiii", 
                spriteid, pose, direction, kind);
  (void) data;
  return rh_tobool(res);
}

/* Calls the eruta_on_start function  */ 
int callrb_on_start() { 
  mrb_value res;
  State * state = state_get();
  Ruby * ruby   = state_ruby(state);
  res           = rh_run_toplevel(ruby, "woe_on_start", "");
  return rh_tobool(res);
}

/* Calls the eruta_on_reload function  */
int callrb_on_reload() { 
  mrb_value res;
  State * state = state_get();
  Ruby * ruby   = state_ruby(state);
  res           = rh_run_toplevel(ruby, "woe_on_reload", "");
  return rh_tobool(res);
}


/* Calls the eruta_on_update function. */
int callrb_on_update(mrb_ruby * self, double dt) {
  mrb_value res;  
  mrb_value mval = mrb_float_value(self, dt);
  res = rh_run_toplevel_args(state_ruby(self), "woe_on_update", 1, &mval);
  return rh_tobool(res);
}

