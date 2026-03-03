## Exemplo 1

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/aoticombr/golang/logger"
)

func main() {

	executablePath, err := os.Executable()
	if err != nil {
		// Lidar com o erro, se necessário
	}
	appRoot := filepath.Dir(executablePath)
	logDir := appRoot //
	fmt.Println(logDir)
	logger, _ := log.NewLogger("INFO", os.Stdout, "[DEVRAIZ]", logDir)
	logger.Info("ler o arquivo")
	logger.Info("Download", "Download", "Download", "Download", "Download", "Download", "Download")
	logger.Info("Descompactar o arquivo")
	logger.Info("ler o arquivo")
	logger.Warning("erro ao ler o arquivo")
	logger.Fatal("erro ao ler o arquivo")

}
```
os logs só serão mostrados por nivel de hierarquia

// **DEBUG** = 0
// **INFO** = 1
// **WARNING** = 2
// **ERROR** = 3
// **CRITICAL** = 4

exemplo se colocar como **ERROR** ele ira mostrar somente **ERROR** ou **CRITICAL** já que o nivel de Error é **3** sempre sera um **>=** em comparação neste caso

Para você que esta programando o ideal é o Debug

use o **logger.Debug()** para identificar onde voce entrou nas rotinas

use o **logger.Info()** para identificar ações

use o **logger.Erro()** para mostrar erros

use o logger.Fatal() para mostrar um erro e **parar a aplicação**, pois o **Fatal(**) como o proprio nome já diz é um **erro gravissimo** para sua aplicação continuar, exemplo se você esta lendo o arquivo **.env** da sua aplicação e não conseguiu isso seria uma falha gravissima ja que sua aplicação depende dele **para funcionar**, se você esta em uma aplicação que possui um ciclo de rodagem use o **logger.Fatal()** somente nesses caso graves já que ele para sua aplicação


