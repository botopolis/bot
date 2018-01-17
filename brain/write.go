package brain

import (
	"errors"
)

func (c *Brain) writeIndices(s Resource) error {
	i, ok := s.(indexer)
	if !ok {
		return nil
	}

	for name, value := range i.Indices() {
		if value == "" {
			continue
		}

		key := i.Namespace() + "_by_" + name + ":" + value
		if err := c.store.Set(key, s.ID()); err != nil {
			return err
		}
	}

	return nil
}

func (c *Brain) Write(s Resource) error {
	if s.ID() == "" || s.Namespace() == "" {
		return errors.New("ID() and Namespace() cannot be empty")
	}
	key := s.Namespace() + ":" + s.ID()
	if err := c.store.Set(key, s); err != nil {
		return err
	}

	return c.writeIndices(s)
}
