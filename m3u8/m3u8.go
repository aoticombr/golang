package m3u8

import (
	"fmt"

	httpaoti "github.com/aoticombr/golang/http"
)

type M3u8 struct {
	url string
	arq []byte
}

func (m3 *M3u8) GetVideoByte(url string) ([]byte, error) {
	m3.url = url
	link := httpaoti.NewHttp()
	link.SetUrl(m3.url)
	link.SetMetodo(httpaoti.M_GET)
	resp, err := link.Send()
	if err != nil {
		return nil, fmt.Errorf("Erro ao obter video: %v", err)
	}
	println(string(resp.Body))
	var (
		Resolucoes []string
		Attribs    []*Attribute
	)
	Resolucoes = extractValuesFromBody(string(resp.Body))
	for _, res := range Resolucoes {
		att, err := NewAttribute(res)
		if err != nil {
			return nil, fmt.Errorf("Erro ao obter video: %v", err)
		}
		Attribs = append(Attribs, att)
	}
	var (
		max_atrib *Attribute
	)

	for _, attrib := range Attribs {
		if max_atrib == nil {
			max_atrib = attrib
		}
		if attrib.Height > max_atrib.Height {
			max_atrib = attrib
		}
	}
	//fmt.Println("Escolhida:", max_atrib.Uri)
	urlnew, err := replacePathInURL(m3.url, max_atrib.Uri)
	if err != nil {
		//fmt.Println("Erro ao concatenar a URL:", err)
		return nil, err
	}
	//fmt.Println(urlnew)
	link2 := httpaoti.NewHttp()
	link2.SetUrl(urlnew)
	link2.SetMetodo(httpaoti.M_GET)
	resp2, err := link2.Send()
	if err != nil {
		return nil, err
	}

	Videos := extractValuesFromBody(string(resp2.Body))
	bt, err := DownloadByte(Videos)
	if err != nil {
		return nil, err
	}
	m3.arq = bt
	return bt, nil
}
func (m3 *M3u8) SaveByteToFile(name string) error {
	return SaveByteToFile(name, m3.arq)
}

func NewM3u8() *M3u8 {
	return &M3u8{}
}
