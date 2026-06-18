unit Unit1;

{
  DelphiAuth (XE2 / CEF4Delphi) - captura do authorization code (OAuth2).

  Versao equivalente a do Delphi 12 (TEdgeBrowser), porem usando o navegador
  Chromium embutido via CEF4Delphi (TChromiumWindow), para rodar no Delphi XE2.

  Fluxo:
    1. Monta a URL de autorizacao e navega o Chromium para ela.
    2. O usuario autentica com a Microsoft dentro da janela.
    3. O Azure redireciona para o CallBackUrl com ?code=...; interceptamos no
       evento OnBeforeBrowse, cancelamos a navegacao (Result := True) e extraimos
       o code da URL.
    4. O code e exibido e pode ser copiado. A API (Go) consome o code direto.

  Observacoes de threading:
    OnBeforeBrowse pode ser chamado em uma thread do CEF (quando
    MultiThreadedMessageLoop = True). Por isso a atualizacao da UI
    (HandleRedirect) e marshalada para a thread principal via TThread.Queue.
}

interface

uses
  Winapi.Windows, Winapi.Messages,
  System.SysUtils, System.Variants, System.Classes, System.StrUtils,
  System.IOUtils,
  Vcl.Graphics, Vcl.Controls, Vcl.Forms, Vcl.Dialogs, Vcl.StdCtrls,
  Vcl.ExtCtrls, Vcl.Clipbrd,
  uCEFChromium, uCEFWindowParent, uCEFChromiumWindow, uCEFTypes, uCEFInterfaces,
  uCEFWinControl, uCEFLinkedWinControlBase, uCEFChromiumCore;

type
  TfrmAuth = class(TForm)
    pnlTop: TPanel;
    lblAuthUrl: TLabel;
    lblClientId: TLabel;
    lblCallback: TLabel;
    lblScope: TLabel;
    lblState: TLabel;
    lblCodeCap: TLabel;
    edtAuthUrl: TEdit;
    edtClientId: TEdit;
    edtCallback: TEdit;
    edtScope: TEdit;
    edtState: TEdit;
    btnLogin: TButton;
    edtCode: TEdit;
    btnCopy: TButton;
    btnClear: TButton;
    pnlStatus: TPanel;
    Timer1: TTimer;
    ChromiumWindow1: TChromiumWindow;
    procedure FormCreate(Sender: TObject);
    procedure btnLoginClick(Sender: TObject);
    procedure btnCopyClick(Sender: TObject);
    procedure btnClearClick(Sender: TObject);
    procedure Timer1Timer(Sender: TObject);
    procedure ChromiumWindow1AfterCreated(Sender: TObject);
    procedure FormShow(Sender: TObject);
    procedure ChromiumWindow1Close(Sender: TObject);
    procedure FormCloseQuery(Sender: TObject; var CanClose: Boolean);
  private
    FReady: Boolean;
    FPendingUrl: string;
    FCode: string;
    FCallbackBase: string; // lido em OnBeforeBrowse (thread do CEF) - usar campo, nao a UI
    FRedirectUri: string;  // passado de OnBeforeBrowse -> HandleRedirect via TThread.Queue
    procedure WMMove(var aMessage: TWMMove); message WM_MOVE;
    procedure WMMoving(var aMessage: TMessage); message WM_MOVING;
    procedure WMEnterMenuLoop(var aMessage: TMessage); message WM_ENTERMENULOOP;
    procedure WMExitMenuLoop(var aMessage: TMessage); message WM_EXITMENULOOP;
    procedure DoHandleRedirect; // roda na thread principal
  protected
    // Variaveis para destruir o form com seguranca (padrao CEF4Delphi)
    FCanClose: Boolean; // True no TChromium.OnClose
    FClosing: Boolean;  // True no CloseQuery
    procedure SetStatus(const Msg: string);
    procedure HandleRedirect(const Uri: string);
    procedure Chromium_OnBeforeBrowse(Sender: TObject; const browser: ICefBrowser;
      const frame: ICefFrame; const request: ICefRequest;
      user_gesture, isRedirect: Boolean; out Result: Boolean);
    procedure Chromium_OnBeforePopup(Sender: TObject; const browser: ICefBrowser;
      const frame: ICefFrame; popup_id: Integer; const targetUrl, targetFrameName: ustring;
      targetDisposition: TCefWindowOpenDisposition; userGesture: Boolean;
      const popupFeatures: TCefPopupFeatures; var windowInfo: TCefWindowInfo;
      var client: ICefClient; var settings: TCefBrowserSettings;
      var extra_info: ICefDictionaryValue; var noJavascriptAccess: Boolean;
      var Result: Boolean);
  public
    function BuildAuthorizationUrl: string;
  end;

var
  frmAuth: TfrmAuth;

implementation

uses
  uCEFApplication;

{$R *.dfm}

// XE2 nao tem System.NetEncoding/TNetEncoding. Helpers proprios (UTF-8):

// Percent-encode mantendo apenas os caracteres "unreserved" (RFC 3986).
function UrlEncode(const S: string): string;
var
  Bytes: TBytes;
  i: Integer;
  b: Byte;
begin
  Result := '';
  Bytes := TEncoding.UTF8.GetBytes(S);
  for i := 0 to High(Bytes) do
  begin
    b := Bytes[i];
    if ((b >= Ord('A')) and (b <= Ord('Z'))) or
       ((b >= Ord('a')) and (b <= Ord('z'))) or
       ((b >= Ord('0')) and (b <= Ord('9'))) or
       (b = Ord('-')) or (b = Ord('_')) or (b = Ord('.')) or (b = Ord('~')) then
      Result := Result + Char(b)
    else
      Result := Result + '%' + IntToHex(b, 2);
  end;
end;

// Decodifica %XX e '+' (espaco), interpretando os bytes como UTF-8.
function UrlDecode(const S: string): string;
var
  Bytes: TBytes;
  i, cnt: Integer;
  c: Char;
begin
  SetLength(Bytes, Length(S));
  cnt := 0;
  i := 1;
  while i <= Length(S) do
  begin
    c := S[i];
    if c = '+' then
    begin
      Bytes[cnt] := Ord(' ');
      Inc(cnt);
      Inc(i);
    end
    else if (c = '%') and (i + 2 <= Length(S)) then
    begin
      Bytes[cnt] := StrToInt('$' + S[i + 1] + S[i + 2]);
      Inc(cnt);
      Inc(i, 3);
    end
    else
    begin
      Bytes[cnt] := Byte(Ord(c) and $FF);
      Inc(cnt);
      Inc(i);
    end;
  end;
  SetLength(Bytes, cnt);
  Result := TEncoding.UTF8.GetString(Bytes);
end;

{ TfrmAuth }

procedure TfrmAuth.FormCreate(Sender: TObject);
begin
  // Valores de exemplo (os mesmos do teste Go). Ajuste conforme necessario.
  edtAuthUrl.Text := 'https://login.microsoftonline.com/cc2dba5e-b8a0-415c-bfe0-ece6b2fff759/oauth2/v2.0/authorize';
  edtClientId.Text := '85dba09e-5061-45ab-b254-d387b43a66ab';
  edtCallback.Text := 'https://mockserver.aoti.com.br/cb';
  edtScope.Text := 'mail.read';
  edtState.Text := 'teste-state-123';

  FReady := False;
  FCanClose := False;
  FClosing := False;

  ChromiumWindow1.ChromiumBrowser.DefaultURL := 'about:blank';
  ChromiumWindow1.ChromiumBrowser.RuntimeStyle := CEF_RUNTIME_STYLE_ALLOY;

  SetStatus('Iniciando navegador...');
end;

procedure TfrmAuth.FormShow(Sender: TObject);
begin
  // Intercepta a navegacao para capturar o redirect com o code.
  ChromiumWindow1.ChromiumBrowser.OnBeforeBrowse := Chromium_OnBeforeBrowse;
  // Bloqueia popups/novas abas (mantem tudo na mesma janela).
  ChromiumWindow1.ChromiumBrowser.OnBeforePopup := Chromium_OnBeforePopup;

  // *MUST* call CreateBrowser. Se o contexto global ainda nao estiver pronto,
  // tenta de novo via Timer.
  if not ChromiumWindow1.CreateBrowser then
    Timer1.Enabled := True;
end;

procedure TfrmAuth.Timer1Timer(Sender: TObject);
begin
  Timer1.Enabled := False;
  if not ChromiumWindow1.CreateBrowser and not ChromiumWindow1.Initialized then
    Timer1.Enabled := True;
end;

procedure TfrmAuth.ChromiumWindow1AfterCreated(Sender: TObject);
begin
  FReady := True;
  SetStatus('Preencha os dados e clique em "Abrir login".');

  // Se o usuario clicou em "Abrir login" antes do navegador ficar pronto.
  if FPendingUrl <> '' then
  begin
    ChromiumWindow1.LoadURL(FPendingUrl);
    FPendingUrl := '';
  end;
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
    'client_id=' + UrlEncode(Trim(edtClientId.Text)) +
    '&response_type=code' +
    '&redirect_uri=' + UrlEncode(Trim(edtCallback.Text)) +
    '&response_mode=query' +
    '&scope=' + UrlEncode(Trim(edtScope.Text));

  if Trim(edtState.Text) <> '' then
    Params := Params + '&state=' + UrlEncode(Trim(edtState.Text));

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
  FCallbackBase := Trim(edtCallback.Text);

  if not FReady then
  begin
    FPendingUrl := BuildAuthorizationUrl;
    SetStatus('Aguardando o navegador iniciar...');
    Exit;
  end;

  SetStatus('Abrindo o login da Microsoft...');
  ChromiumWindow1.LoadURL(BuildAuthorizationUrl);
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
begin
  FCode := '';
  edtCode.Text := '';
  btnCopy.Enabled := False;

  // Limpa o cache (Network.clearBrowserCache) e remove todos os cookies.
  ChromiumWindow1.ChromiumBrowser.ClearCache;
  ChromiumWindow1.ChromiumBrowser.DeleteCookies('', '', True);

  SetStatus('Cache e cookies limpos. Clique em "Abrir login" novamente.');
end;

procedure TfrmAuth.Chromium_OnBeforeBrowse(Sender: TObject;
  const browser: ICefBrowser; const frame: ICefFrame; const request: ICefRequest;
  user_gesture, isRedirect: Boolean; out Result: Boolean);
var
  Uri: string;
begin
  Result := False;
  Uri := request.Url;

  // So intercepta a navegacao para o redirect_uri configurado em btnLoginClick.
  if (FCallbackBase <> '') and StartsText(FCallbackBase, Uri) then
  begin
    // Cancela a navegacao para o callback (a pagina pode nem existir).
    Result := True;
    // HandleRedirect mexe na UI: marshalar para a thread principal.
    FRedirectUri := Uri;
    TThread.Queue(nil, DoHandleRedirect);
  end;
end;

procedure TfrmAuth.DoHandleRedirect;
begin
  HandleRedirect(FRedirectUri);
end;

procedure TfrmAuth.HandleRedirect(const Uri: string);
var
  QueryPart, Pair, Key, Val: string;
  Eq, P, AmpStart, idx: Integer;
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

  // Split manual por '&' (XE2 nao tem TStringHelper.Split). O '&' sentinela no
  // fim garante que o ultimo par tambem seja processado.
  QueryPart := QueryPart + '&';
  AmpStart := 1;
  for idx := 1 to Length(QueryPart) do
  begin
    if QueryPart[idx] <> '&' then
      Continue;

    Pair := Copy(QueryPart, AmpStart, idx - AmpStart);
    AmpStart := idx + 1;
    if Pair = '' then
      Continue;

    Eq := Pos('=', Pair); // 1-based; 0 se nao houver
    if Eq = 0 then
    begin
      Key := UrlDecode(Pair);
      Val := '';
    end
    else
    begin
      Key := UrlDecode(Copy(Pair, 1, Eq - 1));
      Val := UrlDecode(Copy(Pair, Eq + 1, MaxInt));
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

procedure TfrmAuth.Chromium_OnBeforePopup(Sender: TObject;
  const browser: ICefBrowser; const frame: ICefFrame; popup_id: Integer;
  const targetUrl, targetFrameName: ustring;
  targetDisposition: TCefWindowOpenDisposition; userGesture: Boolean;
  const popupFeatures: TCefPopupFeatures; var windowInfo: TCefWindowInfo;
  var client: ICefClient; var settings: TCefBrowserSettings;
  var extra_info: ICefDictionaryValue; var noJavascriptAccess, Result: Boolean);
begin
  // Bloqueia novas abas/janelas; mantem a navegacao na mesma janela.
  Result := (targetDisposition in [CEF_WOD_NEW_FOREGROUND_TAB,
    CEF_WOD_NEW_BACKGROUND_TAB, CEF_WOD_NEW_POPUP, CEF_WOD_NEW_WINDOW]);
end;

procedure TfrmAuth.ChromiumWindow1Close(Sender: TObject);
begin
  FCanClose := True;
  PostMessage(Handle, WM_CLOSE, 0, 0);
end;

procedure TfrmAuth.FormCloseQuery(Sender: TObject; var CanClose: Boolean);
begin
  CanClose := FCanClose;

  if not FClosing then
  begin
    FClosing := True;
    Visible := False;
    ChromiumWindow1.CloseBrowser(True);
  end;
end;

procedure TfrmAuth.WMMove(var aMessage: TWMMove);
begin
  inherited;
  if ChromiumWindow1 <> nil then
    ChromiumWindow1.NotifyMoveOrResizeStarted;
end;

procedure TfrmAuth.WMMoving(var aMessage: TMessage);
begin
  inherited;
  if ChromiumWindow1 <> nil then
    ChromiumWindow1.NotifyMoveOrResizeStarted;
end;

procedure TfrmAuth.WMEnterMenuLoop(var aMessage: TMessage);
begin
  inherited;
  if (aMessage.wParam = 0) and (GlobalCEFApp <> nil) then
    GlobalCEFApp.OsmodalLoop := True;
end;

procedure TfrmAuth.WMExitMenuLoop(var aMessage: TMessage);
begin
  inherited;
  if (aMessage.wParam = 0) and (GlobalCEFApp <> nil) then
    GlobalCEFApp.OsmodalLoop := False;
end;

end.
