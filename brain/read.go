package brain

import "errors"

func (c *Brain) readIndices(s Resource) (id string, err error) {
	i, ok := s.(indexer)
	if !ok {
		return "", errors.New("No lookup value present")
	}

	for name, value := range i.Indices() {
		if value == "" {
			continue
		}

		key := i.Namespace() + "_by_" + name + ":" + value
		if err = c.store.Get(key, &id); err == nil {
			break
		}
	}

	return id, err
}

func (c *Brain) Read(s Resource) error {
	var err error
	if s.Namespace() == "" {
		return errors.New("Namespace() cannot be empty")
	}

	id := s.ID()
	if id == "" {
		if id, err = c.readIndices(s); err != nil {
			return err
		}
	}

	if err = c.store.Get(s.Namespace()+":"+id, s); err != nil {
		return err
	}

	return nil
}
