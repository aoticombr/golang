package bot

import "context"

// IControllerBot é o contrato clássico de um bot: Executar é chamado a cada ciclo.
type IControllerBot interface {
	Executar()
}

// IControllerBotCtx é opcional: se o bot implementar ExecutarCtx, ele é
// chamado no lugar de Executar e recebe o contexto da aplicação, que é
// cancelado quando o serviço para. Bots podem implementar só esta interface.
type IControllerBotCtx interface {
	ExecutarCtx(ctx context.Context)
}
