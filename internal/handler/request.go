package handler

type putUserData struct {
	Blocked     *bool  `json:"blocked"`
	BlockReason string `json:"block_reason"`
}

type getUsersData struct {
	Include []string `form:"include[]"`
}
