program DelphiAuth;

uses
  Vcl.Forms,
  uMainForm in 'uMainForm.pas' {frmAuth};

{$R *.res}

begin
  Application.Initialize;
  Application.MainFormOnTaskbar := True;
  Application.Title := 'DelphiAuth - OAuth2 Authorization Code';
  Application.CreateForm(TfrmAuth, frmAuth);
  Application.Run;
end.
