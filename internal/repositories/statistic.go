package repositories

import (
	"FlexMadeTest/internal/models"
	"context"
	"errors"
	"gorm.io/gorm"
)

// ErrUnsupportedSortingType unsupported sorting type error.
var ErrUnsupportedSortingType = errors.New("unsupported sorting type")

// Statistic repository postgres implementation.
type Statistic struct {
	db *gorm.DB
}

// NewStatistic constructor for NewStatistic.
func NewStatistic(db *gorm.DB) Statistic {
	return Statistic{db: db}
}

// GetStatistic builds request to the database by provided filter and returns models.GetStatisticResultCollection.
func (s Statistic) GetStatistic(ctx context.Context, filter models.GetStatisticFilter) (models.GetStatisticResultCollection, error) {
	col := make(models.GetStatisticResultCollection, 0)

	q := s.db.WithContext(ctx)

	if filter.QueryType == "" {
		q = q.Where(
			q.Where("starts_with(lower(query), lower('select'))").
				Or("starts_with(lower(query), lower('insert'))").
				Or("starts_with(lower(query), lower('update'))").
				Or("starts_with(lower(query), lower('delete'))"),
		)
	} else {
		q = q.Where("starts_with(lower(query), lower(?))", filter.QueryType)
	}

	switch filter.Sorting {
	case "first-fast":
		q = q.Order("max_exec_time ASC")
	case "first-slow":
		q = q.Order("max_exec_time DESC")
	default:
		return nil, ErrUnsupportedSortingType
	}

	q = q.Limit(filter.PerPage).Offset((filter.Page - 1) * filter.PerPage)

	if err := q.Find(&col).Error; err != nil {
		return nil, err
	}

	return col, nil
}
