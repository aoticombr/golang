package main

import (
	"fmt"
	"testing"

	httpaoti "github.com/aoticombr/golang/http"
)

func TestUp_Down(t *testing.T) {
	fmt.Printf("teste")
	link2 := httpaoti.NewHttp()
	link2.SetUrl("https://skylab-api.rocketseat.com.br/journey-nodes/projeto-01")
	link2.SetMetodo(httpaoti.M_GET)
	link2.Request.Header.Authorization = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJhMmNmOTllOC01OTg0LTQ1ZjEtYTU3Yy01YWZmYzc4ZjE5YjYiLCJkYXRhIjp7fSwiaWF0IjoxNzA1MDE0NjA4LCJleHAiOjE3MDUxODc0MDh9.cz1hktCYBGiHxfGePAROPZSTv1IlYHhkMkPCPyDIMEo8jOb1anHEsDrxLj8JzcR5pWspcD_ZVUrTpBdOK2pIYPHhROs_cY4C1LvNEuBJFTqNrRzbLtyF-bRIYEBm77MutsIciIqg5J7epOhNM6NeF1wxyuTcWPgq1dXGOpLnjRruPis1nwLbHkFfqrtdlrkhguJPDynPhK7e-Q6DPz7l99BSvnhkCSKmQh6gpu4J9Fiel7SGCLuEmY7ffyo-rXH8qC5jHLf77_tuF8Zs3kbCQRpuNzOFdA38aXIR4UEiw0a2jYHGZWWycteXdVwrXm9-LyawB1FolUogEro6hGJda9d6saIczKH4oyDxT71iNt1whVkRKay_TfOQ1njA8bm9pG9KzdQyKxA0dwgo7kMQZcypzplwN5kTcHKt6_GDaFnyJuN41KGkopDs-sOnJOW50I-VvlgbhvVVZPJpaPn2BEqlwhEWayc89lOBTptRbDZwRNCRSqfo-NZektNASW_CL4XenEV-yHxsA3CxSY1MJAQaJvUG5KBBBnEzKqupwu_kNVsvi5vSJpZjXEzIM_0BjomiZMi2a3x9oVxL-P1vv7exibGou30H2DHinG8R0P6tdolR_eAoO1CbB00A5H0clfQL5hRAD_BFtSlTlkzG7khAwoOeIUykC37tIfXp-j4"
	link2.Request.Header.AddField("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36")
	link2.Request.Header.AddField("origin", "https://app.rocketseat.com.br")
	link2.Request.Header.AddField("referer", "https://app.rocketseat.com.br/")
	resp2, err := link2.Send()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp2.Body))

	// fmt.Println("Teste")
	// down := m3u8.NewM3u8()
	// byt, err := down.GetVideoByte("https://b-vz-762f4670-e04.tv.pandavideo.com.br/dd0f4f59-f80d-4dfe-864f-e53652f4fea7/playlist.m3u8")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// newUUID := uuid.New()
	// timeString := time.Now().Format("2006-01-02-15-04-05")
	// Name := timeString + newUUID.String() + "playlist.mkv"

	// m3u8.SaveByteToFile(Name, byt)
}
