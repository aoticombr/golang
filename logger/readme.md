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
	logger.Fatal("erro ao ler o arquivo")

}
```
os logs só serão mostrados por nivel de hierarquia

// DEBUG = 0
// INFO = 1
// WARNING = 2
// ERROR = 3
// CRITICAL = 4

exemplo se colocar como "ERROR" ele ira mostrar somente ERROR OU Critital já que o nivel de Error é 3 sempre sera um >= em comparação neste caso
