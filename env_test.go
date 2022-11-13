package godotenv

import (
	"io"
	"strings"
	"testing"
)

func isEqualEnv(t *testing.T, env1 *Env, env2 *Env) bool {
	if len(*env1) != len(*env2) {
		return false
	}

	for k, v := range *env1 {
		debug := v.Comment == nil && (*env2)[k].Comment == nil
		_ = debug

		if !(v.Comment == nil && (*env2)[k].Comment == nil) &&
			*v.Comment != *(*env2)[k].Comment {
			t.Logf("comment mismatch: '%v' vs '%v'\n", *v.Comment, *(*env2)[k].Comment)

			return false
		}

		if v.Data != (*env2)[k].Data {
			t.Logf("data mismatch: '%v' vs '%v'\n", v.Data, (*env2)[k].Data)

			return false
		}

		if v.Quoted != (*env2)[k].Quoted {
			t.Logf("quoted mismatch: %v vs %v\n", v.Quoted, (*env2)[k].Quoted)

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
				"TEST_BOOL": {
					Data:    "true",
					Comment: nil,
					Quoted:  false,
				},
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
				"TEST_BOOL": {
					Data:    "true",
					Comment: nil,
					Quoted:  false,
				},
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
				"TEST_STRING": {
					Data:    "qwerty",
					Comment: nil,
					Quoted:  false,
				},
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
				"TEST_STRING": {
					Data:    "qwerty",
					Comment: nil,
					Quoted:  true,
				},
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
				"TEST_STRING": {
					Data:    "qwerty",
					Comment: nil,
					Quoted:  true,
				},
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
				"TEST_STRING": {
					Data:    "qwerty   \"test\" ",
					Comment: nil,
					Quoted:  true,
				},
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
				"TEST_STRING": {
					Data:    "qwerty   'test' ",
					Comment: nil,
					Quoted:  true,
				},
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
				"TEST_STRING": {
					Data:    `qwerty   \"test\" `,
					Comment: nil,
					Quoted:  true,
				},
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
				"TEST_STRING": {
					Data:    `qwerty   \'test\' `,
					Comment: nil,
					Quoted:  true,
				},
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
				"TEST_STRING": {
					Data: `qwerty   \'test\' `,
					Comment: func() *string {
						comment := " 1342 commented "
						return &comment
					}(),
					Quoted: true,
				},
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
				"TEST_STRING": {
					Data: `qwerty   \'test\' `,
					Comment: func() *string {
						comment := " 1342 comm####ented # comment 2   "
						return &comment
					}(),
					Quoted: true,
				},
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
