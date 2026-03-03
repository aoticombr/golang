package http

type AuthorizationType int

const (
	AT_AutoDetect AuthorizationType = iota
	AT_Basic
	AT_Bearer
	AT_Auth2
	AT_Nenhum
)
