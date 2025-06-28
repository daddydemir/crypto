package coincapadapter

type TokenStrategy string

const (
	QueryStrategy  TokenStrategy = "query"
	HeaderStrategy TokenStrategy = "header"
)
