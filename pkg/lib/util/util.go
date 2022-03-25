package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

func MapJsonToStruct(input interface{}, output interface{}) error {
	decoderConfig := &mapstructure.DecoderConfig{TagName: "json", Result: output}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return errors.New("mapstructure.NewDecoder(decoderConfig) err :" + err.Error())
	}

	derr := decoder.Decode(input)
	if derr != nil {
		return errors.New("decoder.Decode(input) err: " + derr.Error())
	}

	return nil
}

func RawToDecimal(raw string, maxDecimal int) string {
	if raw == "0" {
		return "0"
	}

	raw = strings.TrimLeft(raw, "0")

	decimalString := ""
	if len(raw) <= maxDecimal { // number is less than one
		decimalString = "0."
		for i := 0; i < (maxDecimal - len(raw)); i++ {
			decimalString = decimalString + "0"
		}
		decimalString = decimalString + raw
		decimalString = strings.TrimRight(decimalString, "0")

	} else { // number is greater than one
		numberPart := raw[0:(len(raw) - maxDecimal)]

		decimalPart := raw[(len(raw) - maxDecimal):]
		decimalPart = strings.TrimRight(decimalPart, "0")

		decimalString = numberPart
		if decimalPart != "" {
			decimalString = decimalString + "." + decimalPart
		}
	}

	_, ok := new(big.Float).SetString(decimalString) // check if number is valid
	if !ok {
		return "0"
	}

	return decimalString
}

func DecimalToRaw(decimal string, maxDecimal int) string {
	split := strings.Split(decimal, ".")
	decimal = strings.TrimLeft(decimal, "0")

	if len(split) > 1 {
		decimal = strings.TrimRight(decimal, "0")
		if len(split[1]) > maxDecimal { // reduce decimal count when it is greater than max decimal
			trimmedDecimalPart := split[1][:len(split[1])-(len(split[1])-maxDecimal)]
			decimal = split[0] + "." + trimmedDecimalPart

			return DecimalToRaw(decimal, maxDecimal)
		}
	}

	rawString := ""
	if string(decimal[0]) == "." { // number is less than one
		rawString = decimal[1:]
		for i := 0; i < ((maxDecimal + 1) - len(decimal)); i++ {
			rawString = rawString + "0"
		}
	} else { // number is greater than one
		decimalPart := ""
		if len(split) > 1 {
			decimalPart = strings.TrimRight(split[1], "0")
		}
		rawString = strings.ReplaceAll(decimal, ".", "")
		for i := 0; i < (maxDecimal - len(decimalPart)); i++ {
			rawString = rawString + "0"
		}
	}

	rawString = strings.TrimLeft(rawString, "0")
	_, ok := new(big.Float).SetString(rawString) // check if number is valid
	if !ok {
		return "0"
	}

	return rawString
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

func UniqueStrings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func Microtime() float64 {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	micSeconds := float64(now.Nanosecond()) / 1000000000
	return float64(now.Unix()) + micSeconds
}
