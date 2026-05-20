package domain

type Skill string

const (
	SkillCook      Skill = "cook"
	SkillHouseHelp Skill = "house-help"
	SkillMaid      Skill = "maid"
)

type ConstraintType string

const (
	Hard ConstraintType = "hard"
	Soft ConstraintType = "soft"
)

type Constraint struct {
	Key   string         `json:"key"`
	Value string         `json:"value"`
	Type  ConstraintType `json:"type"`
}

type TimeSlot struct {
	Date  string `json:"date"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type Runner struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Skills     []Skill           `json:"skills"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Slots      []TimeSlot        `json:"slots,omitempty"`
}

type Allocation struct {
	ID                string   `json:"id"`
	RequestID         string   `json:"requestId"`
	RunnerID          string   `json:"runnerId"`
	RunnerName        string   `json:"runnerName"`
	Skill             Skill    `json:"skill"`
	Slot              TimeSlot `json:"slot"`
	SoftScore         int      `json:"softScore"`
	MatchedSoft       []string `json:"matchedSoft,omitempty"`
	UnmatchedSoft     []string `json:"unmatchedSoft,omitempty"`
	MatchedHard       []string `json:"matchedHard,omitempty"`
	ConsideredRunners int      `json:"consideredRunners"`
}

type AddRunnerCommand struct {
	Name       string
	Skills     []Skill
	Attributes map[string]string
	Slots      []TimeSlot
}

type AllocateCommand struct {
	RequestID   string
	Skill       Skill
	Slot        TimeSlot
	Constraints []Constraint
}
