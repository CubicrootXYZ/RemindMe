package handler

type putUserData struct {
	Blocked     *bool  `json:"blocked"`      // user state, if blocked no interaction with the bot is possible
	BlockReason string `json:"block_reason"` // internally displayed reason for a block
}

type getUsersData struct {
	Include []string `form:"include[]"`
}

type postChannelThirdPartyResourceData struct {
	Type        string `json:"type"`         // the type of resource to add
	ResourceURL string `json:"resource_url"` // url to the resource
}
