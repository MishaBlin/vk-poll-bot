package cleanup

import (
	"log"
	pollStruct "mm-polls/internal/lib/types/poll"
	"net/http"
)

type PollProvider interface {
	GetAllPolls() ([]pollStruct.Poll, error)
	DeletePoll(uuid string) error
}

func New(provider PollProvider) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		polls, err := provider.GetAllPolls()
		if err != nil {
			log.Println(err)
			http.Error(writer, "Internal error", http.StatusInternalServerError)
			return
		}

		for _, poll := range polls {
			if !poll.Active {
				err = provider.DeletePoll(poll.ID)
				if err != nil {
					log.Println("Unable to delete poll: ", err)
					http.Error(writer, "Internal error", http.StatusInternalServerError)
					return
				}
			}
		}

		writer.WriteHeader(http.StatusOK)
		log.Println("Cleanup successful")
		if err != nil {
			log.Println("Error encoding response", err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
