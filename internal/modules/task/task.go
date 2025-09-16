package task

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tamboto2000/otaqku-tasks/internal/common"
	"github.com/tamboto2000/otaqku-tasks/internal/dto"
	"github.com/vinovest/sqlx"
)

// Taks statuses
const (
	StatusTODO       = "todo"
	StatusInProgress = "in_progress"
	StatusBlocked    = "blocked"
	StatusDone       = "done"
	StatusAbandoned  = "abandoned"
)

type Task struct {
	ID          int       `db:"id"`
	AccountID   int       `db:"account_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func NewTask(accId int, req dto.CreateTaskRequest) (Task, error) {
	errValidation := common.Error{
		Code:    common.ErrCodeInputValidation,
		Message: "Invalid input",
	}

	// Validate title
	var errTitle common.FieldError
	if len(req.Title) == 0 {
		errTitle.Messages = append(errTitle.Messages, "title can not be empty")
	}

	if len(req.Title) > 100 {
		errTitle.Messages = append(errTitle.Messages, "title can not be longer than 100 characters")
	}

	if len(errTitle.Messages) != 0 {
		errValidation.Fields = append(errValidation.Fields, errTitle)
	}

	if len(errValidation.Fields) != 0 {
		return Task{}, errValidation
	}

	return Task{
		AccountID:   accId,
		Title:       req.Title,
		Description: req.Description,
		Status:      StatusTODO,
	}, nil
}

type TaskList struct {
	Tasks      []Task
	Pagination dto.PaginationMetadata
}

type TaskRepository interface {
	Save(ctx context.Context, task Task) error
}

type PostgreTaskRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewPostgreTaskRepository(db *sqlx.DB, logger *slog.Logger) PostgreTaskRepository {
	return PostgreTaskRepository{db: db, logger: logger}
}

func (repo PostgreTaskRepository) Save(ctx context.Context, task Task) error {
	q := `INSERT INTO tasks (account_id, title, description, status) 
		VALUES (:account_id, :title, :description, :status)`

	_, err := repo.db.NamedExecContext(ctx, q, task)
	if err != nil {
		repo.logger.Error(fmt.Sprintf("error on saving task to database: %v", err), slog.Any("task", task))
		return err
	}

	return nil
}
