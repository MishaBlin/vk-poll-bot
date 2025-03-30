package finish

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
	UpdatePoll(poll *pollStruct.Poll) error
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

		if poll.OwnerID != cmd.UserID {
			log.Println("not owner tried to close poll")

			resp := response.EphemeralResponse{EphemeralText: "Cannot close poll! You are not the owner!"}
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			err = json.NewEncoder(writer).Encode(resp)
			if err != nil {
				log.Println("Error encoding response", err)
				http.Error(writer, "Internal server error", http.StatusInternalServerError)
				return
			}
			return
		}

		if !poll.Active {
			log.Println("tried closing a closed poll")
			http.Error(writer, "Cannot close a closed poll", http.StatusBadRequest)
			return
		}

		poll.Active = false
		err = provider.UpdatePoll(poll)
		if err != nil {
			log.Println("failed to update poll:", err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		summary := util.SummarizePoll(poll)

		upd := response.Update{
			Message: summary,
			Props:   struct{}{},
		}
		resp := response.UpdateResponse{Update: upd}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		log.Printf("ended poll %s\n", poll.Title)
		err = json.NewEncoder(writer).Encode(resp)
		if err != nil {
			log.Println("Error encoding response", err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
