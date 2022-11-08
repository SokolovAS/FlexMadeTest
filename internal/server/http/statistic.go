package http

import (
	"FlexMadeTest/internal/models"
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

//go:generate mockgen -destination=./statistic_mock_storage.go -package=http -source=./statistic.go

type Validator interface {
	Struct(s interface{}) error
}

type StatisticService interface {
	GetStatistic(ctx context.Context, req models.QueriesRequest) (models.StatisticResultCollection, error)
}

// Body represents message for dynamic json responses.
type Body map[string]interface{}

const (
	queryTypeParam = "type"
	sortingParam   = "sorting"
	pageParam      = "page"
	perPageParam   = "per-page"

	defaultPage    = "1"
	defaultPerPage = "20"
	defaultSorting = "first-slow"
)

// NewStatistic constructor for Statistic transport layer.
func NewStatistic(validator Validator, service StatisticService) *Statistic {
	return &Statistic{
		validator: validator,
		service:   service,
	}
}

// Statistic implements transport layer for Fiber application.
type Statistic struct {
	validator Validator
	service   StatisticService
}

// GetQueriesStatistic represents transport layer for.
// [GET] /database/queries
func (t Statistic) GetQueriesStatistic(ctx *fiber.Ctx) error {
	sPage := ctx.Query(pageParam, defaultPage)
	page, err := strconv.Atoi(sPage)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).
			JSON(NewBadRequestError("error parsing page parameter from request", err))
	}

	sPerPage := ctx.Query(perPageParam, defaultPerPage)
	perPage, err := strconv.Atoi(sPerPage)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).
			JSON(NewBadRequestError("error parsing per-page parameter from request", err))
	}

	req := models.QueriesRequest{
		QueryType: ctx.Query(queryTypeParam),
		Sorting:   ctx.Query(sortingParam, defaultSorting),
		Page:      page,
		PerPage:   perPage,
	}

	if err = t.validator.Struct(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).
			JSON(NewBadRequestError("request validation error", err))
	}

	gc := models.QueriesRequest{
		QueryType: req.QueryType,
		Sorting:   req.Sorting,
		Page:      req.Page,
		PerPage:   req.PerPage,
	}

	data, err := t.service.GetStatistic(ctx.Context(), gc)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(NewInternalServerError("error getting database queries statistic", err))
	}

	resp := make(models.QueriesResponse, 0, len(data))
	for _, row := range data {
		resp = append(resp, &models.QueryRow{
			QueryID:           row.QueryID,
			Query:             row.Query,
			MaxExecutionTime:  row.MaxExecutionTime,
			MeanExecutionTime: row.MeanExecutionTime,
		})
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}
