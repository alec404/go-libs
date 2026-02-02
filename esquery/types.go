package esquery

import "github.com/olivere/elastic"

type ESQuery struct {
	MustQuery    []elastic.Query
	MustNotQuery []elastic.Query
	ShouldQuery  []elastic.Query
	Filters      []elastic.Query
	Sorters      []elastic.Sorter
	From         int32 // 分页
	Size         int32
	Indices      []string
	SearchAfter  []interface{}
	Include      []string
	NoTotal      bool // 不需要总数
}

const (
	DefaultFrom     = 0  // 默认页码
	DefaultPageSize = 10 // 默认每页行数
)
