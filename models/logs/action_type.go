package logs

type ActionType int

const (
	ActionAddition ActionType = iota + 1
	ActionOwnerChange
)

func (e ActionType) String() string {
	switch e {
	case ActionAddition:
		return "addition"
	case ActionOwnerChange:
		return "owner_change"
	default:
		return ""
	}
}
