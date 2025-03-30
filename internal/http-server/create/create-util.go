package create

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v6/model"
	pollStruct "mm-polls/internal/lib/types/poll"
	slashRequest "mm-polls/internal/lib/types/slash-request"
	"net/http"
	"os"
	"strings"
)

func getActions(poll *pollStruct.Poll, token string) []*model.PostAction {
	actions := make([]*model.PostAction, 0)
	for i, label := range poll.Options {
		actions = append(actions, buildAction(poll, token, label, "primary", "/vote", i))
	}
	actions = append(actions, buildAction(
		poll, token, "View results", "warning", "/results", -1))

	actions = append(actions, buildAction(
		poll, token, "Close poll", "danger", "/finish", -1))

	return actions
}

func buildAction(poll *pollStruct.Poll, token, label, color, handle string, optionIndex int) *model.PostAction {
	ctx := slashRequest.CommandContext{
		PollID:      poll.ID,
		OwnerID:     poll.OwnerID,
		OptionIndex: optionIndex,
		Token:       token,
	}

	return &model.PostAction{
		Name:  label,
		Style: color,
		Integration: &model.PostActionIntegration{
			URL: os.Getenv("BASE_URL") + handle,
			Context: map[string]interface{}{
				"poll_id":      ctx.PollID,
				"owner_id":     ctx.OwnerID,
				"option_index": ctx.OptionIndex,
				"token":        ctx.Token,
			},
		},
	}
}

func parsePoll(input string) (string, []string, error) {
	parts := strings.Split(input, "|")
	if len(parts) < 2 {
		return "", nil, fmt.Errorf("input must contain a title and at least one option: %s", input)
	}

	title := strings.TrimSpace(parts[0])

	var options []string
	for _, part := range parts[1:] {
		option := strings.TrimSpace(part)
		if option != "" {
			options = append(options, option)
		}
	}

	return title, options, nil
}

func getCommandData(request *http.Request) (*slashRequest.SlashCommandRequest, error) {
	if err := request.ParseForm(); err != nil {
		return nil, err
	}

	cmd := slashRequest.SlashCommandRequest{
		Token:       request.FormValue("token"),
		TeamID:      request.FormValue("team_id"),
		TeamDomain:  request.FormValue("team_domain"),
		ChannelID:   request.FormValue("channel_id"),
		ChannelName: request.FormValue("channel_name"),
		UserID:      request.FormValue("user_id"),
		UserName:    request.FormValue("user_name"),
		Command:     request.FormValue("command"),
		Text:        request.FormValue("text"),
		ResponseURL: request.FormValue("response_url"),
		TriggerID:   request.FormValue("trigger_id"),
	}

	return &cmd, nil
}
