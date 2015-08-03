package world

type Art struct {
    // An art is a special case of technique that requires Essence, 
    // and consumes JP in stead of MP
    Technique
}


type BeingArt struct {
    being       * Being
    art         * Art
    Experience    int
    Level         int
}



