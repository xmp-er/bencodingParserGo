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
	for i := 0; i < len(str); i++ {
		v := string(str[i])
		if v == "i" {
			var temp string = ""
			for {
				v = string(str[i])
				temp += string(v)
				i++
				if v == "e" {
					break
				}
			}
			i--
			t, err := DecodeInt(temp)
			if err != nil {
				return nil, err
			}
			ret = append(ret, t)
		} else if v == "l" {
			start := i
			var temp string = ""
			cnt := 0
			for {
				v = string(str[i])
				temp += string(v)
				if !valid[i] && (v == "l" || v == "d") {
					cnt++
				}
				if !valid[i] && v == "e" {
					cnt--
					if cnt == 0 {
						break
					}
				}
				i++
			}
			end := i
			res, err := DecodeList(temp[1:len(temp)-1], valid[start+1:end])
			if err != nil {
				return nil, err
			}
			ret = append(ret, res)
		} else if v == "d" {
			start := i
			var temp string = ""
			cnt := 0
			for {
				v = string(str[i])
				temp += string(v)
				if !valid[i] && (v == "l" || v == "d") {
					cnt++
				}
				if !valid[i] && v == "e" {
					cnt--
					if cnt == 0 {
						break
					}
				}
				i++
			}
			end := i
			res, err := DecodeDictionary(temp[1:len(temp)-1], valid[start+1:end])
			if err != nil {
				return nil, err
			}
			ret = append(ret, res)
		} else {
			_, err := strconv.Atoi(string(v))
			if err == nil {
				var temp string = ""
				temp_cnt := 0
				o := i + 1
				for {
					if string(str[o]) == ":" {
						break
					}
					temp_cnt++
					v += string(str[o])
					o++
				}
				t, _ := strconv.Atoi(v)
				lim := i + len(v) + t + 1
				for ; i < lim; i++ {
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

func DecodeDictionary(str string, valid []bool) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	k := true
	var key string
	var value interface{}
	for i := 0; i < len(str); i++ {
		v := string(str[i])
		if v == "i" {
			var temp string = ""
			for {
				v = string(str[i])
				temp += string(v)
				i++
				if v == "e" {
					break
				}
			}
			i--
			t, err := DecodeInt(temp)
			if err != nil {
				return nil, err
			}
			if k {
				k = !k
				key = string(t)
				value = nil
			} else if !k {
				k = !k
				value = t
				m[key] = value
			}
		} else if v == "l" {
			var temp string = ""
			cnt := 0
			start := i
			for {
				v = string(str[i])
				temp += string(v)
				if !valid[i] && (v == "l" || v == "d") {
					cnt++
				}
				if !valid[i] && v == "e" {
					cnt--
					if cnt == 0 {
						break
					}
				}
				i++
			}
			end := i
			t, err := DecodeList(temp[1:len(temp)-1], valid[start+1:end])
			if err != nil {
				return nil, err
			}
			if k {
				k = !k
				temp_res := ""
				for _, v := range t {
					temp_res += v.(string)
				}
				key = string(temp_res)
				value = nil
			} else if !k {
				k = !k
				value = t
				m[key] = value
			}
		} else if v == "d" {
			var temp string = ""
			cnt := 0
			start := i
			for {
				v = string(str[i])
				temp += string(v)
				if !valid[i] && (v == "l" || v == "d") {
					cnt++
				}
				if !valid[i] && v == "e" {
					cnt--
					if cnt == 0 {
						break
					}
				}
				i++
			}
			end := i
			t, err := DecodeDictionary(temp[1:len(temp)-1], valid[start+1:end])
			if err != nil {
				return nil, err
			}
			if k {
				k = !k
				temp_res := ""
				for _, v := range t {
					temp_res += v.(string)
				}
				key = string(temp_res)
				value = nil
			} else if !k {
				k = !k
				value = t
				m[key] = value
			}
		} else {
			_, err := strconv.Atoi(string(v))
			if err == nil {
				var temp string = ""
				temp_cnt := 0
				o := i + 1
				for {
					if string(str[o]) == ":" {
						break
					}
					temp_cnt++
					v += string(str[o])
					o++
				}
				t, _ := strconv.Atoi(v)
				lim := i + len(v) + t + 1
				for ; i < lim; i++ {
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
				}
				i--
			}
		}
	}
	return m, nil
}

func Decode(str string, valid []bool) ([]interface{}, error) {
	return DecodeList(str, valid)
}

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
			i--
		} else {
			_, err := strconv.Atoi(string(v))
			if err == nil {
				temp_cnt := 0
				o := i + 1
				for {
					if string(str[o]) == ":" {
						break
					}
					temp_cnt++
					v += string(str[o])
					o++
				}
				t, _ := strconv.Atoi(v)
				lim := i + len(v) + t + 1
				for ; i < lim; i++ {
					(*valid)[i] = true
				}
				i--
			}
		}
	}
}
