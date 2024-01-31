package zinc

import (
	"errors"
	"fmt"
)

// ZincQuery
type ZincQuery struct {
	Query      []byte
	Params     map[string][]string
	SearchType ZincSearchType
}

type ZincSearchType string

const (
	MATCH_QUERY    ZincSearchType = "match"
	MATCHALL_QUERY ZincSearchType = "matchall"

	DEFAULT_PAGE_SIZE = 10
)

var (
	requiredParams = []string{
		"page",
	}

	matchQueryParams = append([]string{
		"term",
	}, requiredParams...)

	matchAllQueryParams = append([]string{}, requiredParams...)
)

// BuildQuery
func BuildQuery(zq ZincQuery) (q string, err error) {
	/*
		if validPage() != nil {
			numResults := DEFAULT_PAGE_SIZE
			n, err := strconv.Atoi(page)
			if err != nil {
				return "", err
			}
			numResults = n

		}
	*/
	switch zq.SearchType {
	case MATCH_QUERY:
		return buildMatchQuery(zq.Query, zq.Params)
	case MATCHALL_QUERY:
		return buildMatchAllQuery(zq.Query, zq.Params)
	}

	return "", nil
}

// getParams
func getParams(params map[string][]string, toGet []string) (p []any, err error) {
	for _, tg := range toGet {
		urlParam := params[tg][0]
		if urlParam == "" {
			return nil, errors.New(fmt.Sprintf("missing param %s", p))
		}
		p = append(p, urlParam)
	}
	return p, nil
}

// buildMatchQuery
func buildMatchQuery(queryTemp []byte, params map[string][]string) (q string, err error) {
	urlParams, err := getParams(params, matchQueryParams)
	if err != nil {
		return "", err
	}
	q = fmt.Sprintf(string(queryTemp), urlParams...)
	return q, err
}

// buildMatchAllQuery
func buildMatchAllQuery(queryTemp []byte, params map[string][]string) (q string, err error) {
	urlParams, err := getParams(params, matchAllQueryParams)
	if err != nil {
		return "", err
	}
	q = fmt.Sprintf(string(queryTemp), urlParams...)
	return q, nil
}
