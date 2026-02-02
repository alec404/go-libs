package esquery

import (
	"fmt"
	"reflect"

	"github.com/olivere/elastic"
)

func NewESQuery(indices []string) *ESQuery {
	return &ESQuery{
		Indices: indices,
		From:    DefaultFrom,
		Size:    DefaultPageSize,
	}
}

func (q *ESQuery) AddMust(query elastic.Query) {
	q.MustQuery = append(q.MustQuery, query)
}

func (q *ESQuery) AddMustNotQuery(query elastic.Query) {
	q.MustNotQuery = append(q.MustNotQuery, query)
}

func (q *ESQuery) AddShouldQuery(query elastic.Query) {
	q.ShouldQuery = append(q.ShouldQuery, query)
}

func (q *ESQuery) AddFilters(query elastic.Query) {
	q.Filters = append(q.Filters, query)
}

func (q *ESQuery) AddSort(fieldSort *elastic.FieldSort) {
	q.Sorters = append(q.Sorters, fieldSort)
}

// 获取第一个排序字段的名称
func (q *ESQuery) GetFirstSorterField() (string, error) {
	// 检查 Sorters 是否为空
	if len(q.Sorters) == 0 {
		return "", fmt.Errorf("sorters is empty")
	}

	// 获取第一个 Sorter
	firstSorter := q.Sorters[0]

	// 类型断言：检查 Sorter 是否为 FieldSort 类型
	fieldSort, ok := firstSorter.(*elastic.FieldSort)
	if !ok {
		return "", fmt.Errorf("first sorter is not a FieldSort type")
	}
	// 使用反射获取 fieldName 的值
	val := reflect.ValueOf(fieldSort).Elem()
	fieldNameField := val.FieldByName("fieldName")

	if !fieldNameField.IsValid() {
		return "", fmt.Errorf("fieldName field not found")
	}

	if fieldNameField.Kind() != reflect.String {
		return "", fmt.Errorf("fieldName is not of type string")
	}

	return fieldNameField.String(), nil
}

func (q *ESQuery) BuildSearchService(client *elastic.Client) *elastic.SearchService {
	searchService := client.Search().
		Index(q.Indices...).
		From(int(q.From)).
		Size(int(q.Size)).
		Query(
			elastic.NewConstantScoreQuery(
				elastic.NewBoolQuery().
					Must(q.MustQuery...).
					MustNot(q.MustNotQuery...).
					Should(q.ShouldQuery...).
					Filter(q.Filters...),
			),
		).
		SortBy(q.Sorters...).
		TrackTotalHits(true)

	if q.NoTotal {
		searchService = searchService.TrackTotalHits(false)
	}

	if len(q.SearchAfter) > 0 {
		searchService = searchService.SearchAfter(q.SearchAfter...)
	}

	if len(q.Include) > 0 {
		searchService = searchService.FetchSourceContext(elastic.NewFetchSourceContext(true).Include(q.Include...))
	}

	return searchService
}
