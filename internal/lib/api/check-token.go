package api

import (
	"fmt"
	"os"
)

func AuthenticateToken(token string) error {
	if token != os.Getenv("COMMAND_TOKEN") {
		return fmt.Errorf("received invalid command token")
	}
	return nil
}
