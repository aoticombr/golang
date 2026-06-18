program Project1;



uses
   {$IFDEF DELPHI16_UP}
  Vcl.Forms, WinApi.Windows,
  {$ELSE}
  Forms, Windows,
  {$ENDIF}
   uCEFApplication,
  Unit1 in 'Unit1.pas' {frmAuth};

{$R *.res}

const
  IMAGE_FILE_LARGE_ADDRESS_AWARE = $0020;



begin



  GlobalCEFApp := TCefApplication.Create;

  if GlobalCEFApp.StartMainProcess then
    begin
      Application.Initialize;
      Application.CreateForm(TfrmAuth, frmAuth);
      Application.Run;
    end;

  GlobalCEFApp.Free;
  GlobalCEFApp := nil;
end.
