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
