package vote

import (
	"encoding/json"
	"fmt"
	"log"
	"mm-polls/internal/lib/api"
	pollStruct "mm-polls/internal/lib/types/poll"
	"mm-polls/internal/lib/types/response"
	slashRequest "mm-polls/internal/lib/types/slash-request"
	"net/http"
)

type PollProvider interface {
	UpdatePoll(poll *pollStruct.Poll) error
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

		if !poll.Active {
			resp := response.EphemeralResponse{EphemeralText: "Cannot vote! Closed poll"}
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)

			log.Println("tried to vote in closed poll")

			err = json.NewEncoder(writer).Encode(resp)
			if err != nil {
				log.Println("Error encoding response", err)
				http.Error(writer, "Internal server error", http.StatusInternalServerError)
				return
			}
			return
		}

		replaceMessage := ""

		optIdx, exists := poll.Voters[cmd.UserID]
		if exists {
			replaceMessage = "Your previous vote is replaced. "
			log.Println("Replacing previous vote on poll")
			poll.Votes[optIdx]--
		}

		poll.Voters[cmd.UserID] = cmd.Context.OptionIndex
		poll.Votes[cmd.Context.OptionIndex]++

		err = provider.UpdatePoll(poll)
		if err != nil {
			log.Println("Error placing vote to db", err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		resp := response.EphemeralResponse{
			EphemeralText: fmt.Sprintf(
				"Thanks for your vote! %sYour vote for %s is recorded!",
				replaceMessage,
				poll.Options[cmd.Context.OptionIndex],
			)}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		log.Printf("created vote for %s\n", poll.Options[cmd.Context.OptionIndex])
		err = json.NewEncoder(writer).Encode(resp)
		if err != nil {
			log.Println("Error encoding response", err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
