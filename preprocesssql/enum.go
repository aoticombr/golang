package preprocesssql

const nullChar byte = '\x00'
const aspasimples byte = 39
const emptstr string = ""
const (
	piCreateParams TPreprocessorInstr = iota
	piCreateMacros
	piExpandParams
	piExpandMacros
	piExpandEscapes
	piParseSQL
	piTransformQuestions
	piTransformEOLs
)
const (
	ncDefault TNameQuoteLevel = 1
	ncSecond  TNameQuoteLevel = 2
	ncThird   TNameQuoteLevel = 3
)
const (
	nsLeft TNameQuoteSide = iota
	nsRight
)
const (
	elDefault TTextEndOfLine = iota
	elWindows
	elUnix
	elMac
)
const (
	mkUnknown TDBMSKind = iota
	mkOracle
	mkMSSQL
	mkMSAccess
	mkMySQL
	mkDB2
	mkSQLAnywhere
	mkAdvantage
	mkInterbase
	mkFirebird
	mkSQLite
	mkPostgreSQL
	mkNexusDB
	mkDataSnap
	mkInformix
	mkTeradata
	mkMongDB
	mkOther
)
const (
	eskText TEscapeKind = iota
	eskString
	eskFloat
	eskDate
	eskTime
	eskDateTime
	eskIdentifier
	eskBoolean
	eskFunction
	eskIF
	eskFI
	eskElse
	eskIIF
	eskEscape
	eskInto
)
const (
	skUnknown TCommandKind = iota
	skSelect
	skSelectForLock
	skSelectForUnLock
	skDelete
	skInsert
	skMerge
	skUpdate
	skCreate
	skAlter
	skDrop
	skStoredProc
	skStoredProcWithCrs
	skStoredProcNoCrs
	skExecute
	skStartTransaction
	skCommit
	skRollback
	skSet
	skSetSchema
	skOther
	skNotResolved
)

const (
	pbByName TParamBindMode = iota
	pbByNumber
)

const (
	ptUnknown TParamType = iota
	ptInput
	ptOutput
	ptInputOutput
	ptResult
)

const (
	prQMark TParamMark = iota
	prName
	prNumber
	prDollar
	prQNumber
)

const (
	efASCII TEscapeFunction = iota
	efLTRIM
	efREPLACE
	efRTRIM
	efABS
	efCOS
	efEXP
	efFLOOR
	efMOD
	efPOWER
	efROUND
	efSIGN
	efSIN
	efSQRT
	efTAN
	efDECODE
	efBIT_LENGTH
	efCHAR
	efCHAR_LENGTH
	efCONCAT
	efINSERT
	efLCASE
	efLEFT
	efLENGTH
	efLOCATE
	efOCTET_LENGTH
	efPOSITION
	efREPEAT
	efRIGHT
	efSPACE
	efSUBSTRING
	efUCASE
	efACOS
	efASIN
	efATAN
	efATAN2
	efCOT
	efCEILING
	efDEGREES
	efLOG
	efLOG10
	efPI
	efRADIANS
	efRANDOM
	efTRUNCATE
	efCURDATE
	efCURTIME
	efNOW
	efDAYNAME
	efDAYOFMONTH
	efDAYOFWEEK
	efDAYOFYEAR
	efEXTRACT
	efHOUR
	efMINUTE
	efMONTH
	efMONTHNAME
	efQUARTER
	efSECOND
	efTIMESTAMPADD
	efTIMESTAMPDIFF
	efWEEK
	efYEAR
	efCATALOG
	efSCHEMA
	efIFNULL
	efIF
	efCONVERT
	efLIMIT
	efNONE
)

const (
	eoQuote TEncodeOption = iota
	eoNormalize
	eoBeautify
)

const (
	npCatalog TNamePart = iota
	npSchema
	npDBLink
	npBaseObject
	npObject
)
