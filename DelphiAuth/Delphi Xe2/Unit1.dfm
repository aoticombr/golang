object frmAuth: TfrmAuth
  Left = 0
  Top = 0
  Caption = 'frmAuth'
  ClientHeight = 701
  ClientWidth = 950
  Color = clBtnFace
  Font.Charset = DEFAULT_CHARSET
  Font.Color = clWindowText
  Font.Height = -11
  Font.Name = 'Tahoma'
  Font.Style = []
  OldCreateOrder = False
  OnCloseQuery = FormCloseQuery
  OnCreate = FormCreate
  OnShow = FormShow
  PixelsPerInch = 96
  TextHeight = 13
  object pnlTop: TPanel
    Left = 0
    Top = 0
    Width = 950
    Height = 225
    Align = alTop
    BevelOuter = bvNone
    Padding.Left = 8
    Padding.Top = 8
    Padding.Right = 8
    TabOrder = 0
    object lblAuthUrl: TLabel
      Left = 8
      Top = 8
      Width = 96
      Height = 13
      Caption = 'AuthUrl (/authorize)'
    end
    object lblClientId: TLabel
      Left = 8
      Top = 56
      Width = 37
      Height = 13
      Caption = 'ClientId'
    end
    object lblCallback: TLabel
      Left = 416
      Top = 56
      Width = 118
      Height = 13
      Caption = 'CallBackUrl (redirect_uri)'
    end
    object lblScope: TLabel
      Left = 8
      Top = 104
      Width = 29
      Height = 13
      Caption = 'Scope'
    end
    object lblState: TLabel
      Left = 416
      Top = 104
      Width = 58
      Height = 13
      Caption = 'State (opc.)'
    end
    object lblCodeCap: TLabel
      Left = 8
      Top = 160
      Width = 81
      Height = 13
      Caption = 'Code capturado:'
    end
    object edtAuthUrl: TEdit
      Left = 9
      Top = 27
      Width = 796
      Height = 21
      TabOrder = 0
    end
    object edtClientId: TEdit
      Left = 8
      Top = 75
      Width = 392
      Height = 21
      TabOrder = 1
    end
    object edtCallback: TEdit
      Left = 416
      Top = 75
      Width = 388
      Height = 21
      TabOrder = 2
    end
    object edtScope: TEdit
      Left = 8
      Top = 123
      Width = 392
      Height = 21
      TabOrder = 3
    end
    object edtState: TEdit
      Left = 416
      Top = 123
      Width = 252
      Height = 21
      TabOrder = 4
    end
    object btnLogin: TButton
      Left = 684
      Top = 121
      Width = 120
      Height = 27
      Caption = 'Abrir login'
      Default = True
      TabOrder = 5
      OnClick = btnLoginClick
    end
    object edtCode: TEdit
      Left = 8
      Top = 179
      Width = 580
      Height = 21
      ReadOnly = True
      TabOrder = 6
    end
    object btnCopy: TButton
      Left = 596
      Top = 178
      Width = 86
      Height = 25
      Caption = 'Copiar'
      Enabled = False
      TabOrder = 7
      OnClick = btnCopyClick
    end
    object btnClear: TButton
      Left = 688
      Top = 178
      Width = 116
      Height = 25
      Caption = 'Limpar cache'
      TabOrder = 8
      OnClick = btnClearClick
    end
  end
  object pnlStatus: TPanel
    Left = 0
    Top = 681
    Width = 950
    Height = 20
    Align = alBottom
    Alignment = taLeftJustify
    BevelOuter = bvNone
    BorderWidth = 1
    Color = clBtnShadow
    ParentBackground = False
    TabOrder = 1
  end
  object ChromiumWindow1: TChromiumWindow
    Left = 0
    Top = 225
    Width = 950
    Height = 456
    Align = alClient
    TabStop = True
    TabOrder = 2
    OnClose = ChromiumWindow1Close
    OnAfterCreated = ChromiumWindow1AfterCreated
    ExplicitLeft = 2
    ExplicitTop = 224
  end
  object Timer1: TTimer
    Enabled = False
    Interval = 300
    OnTimer = Timer1Timer
    Left = 231
    Top = 220
  end
end
