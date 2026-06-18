program Project1;

uses
  Vcl.Forms,
  Unit1 in 'Unit1.pas' {Form1},
  uMainForm in 'uMainForm.pas' {frmAuth};

{$R *.res}

begin
  Application.Initialize;
  Application.MainFormOnTaskbar := True;
  Application.CreateForm(TfrmAuth, frmAuth);
  Application.Run;
end.
