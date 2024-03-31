// This file contains the functions for encoding and decoding bencoded strings.
package processors

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/xmp-er/bpg/validators"
)

// Function to decode a bencoded string value
func DecodeString(str string) (string, error) {
	str = strings.TrimSpace(str)
	s := strings.Split(str, ":")
	if !validators.IsValidString(s) {
		return "", errors.New("string input sent is not valid, corrupted data")
	}
	return s[1], nil
}

// Function to decode an integer from a bencoded string
func DecodeInt(str string) (string, error) {
	if !validators.IsValidInt(str) {
		return "", errors.New("input for integer is not valid, corrupted data")
	}
	return str[1 : len(str)-1], nil
}

// Function to decode a list from a bencoded string
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

// Function to decode a dictionary from a bencoded string
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

// Function to decode a bencoded string
func Decode(str string, valid []bool) ([]interface{}, error) {
	return DecodeList(str, valid)
}

// Function to mark the valid characters in a bencoded string
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

// Function to encode a string in bencoded format
func Encode_string(str string) (string, error) {
	return fmt.Sprint(len(str)) + ":" + str, nil
}

// Function to encode an integer in bencoded format
func Encode_int(i int) (string, error) {
	tmp := fmt.Sprint(i)
	return "i" + tmp + "e", nil
}

// Function to encode a list in bencoded format
func Encode_List(list []interface{}) (string, error) {
	res := "l"
	for _, v := range list {
		_, is_string := v.(string)
		if is_string {
			temp_decoded_int_check, err := strconv.Atoi(v.(string))
			if err == nil {
				v = temp_decoded_int_check
			}
		}
		switch v.(type) {
		case string:
			temp, err := Encode_string(v.(string))
			if err != nil {
				return "", err
			}
			res += temp
		case int:
			temp, err := Encode_int(v.(int))
			if err != nil {
				return "", err
			}
			res += temp
		case []interface{}:
			temp, err := Encode_List(v.([]interface{}))
			if err != nil {
				return "", err
			}
			res += temp
		case map[string]interface{}:
			temp, err := Encode_Dictionary(v.(map[string]interface{}))
			if err != nil {
				return "", err
			}
			res += temp
		}
	}
	return res + "e", nil
}

// Function to encode a dictionary in bencoded format
func Encode_Dictionary(dict map[string]interface{}) (string, error) {
	res := "d"
	for k, v := range dict {
		_, is_string := v.(string)
		if is_string {
			temp_decoded_int_check, err := strconv.Atoi(v.(string))
			if err == nil {
				v = temp_decoded_int_check
			}
		}
		var encoded_value string = ""
		var err error
		switch v.(type) {
		case string:
			encoded_value, err = Encode_string(v.(string))
			if err != nil {
				return "", err
			}
		case int:
			encoded_value, err = Encode_int(v.(int))
			if err != nil {
				return "", err
			}
		case []interface{}:
			encoded_value, err = Encode_List(v.([]interface{}))
			if err != nil {
				return "", err
			}
		case map[string]interface{}:
			encoded_value, err = Encode_Dictionary(v.(map[string]interface{}))
			if err != nil {
				return "", err
			}
		}
		encoded_key, err := Encode_string(k)
		if err != nil {
			return "", err
		}
		res += (encoded_key + encoded_value)
	}
	return res + "e", nil
}
