package utils

import url_utils "net/url"

func UrlParser(URL string, q string) (string, error) {
	url, err := url_utils.Parse(URL)

	if err != nil {
		return "", err
	}
	query := url.Query()

	query.Add("q", q)

	url.RawQuery = query.Encode()

	return url.String(), nil
}
