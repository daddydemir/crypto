package strategy

import "net/http"

type TokenStrategy interface {
	Apply(req *http.Request, token string) error
}

type HeaderTokenStrategy struct{}

func (h HeaderTokenStrategy) Apply(req *http.Request, token string) error {
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}

type QueryTokenStrategy struct {
	ParamName string
}

func (q QueryTokenStrategy) Apply(req *http.Request, token string) error {
	query := req.URL.Query()
	query.Set(q.ParamName, token)
	req.URL.RawQuery = query.Encode()
	return nil
}
