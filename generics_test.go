package generics

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStack(t *testing.T) {
	s := NewStack[int]()

	s.Push(21)
	s.Push(22)
	s.Push(23)

	assert.Equal(t, 23, s.Pop())
	assert.Equal(t, 22, s.Pop())
	assert.Equal(t, 21, s.Pop())
}

func TestNull(t *testing.T) {
	n := NullMap(NullEmpty[int](), func(x int) string {
		return fmt.Sprintf("Hello: %d", x)
	})
	assert.Equal(t, NullEmpty[string](), n)

	n = NullMap(NullNew[int](20), func(x int) string {
		return fmt.Sprintf("Hello: %d", x)
	})
	assert.Equal(t, NullNew[string]("Hello: 20"), n)
}

// User  ...
type User struct {
	ID       int64        `json:"id"`
	Username Null[string] `json:"username"`
}

func TestNullJSON(t *testing.T) {
	data, err := json.Marshal(User{
		ID:       20,
		Username: NullNew("quang tung"),
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, `{"id":20,"username":"quang tung"}`, string(data))

	data, err = json.Marshal(User{
		ID: 20,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, `{"id":20,"username":null}`, string(data))

	data, err = json.Marshal(User{
		ID:       20,
		Username: NullEmpty[string](),
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, `{"id":20,"username":null}`, string(data))
}

func TestNullJSON_Unmarshal(t *testing.T) {
	var user User
	err := json.Unmarshal([]byte(`
{
  "id": 55,
  "username": "hello user"
}
`), &user)
	assert.Equal(t, nil, err)
	assert.Equal(t, User{
		ID:       55,
		Username: NullNew("hello user"),
	}, user)

	err = json.Unmarshal([]byte(`
{
  "id": 55,
  "username": null
}
`), &user)
	assert.Equal(t, nil, err)
	assert.Equal(t, User{
		ID:       55,
		Username: NullEmpty[string](),
	}, user)
}

func buildString(x int) string {
	return fmt.Sprintf("Hello %d", x)
}

func TestSliceMap(t *testing.T) {
	result := SliceMap([]int{2, 3, 4}, buildString)
	assert.Equal(t, []string{
		"Hello 2",
		"Hello 3",
		"Hello 4",
	}, result)
}

func TestGoMapMap(t *testing.T) {
	result := GoMapMap(map[int]int{11: 21, 12: 22}, buildString)
	assert.Equal(t, map[int]string{
		11: "Hello 21",
		12: "Hello 22",
	}, result)
}

func TestResult(t *testing.T) {
	r1 := ResultNew[int](23)
	r2 := ResultAndThen(r1, func(a int) Result[string] {
		return ResultNew(buildString(a))
	})
	assert.Equal(t, ResultNew("Hello 23"), r2)

	r1 = ResultErr[int](errors.New("some error"))
	r2 = ResultAndThen(r1, func(a int) Result[string] {
		return ResultNew(buildString(a))
	})
	assert.Equal(t, ResultErr[string](errors.New("some error")), r2)
}
