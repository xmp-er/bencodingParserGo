package processors

import (
	"errors"
	"strconv"
	"strings"

	"github.com/xmp-er/bencodingParserGo/validators"
)

func DecodeString(str string) (string, error) {
	str = strings.TrimSpace(str)
	s := strings.Split(str, ":")
	if !validators.IsValidString(s) {
		return "", errors.New("string input sent is not valid, corrupted data")
	}
	return s[1], nil
}

func DecodeInt(str string) (string, error) {
	if !validators.IsValidInt(str) {
		return "", errors.New("input for integer is not valid, corrupted data")
	}

	return str[1 : len(str)-1], nil
}

func DecodeList(str string, valid []bool) ([]interface{}, error) {
	var ret []interface{}
	for i, v := range str {
		if v == 'i' {
			var temp string = ""
			for {
				v = rune(str[i])
				temp += string(v)
				i++
				if v == 'e' {
					break
				}
			}
			t, err := DecodeInt(temp)
			if err != nil {
				return nil, err
			}
			ret = append(ret, t)
		} else if v == 'l' {
			var temp string = ""
			cnt := 0
			for {
				v = rune(str[i])
				temp += string(v)
				if !valid[i] && (v == 'l' || v == 'd') {
					cnt++
					i++
				}
				if !valid[i] && v == 'e' {
					cnt--
					if cnt == 0 {
						break
					}
				}
			}
			res, err := DecodeList(temp, valid)
			if err != nil {
				return nil, err
			}
			ret = append(ret, res)
		} else if v == 'd' {
			var temp string = ""
			cnt := 0
			for {
				v = rune(str[i])
				temp += string(v)
				if !valid[i] && (v == 'l' || v == 'd') {
					cnt++
					i++
				}
				if !valid[i] && v == 'e' {
					cnt--
					if cnt == 0 {
						break
					}
				}
			}
			res, err := DecodeDictionary(temp, valid)
			if err != nil {
				return nil, err
			}
			ret = append(ret, res)
		} else {
			t, err := strconv.Atoi(string(v))
			if err == nil {
				var temp string = ""
				for ; i <= (i + t + 1); i++ {
					temp += string(str[i])
				}
				res, err := DecodeString(temp)
				if err != nil {
					return nil, err
				}
				ret = append(ret, res)
				i--
			}
		}

	}
	return ret, nil
}

func DecodeDictionary(str string, valid []bool) (map[interface{}]interface{}, error) {
	m := make(map[interface{}]interface{})
	k := true
	var key interface{}
	var value interface{}
	for i, v := range str {
		if v == 'i' {
			var temp string = ""
			for {
				v = rune(str[i])
				temp += string(v)
				i++
				if v == 'e' {
					break
				}
			}
			t, err := DecodeInt(temp)
			if err != nil {
				return nil, err
			}
			if k {
				k = !k
				key = t
				value = nil
			} else if !k {
				k = !k
				value = t
				m[key] = value
				key = nil
				value = nil
			}
		}
		if v == 'l' {
			var temp string = ""
			cnt := 0
			for {
				v = rune(str[i])
				temp += string(v)
				if !valid[i] && (v == 'l' || v == 'd') {
					cnt++
					i++
				}
				if !valid[i] && v == 'e' {
					cnt--
					if cnt == 0 {
						break
					}
				}
			}
			t, err := DecodeList(temp, valid)
			if err != nil {
				return nil, err
			}
			if k {
				k = !k
				key = t
				value = nil
			} else if !k {
				k = !k
				value = t
				m[key] = value
				key = nil
				value = nil
			}
		}
		if v == 'd' {
			var temp string = ""
			cnt := 0
			for {
				v = rune(str[i])
				temp += string(v)
				if !valid[i] && (v == 'l' || v == 'd') {
					cnt++
					i++
				}
				if !valid[i] && v == 'e' {
					cnt--
					if cnt == 0 {
						break
					}
				}
			}
			t, err := DecodeDictionary(temp, valid)
			if err != nil {
				return nil, err
			}
			if k {
				k = !k
				key = t
				value = nil
			} else if !k {
				k = !k
				value = t
				m[key] = value
				key = nil
				value = nil
			}
		} else {
			t, err := strconv.Atoi(string(v))
			if err == nil {
				var temp string = ""
				for ; i <= (i + t + 1); i++ {
					temp += string(str[i])
				}
				res, err := DecodeString(temp)
				if err != nil {
					return nil, err
				}
				if k {
					k = !k
					key = res
					value = nil
				} else if !k {
					k = !k
					value = res
					m[key] = value
					key = nil
					value = nil
				}
				i--
			}
		}
	}
	return m, nil
}

// func Decode(str string) ([]interface{}, error) {
// 	var result string = ""

// }

func MarkStringAndInts(str string, valid *[]bool) {
	for i := 0; i < len(str); i++ {
		v := string(str[i])
		if v == "i" {
			(*valid)[i] = true
			for {
				v = string(str[i])
				(*valid)[i] = true
				i++
				if v == "e" {
					break
				}
			}
		} else {
			t, err := strconv.Atoi(string(v))
			if err == nil {
				(*valid)[i] = true
				i++
				(*valid)[i] = true
				lim := i + t + 1
				for ; i < lim; i++ {
					(*valid)[i] = true
				}
				i--
			}
		}
	}
}
