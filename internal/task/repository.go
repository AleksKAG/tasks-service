package task

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(task Task) (Task, error) {
	result := r.db.Create(&task)
	return task, result.Error
}

func (r *Repository) Get(id uint) (Task, error) {
	var task Task
	result := r.db.First(&task, id)
	return task, result.Error
}

func (r *Repository) Update(task Task) (Task, error) {
	result := r.db.Model(&task).Updates(Task{
		UserID: task.UserID,
		Title:  task.Title,
		IsDone: task.IsDone,
	})
	if result.Error != nil {
		return Task{}, result.Error
	}
	return task, nil
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Task{}, id).Error
}

func (r *Repository) List(offset, limit int) ([]Task, error) {
	var tasks []Task
	result := r.db.Offset(offset).Limit(limit).Find(&tasks)
	return tasks, result.Error
}

func (r *Repository) ListByUser(userID uint, offset, limit int) ([]Task, error) {
	var tasks []Task
	result := r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&tasks)
	return tasks, result.Error
}
