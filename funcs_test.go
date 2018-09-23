package funcmap

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/gomatic/clock"
	"github.com/stretchr/testify/assert"
)

//
func TestSubstr(t *testing.T) {
	type param struct {
		start, end int
		s, expect  string
	}
	tests := []param{
		{0, 0, "", ""},
		{0, 0, "0123456789abcdef", ""},
		{0, -1, "0123456789abcdef", "0123456789abcde"},
		{0, -2, "0123456789abcdef", "0123456789abcd"},
		{0, -3, "0123456789abcdef", "0123456789abc"},
		{0, -4, "0123456789abcdef", "0123456789ab"},
		{0, -5, "0123456789abcdef", "0123456789a"},
		{0, -6, "0123456789abcdef", "0123456789"},
		{0, -7, "0123456789abcdef", "012345678"},
		{0, -8, "0123456789abcdef", "01234567"},
		{0, -9, "0123456789abcdef", "0123456"},
		{0, -15, "0123456789abcdef", "0"},
		{0, -16, "0123456789abcdef", ""},
		{0, -17, "0123456789abcdef", "0123456789abcde"},
		{1, 0, "0123456789abcdef", "0"},
		{2, -1, "0123456789abcdef", "23456789abcde"},
		{3, -2, "0123456789abcdef", "3456789abcd"},
		{4, -3, "0123456789abcdef", "456789abc"},
		{5, -4, "0123456789abcdef", "56789ab"},
		{6, -5, "0123456789abcdef", "6789a"},
		{7, -6, "0123456789abcdef", "789"},
		{8, -7, "0123456789abcdef", "8"},
		{9, -8, "0123456789abcdef", "8"},
		{10, -9, "0123456789abcdef", "789"},
		{11, -15, "0123456789abcdef", "123456789a"},
		{12, -16, "0123456789abcdef", "0123456789ab"},
		{13, -17, "0123456789abcdef", "de"},
	}
	for _, p := range tests {
		if got := substr(p.start, p.end, p.s); got != p.expect {
			t.Errorf("for:%+v got:%v", p, got)
		}
	}
}

//
func TestIpMath(t *testing.T) {
	tests := map[string][][]string{
		"0.0.0.0": {
			{"0.0.0.0", "_._._._"},
			{"255.0.0.0", "[-1]._._._"},
			{"0.0.0.0", "[0]._._._"},
			{"1.0.0.0", "[+1]._._._"},
			{"0.255.0.0", "_.[-1]._._"},
			{"0.0.0.0", "_.[0]._._"},
			{"0.1.0.0", "_.[+1]._._"},
			{"0.0.255.0", "_._.[-1]._"},
			{"0.0.0.0", "_._.[0]._"},
			{"0.0.1.0", "_._.[+1]._"},
			{"0.0.0.255", "_._._.[-1]"},
			{"0.0.0.0", "_._._.[0]"},
			{"0.0.0.1", "_._._.[+1]"},
		},
		"255.255.255.255": {
			{"255.255.255.255", "_._._._"},
			{"254.255.255.255", "[-1]._._._"},
			{"0.255.255.255", "[0]._._._"},
			{"0.255.255.255", "[+1]._._._"},
			{"255.254.255.255", "_.[-1]._._"},
			{"255.0.255.255", "_.[0]._._"},
			{"255.0.255.255", "_.[+1]._._"},
			{"255.255.254.255", "_._.[-1]._"},
			{"255.255.0.255", "_._.[0]._"},
			{"255.255.0.255", "_._.[+1]._"},
			{"255.255.255.254", "_._._.[-1]"},
			{"255.255.255.0", "_._._.[0]"},
			{"255.255.255.0", "_._._.[+1]"},
			{"7.255.255.255", "[-1,/2,%10]._._._"},
			{"1.255.255.255", "[R]._._._"},
			{"5.255.255.255", "[+2,*5,%10]._._._"},
			{"255.7.255.255", "_.[-1,/2,%10]._._"},
			{"255.192.255.255", "_.[R]._._"},
			{"255.5.255.255", "_.[+2,*5,%10]._._"},
			{"255.255.7.255", "_._.[-1,/2,%10]._"},
			{"255.255.115.255", "_._.[R]._"},
			{"255.255.5.255", "_._.[+2,*5,%10]._"},
			{"255.255.255.7", "_._._.[-1,/2,%10]"},
			{"255.255.255.98", "_._._.[R]"},
			{"255.255.255.5", "_._._.[+2,*5,%10]"},
			{"3.255.255.255", "[+R,*R,%R]._._._"},
			{"255.11.255.255", "_.[+R,*R,%R]._._"},
			{"255.255.38.255", "_._.[+R,*R,%R]._"},
			{"255.255.255.40", "_._._.[+R,*R,%R]"},
		},
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff": {
			{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "_:_:_:_:_:_:_:_"},
			{"fffe:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "[-1]:_:_:_:_:_:_:_"},
			{"0000:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "[0]:_:_:_:_:_:_:_"},
			{"0000:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "[+1]:_:_:_:_:_:_:_"},
			{"ffff:fffe:ffff:ffff:ffff:ffff:ffff:ffff", "_:[-1]:_:_:_:_:_:_"},
			{"ffff:0000:ffff:ffff:ffff:ffff:ffff:ffff", "_:[0]:_:_:_:_:_:_"},
			{"ffff:0000:ffff:ffff:ffff:ffff:ffff:ffff", "_:[+1]:_:_:_:_:_:_"},
			{"ffff:ffff:fffe:ffff:ffff:ffff:ffff:ffff", "_:_:[-1]:_:_:_:_:_"},
			{"ffff:ffff:0000:ffff:ffff:ffff:ffff:ffff", "_:_:[0]:_:_:_:_:_"},
			{"ffff:ffff:0000:ffff:ffff:ffff:ffff:ffff", "_:_:[+1]:_:_:_:_:_"},
			{"ffff:ffff:ffff:fffe:ffff:ffff:ffff:ffff", "_:_:_:[-1]:_:_:_:_"},
			{"ffff:ffff:ffff:0000:ffff:ffff:ffff:ffff", "_:_:_:[0]:_:_:_:_"},
			{"ffff:ffff:ffff:0000:ffff:ffff:ffff:ffff", "_:_:_:[+1]:_:_:_:_"},
			{"ffff:ffff:ffff:ffff:fffe:ffff:ffff:ffff", "_:_:_:_:[-1]:_:_:_"},
			{"ffff:ffff:ffff:ffff:0000:ffff:ffff:ffff", "_:_:_:_:[0]:_:_:_"},
			{"ffff:ffff:ffff:ffff:0000:ffff:ffff:ffff", "_:_:_:_:[+1]:_:_:_"},
			{"ffff:ffff:ffff:ffff:ffff:fffe:ffff:ffff", "_:_:_:_:_:[-1]:_:_"},
			{"ffff:ffff:ffff:ffff:ffff:0000:ffff:ffff", "_:_:_:_:_:[0]:_:_"},
			{"ffff:ffff:ffff:ffff:ffff:0000:ffff:ffff", "_:_:_:_:_:[+1]:_:_"},
			{"ffff:ffff:ffff:ffff:ffff:ffff:fffe:ffff", "_:_:_:_:_:_:[-1]:_"},
			{"ffff:ffff:ffff:ffff:ffff:ffff:0000:ffff", "_:_:_:_:_:_:[0]:_"},
			{"ffff:ffff:ffff:ffff:ffff:ffff:0000:ffff", "_:_:_:_:_:_:[+1]:_"},
			{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffe", "_:_:_:_:_:_:_:[-1]"},
			{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:0000", "_:_:_:_:_:_:_:[0]"},
			{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:0000", "_:_:_:_:_:_:_:[+1]"},
		},
	}

	rand.Seed(0)
	for ip, tests := range tests {
		for _, test := range tests {
			e, m := test[0], test[1]
			v := ip_math(m, ip)
			if e != v {
				t.Errorf("expect:%v for:%v result:%v == %v", e, m, v, e == v)
			}
		}
	}
}

func Test_privateTime_Now(t *testing.T) {
	tests := []struct {
		name  string
		clock clock.Clock
		want  int
	}{
		{name: "format", clock: clock.Format, want: 2006},
		{name: "now", clock: clock.Default, want: time.Now().Year()},
		{name: "epoch", clock: clock.Epoch, want: 1970},
		{name: "now", clock: clock.Default, want: time.Now().Year()},
		{name: "format", clock: clock.Format, want: 2006},
		{name: "epoch", clock: clock.Epoch, want: 1970},
		{name: "now", clock: clock.Default, want: time.Now().Year()},
		{name: "epoch", clock: clock.Epoch, want: 1970},
		{name: "format", clock: clock.Format, want: 2006},
	}
	for i := 0; i < 1; i++ {
		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s %03d", tt.name, i), func(t *testing.T) {
				UseClock(tt.clock)
				now := privateTime.Now()
				assert.Equal(t, now.Year(), tt.want, fmt.Sprintf("%s : %s : %s", tt.name, string(tt.clock), now.String()))
			})
		}
	}
}
