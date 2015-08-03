package world


type Direction  string

type Exit struct {
    Direction 
    ToRoomID    int
    toRoom    * Room         
}

type Room struct {
    Entity
    Exits   map[Direction]Exit
}

