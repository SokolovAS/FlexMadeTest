package repositories

import (
	"FlexMadeTest/internal/models"
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func TestStatistic_GetStatistic(t *testing.T) {
	successRes := models.StatisticResultCollection{
		{
			QueryID:           1,
			Query:             "TEST UPDATE 1",
			MaxExecutionTime:  33,
			MeanExecutionTime: 12,
		},
		{
			QueryID:           2,
			Query:             "TEST UPDATE 2",
			MaxExecutionTime:  22,
			MeanExecutionTime: 8,
		},
	}

	testError := errors.New("test error")

	type args struct {
		ctx    context.Context
		filter models.StatisticFilter
	}
	tests := map[string]struct {
		mock    func(mock sqlmock.Sqlmock)
		args    args
		want    models.StatisticResultCollection
		wantErr error
	}{
		"error.unsupported sorting type": {
			mock: func(mock sqlmock.Sqlmock) {
				return
			},
			args: args{
				filter: models.StatisticFilter{
					QueryType: "select",
					Sorting:   "bad-sorting-type",
					Page:      1,
					PerPage:   2,
				},
			},
			want:    nil,
			wantErr: ErrUnsupportedSortingType,
		},
		"error.query execution": {
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pg_stat_statements" WHERE starts_with(lower(query), lower($1)) ORDER BY max_exec_time ASC LIMIT 2`)).
					WithArgs("select").
					WillReturnError(testError)
			},
			args: args{
				filter: models.StatisticFilter{
					QueryType: "select",
					Sorting:   "first-fast",
					Page:      1,
					PerPage:   2,
				},
			},
			want:    nil,
			wantErr: testError,
		},
		"success.with type": {
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pg_stat_statements" WHERE starts_with(lower(query), lower($1)) ORDER BY max_exec_time DESC LIMIT 2 OFFSET 2`)).
					WithArgs("update").
					WillReturnRows(
						sqlmock.NewRows([]string{"queryid", "query", "max_exec_time", "mean_exec_time"}).
							AddRow(1, "TEST UPDATE 1", 33, 12).
							AddRow(2, "TEST UPDATE 2", 22, 8),
					)
			},
			args: args{
				filter: models.StatisticFilter{
					QueryType: "update",
					Sorting:   "first-slow",
					Page:      2,
					PerPage:   2,
				},
			},
			want:    successRes,
			wantErr: nil,
		},
		"success.without type": {
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pg_stat_statements" WHERE (starts_with(lower(query), lower('select')) OR starts_with(lower(query), lower('insert')) OR starts_with(lower(query), lower('update')) OR starts_with(lower(query), lower('delete'))) ORDER BY max_exec_time DESC LIMIT 2`)).
					WillReturnRows(
						sqlmock.NewRows([]string{"queryid", "query", "max_exec_time", "mean_exec_time"}).
							AddRow(1, "TEST UPDATE 1", 33, 12).
							AddRow(2, "TEST UPDATE 2", 22, 8),
					)
			},
			args: args{
				filter: models.StatisticFilter{
					QueryType: "",
					Sorting:   "first-slow",
					Page:      1,
					PerPage:   2,
				},
			},
			want:    successRes,
			wantErr: nil,
		},
	}
	for caseName, tt := range tests {
		t.Run(caseName, func(t *testing.T) {

			con, db, mock := NewGormDBMock(t)
			defer con.Close()

			tt.mock(mock)

			got, err := NewStatistic(db).GetStatistic(tt.args.ctx, tt.args.filter)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func NewGormDBMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected on sqlmock.New", err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock",
		DriverName:           "postgres",
		Conn:                 dbConn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("Failed to open gorm v2 db, got error: %v", err)
	}

	return dbConn, db, mock
}
