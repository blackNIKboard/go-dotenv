package godotenv

import (
	"io"
	"slices"
	"strings"
	"testing"
)

func isEqualEnv(t *testing.T, env1 *Env, env2 *Env) bool {
	if len(env1.Keys) != len(env2.Keys) {
		return false
	}

	if !slices.Equal(env1.Keys, env2.Keys) {
		return false
	}

	for _, k := range env1.Keys {
		v := env1.Data[k]

		debug := v.Comment == nil && (*env2).Data[k].Comment == nil
		_ = debug

		if !(v.Comment == nil && (*env2).Data[k].Comment == nil) &&
			*v.Comment != *(*env2).Data[k].Comment {
			t.Logf("comment mismatch: '%v' vs '%v'\n", *v.Comment, *(*env2).Data[k].Comment)

			return false
		}

		if v.Data != (*env2).Data[k].Data {
			t.Logf("data mismatch: '%v' vs '%v'\n", v.Data, (*env2).Data[k].Data)

			return false
		}

		if v.Quoted != (*env2).Data[k].Quoted {
			t.Logf("quoted mismatch: %v vs %v\n", v.Quoted, (*env2).Data[k].Quoted)

			return false
		}
	}

	return true
}

func TestEnv_Read(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		e       *Env
		args    args
		wantErr bool
		want    *Env
	}{
		{
			"Single bool true",
			&Env{},
			args{
				strings.NewReader(`TEST_BOOL=true`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_BOOL": {
						Data:    "true",
						Comment: nil,
						Quoted:  false,
					},
				},
				Keys: []string{"TEST_BOOL"},
			},
		},
		{
			"Single bool true",
			&Env{},
			args{
				strings.NewReader(`TEST_BOOL=true`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{"TEST_BOOL": {
					Data:    "true",
					Comment: nil,
					Quoted:  false,
				},
				},
				Keys: []string{"TEST_BOOL"},
			},
		},
		{
			"Single string unquoted",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING=qwerty`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "qwerty",
						Comment: nil,
						Quoted:  false,
					}},
				Keys: []string{"TEST_STRING"},
			},
		},
		{
			"Single string simple quoted 1",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING="qwerty"`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "qwerty",
						Comment: nil,
						Quoted:  true,
					},
				},
				Keys: []string{"TEST_STRING"},
			},
		},
		{
			"Single string simple quoted 2",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING='qwerty'`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "qwerty",
						Comment: nil,
						Quoted:  true,
					},
				},
				Keys: []string{"TEST_STRING"},
			},
		},
		{
			"Single string complex quoted 1",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING='qwerty   "test" '`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "qwerty   \"test\" ",
						Comment: nil,
						Quoted:  true,
					},
				},
				Keys: []string{"TEST_STRING"},
			},
		},
		{
			"Single string complex quoted 2",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING="qwerty   'test' "`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "qwerty   'test' ",
						Comment: nil,
						Quoted:  true,
					},
				},
				Keys: []string{"TEST_STRING"},
			},
		},
		{
			"Single string complex quoted 3",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING="qwerty   \"test\" "`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    `qwerty   \"test\" `,
						Comment: nil,
						Quoted:  true,
					},
				},
				Keys: []string{"TEST_STRING"},
			},
		},
		{
			"Single string complex quoted 4",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING='qwerty   \'test\' '`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    `qwerty   \'test\' `,
						Comment: nil,
						Quoted:  true,
					},
				},
				Keys: []string{"TEST_STRING"},
			},
		},
		{
			"Single string complex quoted 5 + simple comment",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING='qwerty   \'test\' ' # 1342 commented `),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING": {
						Data: `qwerty   \'test\' `,
						Comment: func() *string {
							comment := " 1342 commented "
							return &comment
						}(),
						Quoted: true,
					},
				},
				Keys: []string{"TEST_STRING"},
			},
		},
		{
			"Single string complex quoted 6 + complex comment",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING='qwerty   \'test\' '# 1342 comm####ented # comment 2   `),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING": {
						Data: `qwerty   \'test\' `,
						Comment: func() *string {
							comment := " 1342 comm####ented # comment 2   "
							return &comment
						}(),
						Quoted: true,
					},
				},
				Keys: []string{"TEST_STRING"},
			},
		},
		{
			"Single string with path data quoted",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING='/data/test/'`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    `/data/test/`,
						Comment: nil,
						Quoted:  true,
					},
				},
				Keys: []string{"TEST_STRING"},
			},
		},
		{
			"Single string with path data unquoted",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING=/data/test/`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    `/data/test/`,
						Comment: nil,
						Quoted:  false,
					},
				},
				Keys: []string{"TEST_STRING"},
			},
		},
		{
			"Preserve order 1",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING1=test
TEST_STRING=/data/test/`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING1": {
						Data:    `test`,
						Comment: nil,
						Quoted:  false,
					},
					"TEST_STRING": {
						Data:    `/data/test/`,
						Comment: nil,
						Quoted:  false,
					},
				},
				Keys: []string{"TEST_STRING1", "TEST_STRING"},
			},
		},
		{
			"Preserve order 2",
			&Env{},
			args{
				strings.NewReader(`TEST_STRING=/data/test/
TEST_STRING1=test`),
			},
			false,
			&Env{
				Data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    `/data/test/`,
						Comment: nil,
						Quoted:  false,
					},
					"TEST_STRING1": {
						Data:    `test`,
						Comment: nil,
						Quoted:  false,
					},
				},
				Keys: []string{"TEST_STRING", "TEST_STRING1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Read(tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("Env.Read() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Log("comparing data")
				if !isEqualEnv(t, tt.e, tt.want) {
					t.Errorf("Env.Read() result unexpected, want '%v' vs got '%v'", *tt.e, *tt.want)
				}
			}
		})
	}
}

//  1342 comm####ented # comment 2   ]
//  1342 comm####ented # comment 2   '
