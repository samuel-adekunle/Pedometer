package steps

type Step struct {
	UserName    string         `json:"user_name"`
	Challenges  []string       `json:"challenges"`
	Count       map[string]int `json:"count"`
	DailyTarget int            `json:"daily_target"`
}

type StepRequestBody struct {
	UserName    string         `json:"user_name"`
	Count       map[string]int `json:"count"`
	DailyTarget int            `json:"daily_target"`
}
