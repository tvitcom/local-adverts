package util

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	// "strconv"
	"github.com/tvitcom/local-adverts/bkp/swearfilter"
)

func SwearDetector(message string) ([]string, error) {
	var swears = []string{`аеб`,`аёб`,`банамат`,`бля`,`ETC IN FULL VERSION`,`целка`,`чмо`,`чмыр`,`шалав`,`шлюх`,`шлюшк`,`ъеб`,`ьеб`,`ябывает`}
	filter := swearfilter.NewSwearFilter(true, swears...)
	return filter.Check(message)
}

func IsBadWords(str string, bads string) bool {
	// bads := "dsa,dri,def,da,dass,sdd"
	for _, v := range strings.Split(bads, ",") {
		re := regexp.MustCompile(`.?`+ v + `.?`)
		res := re.Match([]byte(str))
		if res {
			return true
		}
	}
	return false
}

//TODO use that: regexp.QuoteMeta(`Escaping symbols like: .+*?()|[]{}^$`)
func Addcslashes(s string) string {
	var	slashedCharList string = "_*~'`<>#|\\"
	var result []rune
	for _, ch := range []rune(s) {
		for _, v := range []rune(slashedCharList) {
			if ch == v {
				result = append(result, '\\')
			}
		}
		result = append(result, ch)
	}
	return string(result)
}

func PhoneNormalisation(n string) string {
	//remove `()-+:\b` in phone number string
	r := strings.NewReplacer("-", "", "(", "", ")", "", "+", "", " ", "", ":", "")
	return r.Replace(n)
}

func ExtractTitle(str string, limitChars, limitWords int) string {
  // Validation and Triming
  prefix := "Продам"
  str = strings.TrimSpace(strings.ToValidUTF8(str,""))
  var fields []string
  fields = strings.Fields(str)
  if strings.EqualFold(fields[0], prefix) {
      fields = fields[1:]
  }
  if foldedLen := len(fields); foldedLen > limitWords {
    fields[0] = strings.Title(fields[0])
    fields = fields[:limitWords]
  }
  return Substr(strings.Join(fields, " "), 0, limitChars)
}

func Substr(input string, start int, length int) string {
    asRunes := []rune(input)
    if start >= len(asRunes) {
        return ""
    }
    if start+length > len(asRunes) {
        length = len(asRunes) - start
    }
    return string(asRunes[start : start+length])
}

func ExtractDigitsString(str string) string {
	var digits []rune
	for _, c := range str {
	    if unicode.IsDigit(c) {
	        digits = append(digits, c)
	    }
	}
	return string(digits)
}

func Stringer(raw interface{}) string {
    return fmt.Sprint(raw)
}
