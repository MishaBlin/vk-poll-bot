package results

import (
	"encoding/json"
	"log"
	"mm-polls/internal/lib/api"
	pollStruct "mm-polls/internal/lib/types/poll"
	"mm-polls/internal/lib/types/response"
	slashRequest "mm-polls/internal/lib/types/slash-request"
	"mm-polls/internal/lib/util"
	"net/http"
)

type PollProvider interface {
	GetPoll(uuid string) (*pollStruct.Poll, error)
}

func New(provider PollProvider) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var cmd slashRequest.SlashCommandRequest
		if err := json.NewDecoder(request.Body).Decode(&cmd); err != nil {
			log.Println("failed to decode slash command:", err)
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}

		err := api.AuthenticateToken(cmd.Context.Token)
		if err != nil {
			log.Println(err)
			http.Error(writer, "Invalid token", http.StatusUnauthorized)
			return
		}

		poll, err := provider.GetPoll(cmd.Context.PollID)
		if err != nil {
			log.Println("failed to get poll from db:", err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		summary := util.SummarizePoll(poll)

		resp := response.EphemeralResponse{EphemeralText: summary}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		log.Printf("results sent for %s\n", poll.Title)
		err = json.NewEncoder(writer).Encode(resp)
		if err != nil {
			log.Println("Error encoding response", err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
