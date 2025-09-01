package task

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(t Task) (Task, error) {
	result := r.db.Create(&t)
	return t, result.Error
}

func (r *Repository) Get(id uint32) (Task, error) {
	var t Task
	result := r.db.First(&t, id)
	return t, result.Error
}

func (r *Repository) Update(t Task) (Task, error) {
	result := r.db.Model(&Task{}).Where("id = ?", t.ID).Updates(map[string]interface{}{
		"title": t.Title,
	})
	if result.Error != nil {
		return Task{}, result.Error
	}
	return t, nil
}

func (r *Repository) Delete(id uint32) error {
	return r.db.Delete(&Task{}, id).Error
}

func (r *Repository) List(offset, limit int) ([]Task, error) {
	var tasks []Task
	result := r.db.Offset(offset).Limit(limit).Find(&tasks)
	return tasks, result.Error
}
