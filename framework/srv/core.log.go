package srv

import "github.com/aoticombr/golang/lib"

func (srv *CoreSrv) LogDebug(v ...interface{}) {
	lib.NewLog().Debug(srv.Srv.Name, v...)
}
func (srv *CoreSrv) LogInfo(v ...interface{}) {
	lib.NewLog().Info(srv.Srv.Name, v...)
}
func (srv *CoreSrv) LogError(v ...interface{}) {
	lib.NewLog().Error(srv.Srv.Name, v...)
}
func (srv *CoreSrv) LogWarning(v ...interface{}) {
	lib.NewLog().Warning(srv.Srv.Name, v...)
}
func (srv *CoreSrv) LogCritical(v ...interface{}) {
	lib.NewLog().Critical(srv.Srv.Name, v...)
}
func (srv *CoreSrv) LogFatal(v ...interface{}) {
	lib.NewLog().Fatal(srv.Srv.Name, v...)
}
