package world


type Technique struct {
    Entity
    Kind        int
    Effect      int
    Level       int
    Cost        int
    SkillID     int
    skill     * Skill
    onUse       func (me * Technique, caster * Being, targets ...*Being) (bool)
}


type BeingTechnique struct {
    being       * Being
    technique   * Technique
    Experience    int
    Level         int
}

/* An exploit is a special technique that can only be used a few times 
 * per day.
 */
type Exploit struct {
    // it's a special case of a technique
    Technique
}

type BeingExploit struct {
    being       * Being
    exploit     * Exploit
    Experience    int
    Level         int
    // How many times the exploit may be used.
    Uses          Vital
}



