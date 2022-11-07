package services

import (
	"FlexMadeTest/internal/models"
	"context"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestStatistic_GetStatistic(t *testing.T) {

	mockCtrl := gomock.NewController(t)

	testError := errors.New("test error")

	testCtx := context.Background()
	testReq := models.QueriesRequest{
		QueryType: "select",
		Sorting:   "asc",
		Page:      1,
		PerPage:   2,
	}

	type fields struct {
		repository StatisticRepository
	}
	type args struct {
		ctx context.Context
		req models.QueriesRequest
	}
	tests := map[string]struct {
		fields  fields
		args    args
		want    models.StatisticResultCollection
		wantErr error
	}{
		"error.getting statistic": {
			fields: fields{
				repository: func() StatisticRepository {
					repo := NewMockStatisticRepository(mockCtrl)
					repo.EXPECT().GetStatistic(testCtx, models.StatisticFilter{
						QueryType: "select",
						Sorting:   "asc",
						Page:      1,
						PerPage:   2,
					}).Return(nil, testError)
					return repo
				}(),
			},
			args: args{
				ctx: testCtx,
				req: testReq,
			},
			want:    nil,
			wantErr: testError,
		},
		"success": {
			fields: fields{
				repository: func() StatisticRepository {
					repo := NewMockStatisticRepository(mockCtrl)
					repo.EXPECT().GetStatistic(testCtx, models.StatisticFilter{
						QueryType: "select",
						Sorting:   "asc",
						Page:      1,
						PerPage:   2,
					}).Return(models.StatisticResultCollection{
						{
							QueryID:           1,
							Query:             "TEST Q1",
							MaxExecutionTime:  2.22,
							MeanExecutionTime: 1.86,
						},
						{
							QueryID:           2,
							Query:             "TEST Q2",
							MaxExecutionTime:  3.56,
							MeanExecutionTime: 2.92,
						},
					}, nil)

					return repo
				}(),
			},
			args: args{
				ctx: testCtx,
				req: testReq,
			},
			want: models.StatisticResultCollection{
				{
					QueryID:           1,
					Query:             "TEST Q1",
					MaxExecutionTime:  2.22,
					MeanExecutionTime: 1.86,
				},
				{
					QueryID:           2,
					Query:             "TEST Q2",
					MaxExecutionTime:  3.56,
					MeanExecutionTime: 2.92,
				},
			},
			wantErr: nil,
		},
	}
	for caseName, tt := range tests {
		t.Run(caseName, func(t *testing.T) {
			got, err := NewStatistic(tt.fields.repository).GetStatistic(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
