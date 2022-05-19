package utils

import (
	"errors"
	"math/rand"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "test1"},
		{name: "test2"},
		{name: "test3"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			f := func() error {
				num := rand.Int()
				v, err := exec(int(num))
				t.Log(v, err)
				return err

			}
			opts := []Option{}
			opts = append(opts, WithRetryCount(2), WithDelayTime(3*time.Second))
			err := Do(f, opts...)
			if err != nil {
				t.Log(err)
			}
		})
	}
}

func exec(x int) (int, error) {
	if x%13 == 1 {
		return x, errors.New("exec error")
	}
	return x, nil
}
