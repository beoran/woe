package world


// An entity is anything that can exist in a World
type Entity struct {
    ID                  string
    Name                string
    ShortDescription    string
    Description         string
    Aliases           []string
}

