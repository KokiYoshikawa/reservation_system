package repository

type Repository struct{}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) GetRootMessage() string {
	return "reservation system backend is running"
}

func (r *Repository) GetHealthStatus() string {
	return "ok"
}
