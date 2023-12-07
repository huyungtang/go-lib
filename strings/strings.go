package strings

import (
	"fmt"
	"math/rand"
	"regexp"
	base "strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Code128A
// ****************************************************************************************************************************************
func Code128A(s string) string {
	bs := make([]rune, 1, len(s)+3)
	bs[0] = 203
	ck := 0
	for i, j := range s {
		if j >= 32 {
			ck += int(j-32) * (i + 1)
		} else {
			ck += int(j+64) * (i + 1)
		}
		bs = append(bs, j)
	}

	if c := ck % 103; c < 95 {
		ck = c + 32
	} else {
		ck = c + 100
	}

	bs = append(bs, rune(ck), 206)

	return string(bs)
}

// Currency
// ****************************************************************************************************************************************
func Currency(format string, val interface{}) string {
	return message.NewPrinter(language.English).Sprintf(format, val)
}

// Format
// ****************************************************************************************************************************************
func Format(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// Find
// ****************************************************************************************************************************************
func Find(str, pattern string) string {
	return regexp.MustCompile(pattern).FindString(str)
}

// Fixed
// ****************************************************************************************************************************************
func Fixed(str string, ml int) string {
	out := make([]rune, 0, len([]rune(str)))
	cl := 0
	for _, j := range str {
		if j == 10 {
			cl = 0
		} else {
			bs := []byte(string(j))
			if l := len(bs); cl+l > ml {
				out = append(out, 10)
				cl = l
			} else {
				cl += l
			}
		}

		out = append(out, j)
	}

	return string(out)
}

// HasPrefix
// ****************************************************************************************************************************************
func HasPrefix(str, pre string) bool {
	return base.HasPrefix(str, pre)
}

// HasSuffix
// ****************************************************************************************************************************************
func HasSuffix(str, suf string) bool {
	return base.HasSuffix(str, suf)
}

// Index
// ****************************************************************************************************************************************
func Index(s string, sp string) int {
	return base.Index(s, sp)
}

// Join
// ****************************************************************************************************************************************
func Join(strs []string, sep string) string {
	return base.Join(strs, sep)
}

// LastIndex
// ****************************************************************************************************************************************
func LastIndex(s, ss string) int {

	return base.LastIndex(s, ss)
}

// OmitEmpty
// ****************************************************************************************************************************************
func OmitEmpty(strs []string) []string {
	return omit(strs, func(s string) bool {
		return s == ""
	})
}

// Omit
// ****************************************************************************************************************************************
func Omit(strs []string, pattern string) []string {
	r := regexp.MustCompile(pattern)

	return omit(strs, func(s string) bool {
		return r.MatchString(s)
	})
}

// Parse
// ****************************************************************************************************************************************
func Parse(str, format string, args ...interface{}) (err error) {
	_, err = fmt.Sscanf(str, format, args...)

	return
}

// Random
// ****************************************************************************************************************************************
func Random(size int) string {
	rs := make([]rune, size)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; {
		if r := (rd.Intn(74) + 48); (r >= 48 && r <= 57) || (r >= 65 && r <= 90) || (r >= 97 && r <= 122) {
			rs[i] = rune(r)
			i++
		}
	}

	return string(rs)
}

// Replace
// ****************************************************************************************************************************************
func Replace(str, old, new string) string {
	return base.ReplaceAll(str, old, new)
}

// Reverse
// ****************************************************************************************************************************************
func Reverse(str string) string {
	r := []rune(str)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	return string(r)
}

// Split
// ****************************************************************************************************************************************
func Split(s, sep string) []string {
	return base.Split(s, sep)
}

// ToLower
// ****************************************************************************************************************************************
func ToLower(str string) string {
	return base.ToLower(str)
}

// ToUpper
// ****************************************************************************************************************************************
func ToUpper(str string) string {
	return base.ToUpper(str)
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// omit ***********************************************************************************************************************************
func omit(strs []string, fn func(string) bool) (rtn []string) {
	rtn = make([]string, 0, len(strs))
	for _, v := range strs {
		if !fn(v) {
			rtn = append(rtn, v)
		}
	}

	return
}
