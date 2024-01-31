package zinc

import (
	"errors"
	"fmt"
)

// ZincQuery
type ZincQuery struct {
	Params     map[string][]string
	SearchType ZincSearchType
}

type ZincSearchType string

const (
	MATCH_QUERY    ZincSearchType = "match"
	MATCHALL_QUERY ZincSearchType = "matchall"
)

var (
	requiredParams = []string{
		"page",
	}

	matchQueryParams = append([]string{
		"term",
	}, requiredParams...)

	matchAllQueryParams = append([]string{}, requiredParams...)

	matchAllQueryTempl = []byte(`{
        "search_type": "matchall",
        "from": 0,
        "max_results": %v,
        "_source": []
    }`)

	matchQueryTempl = []byte(`{
        "search_type": "match",
        "query": {
            "term": "%v",
            "field": "_all"
        },
        "sort_fields": ["-@timestamp"],
        "from": 0,
        "max_results": %v,
        "_source": [ ]
    }`)
)

// BuildQuery
func BuildQuery(zq ZincQuery) (q string, err error) {
	switch zq.SearchType {
	case MATCH_QUERY:
		return buildMatchQuery(zq.Params)
	case MATCHALL_QUERY:
		return buildMatchAllQuery(zq.Params)
	}

	return "", errors.New(fmt.Sprintf("invalid SearchType %s", zq.SearchType))
}

// getParams
func getParams(params map[string][]string, toGet []string) (p []any, err error) {
	for _, tg := range toGet {
		if params[tg] == nil {
			return nil, errors.New(fmt.Sprintf("missing param %s", tg))
		}
		urlParam := params[tg][0]
		p = append(p, urlParam)
	}
	return p, nil
}

// buildMatchQuery
func buildMatchQuery(params map[string][]string) (q string, err error) {
	urlParams, err := getParams(params, matchQueryParams)
	if err != nil {
		return "", err
	}
	q = fmt.Sprintf(string(matchQueryTempl), urlParams...)
	return q, err
}

// buildMatchAllQuery
func buildMatchAllQuery(params map[string][]string) (q string, err error) {
	urlParams, err := getParams(params, matchAllQueryParams)
	if err != nil {
		return "", err
	}
	q = fmt.Sprintf(string(matchAllQueryTempl), urlParams...)
	return q, nil
}
