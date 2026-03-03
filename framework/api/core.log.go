package api

import lib "github.com/aoticombr/golang/lib"

func (api *CoreApi) LogDebug(v ...interface{}) {
	lib.NewLog().Debug(api.Api.Name, v...)
}
func (api *CoreApi) LogInfo(v ...interface{}) {
	lib.NewLog().Info(api.Api.Name, v...)
}
func (api *CoreApi) LogError(v ...interface{}) {
	lib.NewLog().Error(api.Api.Name, v...)
}
func (api *CoreApi) LogWarning(v ...interface{}) {
	lib.NewLog().Warning(api.Api.Name, v...)
}
func (api *CoreApi) LogCritical(v ...interface{}) {
	lib.NewLog().Critical(api.Api.Name, v...)
}
func (api *CoreApi) LogFatal(v ...interface{}) {
	lib.NewLog().Fatal(api.Api.Name, v...)
}
