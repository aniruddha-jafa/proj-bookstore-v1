package security

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTokenGeneration(t *testing.T) {
	var tests = []struct {
		userId   uuid.UUID
		secret   string
		duration time.Duration
		now      time.Time
		//
		err  error
		want string
	}{
		{
			uuid.MustParse("1ccd635a-d484-4900-bc9a-e1881c54ff9b"),
			"abcde12345",
			time.Minute * 1,
			// Using a date in the future
			time.Date(2100, time.December, 15, 0, 0, 0, 0, time.UTC),
			nil,
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiIxY2NkNjM1YS1kNDg0LTQ5MDAtYmM5YS1lMTg4MWM1NGZmOWIiLCJleHAiOjQxMzI1MTIwNjAsImlhdCI6NDEzMjUxMjAwMH0.WcaslLIGa1TxF8MxI8zbOZ-E9yfyJt217GyGrbsny6o",
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s %s %v", tt.userId, tt.secret, tt.duration)
		t.Run(testname, func(t *testing.T) {
			got, err := MakeJwt(tt.userId, tt.secret, tt.duration, tt.now)
			if err != nil {
				t.Errorf("unexpected error for: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("got: %s, want: %s", got, tt.want)
			}
			userId, err := ValidateJwt(got, tt.secret)
			if err != nil {
				t.Errorf("unable to validate JWT %v", err)
			}
			if userId != tt.userId {
				t.Errorf("user id mismatch: got %v, want: %v", userId, tt.userId)
			}
		})
	}
}
