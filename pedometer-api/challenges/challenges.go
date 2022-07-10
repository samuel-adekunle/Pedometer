package challenges

type Challenges struct {
	ChallengeName string         `json:"challenge_name"`
	Current       int            `json:"current"`
	Target        int            `json:"target"`
	Steps         map[string]int `json:"steps"`
	StartDate     string         `json:"start_date"`
	EndDate       string         `json:"end_date"`
}

type ChallengeRequestBody struct {
	ChallengeName string `json:"challenge_name"`
	Target        int    `json:"target"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
}
