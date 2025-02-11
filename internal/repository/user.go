package repository

import (
	"context"
	"errors"

	"github.com/PosokhovVadim/stawberry/internal/app/apperror"
	"github.com/PosokhovVadim/stawberry/internal/domain/entity"
	"github.com/PosokhovVadim/stawberry/internal/repository/model"
	"gorm.io/gorm"
)

const transactionContextKey = "utx"

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByEmail - возвращает найденного по email пользователя.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (founded entity.User, err error) {
	var userModel model.User
	if err = r.conn(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, apperror.ErrAuthUserNotFound
		}
		return entity.User{}, &apperror.AuthError{
			Code:    apperror.DatabaseError,
			Message: "failed to get user by email",
			Err:     err,
		}
	}

	return model.ConvertUserToEntity(userModel), nil
}

// Insert - сохраняет пользователя в бд.
func (r *UserRepository) Save(ctx context.Context, entity entity.Register) (id uint, err error) {
	userModel := model.ConvertUserFromSvc(entity)
	if err = r.conn(ctx).Create(&userModel).Error; err != nil {
		if isDuplicateError(err) {
			return 0, &apperror.AuthError{
				Code:    apperror.DuplicateError,
				Message: "user already exists",
				Err:     err,
			}
		}
		return 0, &apperror.AuthError{
			Code:    apperror.DatabaseError,
			Message: "failed to create user",
			Err:     err,
		}
	}

	return userModel.ID, nil
}

func (r *UserRepository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, transactionContextKey, tx))
	})
}

func (r *UserRepository) conn(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionContextKey).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}

	return r.db.WithContext(ctx)
}
