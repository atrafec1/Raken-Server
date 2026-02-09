package dto

type PayrollEntry struct {
	EmployeeCode   string
	CurrentDate    string
	CraftLevel     string
	JobNumber      string
	Phase          string
	CostCode       string
	ChangeOrder    string
	RegularHours   float64
	OvertimeHours  float64
	PremiumHours   float64
	Day            int
	EquipmentCode  string
	DownFlag       string
	SpecialPayType string
	SpecialPayCode string
	SpecialUnits   float64
	SpecialRate    float64
}
