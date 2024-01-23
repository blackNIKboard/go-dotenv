package godotenv

import (
	"io"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func isEqualEnv(t *testing.T, env1 *Env, env2 *Env) bool {
	if len(env1.keys) != len(env2.keys) {
		return false
	}

	if !slices.Equal(env1.keys, env2.keys) {
		return false
	}

	for _, k := range env1.keys {
		v := env1.data[k]

		debug := v.Comment == nil && (*env2).data[k].Comment == nil
		_ = debug

		if !(v.Comment == nil && (*env2).data[k].Comment == nil) &&
			*v.Comment != *(*env2).data[k].Comment {
			t.Logf("comment mismatch: '%v' vs '%v'\n", *v.Comment, *(*env2).data[k].Comment)

			return false
		}

		if v.Data != (*env2).data[k].Data {
			t.Logf("data mismatch: '%v' vs '%v'\n", v.Data, (*env2).data[k].Data)

			return false
		}

		if v.Quoted != (*env2).data[k].Quoted {
			t.Logf("quoted mismatch: %v vs %v\n", v.Quoted, (*env2).data[k].Quoted)

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
				data: map[string]EnvEntry{
					"TEST_BOOL": {
						Data:    "true",
						Comment: nil,
						Quoted:  false,
					},
				},
				keys: []string{"TEST_BOOL"},
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
				data: map[string]EnvEntry{"TEST_BOOL": {
					Data:    "true",
					Comment: nil,
					Quoted:  false,
				},
				},
				keys: []string{"TEST_BOOL"},
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
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "qwerty",
						Comment: nil,
						Quoted:  false,
					}},
				keys: []string{"TEST_STRING"},
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
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "qwerty",
						Comment: nil,
						Quoted:  true,
					},
				},
				keys: []string{"TEST_STRING"},
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
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "qwerty",
						Comment: nil,
						Quoted:  true,
					},
				},
				keys: []string{"TEST_STRING"},
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
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "qwerty   \"test\" ",
						Comment: nil,
						Quoted:  true,
					},
				},
				keys: []string{"TEST_STRING"},
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
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "qwerty   'test' ",
						Comment: nil,
						Quoted:  true,
					},
				},
				keys: []string{"TEST_STRING"},
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
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    `qwerty   \"test\" `,
						Comment: nil,
						Quoted:  true,
					},
				},
				keys: []string{"TEST_STRING"},
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
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    `qwerty   \'test\' `,
						Comment: nil,
						Quoted:  true,
					},
				},
				keys: []string{"TEST_STRING"},
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
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data: `qwerty   \'test\' `,
						Comment: func() *string {
							comment := " 1342 commented "
							return &comment
						}(),
						Quoted: true,
					},
				},
				keys: []string{"TEST_STRING"},
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
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data: `qwerty   \'test\' `,
						Comment: func() *string {
							comment := " 1342 comm####ented # comment 2   "
							return &comment
						}(),
						Quoted: true,
					},
				},
				keys: []string{"TEST_STRING"},
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
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    `/data/test/`,
						Comment: nil,
						Quoted:  true,
					},
				},
				keys: []string{"TEST_STRING"},
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
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    `/data/test/`,
						Comment: nil,
						Quoted:  false,
					},
				},
				keys: []string{"TEST_STRING"},
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
				data: map[string]EnvEntry{
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
				keys: []string{"TEST_STRING1", "TEST_STRING"},
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
				data: map[string]EnvEntry{
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
				keys: []string{"TEST_STRING", "TEST_STRING1"},
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

func TestEnv_Add(t *testing.T) {
	type args struct {
		key     string
		value   string
		quoted  bool
		comment *string
	}
	tests := []struct {
		name    string
		env     *Env
		args    args
		want    *Env
		wantErr bool
	}{
		{
			name: "Add string to empty Env",
			env:  &Env{},
			args: args{
				key:     "TEST_STRING",
				value:   "test",
				quoted:  false,
				comment: nil,
			},
			want: &Env{
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "test",
						Quoted:  false,
						Comment: nil,
					},
				},
				keys: []string{"TEST_STRING"},
			},
			wantErr: false,
		},
		{
			name: "Modify existing string in Env",
			env: &Env{
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "test",
						Quoted:  false,
						Comment: nil,
					},
				},
				keys: []string{"TEST_STRING"},
			},
			args: args{
				key:     "TEST_STRING",
				value:   "test1",
				quoted:  false,
				comment: nil,
			},
			want: &Env{
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "test1",
						Quoted:  false,
						Comment: nil,
					},
				},
				keys: []string{"TEST_STRING"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.env
			if err := e.Add(tt.args.key, tt.args.value, tt.args.quoted, tt.args.comment); (err != nil) != tt.wantErr {
				t.Errorf("Env.Add() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !isEqualEnv(t, e, tt.want) {
				t.Errorf("Environments are not equal: want %v; got %v", *tt.want, *e)
			}
		})
	}
}

func TestEnv_Get(t *testing.T) {
	tests := []struct {
		name    string
		env     *Env
		request string
		want    string
		wantErr bool
	}{
		{
			name:    "get from empty",
			env:     &Env{},
			request: "TEST_STRING",
			want:    "",
			wantErr: true,
		},
		{
			name: "get simple sting",
			env: &Env{
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    "test",
						Quoted:  false,
						Comment: nil,
					},
				},
				keys: []string{"TEST_STRING"},
			},
			request: "TEST_STRING",
			want:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.env
			got, err := e.Get(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Env.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Delete(t *testing.T) {
	tests := []struct {
		name    string
		env     *Env
		request string
		want    *Env
		wantErr bool
	}{
		{
			name:    "delete nonexistant",
			env:     &Env{},
			request: "TEST_STRING",
			want:    &Env{},
			wantErr: true,
		},
		{
			name: "delete existing, first of two",
			env: &Env{
				data: map[string]EnvEntry{
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
				keys: []string{"TEST_STRING", "TEST_STRING1"},
			},
			request: "TEST_STRING",
			want: &Env{
				data: map[string]EnvEntry{
					"TEST_STRING1": {
						Data:    `test`,
						Comment: nil,
						Quoted:  false,
					},
				},
				keys: []string{"TEST_STRING1"},
			},
			wantErr: false,
		},
		{
			name: "delete existing, second of two",
			env: &Env{
				data: map[string]EnvEntry{
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
				keys: []string{"TEST_STRING", "TEST_STRING1"},
			},
			request: "TEST_STRING1",
			want: &Env{
				data: map[string]EnvEntry{
					"TEST_STRING": {
						Data:    `/data/test/`,
						Comment: nil,
						Quoted:  false,
					},
				},
				keys: []string{"TEST_STRING"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.env
			if err := e.Delete(tt.request); (err != nil) != tt.wantErr {
				t.Errorf("Env.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !isEqualEnv(t, e, tt.want) {
				t.Errorf("Environments are not equal: want %v; got %v", *tt.want, *e)
			}
		})
	}
}

func TestEnv_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		env     *Env
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "Get from empty env",
			env:     &Env{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Get all 1",
			env: &Env{
				data: map[string]EnvEntry{
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
				keys: []string{"TEST_STRING", "TEST_STRING1"},
			},
			want: map[string]string{
				"TEST_STRING":  "/data/test/",
				"TEST_STRING1": "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.env
			got, err := e.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
