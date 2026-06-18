unit uMainForm;

{
  DelphiAuth - captura do authorization code (OAuth2 Authorization Code flow).

  Abre o endpoint de autorizacao da Microsoft (AuthUrl) numa janela embutida
  (TEdgeBrowser / WebView2), no estilo do Postman. Apos o login, o servidor
  redireciona para o CallBackUrl com ?code=...; interceptamos essa navegacao no
  evento NavigationStarting, cancelamos antes de carregar a pagina de callback
  (que pode nem existir) e extraimos o code da query string.

  O objetivo deste app e apenas CAPTURAR e EXIBIR o code - a troca por token e
  feita pela API (componente http em Go), que consome o code diretamente.

  Requisitos:
    - Delphi 11/12 (TEdgeBrowser disponivel)
    - Microsoft Edge WebView2 Runtime instalado (padrao no Windows 11)
    - WebView2Loader.dll ao lado do executavel (vem com o Delphi)
    - O CallBackUrl precisa estar registrado como redirect URI no app do Azure AD
}

interface

uses
  Winapi.Windows, Winapi.Messages, Winapi.WebView2,
  System.SysUtils, System.Variants, System.Classes, System.StrUtils,
  System.IOUtils, System.NetEncoding,
  Vcl.Graphics, Vcl.Controls, Vcl.Forms, Vcl.Dialogs, Vcl.StdCtrls,
  Vcl.ExtCtrls, Vcl.Clipbrd, Vcl.Edge, Winapi.ActiveX;

type
  TfrmAuth = class(TForm)
    pnlTop: TPanel;
    lblAuthUrl: TLabel;
    edtAuthUrl: TEdit;
    lblClientId: TLabel;
    edtClientId: TEdit;
    lblCallback: TLabel;
    edtCallback: TEdit;
    lblScope: TLabel;
    edtScope: TEdit;
    lblState: TLabel;
    edtState: TEdit;
    btnLogin: TButton;
    lblCodeCap: TLabel;
    edtCode: TEdit;
    btnCopy: TButton;
    btnClear: TButton;
    pnlStatus: TPanel;
    EdgeBrowser1: TEdgeBrowser;
    procedure FormCreate(Sender: TObject);
    procedure btnLoginClick(Sender: TObject);
    procedure btnCopyClick(Sender: TObject);
    procedure btnClearClick(Sender: TObject);
    procedure EdgeBrowser1NavigationStarting(Sender: TCustomEdgeBrowser;
      Args: TNavigationStartingEventArgs);
    procedure EdgeBrowser1CreateWebViewCompleted(Sender: TCustomEdgeBrowser;
      AResult: HRESULT);
  private
    FReady: Boolean;
    FPendingUrl: string;
    FCode: string;
    function BuildAuthorizationUrl: string;
    procedure HandleRedirect(const Uri: string);
    procedure SetStatus(const Msg: string);
  public
  end;

var
  frmAuth: TfrmAuth;

implementation

{$R *.dfm}

procedure TfrmAuth.FormCreate(Sender: TObject);
begin
  // Pasta de dados do WebView2 em local gravavel (evita erro ao lado do .exe).
  EdgeBrowser1.UserDataFolder := TPath.Combine(TPath.GetTempPath, 'DelphiAuthWV2');

  // Valores de exemplo (os mesmos do teste Go). Ajuste conforme necessario.
  edtAuthUrl.Text := 'https://login.microsoftonline.com/cc2dba5e-b8a0-415c-bfe0-ece6b2fff759/oauth2/v2.0/authorize';
  edtClientId.Text := '85dba09e-5061-45ab-b254-d387b43a66ab';
  edtCallback.Text := 'https://mockserver.aoti.com.br/cb';
  edtScope.Text := 'mail.read';
  edtState.Text := 'teste-state-123';

  // O TEdgeBrowser NAO inicializa sozinho ao chamar Navigate: e preciso criar o
  // WebView2 explicitamente. Sem isto a navegacao fica adiada e a pagina nunca carrega.
  SetStatus('Iniciando WebView2...');
  EdgeBrowser1.CreateWebView;
end;

procedure TfrmAuth.EdgeBrowser1CreateWebViewCompleted(Sender: TCustomEdgeBrowser;
  AResult: HRESULT);
begin
  if Succeeded(AResult) then
  begin
    FReady := True;
    SetStatus('Preencha os dados e clique em "Abrir login".');
    // Se o usuario ja clicou em "Abrir login" antes do WebView2 ficar pronto.
    if FPendingUrl <> '' then
    begin
      EdgeBrowser1.Navigate(FPendingUrl);
      FPendingUrl := '';
    end;
  end
  else
    SetStatus('Falha ao iniciar o WebView2 (HRESULT=0x' + IntToHex(AResult, 8) +
      '). Verifique se o Edge WebView2 Runtime esta instalado e se o WebView2Loader.dll ' +
      'esta ao lado do .exe.');
end;

procedure TfrmAuth.SetStatus(const Msg: string);
begin
  pnlStatus.Caption := '  ' + Msg;
end;

function TfrmAuth.BuildAuthorizationUrl: string;
var
  Params: string;
begin
  // response_mode=query garante o code como parametro de query no redirect GET.
  Params :=
    'client_id=' + TNetEncoding.URL.Encode(Trim(edtClientId.Text)) +
    '&response_type=code' +
    '&redirect_uri=' + TNetEncoding.URL.Encode(Trim(edtCallback.Text)) +
    '&response_mode=query' +
    '&scope=' + TNetEncoding.URL.Encode(Trim(edtScope.Text));

  if Trim(edtState.Text) <> '' then
    Params := Params + '&state=' + TNetEncoding.URL.Encode(Trim(edtState.Text));

  Result := Trim(edtAuthUrl.Text) + '?' + Params;
end;

procedure TfrmAuth.btnLoginClick(Sender: TObject);
begin
  if Trim(edtAuthUrl.Text) = '' then
  begin
    SetStatus('Informe a AuthUrl.');
    Exit;
  end;
  if Trim(edtCallback.Text) = '' then
  begin
    SetStatus('Informe o CallBackUrl (redirect_uri).');
    Exit;
  end;

  FCode := '';
  edtCode.Text := '';
  btnCopy.Enabled := False;

  if not FReady then
  begin
    // WebView2 ainda inicializando: guarda a URL e navega ao concluir.
    FPendingUrl := BuildAuthorizationUrl;
    SetStatus('Aguardando o WebView2 iniciar...');
    Exit;
  end;

  SetStatus('Abrindo o login da Microsoft...');
  EdgeBrowser1.Navigate(BuildAuthorizationUrl);
end;

procedure TfrmAuth.EdgeBrowser1NavigationStarting(Sender: TCustomEdgeBrowser;
  Args: TNavigationStartingEventArgs);
var
  Uri: string;
  P: PWideChar;
begin
  // O wrapper TNavigationStartingEventArgs nao expoe a URI diretamente; lemos da
  // interface COM. Get_uri aloca a string, que deve ser liberada com CoTaskMemFree.
  P := nil;
  if Failed(Args.ArgsInterface.Get_uri(P)) or (P = nil) then
    Exit;
  Uri := string(P);
  CoTaskMemFree(P);

  // So intercepta quando a navegacao for para o redirect_uri configurado.
  if StartsText(Trim(edtCallback.Text), Uri) then
  begin
    // Cancela a navegacao para o callback (a pagina pode nem existir) e
    // extrai o code da propria URL pretendida.
    Args.ArgsInterface.Set_Cancel(1);
    HandleRedirect(Uri);
  end;
end;

procedure TfrmAuth.HandleRedirect(const Uri: string);
var
  QueryPart, Pair, Key, Val: string;
  Pairs: TArray<string>;
  Eq, P: Integer;
  Code, State, ErrCode, ErrDesc: string;
begin
  P := Pos('?', Uri);
  if P = 0 then
  begin
    SetStatus('Redirect recebido sem query string: ' + Uri);
    Exit;
  end;

  QueryPart := Copy(Uri, P + 1, MaxInt);

  // Remove fragmento (#...) se houver.
  P := Pos('#', QueryPart);
  if P > 0 then
    QueryPart := Copy(QueryPart, 1, P - 1);

  Pairs := QueryPart.Split(['&']);
  for Pair in Pairs do
  begin
    if Pair = '' then
      Continue;

    Eq := Pair.IndexOf('='); // 0-based; -1 se nao houver
    if Eq < 0 then
    begin
      Key := TNetEncoding.URL.Decode(Pair);
      Val := '';
    end
    else
    begin
      Key := TNetEncoding.URL.Decode(Pair.Substring(0, Eq));
      Val := TNetEncoding.URL.Decode(Pair.Substring(Eq + 1));
    end;

    if SameText(Key, 'code') then
      Code := Val
    else if SameText(Key, 'state') then
      State := Val
    else if SameText(Key, 'error') then
      ErrCode := Val
    else if SameText(Key, 'error_description') then
      ErrDesc := Val;
  end;

  if ErrCode <> '' then
  begin
    SetStatus('Erro retornado pelo IdP: ' + ErrCode + ' - ' + ErrDesc);
    Exit;
  end;

  if Code = '' then
  begin
    SetStatus('Callback sem o parametro "code".');
    Exit;
  end;

  // Valida o state (protecao CSRF) quando foi enviado.
  if (Trim(edtState.Text) <> '') and (State <> Trim(edtState.Text)) then
  begin
    SetStatus('ATENCAO: state divergente (possivel CSRF). Recebido: "' + State + '".');
    Exit;
  end;

  FCode := Code;
  edtCode.Text := Code;
  btnCopy.Enabled := True;
  SetStatus('Code capturado com sucesso! Copie e entregue para a API.');
end;

procedure TfrmAuth.btnCopyClick(Sender: TObject);
begin
  if FCode = '' then
  begin
    SetStatus('Nenhum code capturado ainda.');
    Exit;
  end;
  Clipboard.AsText := FCode;
  SetStatus('Code copiado para a area de transferencia.');
end;

procedure TfrmAuth.btnClearClick(Sender: TObject);
var
  Folder: string;
  i: Integer;
  Deleted: Boolean;
begin
  Folder := EdgeBrowser1.UserDataFolder;

  FCode := '';
  edtCode.Text := '';
  btnCopy.Enabled := False;
  FReady := False;
  FPendingUrl := '';
  SetStatus('Limpando cache e sessao do navegador...');

  // 1) Fecha o WebView2 para liberar os arquivos da UserDataFolder.
  EdgeBrowser1.CloseWebView;

  // 2) Apaga a pasta de dados inteira (cache de disco, cookies, localStorage,
  //    IndexedDB, etc.). O processo msedgewebview2.exe pode demorar a liberar os
  //    arquivos, entao tentamos algumas vezes antes de desistir.
  Deleted := True;
  if (Folder <> '') and TDirectory.Exists(Folder) then
  begin
    Deleted := False;
    for i := 1 to 20 do
    begin
      try
        TDirectory.Delete(Folder, True);
        Deleted := True;
        Break;
      except
        Sleep(150);
        Application.ProcessMessages;
      end;
    end;
  end;

  // 3) Recria o WebView2 do zero (OnCreateWebViewCompleted reabilita o uso).
  EdgeBrowser1.CreateWebView;

  if Deleted then
    SetStatus('Cache e sessao limpos. Reiniciando WebView2...')
  else
    SetStatus('Sessao reiniciada, mas alguns arquivos de cache estavam em uso. ' +
      'Feche o app e apague a pasta manualmente se precisar de limpeza total: ' + Folder);
end;

end.
