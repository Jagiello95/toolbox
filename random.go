package toolbox

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

type Random struct{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func (r *Random) RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func (r *Random) RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func (r *Random) RandomOwner() string {
	return r.RandomString(6)
}

// RandomMoney generates a random amount of money
func (r *Random) RandomMoney() int64 {
	return r.RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func (r *Random) RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomEmail generates a random email
func (r *Random) RandomEmail() string {
	return fmt.Sprintf("%s@email.com", r.RandomString(6))
}
