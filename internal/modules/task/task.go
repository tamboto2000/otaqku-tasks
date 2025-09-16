package task

import (
	"context"
	"database/sql"
	"errors"
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
	GetByAccountID(ctx context.Context, accId int, paginate dto.Pagination) (TaskList, error)
	GetByAccountIDAndID(ctx context.Context, accId, id int) (Task, error)
	DeleteByAccountIDAndID(ctx context.Context, accId, id int) error
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

func (repo PostgreTaskRepository) GetByAccountID(ctx context.Context, accId int, paginate dto.Pagination) (TaskList, error) {
	q := `SELECT id, title, status, created_at, updated_at FROM tasks WHERE account_id = $1 AND deleted_at IS NULL LIMIT $2 OFFSET $3`

	var taskList TaskList
	rows, err := repo.db.QueryxContext(ctx, q, accId, paginate.PageSize, common.GetOffset(paginate.Page, paginate.PageSize))
	if err != nil {
		if err == sql.ErrNoRows {
			return TaskList{}, common.ErrNotFound
		}

		repo.logger.Error(fmt.Sprintf("error on fetching list of tasks: %v", err), slog.Int("account_id", accId))
		return TaskList{}, err
	}

	for rows.Next() {
		var task Task
		if err := rows.StructScan(&task); err != nil {
			repo.logger.Error(fmt.Sprintf("error on scanning task from database to struct: %v", err))
			return TaskList{}, err
		}

		taskList.Tasks = append(taskList.Tasks, task)
	}

	// Get the total count of account's tasks
	q = `SELECT COUNT(id) FROM tasks WHERE account_id = $1 AND deleted_at IS NULL`
	row := repo.db.QueryRowContext(ctx, q, accId)
	var totalCount int
	if err := row.Scan(&totalCount); err != nil {
		repo.logger.Error(fmt.Sprintf("error on fetching account's total task count: %v", err), slog.Int("account_id", accId))
		return TaskList{}, err
	}

	pageMetaData := dto.PaginationMetadata{
		Pagination: paginate,
		Total:      totalCount,
	}

	taskList.Pagination = pageMetaData

	return taskList, nil
}

func (repo PostgreTaskRepository) GetByAccountIDAndID(ctx context.Context, accId, id int) (Task, error) {
	q := `SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE account_id = $1 AND id = $2 AND deleted_at IS NULL`

	var task Task
	row := repo.db.QueryRowxContext(ctx, q, accId, id)
	if err := row.StructScan(&task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, common.ErrNotFound
		}

		repo.logger.Error(fmt.Sprintf("error on fetching a sigle task: %v", err), slog.Int("account_id", accId), slog.Int("id", id))
		return Task{}, err
	}

	return task, nil
}

func (repo PostgreTaskRepository) DeleteByAccountIDAndID(ctx context.Context, accId, id int) error {
	q := `UPDATE tasks SET deleted_at = CURRENT_TIMESTAMP WHERE account_id = $1 AND id = $2`

	_, err := repo.db.ExecContext(ctx, q, accId, id)
	if err != nil {
		repo.logger.Error(fmt.Sprintf("error on deleting a task: %v", err), slog.Int("id", id))
		return err
	}

	return nil
}
