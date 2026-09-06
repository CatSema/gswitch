package i18n

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"strconv"
	"strings"
)

// Reject duplicate object members instead of silently keeping the last translation.
func uniqueJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := jsonValue(decoder); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return errors.New("trailing JSON data")
	}
	return nil
}

func jsonValue(d *json.Decoder) error {
	token, err := d.Token()
	if err != nil {
		return err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	seen := map[string]bool{}
	for d.More() {
		if delimiter == '{' {
			token, err := d.Token()
			if err != nil {
				return err
			}
			key, ok := token.(string)
			if !ok || seen[key] {
				return fmt.Errorf("duplicate or invalid JSON key %v", token)
			}
			seen[key] = true
		}
		if err := jsonValue(d); err != nil {
			return err
		}
	}
	_, err = d.Token()
	return err
}

// UI messages only use scalar printf parameters. Widths, precision and stars
// are deliberately rejected so translators cannot introduce unsafe formatting.
func parameters(text string) (map[int]byte, error) {
	result := map[int]byte{}
	next := 1
	for i := 0; i < len(text); i++ {
		if text[i] != '%' {
			continue
		}
		i++
		if i < len(text) && text[i] == '%' {
			continue
		}
		if i < len(text) && text[i] == '[' {
			end := strings.IndexByte(text[i:], ']')
			if end < 0 {
				return nil, errors.New("unclosed parameter index")
			}
			n, err := strconv.Atoi(text[i+1 : i+end])
			if err != nil || n < 1 {
				return nil, errors.New("invalid parameter index")
			}
			next = n
			i += end + 1
		}
		if i >= len(text) || !strings.ContainsRune("sdv", rune(text[i])) {
			return nil, fmt.Errorf("unsupported format in %q", text)
		}
		if old, ok := result[next]; ok && old != text[i] {
			return nil, errors.New("conflicting parameter types")
		}
		result[next] = text[i]
		next++
	}
	return result, nil
}

func sameParameters(base, translation string) error {
	expected, err := parameters(base)
	if err != nil {
		return err
	}
	actual, err := parameters(translation)
	if err != nil {
		return err
	}
	if !maps.Equal(expected, actual) {
		return errors.New("translation parameters do not match English")
	}
	return nil
}
