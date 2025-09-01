package task

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTask(t Task) (Task, error) {
	return s.repo.Create(t)
}

func (s *Service) GetTask(id uint32) (Task, error) {
	return s.repo.Get(id)
}

func (s *Service) UpdateTask(t Task) (Task, error) {
	return s.repo.Update(t)
}

func (s *Service) DeleteTask(id uint32) error {
	return s.repo.Delete(id)
}

func (s *Service) ListTasks(page, pageSize int) ([]Task, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}
