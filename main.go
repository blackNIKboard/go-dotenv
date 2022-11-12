package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gijsbers/go-pcre"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	fileIn, err := os.OpenFile("text.env", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer fileIn.Close()

	fileOut, err := os.OpenFile("text.out.env", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer fileIn.Close()

	env, err := Read(fileIn)
	if err != nil {
		panic(err)
	}

	err = env.Write(fileOut)
	if err != nil {
		panic(err)
	}

	spew.Dump(env)
}

type Env map[string]EnvEntry

type EnvEntry struct {
	Data    string
	Comment *string
	Quoted  bool
}

func (e *Env) Write(writer io.Writer) error {
	for key, entry := range *e {
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

func Read(reader io.Reader) (Env, error) {
	var (
		result = Env{}
		re     = pcre.MustCompile(`[A-Z,_]+=((["'])(?:[^\2\\]|\\.)*?\2|\w+)`, 0)
	)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var (
			entry   EnvEntry
			matcher = re.MatcherString(scanner.Text(), 0)
		)

		entry.Comment = func() *string {
			split := strings.SplitN(scanner.Text(), "#", 2)

			if len(split) > 1 {
				return &split[1]
			}

			return nil
		}()

		if !matcher.Matches() {
			continue
		}

		rawEntry := strings.Split(matcher.GroupString(0), "=")

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

		result[rawEntry[0]] = entry
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
