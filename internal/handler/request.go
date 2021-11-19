package handler

type putUserData struct {
	Blocked     *bool  `json:"blocked"`      // user state, if blocked no interaction with the bot is possible
	BlockReason string `json:"block_reason"` // internally displayed reason for a block
}

type getUsersData struct {
	Include []string `form:"include[]"`
}
