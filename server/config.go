package main

import (
	"fmt"
)

type Configuration struct {
	BASE_URL     string
	SALT         string
	CallBack_URL string
}

func (c *Configuration) IsValid() error {
	if len(c.BASE_URL) == 0 {
		return fmt.Errorf("BASE URL is not configured.")
	} else if len(c.SALT) == 0 {
		return fmt.Errorf("SALT is not configured.")
	} else if len(c.CallBack_URL) == 0 {
		return fmt.Errorf("Callback URL is not configured.")
	}

	return nil
}
