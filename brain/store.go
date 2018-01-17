package brain

// Resource is a storable interface in the brain
// type User struct {
//   id       string
//   Name     string
//   Username string `store:"index"`
// }
//
// func (u User) ID() string {
//   return string(u.id)
// }
//
// func (u User) Namespace() string {
//   return "app:user"
// }
//
// // optionally:
// func (u User) Indices() map[string]string {
//   return map[string]string{
//     "name": u.Name,
//   }
// }
type Resource interface {
	ID() string
	Namespace() string
}

type indexer interface {
	Resource
	Indices() map[string]string
}

type Store interface {
	Get(key string, i interface{}) error
	Set(key string, i interface{}) error
}

// Brain is the store
type Brain struct{ store Store }

// New creates a new connection to Redis
func New(s Store) *Brain {
	return &Brain{store: s}
}
