package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/PosokhovVadim/stawberry/internal/app/apperror"
	"github.com/PosokhovVadim/stawberry/internal/domain/entity"
	"github.com/PosokhovVadim/stawberry/internal/repository/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const transactionContextKey = "stx"

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context, tokenID uuid.UUID) (founded entity.Session, err error) {
	var m model.Session
	if err = r.conn(ctx).Where("token_id = ?", tokenID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Session{}, apperror.ErrSessionNotFound
		}
		return entity.Session{}, apperror.ErrSessionDatabaseError.Internal(err)
	}

	return model.ConvertSessionToEntity(m), nil
}

func (r *Repository) GetForUpdate(ctx context.Context, tokenID uuid.UUID) (founded entity.Session, err error) {
	var m model.Session

	err = r.conn(ctx).
		Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"}).
		Where("token_id = ?", tokenID).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Session{}, apperror.ErrSessionNotFound
		}
		return entity.Session{}, apperror.ErrSessionDatabaseError.Internal(err)
	}

	return model.ConvertSessionToEntity(m), nil
}

func (r *Repository) Create(ctx context.Context, session entity.Session) (id uint, err error) {
	var m = model.ConvertSessionFromEntity(session)

	if err = r.conn(ctx).Create(&m).Error; err != nil {
		return 0, fmt.Errorf("failed to create session: %w", err)
	}

	return m.ID, nil
}

func (r *Repository) Update(ctx context.Context, session entity.Session) (err error) {
	m := model.ConvertSessionFromEntity(session)

	if err = r.conn(ctx).Save(&m).Error; err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

func (r *Repository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, transactionContextKey, tx))
	})
}

func (r *Repository) conn(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionContextKey).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}

	return r.db.WithContext(ctx)
}
