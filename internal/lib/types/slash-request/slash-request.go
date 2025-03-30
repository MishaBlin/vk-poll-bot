package slash_request

type SlashCommandRequest struct {
	Token       string         `json:"token"`
	TeamID      string         `json:"team_id"`
	TeamDomain  string         `json:"team_domain"`
	ChannelID   string         `json:"channel_id"`
	ChannelName string         `json:"channel_name"`
	UserID      string         `json:"user_id"`
	UserName    string         `json:"user_name"`
	Command     string         `json:"command"`
	Text        string         `json:"text"`
	ResponseURL string         `json:"response_url"`
	TriggerID   string         `json:"trigger_id"`
	Context     CommandContext `json:"context"`
}

type CommandContext struct {
	PollID      string `json:"poll_id"`
	OwnerID     string `json:"owner_id"`
	OptionIndex int    `json:"option_index"`
	Token       string `json:"token"`
}
