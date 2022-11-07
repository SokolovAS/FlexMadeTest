package services

import (
	"FlexMadeTest/internal/models"
	"FlexMadeTest/internal/repositories"
	"context"
)

// Statistic service implementation.
type Statistic struct {
	repository *repositories.Statistic
}

// NewStatistic constructor produces Statistic.
func NewStatistic(repository *repositories.Statistic) *Statistic {
	return &Statistic{repository: repository}
}

// GetStatistic receives dto.GetQueriesReq and calls repository to get statistic. Returns dto.QueryStatisticCollection.
func (s Statistic) GetStatistic(ctx context.Context, req models.GetQueriesRequest) (models.GetStatisticResultCollection, error) {

	filter := models.GetStatisticFilter{
		QueryType: req.QueryType,
		Sorting:   req.Sorting,
		Page:      req.Page,
		PerPage:   req.PerPage,
	}

	statisticCollection, err := s.repository.GetStatistic(ctx, filter)
	if err != nil {
		return nil, err
	}

	collection := make(models.GetStatisticResultCollection, 0, len(statisticCollection))
	for _, item := range statisticCollection {
		collection = append(collection, &models.GetStatisticResultRow{
			QueryID:           item.QueryID,
			Query:             item.Query,
			MaxExecutionTime:  item.MaxExecutionTime,
			MeanExecutionTime: item.MeanExecutionTime,
		})
	}

	return collection, nil
}
