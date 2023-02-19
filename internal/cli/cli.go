package cli

import (
	"belby/internal/cli/handlers"
	"errors"
)

func Handle(args []string) error {
	if len(args) > 2 {
		return errors.New("wrong args length")
	}

	if len(args) == 1 {
		return handlers.Run()
	}

	switch args[1] {
	case "run":
		return handlers.Run()
	case "migrate":
		return handlers.Migrate()
	}

	return errors.New("wrong args\navailable args: run, migrate")
}
