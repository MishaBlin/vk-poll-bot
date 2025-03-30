package create

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mattermost/mattermost-server/v6/model"
	"log"
	"mm-polls/internal/lib/api"
	pollStruct "mm-polls/internal/lib/types/poll"
	"net/http"
)

type PollProvider interface {
	CreatePoll(poll *pollStruct.Poll) error
}

func New(provider PollProvider) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		cmd, err := getCommandData(request)
		if err != nil {
			log.Println("error decoding request", err)
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}

		err = api.AuthenticateToken(cmd.Token)
		if err != nil {
			log.Println(err)
			http.Error(writer, "Invalid token", http.StatusUnauthorized)
			return
		}

		title, opts, err := parsePoll(cmd.Text)
		if err != nil {
			log.Println("Error parsing poll", err)
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}

		poll := &pollStruct.Poll{
			ID:      uuid.New().String(),
			Title:   title,
			OwnerID: cmd.UserID,
			Options: opts,
			Votes:   make([]int, len(opts)),
			Voters:  make(map[string]int),
			Active:  true,
		}

		err = provider.CreatePoll(poll)
		if err != nil {
			log.Println("Error creating poll", err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		actions := getActions(poll, cmd.Token)

		resp := model.CommandResponse{
			Text:         "*A new poll is started!*",
			ResponseType: "in_channel",
			ChannelId:    cmd.ChannelID,
			Attachments:  []*model.SlackAttachment{{Title: poll.Title, Actions: actions}},
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		log.Printf("created poll %s\n", poll.ID)
		err = json.NewEncoder(writer).Encode(resp)
		if err != nil {
			log.Println("Error encoding response", err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
