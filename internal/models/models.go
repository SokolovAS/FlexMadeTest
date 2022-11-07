package models

// StatisticFilter request parameter represents filter.
type StatisticFilter struct {
	QueryType string
	Sorting   string
	Page      int
	PerPage   int
}

// StatisticResultRow model for mapping database result.
type StatisticResultRow struct {
	QueryID           int64   `gorm:"column:queryid"`
	Query             string  `gorm:"column:query"`
	MaxExecutionTime  float64 `gorm:"column:max_exec_time"`
	MeanExecutionTime float64 `gorm:"column:mean_exec_time"`
}

// TableName provides mode table name.
func (g StatisticResultRow) TableName() string {
	return "pg_stat_statements"
}

// StatisticResultCollection collection of StatisticResultRow
type StatisticResultCollection []*StatisticResultRow

type QueriesRequest struct {
	QueryType string `validate:"omitempty,oneof=select insert update delete"`
	Sorting   string `validate:"oneof=first-slow first-fast,required"`
	Page      int    `validate:"gt=0,required"`
	PerPage   int    `validate:"gt=0,required"`
}

type QueriesResponse []*QueryRow

type QueryRow struct {
	QueryID           int64   `json:"query_id"`
	Query             string  `json:"query"`
	MaxExecutionTime  float64 `json:"max_execution_time"`
	MeanExecutionTime float64 `json:"mean_execution_time"`
}
