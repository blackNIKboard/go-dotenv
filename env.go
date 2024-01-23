package godotenv

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func (e *Env) Write(writer io.Writer) error {
	for _, key := range e.Keys {
		entry := e.Data[key]

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
			`[A-Z0-9,_]+=(("(?:[^"\\]|\\.)*?"|([^'"])+)|('(?:[^'\\]|\\.)*?'|([^'"])+))`,
		)
	)

	e.Data = map[string]EnvEntry{}

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

		(*e).Data[rawEntry[0]] = entry
		(*e).Keys = append((*e).Keys, rawEntry[0])
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
