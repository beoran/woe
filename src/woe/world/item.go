package world

type Item struct {
    Entity
    Quality     int
    Price       int
    Type        int
    Equippable  int
}

type ItemPointer struct {
    ID     ID
    item * Item
}




