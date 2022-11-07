package http

import (
	"FlexMadeTest/internal/models"
	"errors"
	"github.com/go-playground/assert/v2"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestStatistic_GetQueriesStatistic(t *testing.T) {
	ctrl := gomock.NewController(t)
	testServiceError := errors.New("test service error")

	type fields struct {
		validator Validator
		service   StatisticService
	}
	tests := map[string]struct {
		route         string
		expStatusCode int
		fields        fields
	}{
		"error.wrong page param": {
			route:         "/database/queries?type=update&sorting=first-slow&page=wrong-page-param&per-page=2",
			expStatusCode: http.StatusBadRequest,
			fields:        fields{},
		},
		"error.wrong per-page param": {
			route:         "/database/queries?type=update&sorting=first-slow&page=1&per-page=wrong-per-page-param",
			expStatusCode: http.StatusBadRequest,
			fields:        fields{},
		},
		"error.struct validator error": {
			route:         "/database/queries?type=wrong-type&sorting=wrong-sorting&page=1&per-page=2",
			expStatusCode: http.StatusBadRequest,
			fields: fields{
				validator: func() Validator {
					mock := NewMockValidator(ctrl)

					mock.EXPECT().Struct(&models.GetQueriesRequest{
						QueryType: "wrong-type",
						Sorting:   "wrong-sorting",
						Page:      1,
						PerPage:   2,
					}).Return(validator.ValidationErrors{
						fieldError{
							field: "test1field",
							error: "test1error",
						},
						fieldError{
							field: "test2field",
							error: "test2error",
						},
					})

					return mock
				}(),
			},
		},
		"error.get statistic": {
			route:         "/database/queries?type=select&sorting=first-slow&page=1&per-page=2",
			expStatusCode: http.StatusInternalServerError,
			fields: fields{
				validator: func() Validator {
					mock := NewMockValidator(ctrl)

					mock.EXPECT().Struct(&models.GetQueriesRequest{
						QueryType: "select",
						Sorting:   "first-slow",
						Page:      1,
						PerPage:   2,
					}).Return(nil)

					return mock
				}(),
				service: func() StatisticService {
					mock := NewMockStatisticService(ctrl)
					mock.EXPECT().GetStatistic(gomock.Any(), models.GetQueriesRequest{
						QueryType: "select",
						Sorting:   "first-slow",
						Page:      1,
						PerPage:   2,
					}).Return(nil, testServiceError)

					return mock
				}(),
			},
		},
		"success": {
			route:         "/database/queries?type=select&sorting=first-slow&page=1&per-page=2",
			expStatusCode: http.StatusOK,
			fields: fields{
				validator: func() Validator {
					mock := NewMockValidator(ctrl)

					mock.EXPECT().Struct(&models.GetQueriesRequest{
						QueryType: "select",
						Sorting:   "first-slow",
						Page:      1,
						PerPage:   2,
					}).Return(nil)

					return mock
				}(),
				service: func() StatisticService {
					mock := NewMockStatisticService(ctrl)
					mock.EXPECT().GetStatistic(gomock.Any(), models.GetQueriesRequest{
						QueryType: "select",
						Sorting:   "first-slow",
						Page:      1,
						PerPage:   2,
					}).Return(models.GetStatisticResultCollection{
						{
							QueryID:           1,
							Query:             "TEST Q1",
							MaxExecutionTime:  5.55,
							MeanExecutionTime: 4.25,
						},
						{
							QueryID:           2,
							Query:             "TEST Q2",
							MaxExecutionTime:  4.34,
							MeanExecutionTime: 4.22,
						},
					}, nil)

					return mock
				}(),
			},
		},
	}

	for caseName, tt := range tests {
		t.Run(caseName, func(t *testing.T) {
			app := fiber.New()
			app.Get("/database/queries", NewStatistic(tt.fields.validator, tt.fields.service).GetQueriesStatistic)
			resp, err := app.Test(httptest.NewRequest(
				http.MethodGet,
				tt.route,
				nil,
			))

			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expStatusCode, resp.StatusCode)
		})
	}
}

type fieldError struct {
	field string
	error string
}

func (f fieldError) Tag() string {
	return ""
}

func (f fieldError) ActualTag() string {
	return ""
}

func (f fieldError) Namespace() string {
	return ""
}

func (f fieldError) StructNamespace() string {
	return ""
}

func (f fieldError) Field() string {
	return f.field
}

func (f fieldError) StructField() string {
	return ""
}

func (f fieldError) Value() interface{} {
	return ""
}

func (f fieldError) Param() string {
	return ""
}

func (f fieldError) Kind() reflect.Kind {
	return reflect.Kind(1)
}

func (f fieldError) Type() reflect.Type {
	return nil
}

func (f fieldError) Translate(_ ut.Translator) string {
	return ""
}

func (f fieldError) Error() string {
	return f.error
}
