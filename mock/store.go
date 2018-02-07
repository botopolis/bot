package mock

// Store is a mock bot.Store
type Store struct {
	// Called by Get
	GetFunc func(string, interface{}) error
	// Called by Set
	SetFunc func(string, interface{}) error
	// Called by Delete
	DeleteFunc func(string) error
}

// NewStore returns a NOOP Store
func NewStore() *Store {
	return &Store{
		GetFunc:    func(s string, i interface{}) error { return nil },
		SetFunc:    func(s string, i interface{}) error { return nil },
		DeleteFunc: func(s string) error { return nil },
	}
}

// Get delegates to Chat.GetFunc
func (s Store) Get(k string, i interface{}) error { return s.GetFunc(k, i) }

// Set delegates to Chat.SetFunc
func (s Store) Set(k string, i interface{}) error { return s.SetFunc(k, i) }

// Delete delegates to Chat.DeleteFunc
func (s Store) Delete(k string) error { return s.DeleteFunc(k) }
