package zei

const (
	authNone authMode = iota
	authBasic
	authBearer
)

type authMode int
