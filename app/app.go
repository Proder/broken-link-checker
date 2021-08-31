package app

import "broken-link-checker/app/config"

func Run() error {
	// Get the server settings
	_ = config.Get()

	return nil
}
