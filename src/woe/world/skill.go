package world

type Skill struct {
    Entity
    Kind              int
    
    derived           func (m * Being)  int 
    talent          * int 
}

type BeingSkill struct {
    being       * Being
    technique   * Skill
    Experience    int
    Level         int
}

