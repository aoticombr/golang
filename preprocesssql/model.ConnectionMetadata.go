package preprocesssql

import (
	"fmt"
)

type TConnectionMetadata struct {
}

func (cm *TConnectionMetadata) GetNameQuoteChar(AQuote TNameQuoteLevel, ASide TNameQuoteSide) byte {

	if AQuote == ncDefault {
		return '"'
	} else {
		return 0
	}
}
func (cm *TConnectionMetadata) AddEscapeSequenceArgs(ASeq *TEscapeData) string {
	var i int

	Result := ""
	for i = 0; i < len(ASeq.Args); i++ {
		if i > 0 {
			Result += ", "
		}
		Result += ASeq.Args[i]
	}
	return Result
}
func (cm *TConnectionMetadata) InternalEscapeBoolean(AStr string) string {
	return AStr
}
func (cm *TConnectionMetadata) InternalEscapeDate(AStr string) string {
	return AnsiQuotedStr(AStr, aspasimples)
}
func (cm *TConnectionMetadata) InternalEscapeTime(AStr string) string {
	return AnsiQuotedStr(AStr, aspasimples)
}
func (cm *TConnectionMetadata) InternalEscapeDateTime(AStr string) string {
	return AnsiQuotedStr(AStr, aspasimples)
}
func (cm *TConnectionMetadata) InternalEscapeFloat(AStr string) string {
	return AStr
}
func (cm *TConnectionMetadata) InternalEscapeString(AStr string) string {
	return AnsiQuotedStr(AStr, aspasimples)
}
func (cm *TConnectionMetadata) InternalEscapeEscape(AEscape byte, AStr string) string {
	return AStr + " ESCAPE " + QuotedStr(string(AEscape))
}
func (cm *TConnectionMetadata) InternalEscapeInto(AStr string) string {
	return ""
}
func (cm *TConnectionMetadata) InternalEscapeFunction(ASeq *TEscapeData) string {
	Result := ASeq.Name
	if len(ASeq.Args) > 0 {
		Result += "(" + cm.AddEscapeSequenceArgs(ASeq) + ")"
	}
	return Result
}
func (cm *TConnectionMetadata) EscapeFuncToID(ASeq *TEscapeData) error {
	var (
		sName string
		eFunc TEscapeFunction
	)

	sName = ASeq.Name
	// character
	if sName == "ASCII" {
		eFunc = efASCII
	} else if sName == "LTRIM" {
		eFunc = efLTRIM
	} else if sName == "REPLACE" {
		eFunc = efREPLACE
	} else if sName == "RTRIM" {
		eFunc = efRTRIM
	} else if sName == "DECODE" {
		eFunc = efDECODE
	} else if sName == "BIT_LENGTH" {
		eFunc = efBIT_LENGTH
	} else if sName == "CHAR" {
		eFunc = efCHAR
	} else if (sName == "CHAR_LENGTH") || (sName == "CHARACTER_LENGTH") {
		eFunc = efCHAR_LENGTH
	} else if sName == "CONCAT" {
		eFunc = efCONCAT
	} else if sName == "INSERT" {
		eFunc = efINSERT
	} else if sName == "LCASE" {
		eFunc = efLCASE
	} else if sName == "LEFT" {
		eFunc = efLEFT
	} else if sName == "LENGTH" {
		eFunc = efLENGTH
	} else if sName == "LOCATE" {
		eFunc = efLOCATE
	} else if sName == "OCTET_LENGTH" {
		eFunc = efOCTET_LENGTH
	} else if sName == "POSITION" {
		eFunc = efPOSITION
	} else if sName == "REPEAT" {
		eFunc = efREPEAT
	} else if sName == "RIGHT" {
		eFunc = efRIGHT
	} else if sName == "SPACE" {
		eFunc = efSPACE
	} else if sName == "SUBSTRING" {
		eFunc = efSUBSTRING
	} else if sName == "UCASE" {
		eFunc = efUCASE
		// numeric
	} else if sName == "ACOS" {
		eFunc = efACOS
	} else if sName == "ASIN" {
		eFunc = efASIN
	} else if sName == "ATAN" {
		eFunc = efATAN
	} else if sName == "CEILING" {
		eFunc = efCEILING
	} else if sName == "DEGREES" {
		eFunc = efDEGREES
	} else if sName == "LOG" {
		eFunc = efLOG
	} else if sName == "LOG10" {
		eFunc = efLOG10
	} else if sName == "PI" {
		eFunc = efPI
	} else if sName == "RADIANS" {
		eFunc = efRADIANS
	} else if sName == "RANDOM" {
		eFunc = efRANDOM
	} else if sName == "TRUNCATE" {
		eFunc = efTRUNCATE
	} else if sName == "ABS" {
		eFunc = efABS
	} else if sName == "COS" {
		eFunc = efCOS
	} else if sName == "EXP" {
		eFunc = efEXP
	} else if sName == "FLOOR" {
		eFunc = efFLOOR
	} else if sName == "MOD" {
		eFunc = efMOD
	} else if sName == "POWER" {
		eFunc = efPOWER
	} else if sName == "ROUND" {
		eFunc = efROUND
	} else if sName == "SIGN" {
		eFunc = efSIGN
	} else if sName == "SIN" {
		eFunc = efSIN
	} else if sName == "SQRT" {
		eFunc = efSQRT
	} else if sName == "TAN" {
		eFunc = efTAN
		// date and time
	} else if (sName == "CURRENT_DATE") || (sName == "CURDATE") {
		eFunc = efCURDATE
	} else if (sName == "CURRENT_TIME") || (sName == "CURTIME") {
		eFunc = efCURTIME
	} else if (sName == "CURRENT_TIMESTAMP") || (sName == "NOW") {
		eFunc = efNOW
	} else if sName == "DAYNAME" {
		eFunc = efDAYNAME
	} else if sName == "DAYOFMONTH" {
		eFunc = efDAYOFMONTH
	} else if sName == "DAYOFWEEK" {
		eFunc = efDAYOFWEEK
	} else if sName == "DAYOFYEAR" {
		eFunc = efDAYOFYEAR
	} else if sName == "EXTRACT" {
		eFunc = efEXTRACT
	} else if sName == "HOUR" {
		eFunc = efHOUR
	} else if sName == "MINUTE" {
		eFunc = efMINUTE
	} else if sName == "MONTH" {
		eFunc = efMONTH
	} else if sName == "MONTHNAME" {
		eFunc = efMONTHNAME
	} else if sName == "QUARTER" {
		eFunc = efQUARTER
	} else if sName == "SECOND" {
		eFunc = efSECOND
	} else if sName == "TIMESTAMPADD" {
		eFunc = efTIMESTAMPADD
	} else if sName == "TIMESTAMPDIFF" {
		eFunc = efTIMESTAMPDIFF
	} else if sName == "WEEK" {
		eFunc = efWEEK
	} else if sName == "YEAR" {
		eFunc = efYEAR
		// system
	} else if sName == "CATALOG" {
		eFunc = efCATALOG
	} else if sName == "SCHEMA" {
		eFunc = efSCHEMA
	} else if sName == "IFNULL" {
		eFunc = efIFNULL
	} else if (sName == "IF") || (sName == "IIF") {
		eFunc = efIF
	} else if sName == "LIMIT" {
		eFunc = efLIMIT
		// convert
	} else if sName == "CONVERT" {
		eFunc = efCONVERT
	} else {
		eFunc = efNONE
		// unsupported ATAN2, COT
		return fmt.Errorf("eFunc are not supported")
	}
	ASeq.Func = eFunc
	return nil
}
func (cm *TConnectionMetadata) GetNameParts() TNameParts {
	return TNameParts{npBaseObject, npObject}
}
func (cm *TConnectionMetadata) EncodeObjName(AParsedName TParsedName, ACommand IFDPhysCommand, AOpts TEncodeOptions) string {
	var (
		rName  TParsedName
		Result string
		//eParts TNameParts
	)
	//eParts = cm.GetNameParts()
	rName = AParsedName
	if In(eoBeautify, AOpts) {
		// if !In(npCatalog, eParts) || In(npCatalog, FConnectionObj.RemoveDefaultMeta) || (rName.FCatalog == "*") ||
		// 	(AnsiCompareText(FConnectionObj.DefaultCatalog, rName.FCatalog) == 0) {
		// 	rName.FCatalog = emptstr
		// }
		// if !In(npSchema, eParts) || In(npSchema, FConnectionObj.RemoveDefaultMeta) || (rName.FSchema == "*") ||
		// 	(AnsiCompareText(FConnectionObj.DefaultSchema, rName.FSchema) == 0) {
		// 	rName.FSchema == emptstr
		// }
		if rName.FBaseObject == "*" {
			rName.FBaseObject = emptstr
		}
		if rName.FObject == "*" {
			rName.FObject = emptstr
		}
		// if !In(eoQuote, AOpts) {
		// 	rName.FCatalog = QuoteNameIfReq(rName.FCatalog, npCatalog)
		// 	rName.FSchema = QuoteNameIfReq(rName.FSchema, npSchema)
		// 	rName.FBaseObject = QuoteNameIfReq(rName.FBaseObject, npBaseObject)
		// 	rName.FObject = QuoteNameIfReq(rName.FObject, npObject)
		// }
	}
	// else if In(eoNormalize, AOpts) {
	// 	rName.FCatalog = NormObjName(rName.FCatalog, npCatalog)
	// 	rName.FSchema = NormObjName(rName.FSchema, npSchema)
	// 	rName.FBaseObject = NormObjName(rName.FBaseObject, npBaseObject)
	// 	rName.FObject = NormObjName(rName.FObject, npObject)
	// }
	// if In(eoQuote, AOpts) {
	// 	rName.FCatalog = QuoteObjName(rName.FCatalog, npCatalog)
	// 	rName.FSchema = QuoteObjName(rName.FSchema, npSchema)
	// 	rName.FBaseObject = QuoteObjName(rName.FBaseObject, npBaseObject)
	// 	rName.FObject = QuoteObjName(rName.FObject, npObject)
	// }
	//Result = InternalEncodeObjName(rName, ACommand)
	Result = rName.FBaseObject
	return Result
}

func (cm *TConnectionMetadata) TranslateEscapeSequence(seq *TEscapeData) (string, error) {
	switch seq.Kind {
	case eskFloat:
		return cm.InternalEscapeFloat(UnQuoteBase(seq.Args[0], aspasimples, aspasimples)), nil
	case eskDate:
		return cm.InternalEscapeDate(UnQuoteBase(seq.Args[0], aspasimples, aspasimples)), nil
	case eskTime:
		return cm.InternalEscapeTime(UnQuoteBase(seq.Args[0], aspasimples, aspasimples)), nil
	case eskDateTime:
		return cm.InternalEscapeDateTime(UnQuoteBase(seq.Args[0], aspasimples, aspasimples)), nil
	// case eskIdentifier:
	// 	rName := UnQuoteBase(seq.Args[0], 39, 39)
	// 	return cm.EncodeObjName(rName, nil, TEncodeOptions{eoQuote}), nil
	case eskBoolean:
		return cm.InternalEscapeBoolean(UnQuoteBase(seq.Args[0], aspasimples, aspasimples)), nil
	case eskString:
		return cm.InternalEscapeString(UnQuoteBase(seq.Args[0], aspasimples, aspasimples)), nil
	// case eskFunction:
	// 	err := cm.EscapeFuncToID(seq)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	// LIMIT has a syntax of a function, but it is special escape sequence,
	// 	// processed internally by command preprocessor and TPhysCommand
	// 	if seq.Func == efLIMIT {
	// 		return "", nil
	// 	} else {
	// 		defer func() {
	// 			if r := recover(); r != nil {
	// 				if e, ok := r.(Exception); ok {
	// 					return "", e
	// 				}

	// 			}
	// 		}()

	// 		return cm.InternalEscapeFunction(seq), nil
	// 	}
	// case eskIIF:
	// 	i := 0
	// 	var s string
	// 	var lCurrent bool
	// 	result := ""
	// 	for i < len(seq.Args) && i&1 != 1 {
	// 		s = strings.TrimSpace(seq.Args[i])
	// 		lCurrent = false
	// 		if IsRDBMSKind(s, &lCurrent) {
	// 			if lCurrent {
	// 				result = seq.Args[i+1]
	// 				break
	// 			}
	// 		} else if s != "" {
	// 			result = seq.Args[i+1]
	// 			break
	// 		}
	// 		i += 2
	// 	}
	// 	if i == len(seq.Args)-1 {
	// 		result = seq.Args[i]
	// 	}
	// 	return result
	// case eskIF:
	// 	s := strings.TrimSpace(seq.Args[0])
	// 	var lCurrent bool
	// 	if IsRDBMSKind(s, &lCurrent) {
	// 		if lCurrent {
	// 			return "True"
	// 		}
	// 	} else if s != "" {
	// 		return "True"
	// 	}
	// 	return ""
	case eskFI:
		return "", nil
	case eskEscape:
		return cm.InternalEscapeEscape(seq.Args[0][1], seq.Args[1]), nil
	case eskInto:
		return cm.InternalEscapeInto(seq.Args[0]), nil
	default:
		return "", nil
	}
}

func NewConnectionMetadata() *TConnectionMetadata {
	return &TConnectionMetadata{}
}
