package http

type AuthorizationType int

const (
	AutoDetect AuthorizationType = iota
	Basic
	Bearer
	Nenhum
)
