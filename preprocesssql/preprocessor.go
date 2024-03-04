package preprocesssql

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type TPhysPreprocessor struct {
	Source      string
	Destination string
	Kind        TDBMSKind
	//-----------------------------
	Params       *TParams
	MacrosUpd    *TMacros
	MacrosRead   *TMacros
	Markers      TStrings
	Instrs       TPreprocessorInstrs
	ConnMetadata *TConnectionMetadata
	Nested       bool

	//------------------------
	fNameDelims1 TListByte
	fNameDelims2 TListByte
	//------------------------
	fInNames        TNameQuoteLevels
	fLineSeparator  TTextEndOfLine
	fSQLCommandKind TCommandKind
	fParamMark      TParamMark
	//------------------------
	fSQLFromValue string
	//------------------------
	fSourceLen        int
	fSrcIndex         int
	fCommitedIndex    int
	fDestinationIndex int
	fBraceLevel       int
	fEscapeLevel      int
	fParamCount       int
	fSkipEscapes      int
	fSQLOrderByPos    int
	fSQLValuesPos     int
	fSQLLimitSkip     int
	fSQLLimitRows     int
	//------------------------
	fCh     byte
	fPrevCh byte
	//------------------------
	fInComment1         bool
	fInComment2         bool
	fInStr1             bool
	fInStr2             bool
	fInStr3             bool
	fInMySQLConditional bool
	fInProgramBlock     bool
	fWasIntoEscape      bool

	fInIntoEscape bool
	//----

}

func NewPhysPreprocessor() *TPhysPreprocessor {
	return &TPhysPreprocessor{
		fLineSeparator:      elDefault,
		Kind:                mkOracle,
		Destination:         "",
		ConnMetadata:        NewConnectionMetadata(),
		Params:              NewParams(),
		MacrosUpd:           NewMacros(),
		MacrosRead:          NewMacros(),
		fSQLFromValue:       "",
		fCh:                 byte(nullChar),
		fPrevCh:             byte(nullChar),
		fSQLOrderByPos:      0,
		fSQLValuesPos:       0,
		fSrcIndex:           -1,
		fCommitedIndex:      0,
		fEscapeLevel:        0,
		fBraceLevel:         0,
		fParamCount:         0,
		fSQLCommandKind:     skUnknown,
		fSQLLimitSkip:       0,
		fSQLLimitRows:       -1,
		fInComment1:         false,
		fInMySQLConditional: false,
		fInComment2:         false,
		fInStr1:             false,
		fInStr2:             false,
		fInStr3:             false,
		fInProgramBlock:     false,
		fWasIntoEscape:      false,
		fInNames:            NewTNameQuoteLevels(),
	}
}
func (pp *TPhysPreprocessor) PushWriter() {
	pp.fDestinationIndex = 1
	pp.fCommitedIndex = pp.fSrcIndex
}
func (pp *TPhysPreprocessor) ProcessIdentifier(ADotAllowed bool, AIsQuoted *bool) (string, error) {
	var (
		aBuff [256]byte
		i     int
	//	eQuote TNameQuoteLevel
	)

	ProcessQuotedDelim := func(ADelim1, ADelim2 byte) (string, bool, error) {
		result := ""
		if ADelim1 != 0 && ADelim1 != ' ' && pp.fCh == ADelim1 {
			*AIsQuoted = true
			i := 0
			for {
				i++
				if i == 256 {
					return "", false, fmt.Errorf("ProcessIdentifier: too long identifier")
				}
				aBuff[i-1] = pp.GetChar()
				if aBuff[i-1] == 0 || aBuff[i-1] == ADelim2 {
					break
				}
			}
			result = string(aBuff[:i])
			return result, true, nil
		}
		return result, false, nil
	}

	Result := ""
	*AIsQuoted = false
	i = -1
	if pp.ConnMetadata != nil {
		pp.GetChar()
		NQL := TNameQuoteLevels{ncDefault, ncSecond, ncThird}
		for _, eQuote := range NQL {
			_, found, err := ProcessQuotedDelim(
				pp.ConnMetadata.GetNameQuoteChar(eQuote, nsLeft),
				pp.ConnMetadata.GetNameQuoteChar(eQuote, nsRight))
			if err != nil {
				return "", err
			}
			if found {
				break
			}
		}
		pp.PutBack()
	}

	for {
		i++
		if i == 256 {
			return "", fmt.Errorf("ProcessIdentifier: too long identifier")
		}
		aBuff[i] = pp.GetChar()
		if !(In(aBuff[i], []byte{
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'#', '$', '_'}) || (ADotAllowed && aBuff[i] == '.')) && !unicode.IsLetter(rune(aBuff[i])) {
			break
		}
	}
	pp.PutBack()
	Result = string(aBuff[:i])
	return strings.ToUpper(Result), nil
}

func (pp *TPhysPreprocessor) UnCommit(AChars int) {
	pp.fCommitedIndex -= AChars
}
func (pp *TPhysPreprocessor) Commit(ASkip int) {
	var iLen int
	iLen = pp.fSrcIndex - pp.fCommitedIndex + ASkip
	if pp.fCommitedIndex+iLen >= pp.fSourceLen {
		iLen = pp.fSourceLen - pp.fCommitedIndex
	}
	if iLen > 0 {
		for pp.fDestinationIndex+iLen-1 > len(pp.Destination) {
			newLen := len(pp.Destination) * 2
			if newLen < pp.fDestinationIndex+iLen {
				newLen = pp.fDestinationIndex + iLen
			}
			newDest := make([]byte, newLen)
			copy(newDest, pp.Destination)
			pp.Destination = string(newDest)
		}
		pp.Destination = pp.Source[pp.fCommitedIndex+1 : pp.fSrcIndex-1]
		pp.fDestinationIndex += iLen
	}
	pp.fCommitedIndex = pp.fSrcIndex
}
func (pp *TPhysPreprocessor) ProcessParam() error {
	var (
		sName, sSubst string
		lIsQuoted     *bool
		iPar          int
		cPrevCh, cCh  byte
		err           error
		oPar          *TParam
	)
	cPrevCh = pp.fPrevCh
	cCh = pp.GetChar()

	BuildSubstName := func() error {
		return nil
	}
	//escrever depois
	if In(cCh, []byte{'=', ' ', 13, 10, 9, 0}) || In(cPrevCh, pp.fNameDelims1) || In(cPrevCh, pp.fNameDelims2) ||
		In(cPrevCh, []byte{
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'#', '$', '_'}) {
		// skip:
		// - PL/SQL assignment operator - :=
		// - TSQL label - name:
		// - Informix - catalog:schema.object
		// - if before ':' is an identifier
	} else if cCh == ':' {
		if In(pp.Kind, []TDBMSKind{mkPostgreSQL, mkMSSQL, mkAdvantage, mkInformix}) ||
			(pp.Kind == mkSQLite) && pp.Nested {
			// skip:
			// - PostgreSQL, Informix ::
			// - MSSQL SELECT ... FROM ::fn_xxxx()
			// - ADS SELECT ::conn.TransactionCount FROM system.iota
			// - SQLite SELECT xxxx AS "nm&1", &1=::type
		} else {
			pp.Commit(-1)
		}
	} else {
		pp.Commit(-2)
		pp.PutBack()
		lIsQuoted = new(bool)
		*lIsQuoted = false
		if pp.Kind == mkOracle {
			sName, err = pp.ProcessIdentifier(false, lIsQuoted)
			if err != nil {
				return err
			}
			// skip:
			// - Oracle triggers :new. / :old.
			if pp.GetChar() == '.' {
				pp.UnCommit(2)
				pp.Commit(0)
				return nil
			} else {
				pp.PutBack()
			}
		} else {
			sName, err = pp.ProcessIdentifier(true, lIsQuoted)
			if err != nil {
				return err
			}
		}

		iPar = -1
		oPar = nil
		sSubst = sName

		if In(piCreateParams, pp.Instrs) {
			iPar = pp.Params.IndexOf(sName)
			if (pp.Params.BindMode == pbByNumber) || (iPar == -1) {
				BuildSubstName()
				if (oPar == nil) || (pp.Params.BindMode == pbByNumber) {
					oPar = pp.Params.NewParam()
					oPar.Name = sName
					if pp.Params.BindMode == pbByNumber {
						oPar.Position = pp.Params.Count()
					}
					oPar.IsCaseSensitive = *lIsQuoted
				}
			} else {
				sSubst = ":" + sSubst
				oPar = pp.Params.Items[iPar]
			}
			if In(oPar.ParamType, TParamTypes{ptUnknown, ptInput}) {
				if pp.fInIntoEscape {
					if (iPar == -1) && (pp.Kind == mkFirebird) {
						oPar.ParamType = ptOutput
					} else {
						oPar.ParamType = ptInputOutput
					}

				}
			}
		} else if pp.fParamMark == prName {
			// The following is needed, when piCreateParams is not included, but a SQL
			// will be written in any case. So, we should take into account ParamNameMaxLength
			// and lIsQuoted for prName markers.
			oPar = pp.Params.FindParam(sName)
			BuildSubstName()
		} else {
			sSubst = ":" + sSubst
		}

		if In(piExpandParams, pp.Instrs) {
			pp.Params.Markers.Add(sName)
			switch pp.fParamMark {
			case prQMark:
				sSubst = "?"
			case prNumber:
				{
					//Inc(FParamCount)
					sSubst = ":" + IntToStr(pp.Params.Count())
				}
			case prDollar:
				if pp.Params.BindMode == pbByNumber {
					//Inc(FParamCount)
					sSubst = "$" + IntToStr(pp.Params.Count())
				} else {
					if oPar == nil {
						oPar = pp.Params.FindParam(sName)
					}

					if oPar != nil {
						sSubst = "$" + IntToStr(oPar.Index+1)
					} else {
						sSubst = string(nullChar)
					}

				}
			case prQNumber:
				{
					// Inc(FParamCount);
					sSubst = "?" + IntToStr(pp.Params.Count())
				}
			}
		}
		pp.WriteStr(sSubst)
	}
	return nil
}
func (pp *TPhysPreprocessor) WriteStr(AStr string) {
	var iLen int
	iLen = len(AStr)
	if iLen > 0 {
		// for pp.fDestinationIndex+iLen-1 > len(pp.Destination) {
		// 	newLen := len(pp.Destination) * 2
		// 	if newLen < pp.fDestinationIndex+iLen {
		// 		newLen = pp.fDestinationIndex + iLen
		// 	}
		// 	newDest := make([]byte, newLen)
		// 	copy(newDest, pp.Destination)
		// 	pp.Destination = string(newDest)
		// }
		//pp.Destination = AStr[pp.fCommitedIndex+1 : pp.fSrcIndex-1]
		pp.Destination += AStr
		pp.fDestinationIndex += iLen
	}
	pp.fCommitedIndex = pp.fSrcIndex
}
func (pp *TPhysPreprocessor) Missed(AStr string) error {
	err := NewFDException(AStr)
	return err
}
func (pp *TPhysPreprocessor) SkipWS() {
	var ch byte
	for {
		ch = pp.GetChar()
		if ch > ' ' || ch == '\x00' {
			break
		}
	}
	if ch != '\x00' {
		pp.PutBack()
	}
}
func (pp *TPhysPreprocessor) ProcessQuestion() {
	//escrever depois
}
func (pp *TPhysPreprocessor) ProcessQuoteTag() {
	//escrever depois
}
func (pp *TPhysPreprocessor) ProcessMacro(AFirstCh byte) error {
	var (
		sName  string
		sRes   string
		oMacro *TMacro
		//lIsRaw,
		lProcessRes bool
		lIsQuoted   *bool
		//i                   int
		oPP *TPhysPreprocessor
		cCh byte
		err error
	)
	lIsQuoted = new(bool)
	//	lIsRaw = (AFirstCh == '!')
	cCh = pp.GetChar()
	// if In(cCh, []byte{'=', '<', '>', '&', '''', ' ', #13, #10, #9, #0}) {
	if In(cCh, []byte{'=', '<', '>', '&', aspasimples, ' ', 13, 10, 9, 0}) {
		// skip:
		// - !=, !<, !>, &=, &&, & operators
		// - } of string literal
		// - delimiters
		if cCh == aspasimples {
			pp.PutBack()
		}
	} else if cCh == AFirstCh {
		pp.Commit(-1)
	} else {
		if In(piExpandMacros, pp.Instrs) {
			pp.Commit(-2)
		}
		pp.PutBack()

		*lIsQuoted = false
		sName, err = pp.ProcessIdentifier(false, lIsQuoted)
		if err != nil {
			return err
		}
		if (pp.MacrosUpd != nil) && In(piCreateMacros, pp.Instrs) {
			oMacro = pp.MacrosUpd.FindMacro(sName)
			if oMacro == nil {
				oMacro = pp.MacrosUpd.NewMacro()
				oMacro.Name = sName
				// if lIsRaw {
				// 	oMacro.DataType = mdRaw

				// } else {
				// 	oMacro.DataType = mdIdentifier
				// }
			}
		} else {
			oMacro = nil
		}

		if In(piExpandMacros, pp.Instrs) {
			if (pp.MacrosUpd != pp.MacrosRead) || (oMacro == nil) && !In(piCreateMacros, pp.Instrs) {
				oMacro = pp.MacrosRead.FindMacro(sName)
			}

			if oMacro != nil {
				sRes = oMacro.SQL
				lProcessRes = false
				for i := 0; i < len(sRes); i++ {
					if In(sRes[i], []byte{'!', '&', ':', '{'}) {
						lProcessRes = true
						break
					}
				}
				if lProcessRes {
					oPP = NewPhysPreprocessor()
					oPP.Nested = true
					oPP.ConnMetadata = pp.ConnMetadata
					oPP.Params = pp.Params
					oPP.MacrosUpd = pp.MacrosUpd
					oPP.MacrosRead = pp.MacrosRead
					oPP.Instrs = pp.Instrs
					oPP.Instrs.Remove(piParseSQL)
					//oPP.DesignMode = pp.DesignMode
					oPP.Source = sRes
					oPP.Execute()

					sRes = oPP.Destination

					oPP = nil
				}
				pp.WriteStr(sRes)
			} else {
				pp.WriteStr("")
			}

		}
	}
	return nil
}
func (pp *TPhysPreprocessor) TranslateEscape(aEscape *TEscapeData) (string, error) {
	result, err := pp.ConnMetadata.TranslateEscapeSequence(aEscape)
	if err != nil {
		return "", err
	}
	if aEscape.Kind == eskFunction && aEscape.Func == efLIMIT {
		if len(aEscape.Args) == 2 {
			pp.fSQLLimitSkip, _ = strconv.Atoi(aEscape.Args[0])
			pp.fSQLLimitRows, _ = strconv.Atoi(aEscape.Args[1])
		} else {
			pp.fSQLLimitRows, _ = strconv.Atoi(aEscape.Args[0])
		}
	}
	return result, nil
}
func (pp *TPhysPreprocessor) ProcessEscape() (TEscapeKind, error) {
	var (
		sKind string
		rEsc  *TEscapeData
		err   error
	)
	rEsc = NewEscapedData()
	// check for GUID and $DEFINE
	iPrevSrcIndex := pp.fSrcIndex
	iCnt := 0
	ch := pp.GetChar()
	if strings.ContainsAny(string(ch), "0123456789abcdefABCDEF-$") {
		iCnt++
	}

	if iCnt > 3 {
		pp.Commit(0)
		return eskText, nil
	}
	pp.fSrcIndex = iPrevSrcIndex

	// it is rather escape sequence
	pp.Commit(-1)
	pp.fEscapeLevel++
	pp.SkipWS()
	lTemp := false
	value, err := pp.ProcessIdentifier(false, &lTemp)
	sKind = strings.ToUpper(value)
	pp.SkipWS()
	switch sKind {
	case "E":
		rEsc.Kind = eskFloat
		rEsc.Args = make([]string, 1)
		value, err = pp.ProcessCommand()
		if err != nil {
			return eskText, err
		}
		rEsc.Args[0] = value
	case "D":
		rEsc.Kind = eskDate
		rEsc.Args = make([]string, 1)
		value, err = pp.ProcessCommand()
		if err != nil {
			return eskText, err
		}
		rEsc.Args[0] = value
	case "T":
		rEsc.Kind = eskTime
		rEsc.Args = make([]string, 1)
		value, err = pp.ProcessCommand()
		if err != nil {
			return eskText, err
		}
		rEsc.Args[0] = value
	case "DT":
		rEsc.Kind = eskDateTime
		rEsc.Args = make([]string, 1)
		value, err = pp.ProcessCommand()
		if err != nil {
			return eskText, err
		}
		rEsc.Args[0] = value
	case "ID":
		rEsc.Kind = eskIdentifier
		rEsc.Args = make([]string, 1)
		value, err = pp.ProcessCommand()
		if err != nil {
			return eskText, err
		}
		rEsc.Args[0] = value
	case "L":
		rEsc.Kind = eskBoolean
		rEsc.Args = make([]string, 1)
		value, err = pp.ProcessCommand()
		if err != nil {
			return eskText, err
		}
		rEsc.Args[0] = value
	case "S":
		rEsc.Kind = eskString
		rEsc.Args = make([]string, 1)
		pp.PushWriter()
		for {
			ch := pp.GetChar()
			switch ch {
			case '!', '&':
				if InIn(pp.Instrs, TPreprocessorInstrs{piExpandMacros, piCreateMacros}) {
					pp.ProcessMacro(pp.fCh)
					pp.GetChar()
				}
			case '\\':
				pp.Commit(-1)
				pp.GetChar()
				if pp.fCh == '}' {
					pp.GetChar()
				}
			}
			if pp.fCh == '}' || pp.fCh == 0 {
				break
			}
		}
		if pp.fCh != 0 {
			pp.PutBack()
		}
		rEsc.Args[0] = pp.PopWriter()
	case "ESCAPE":
		rEsc.Kind = eskEscape
		pp.SkipWS()
		pp.GetChar()
		if pp.fCh != '\'' {
			pp.Missed("'")
		}
		rEsc.Args = make([]string, 2)
		rEsc.Args[0] = string(pp.GetChar())
		pp.GetChar()
		if pp.fCh != '\'' {
			pp.Missed("'")
		}
		pp.SkipWS()
		value, err = pp.ProcessCommand()
		if err != nil {
			return eskText, err
		}
		rEsc.Args[1] = value
	case "INTO", "RETURNING_VALUES", "RETURNING":
		if pp.fInIntoEscape {
			pp.Missed("}")
		}
		rEsc.Kind = eskInto
		rEsc.Args = make([]string, 1)
		pp.fInIntoEscape = true
		defer func() {
			pp.fInIntoEscape = false
			pp.fWasIntoEscape = true
		}()
		value, err = pp.ProcessCommand()
		if err != nil {
			return eskText, err
		}
		rEsc.Args[0] = value
	case "IF":
		rEsc.Kind = eskIF
		rEsc.Args = make([]string, 1)
		ePrevInstrs := pp.Instrs
		pp.Instrs.Removes(TPreprocessorInstrs{piExpandParams, piCreateParams})
		defer func() {
			pp.Instrs = ePrevInstrs
		}()
		value, err = pp.ProcessCommand()
		if err != nil {
			return eskText, err
		}
		rEsc.Args[0] = value
	case "FI":
		rEsc.Kind = eskFI
	case "IIF":
		rEsc.Kind = eskIIF
		pp.GetChar()
		if pp.fCh != '(' {
			pp.Missed("(")
		}
		pp.fBraceLevel++
		for {
			pp.GetChar()
			if pp.fCh == ')' || pp.fCh == 0 {
				break
			}
			if pp.fCh != ')' {
				value, err = pp.ProcessCommand()
				if err != nil {
					return eskText, err
				}
				rEsc.Args = append(rEsc.Args, value)
				pp.GetChar()
				if pp.fCh != ')' {
					value, err = pp.ProcessCommand()
					if err != nil {
						return eskText, err
					}
					rEsc.Args = append(rEsc.Args, value)
					pp.GetChar()
				}
			}
		}
		if pp.fCh == ')' {
			pp.fBraceLevel--
		}
		if rEsc.Args[len(rEsc.Args)-1] == "" {
			rEsc.Args = rEsc.Args[:len(rEsc.Args)-1]
		}
		pp.SkipWS()
	case "STATIC":
		rEsc.Kind = eskText
		pp.WriteStr("{static}")
	default:
		rEsc.Kind = eskFunction
		if sKind == "FN" {
			rEsc.Name, err = pp.ProcessIdentifier(false, &lTemp)
			if err != nil {
				return eskText, err
			}
		} else {
			rEsc.Name = sKind
		}
		if rEsc.Name == "" {
			//return eskText, pp.ErrorEmptyName()
			return eskText, fmt.Errorf("ProcessEscape: empty name")
		}
		pp.SkipWS()
		pp.GetChar()
		if pp.fCh != '(' {
			pp.Missed("(")
		}
		pp.fBraceLevel++
		for {
			pp.GetChar()
			if pp.fCh != ')' {
				pp.PutBack()
				value, err := pp.ProcessCommand()
				if err != nil {
					return eskText, err
				}
				rEsc.Args = append(rEsc.Args, value)
				pp.GetChar()
			}
			if pp.fCh == ')' || pp.fCh == 0 {
				break
			}
		}
		if pp.fCh == ')' {
			pp.fBraceLevel--
		}
		pp.SkipWS()
	}
	if pp.GetChar() != '}' {
		pp.Missed("}")
	}
	pp.fEscapeLevel--
	if rEsc.Kind == eskIF {
		value, err := pp.TranslateEscape(rEsc)
		if err != nil {
			return eskText, err
		}
		if value == "" {
			ePrevInstrs := pp.Instrs
			pp.Instrs.Removes(TPreprocessorInstrs{piExpandMacros, piExpandParams, piCreateParams})
			defer func() {
				pp.fSkipEscapes--
				pp.Instrs = ePrevInstrs
			}()
			pp.fSkipEscapes++
			pp.ProcessCommand()
		} else {
			value, err := pp.TranslateEscape(rEsc)
			if err != nil {
				return eskText, err
			}
			pp.WriteStr(value)
		}
	} else if pp.fSkipEscapes == 0 {
		value, err := pp.TranslateEscape(rEsc)
		if err != nil {
			return eskText, err
		}
		pp.WriteStr(value)
	}
	return rEsc.Kind, nil
}
func (pp *TPhysPreprocessor) PutBack() {
	pp.fSrcIndex--
	pp.fPrevCh = pp.Source[pp.fSrcIndex]
	pp.fCh = pp.Source[pp.fSrcIndex]
}
func (pp *TPhysPreprocessor) GetChar() byte {
	var (
		ret byte
	)
	pp.fSrcIndex++
	if pp.fSrcIndex > pp.fSourceLen-1 {
		ret = 0
	} else {
		ret = pp.Source[pp.fSrcIndex]
	}
	pp.fPrevCh = pp.fCh
	pp.fCh = ret
	return ret
}
func (pp *TPhysPreprocessor) PopWriter() string {
	pp.Commit(0)

	return ""
	//return Copy(pp.Destination, 1, pp.fDestinationIndex-1)
}
func (pp *TPhysPreprocessor) ProcessCommand() (string, error) {
	var (
		iEnterBraceLevel int
	)
	pp.PushWriter()
	iEnterBraceLevel = pp.fBraceLevel

OuterLoop1:
	for {
		pp.GetChar()
		switch pp.fCh {
		case '}':
			if In(piExpandEscapes, pp.Instrs) &&
				!pp.fInComment1 && !pp.fInComment2 && !pp.fInStr1 && !pp.fInStr2 && (len(pp.fInNames) == 0) {
				break OuterLoop1
			}
		case '(':
			if !pp.fInComment1 && !pp.fInComment2 && !pp.fInStr1 && !pp.fInStr2 && (len(pp.fInNames) == 0) {
				pp.fBraceLevel++
			}
		case ')':
			if !pp.fInComment1 && !pp.fInComment2 && !pp.fInStr1 && !pp.fInStr2 && (len(pp.fInNames) == 0) {
				if (pp.fEscapeLevel > 0) && (pp.fBraceLevel == iEnterBraceLevel) {
					break OuterLoop1
				}
				pp.fBraceLevel--
			}
		case '\\':
			if !(In(piExpandEscapes, pp.Instrs) || pp.fInComment1 || pp.fInComment2 || pp.fInStr1 || pp.fInStr2 || (len(pp.fInNames) != 0)) {
				pp.GetChar()
			} else {
				pp.Commit(-1)
				pp.GetChar()
			}
		case '/':
			pp.GetChar()
			if !pp.fInComment1 && !pp.fInComment2 && !pp.fInStr1 && !pp.fInStr2 && !pp.fInStr3 && (len(pp.fInNames) == 0) && (pp.fCh == '*') {
				if pp.Kind == mkMySQL {
					pp.GetChar()
					if pp.fCh == '!' {
						pp.fInMySQLConditional = true
					} else {
						pp.PutBack()
						pp.fInComment1 = true
					}
				} else {
					pp.fInComment1 = true
				}
			} else {
				pp.PutBack()
			}
		case '*':
			pp.GetChar()
			if !pp.fInComment2 && !pp.fInStr1 && !pp.fInStr2 && !pp.fInStr3 && (len(pp.fInNames) == 0) && (pp.fCh == '/') {
				if pp.fInMySQLConditional {
					pp.fInMySQLConditional = false
				} else {
					pp.fInComment1 = false
				}
			} else {
				pp.PutBack()
			}
		case '-':
			pp.GetChar()
			if !pp.fInComment1 && !pp.fInStr1 && !pp.fInStr2 && !pp.fInStr3 && (len(pp.fInNames) == 0) && (pp.fCh == '-') {
				pp.fInComment2 = true
			} else {
				pp.PutBack()
			}
		case '\'':
			pp.GetChar()
			if !pp.fInComment1 && !pp.fInComment2 && (len(pp.fInNames) == 0) && !pp.fInStr2 && !pp.fInStr3 &&
				((pp.fCh != '\'') || !pp.fInStr1) {
				pp.PutBack()
				pp.fInStr1 = !pp.fInStr1
			}
		case '\r', '\n':
			if !pp.fInComment1 && pp.fInComment2 {
				pp.fInComment2 = false
			}
			if In(piTransformEOLs, pp.Instrs) {
				switch pp.fLineSeparator {
				case elUnix:
					if pp.fCh == '\r' {
						pp.Commit(-1)
						if pp.GetChar() != '\n' {
							pp.PutBack()
							pp.WriteStr("\n")
						}
					} else {
						pp.Commit(0)
						if pp.GetChar() == '\r' {
							pp.Commit(-1)
						} else {
							pp.PutBack()
						}
					}
				case elMac:
					if pp.fCh == '\n' {
						pp.Commit(-1)
						if pp.GetChar() != '\r' {
							pp.PutBack()
							pp.WriteStr("\r")
						}
					} else {
						pp.Commit(0)
						if pp.GetChar() == '\n' {
							pp.Commit(-1)
						} else {
							pp.PutBack()
							pp.WriteStr("\n")
						}
					}
				case elWindows:
					if pp.fCh == '\n' {
						pp.Commit(-1)
						if pp.GetChar() != '\r' {
							pp.PutBack()
						} else {
							pp.Commit(-1)
						}
						pp.WriteStr("\r\n")
					} else {
						pp.Commit(0)
						if pp.GetChar() == '\n' {
							pp.Commit(0)
						} else {
							pp.PutBack()
							pp.WriteStr("\n")
						}
					}
				}
			}
		case ':':
			if !pp.fInComment1 && !pp.fInComment2 && !pp.fInStr1 && !pp.fInStr2 &&
				!pp.fInStr3 && (len(pp.fInNames) == 0) &&
				!pp.fInProgramBlock && (InIn(TPreprocessorInstrs{piExpandParams, piCreateParams}, pp.Instrs)) {
				pp.ProcessParam()
			}
		case '?':
			if !pp.fInComment1 && !pp.fInComment2 && !pp.fInStr1 && !pp.fInStr2 &&
				!pp.fInStr3 && (len(pp.fInNames) == 0) &&
				(In(piTransformQuestions, pp.Instrs) && (InIn(TPreprocessorInstrs{piExpandParams, piCreateParams}, pp.Instrs))) {
				pp.ProcessQuestion()
			}
		case '{':
			if !pp.fInComment1 && !pp.fInComment2 && !pp.fInStr1 && !pp.fInStr2 && (len(pp.fInNames) == 0) &&
				In(piExpandEscapes, pp.Instrs) {

				esk, err := pp.ProcessEscape()
				if err != nil {
					return "", err
				}
				if esk == eskFI {
					if pp.fCh == '}' {
						pp.GetChar()
					}
					break OuterLoop1
				}
			}
		case ',':
			if !pp.fInComment1 && !pp.fInComment2 && !pp.fInStr1 && !pp.fInStr2 &&
				(iEnterBraceLevel > 0) && (iEnterBraceLevel == pp.fBraceLevel) {
				break OuterLoop1
			}
		case '!', '&':
			if (InIn(TPreprocessorInstrs{piExpandMacros, piCreateMacros}, pp.Instrs)) {
				pp.ProcessMacro(pp.fCh)
			}
		case '$':
			if pp.Kind == mkPostgreSQL {
				pp.ProcessQuoteTag()
			}
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
			'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
			'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			if (pp.ConnMetadata != nil) && In(pp.Kind, TDBMSKinds{mkInterbase, mkFirebird}) &&
				!pp.fInComment1 && !pp.fInComment2 && !pp.fInStr1 && !pp.fInStr2 && (len(pp.fInNames) == 0) &&
				!pp.fInProgramBlock {
				pp.fInProgramBlock =
					In(pp.fCh, []byte{'B', 'b'}) &&
						In(pp.GetChar(), []byte{'E', 'e'}) &&
						In(pp.GetChar(), []byte{'G', 'g'}) &&
						In(pp.GetChar(), []byte{'I', 'i'}) &&
						In(pp.GetChar(), []byte{'N', 'n'}) &&
						!In(pp.GetChar(), []byte{
							'0', '1', '2', '3', '4', '5', '6', '7', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g',
							'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
							'w', 'x', 'y', 'z', 'A', 'B',
							'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
							'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '#', '$', '_'})
				for In(pp.fCh, []byte{
					'0', '1', '2', '3', '4', '5', '6', '7', '9', 'a', 'b',
					'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
					'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D',
					'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
					'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '#', '$', '_'}) {
					pp.GetChar()
				}
				pp.PutBack()
			}
		case ';':
			if !pp.fInComment1 && !pp.fInComment2 && !pp.fInStr1 && !pp.fInStr2 && !pp.fInMySQLConditional && (pp.fBraceLevel == 0) {
				pp.fParamCount = 0
			}
		default:

			if (pp.fCh == '"') && !In(byte('"'), pp.fNameDelims1) {
				pp.GetChar()
				if !pp.fInComment1 && !pp.fInComment2 && (len(pp.fInNames) == 0) && !pp.fInStr1 && !pp.fInStr3 &&
					((byte(pp.fCh) != '"') || !pp.fInStr2) {
					pp.PutBack()
					pp.fInStr2 = !pp.fInStr2
				}
			} else if (In(byte(pp.fCh), pp.fNameDelims1) || In(byte(pp.fCh), pp.fNameDelims2)) &&
				!pp.fInComment1 && !pp.fInComment2 && !pp.fInStr1 && !pp.fInStr2 && !pp.fInStr3 {

				listaNQL := TNameQuoteLevels{ncDefault, ncSecond, ncThird}

			OuterLoop2:
				for _, eQuote := range listaNQL {
					cQuote1 := pp.ConnMetadata.GetNameQuoteChar(eQuote, nsLeft)
					cQuote2 := pp.ConnMetadata.GetNameQuoteChar(eQuote, nsRight)
					if pp.fCh == cQuote1 {
						if cQuote1 == cQuote2 {
							if In(eQuote, pp.fInNames) {
								pp.fInNames = Exclude(pp.fInNames, eQuote)
							} else {
								pp.fInNames = Include(pp.fInNames, eQuote)
							}
						} else {
							Include(pp.fInNames, eQuote)
						}
						break
					} else if pp.fCh == cQuote2 {
						if cQuote1 == cQuote2 {
							if In(eQuote, pp.fInNames) {
								pp.fInNames = Exclude(pp.fInNames, eQuote)
							} else {
								pp.fInNames = Include(pp.fInNames, eQuote)
							}
						} else {
							pp.fInNames = Exclude(pp.fInNames, eQuote)
						}
						break OuterLoop2
					}
				}
			}
		}
		if pp.fCh == nullChar {
			break OuterLoop1
		}
	}
	return pp.PopWriter(), nil
}
func (pp *TPhysPreprocessor) Execute() error {
	if pp.Source == "" {
		return nil
	}
	pp.Instrs = Include(pp.Instrs, piTransformQuestions)
	pp.fNameDelims1.Clear()
	pp.fNameDelims2.Clear()

	if pp.ConnMetadata != nil {
		lista := TNameQuoteLevels{ncDefault, ncSecond, ncThird}
		for _, eQuote := range lista {
			pp.fNameDelims1 = Include(pp.fNameDelims1, pp.ConnMetadata.GetNameQuoteChar(eQuote, nsLeft))
			pp.fNameDelims2 = Include(pp.fNameDelims2, pp.ConnMetadata.GetNameQuoteChar(eQuote, nsRight))
		}
	}
	pp.fDestinationIndex = 1
	pp.Markers.Clear()
	pp.fSourceLen = len(pp.Source)
	value, err := pp.ProcessCommand()
	if err != nil {
		return err
	}
	pp.Destination = value
	return nil
}

/*
PreprocessSQL:

	PreprocessSQL(commandtext string)(Params, MacrosUpd, MacrosRead,Errors)
*/
func PreprocessSQL(commandtext string, ACreateParams, ACreateMacros, AExpandMacros, AExpandEscape, AParseSQL bool) (*TParams, *TMacros, *TMacros, error) {
	oPrep := NewPhysPreprocessor()
	oPrep.Source = commandtext
	oPrep.Params.Clear()
	oPrep.MacrosUpd.Clear()
	oPrep.MacrosRead.Clear()
	oPrep.Instrs.Clear()
	if ACreateParams {
		oPrep.Instrs.Add(piCreateParams)
	}
	if ACreateMacros {
		oPrep.Instrs.Add(piCreateMacros)
	}
	if AExpandMacros {
		oPrep.Instrs.Add(piExpandMacros)
	}
	if AExpandEscape {
		oPrep.Instrs.Add(piExpandEscapes)
	}
	if AParseSQL {
		oPrep.Instrs.Add(piParseSQL)
	}
	oPrep.Execute()
	// for _, par := range oPrep.Params.Items {
	// 	fmt.Println("param.Name:", par.Name)
	// }
	// for _, mac := range oPrep.MacrosUpd.Items {
	// 	fmt.Println("macro.Name:", mac.Name)
	// }
	return oPrep.Params, oPrep.MacrosUpd, oPrep.MacrosRead, nil
}

func main() {
	sql := `select
   'S' AS GERAR_ENTIDADE_INTEGRADA
  ,NULL AS ID_ENTIDADE_INTEGRADA
  ,os.cod_empresa
  ,nvl(os.numero_os_fabrica, os.numero_os) as numero_os_fabrica
  ,trunc(os.data_encerrada) as data_encerrada
  ,os.hora_encerrada as hora_encerrada
  from os
  where 1 = 1
    &macroxxx
	and os.cod_empresa = :COD_EMPRESA
	and os.numero_os > 0
	and os.status_os = 1
	and nvl(os.orcamento, 'N') = 'N'
	and rownum <= 10 
	and trunc(os.data_encerrada) >= trunc(nvl('hh:mm:ss',sysdate))
	and trunc(os.data_encerrada) >= trunc(nvl(:DT_INICIO,sysdate))
	and trunc(os.data_encerrada) < trunc(sysdate + 1)
	 and trunc(os.data_encerrada) = :data_encerrada
	and not exists(SELECT 1 
				  FROM FAB_EI_OP5_MB_MGT a1
				  where a1.cod_empresa = os.cod_empresa
					and a1.numero_os_fabrica = nvl(os.numero_os_fabrica, os.numero_os)
					and trunc(a1.data_encerrada) = trunc(os.data_encerrada )
					and a1.hora_encerrada = os.hora_encerrada)`

	//sql = "select * from dual where id_xxx = :id_sss and xxx = :bbb_ff and yyy =:ccc_uu &hhhh');"
	Params, MacrosUpd, MacrosRead, err := PreprocessSQL(sql, true, true, true, true, true)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, p := range Params.Items {
			fmt.Println("param.Name:", p.Name)
		}
		for _, m := range MacrosUpd.Items {
			fmt.Println("macro.Name:", m.Name)
		}
		for _, m := range MacrosRead.Items {
			fmt.Println("macro.Name:", m.Name)
		}
	}

}
