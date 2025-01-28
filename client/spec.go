package client

import "errors"

type Spec struct {
	Key string `json:"key"`
}

func (s *Spec) SetDefaults() {
}

func (s *Spec) Validate() error {
	if s.Key == "" {
		return errors.New("key is required")
	}
	return nil
}
