package utils

import url_utils "net/url"

func UrlParser(URL string, q map[string]string) (string, error) {
	url, err := url_utils.Parse(URL)

	if err != nil {
		return "", err
	}
	query := url.Query()

	for k, v := range q {
		query.Add(k, v)
	}

	url.RawQuery = query.Encode()

	return url.String(), nil
}
