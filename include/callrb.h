#ifndef collide_H_INCLUDED
#define collide_H_INCLUDED

#include "state.h"

enum CollisionKinds_ {
  COLLIDE_BEGIN     = 1,
  COLLIDE_COLLIDING = 2,
  COLLIDE_END       = 3
};

int callrb_sprite_event(SpriteState * spritestate, int kind, void * data);

int callrb_on_start();

int callrb_on_reload();

int callrb_on_update(State * self);


#endif




