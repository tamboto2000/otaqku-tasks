package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/mail"
	"regexp"

	"github.com/tamboto2000/otaqku-tasks/internal/common"
	"github.com/tamboto2000/otaqku-tasks/internal/dto"
	"github.com/vinovest/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var userNameRegex = regexp.MustCompile(`^[A-Za-z]+(?:[ '-][A-Za-z]+)*$`)

type Account struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password []byte `db:"password"`
}

func NewAccount(req dto.CreateAccountRequest) (Account, error) {
	errValidation := common.Error{
		Code:    common.ErrCodeInputValidation,
		Message: "invalid input",
	}

	// Validate name
	errName := common.FieldError{Name: "name"}
	if len(req.Name) == 0 {
		errName.Messages = append(errName.Messages, "name can not be empty")
	}

	if len(req.Name) > 100 {
		errName.Messages = append(errName.Messages, "name length can not be greater than 100")
	}

	if !userNameRegex.MatchString(req.Name) {
		errName.Messages = append(errName.Messages, "name can only contain letters, spaces, hyphens, and apostrophes")
	}

	if len(errName.Messages) != 0 {
		errValidation.Fields = append(errValidation.Fields, errName)
	}

	// Validate email
	errEmail := common.FieldError{Name: "email"}
	if len(req.Email) == 0 {
		errEmail.Messages = append(errEmail.Messages, "email can not be empty")
	}

	if len(req.Email) > 100 {
		errEmail.Messages = append(errEmail.Messages, "email length can not be greater than 100")
	}

	emailAddr, err := mail.ParseAddress(req.Email)
	if err != nil {
		errEmail.Messages = append(errEmail.Messages, "email is not valid")
	}

	req.Email = emailAddr.Address

	if len(errEmail.Messages) != 0 {
		errValidation.Fields = append(errValidation.Fields, errEmail)
	}

	// Validate password
	errPwd := common.FieldError{Name: "password"}
	if len(req.Password) == 0 {
		errPwd.Messages = append(errPwd.Messages, "password can not be empty")
	}

	if len(req.Password) > 100 {
		errPwd.Messages = append(errPwd.Messages, "password length can not be greater than 100")
	}

	if len(errPwd.Messages) != 0 {
		errValidation.Fields = append(errValidation.Fields, errPwd)
	}

	if len(errValidation.Fields) != 0 {
		return Account{}, errValidation
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, fmt.Errorf("build account error: %v", err)
	}

	acc := Account{
		Name:     req.Name,
		Email:    req.Email,
		Password: pwdHash,
	}

	return acc, nil
}

func (acc Account) MatchPassword(pwd string) error {
	return bcrypt.CompareHashAndPassword(acc.Password, []byte(pwd))
}

type AccountRepository interface {
	Save(ctx context.Context, acc Account) error
	GetByEmail(ctx context.Context, email string) (Account, error)
	IsExistsByEmail(ctx context.Context, email string) (bool, error)
}

type PostgreAccountRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewPostgreAccountRepository(db *sqlx.DB, logger *slog.Logger) PostgreAccountRepository {
	return PostgreAccountRepository{db: db, logger: logger}
}

func (repo PostgreAccountRepository) Save(ctx context.Context, acc Account) error {
	q := `INSERT INTO accounts (name, email, password) VALUES ($1, $2, $3)`
	_, err := repo.db.ExecContext(ctx, q, acc.Name, acc.Email, acc.Password)
	if err != nil {
		repo.logger.Error(fmt.Sprintf("error on saving account to database: %v", err))
		return err
	}

	return nil
}

func (repo PostgreAccountRepository) GetByEmail(ctx context.Context, email string) (Account, error) {
	q := `SELECT id, name, email, password FROM accounts WHERE email = $1`
	row := repo.db.QueryRowxContext(ctx, q, email)
	var acc Account
	if err := row.StructScan(&acc); err != nil {
		if err == sql.ErrNoRows {
			return Account{}, common.ErrNotFound
		}

		repo.logger.Error(fmt.Sprintf("error on fetching account by email: %v", err))
		return Account{}, err
	}

	return acc, nil
}

func (repo PostgreAccountRepository) IsExistsByEmail(ctx context.Context, email string) (bool, error) {
	q := `SELECT id FROM accounts WHERE email = $1`
	var id int
	row := repo.db.QueryRowContext(ctx, q, email)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		repo.logger.Error(fmt.Sprintf("error on checking account existence by email: %v", err), slog.String("email", email))
		return false, err
	}

	return true, nil
}
