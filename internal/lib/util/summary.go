package util

import (
	"fmt"
	"math"
	pollStruct "mm-polls/internal/lib/types/poll"
)

func SummarizePoll(poll *pollStruct.Poll) string {
	totalVotes := 0
	for _, voteCount := range poll.Votes {
		totalVotes += voteCount
	}

	var message string
	if poll.Active {
		message = fmt.Sprintf("Current summary for poll: %s", poll.Title)
	} else {
		message = fmt.Sprintf("Closed summary for poll: %s", poll.Title)
	}

	message += "\n\n"

	message += "| Option | Votes"
	message += " | Percent |\n"
	message += "| :----- | -----: | -----: |\n"

	for i, option := range poll.Options {
		voteCount := 0
		if i < len(poll.Votes) {
			voteCount = poll.Votes[i]
		}

		var percentage int
		if totalVotes > 0 {
			percentage = int(math.Round(100 * float64(voteCount) / float64(totalVotes)))
		}

		message += fmt.Sprintf("| %s | %d | %d%% |\n", option, voteCount, percentage)
	}

	return message
}
