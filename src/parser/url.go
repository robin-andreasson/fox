package parser

import (
	"regexp"
	"strings"
)

func Url(path string, h_path string, h_path_rex string, h_params [][]string) (bool, map[string]string, map[string]string) {

	path, query_string, found := strings.Cut(path, "?")

	params := map[string]string{}

	if path != h_path {

		if params = getParams(path, h_path, h_path_rex, h_params); params == nil {

			return false, nil, nil
		}

	}

	queries := getQueryStrings(query_string, found)

	return true, queries, params
}

func getParams(path string, h_path string, h_path_rex string, h_params [][]string) map[string]string {

	params := map[string]string{}

	rex := regexp.MustCompile(h_path_rex)

	result := rex.FindStringSubmatch(path)

	if len(result) == 0 || result[0] != path {

		//if no params, then
		if rex.FindString(path) != "" {
			return params
		}

		return nil
	}

	for i := 1; i < len(result); i++ {

		key := urldecode(h_params[i-1][0])
		value := urldecode(result[i])

		if len(h_params[i-1]) > 1 {
			rex := regexp.MustCompile("^" + h_params[i-1][1] + "$")

			result := rex.FindString(value)

			if result == "" {
				return nil
			}
		}

		params[key] = value
	}

	return params
}

func getQueryStrings(query_string string, found bool) map[string]string {

	if !found || query_string == "" {
		return map[string]string{}
	}

	rex := regexp.MustCompile("=|&")

	query_seg := rex.Split(query_string, -1)

	queries := map[string]string{}

	length := len(query_seg)

	if length%2 != 0 {
		length--
	}

	for i := 0; i < length; i += 2 {

		key := urldecode(query_seg[i])
		value := urldecode(query_seg[i+1])

		queries[key] = value
	}

	return queries

}
