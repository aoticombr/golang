package main

import (
	"fmt"
	"time"

	"github.com/aoticombr/golang/http"
)

type ReadSocket struct {
}

func (rs *ReadSocket) Read(messageType int, body []byte, err error) {
	fmt.Println("-------Read-------")
	fmt.Println(time.Now())
	fmt.Println("ReadSocket.read")
	fmt.Println("messageType:", messageType)
	fmt.Println("body:", string(body))
	fmt.Println("err:", err)
}
func (rs *ReadSocket) Error(msg string) {
	fmt.Println("-------Error-------")
	fmt.Println(time.Now())
	fmt.Println("msg:", msg)
}
func (rs *ReadSocket) Msg(msg string) {
	fmt.Println("-------Msg-------")
	fmt.Println(time.Now())
	fmt.Println("msg:", msg)
}
func main() {
	var rs *ReadSocket
	rs = &ReadSocket{}
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.Request.Header.AddField("x-vectury-dealer", "07600973")
	cp.Authorization = `eyJraWQiOiJEVXBTQWFUREtQRHphS19mV0NJcldxUUVOVTQ4bzM2ZXR5ZlV5eG9PaFY0IiwiYWxnIjoiUlMyNTYifQ.eyJ2ZXIiOjEsImp0aSI6IkFULklIUGRpQ0JtbkhxX0pjcDNSeG56bHY1UkozR0pjbDFudXJjVmlwckxMcVUiLCJpc3MiOiJodHRwczovL3Nzby11YXQucmVuYXVsdC5jb20vb2F1dGgyL2F1c3R3b2VzaGJDa1BKeXcxNDE2IiwiYXVkIjoiaHR0cHM6Ly9hcGlzLnJlbmF1bHQuY29tIiwiaWF0IjoxNjk4MTc0NzY1LCJleHAiOjE2OTgxNzgzNjUsImNpZCI6Imlybi03MDcyNV91YXRfcGtqd3Rfb3VxZHN2Z2NqczNxIiwic2NwIjpbImFwaXMuZGVmYXVsdCIsImRmdC12MS5kZWFsZXJzLWRvd25sb2FkIiwiZGZ0LXYxLmRlYWxlcnMtdXBsb2FkIl0sInN1YiI6Imlybi03MDcyNV91YXRfcGtqd3Rfb3VxZHN2Z2NqczNxIiwiaXJuIjoiSVJOLTcwNzI1IiwiY2VydC11aWQiOiJpcm4tNzA3MjVfdWF0X3Brand0In0.J66bg1F8UpHJOhfus28aNs3_YmRvhaU0Y7KVpEPVbsjtr8mkuNS3ulMVvZp9ba3yw94AMMT_aolWM2qgzRWvas9ugYNebBJez_B5Nj8SPdQF34HJxC_FEINCfd6IqRDF-NhlaMpbVgvycc8DwzM2Jq_5YZ2P1VQryOtze07iXSbGm7HGwRGjX_e7_0nrsQ0P5_AxHmPHWSVZRMJukl4dA6LjwViyD5U8ZCkwSkcfPXJfv37cg693iTldW098ZNJdAM0yBr1dj0Ig_vR_LCFc3YPWLi_KfwpfHPgB5xFClUOUuMXW18KX7-xi1gkL2gBIibb3WGcg8pDxP6Gl3Syhjg`
	cp.AuthorizationType = http.AT_Bearer
	cp.SetUrl("ws://localhost:3030/")

	cp.Metodo = http.M_POST
	cp.EncType = http.ET_WEB_SERVICE

	cp.OnSend = rs
	fmt.Println("###############ini###################")
	err := cp.Conectar()
	fmt.Println("###############fim###################")
	if err != nil {
		panic(err)
	}
	select {}
}
