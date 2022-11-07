package services

import (
	"FlexMadeTest/internal/models"
	"context"
)

//go:generate mockgen -destination=./statistic_mock_test.go -package=services -source=./statistic.go
type StatisticRepository interface {
	GetStatistic(ctx context.Context, filter models.StatisticFilter) (models.StatisticResultCollection, error)
}

// Statistic service implementation.
type Statistic struct {
	repository StatisticRepository
}

// NewStatistic constructor produces Statistic.
func NewStatistic(repository StatisticRepository) Statistic {
	return Statistic{repository: repository}
}

// GetStatistic receives dto.GetQueriesReq and calls repository to get statistic. Returns dto.QueryStatisticCollection.
func (s Statistic) GetStatistic(ctx context.Context, req models.QueriesRequest) (models.StatisticResultCollection, error) {

	filter := models.StatisticFilter{
		QueryType: req.QueryType,
		Sorting:   req.Sorting,
		Page:      req.Page,
		PerPage:   req.PerPage,
	}

	statisticCollection, err := s.repository.GetStatistic(ctx, filter)
	if err != nil {
		return nil, err
	}

	collection := make(models.StatisticResultCollection, 0, len(statisticCollection))
	for _, item := range statisticCollection {
		collection = append(collection, &models.StatisticResultRow{
			QueryID:           item.QueryID,
			Query:             item.Query,
			MaxExecutionTime:  item.MaxExecutionTime,
			MeanExecutionTime: item.MeanExecutionTime,
		})
	}

	return collection, nil
}
