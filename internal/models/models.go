package models

// GetStatisticFilter request parameter represents filter.
type GetStatisticFilter struct {
	QueryType string
	Sorting   string
	Page      int
	PerPage   int
}

// GetStatisticResultRow model for mapping database result.
type GetStatisticResultRow struct {
	QueryID           int64   `gorm:"column:queryid"`
	Query             string  `gorm:"column:query"`
	MaxExecutionTime  float64 `gorm:"column:max_exec_time"`
	MeanExecutionTime float64 `gorm:"column:mean_exec_time"`
}

// TableName provides mode table name.
func (g GetStatisticResultRow) TableName() string {
	return "pg_stat_statements"
}

// GetStatisticResultCollection collection of GetStatisticResultRow
type GetStatisticResultCollection []*GetStatisticResultRow

type GetQueriesRequest struct {
	QueryType string `validate:"omitempty,oneof=select insert update delete"`
	Sorting   string `validate:"oneof=first-slow first-fast,required"`
	Page      int    `validate:"gt=0,required"`
	PerPage   int    `validate:"gt=0,required"`
}

type GetQueriesResponse []*QueryRow

type QueryRow struct {
	QueryID           int64   `json:"query_id"`
	Query             string  `json:"query"`
	MaxExecutionTime  float64 `json:"max_execution_time"`
	MeanExecutionTime float64 `json:"mean_execution_time"`
}
