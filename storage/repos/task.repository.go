package repos

import (
	"clean-rest-arch/internal/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task *models.TaskEntity) (uint, error)
	GetTask(userId uint, taskId uint) (*models.TaskEntity, error)
	GetAllUserTasks(userId uint) ([]*models.TaskEntity, error)
	UpdateTask(task *models.TaskEntity) error
	DeleteTask(userId uint, taskId uint) error
}

type taskRepo struct {
	database *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepo{database: db}
}

func (r *taskRepo) CreateTask(task *models.TaskEntity) (uint, error) {
	const op = "storage.repos.CreateTask"

	var maxTaskID uint
	r.database.Model(&models.TaskEntity{}).
		Select("MAX(task_id)").
		Where("user_id = ?", task.UserId).
		Scan(&maxTaskID)
	task.TaskId = maxTaskID + 1

	result := r.database.Create(task)
	if result.Error != nil {
		return 0, fmt.Errorf("%s: %w", op, result.Error)
	}

	return task.Id, nil
}

func (r *taskRepo) GetTask(userId uint, taskId uint) (*models.TaskEntity, error) {
	const op = "storage.repos.GetTask"

	var task models.TaskEntity
	result := r.database.Where("user_id = ? AND task_id = ?", userId, taskId).First(&task)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, fmt.Errorf("%s: %w", op, result.Error)
	}

	return &task, nil
}

func (r *taskRepo) GetAllUserTasks(userId uint) ([]*models.TaskEntity, error) {
	const op = "storage.repos.GetAllUserTasks"

	var tasks []*models.TaskEntity
	result := r.database.Where("user_id = ?", userId).Find(&tasks)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, fmt.Errorf("%s: %w", op, result.Error)
	}

	return tasks, nil
}

func (r *taskRepo) UpdateTask(task *models.TaskEntity) error {
	const op = "storage.repos.UpdateTask"

	var existTask models.TaskEntity
	result := r.database.Where("user_id = ? AND task_id = ?", task.UserId, task.TaskId).First(&existTask)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return result.Error
		}
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	existTask.Title = task.Title
	existTask.Data = task.Data
	existTask.UpdatedAt = time.Now()

	result = r.database.Save(existTask)
	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	return nil
}

func (r *taskRepo) DeleteTask(userId uint, taskId uint) error {
	const op = "storage.repos.DeleteTask"

	result := r.database.Where(&models.TaskEntity{UserId: userId, TaskId: taskId}).Delete(&models.TaskEntity{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	return nil
}
