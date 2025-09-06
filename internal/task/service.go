package task

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTask(task Task) (Task, error) {
	return s.repo.Create(task)
}

func (s *Service) GetTask(id uint) (Task, error) {
	return s.repo.Get(id)
}

func (s *Service) UpdateTask(id uint, task Task) (Task, error) {
	existing, err := s.repo.Get(id)
	if err != nil {
		return Task{}, err
	}
	existing.UserID = task.UserID
	existing.Title = task.Title
	existing.IsDone = task.IsDone
	return s.repo.Update(existing)
}

func (s *Service) DeleteTask(id uint) error {
	return s.repo.Delete(id)
}

func (s *Service) ListTasks(page, pageSize int) ([]Task, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}

func (s *Service) ListTasksByUser(userID uint, page, pageSize int) ([]Task, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListByUser(userID, offset, pageSize)
}
