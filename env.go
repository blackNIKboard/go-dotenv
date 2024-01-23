package godotenv

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	ErrNotFound        = fmt.Errorf("key not found in env")
	ErrKeyNotMatched   = fmt.Errorf("key do not match the regexp")
	ErrValueNotMatched = fmt.Errorf("value do not match the regexp")
)

var regexKey = `[A-Z0-9,_]+`
var regexValue = `(("(?:[^"\\]|\\.)*?"|([^'"])+)|('(?:[^'\\]|\\.)*?'|([^'"])+))`

func (e *Env) Write(writer io.Writer) error {
	for _, key := range e.keys {
		entry := e.data[key]

		if _, err := writer.Write(
			[]byte(fmt.Sprintf("%s=%s\n",
				key,
				func() string {
					result := ""

					if entry.Quoted {
						result += "\"" + entry.Data + "\""
					} else {
						result += entry.Data
					}

					if entry.Comment != nil {
						result += "    #" + *entry.Comment
					}

					return result
				}(),
			))); err != nil {
			return err
		}
	}

	return nil
}

func (e *Env) Read(reader io.Reader) error {
	var (
		re = regexp.MustCompile(
			regexKey + "=" + regexValue,
		)
	)

	e.data = map[string]EnvEntry{}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var (
			entry EnvEntry
		)

		entry.Comment = func() *string {
			split := strings.SplitN(scanner.Text(), "#", 2)

			if len(split) > 1 {
				return &split[1]
			}

			return nil
		}()

		if !re.MatchString(scanner.Text()) {
			continue
		}

		rawEntry := strings.Split(re.FindString(scanner.Text()), "=")

		fmt.Println(rawEntry)

		entry.Data = rawEntry[1]

		entry.Quoted = func() bool {
			length := len(entry.Data)

			if length < 2 {
				return false
			}

			if !(entry.Data[0] == entry.Data[length-1]) {
				return false
			}

			if !(string(entry.Data[0]) == "'" || string(entry.Data[0]) == "\"") {
				return false
			}

			entry.Data = entry.Data[1 : length-1]

			return true
		}()

		(*e).data[rawEntry[0]] = entry
		(*e).keys = append((*e).keys, rawEntry[0])
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (e *Env) Add(key, value string, quoted bool, comment *string) error {
	reKey := regexp.MustCompile(regexKey)
	reValue := regexp.MustCompile(regexValue)

	if !reKey.MatchString(key) {
		return ErrKeyNotMatched
	}

	if !reValue.MatchString(value) {
		return ErrValueNotMatched
	}

	if e.data == nil {
		e.data = map[string]EnvEntry{}
	}

	e.data[key] = EnvEntry{
		Data:    value,
		Quoted:  quoted,
		Comment: comment,
	}

	e.keys = append(e.keys, key)

	return nil
}

func (e *Env) Get(key string) (string, error) {
	if e.data == nil {
		return "", ErrNotFound
	}

	value, ok := e.data[key]
	if !ok {
		return "", ErrNotFound
	}

	return value.Data, nil
}

func (e *Env) Delete(key string) error {
	_, err := e.Get(key)
	if err != nil {
		return err
	}

	delete(e.data, key)

	for i, v := range e.keys {
		if v == key {
			e.keys = append(e.keys[:i], e.keys[i+1:]...)
		}
	}

	return nil
}
