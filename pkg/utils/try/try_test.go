package try

import (
	"io"
	"testing"

	"github.com/jopbrown/gobase/log"
	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	times := 0
	err := Do(func() error {
		times++
		log.Infof("try %d times", times)
		return io.EOF
	}, Option().SetLimitTimes(3))

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 3, times)

	times = 0
	err = Do(func() error {
		times++
		log.Infof("try %d times", times)
		if times == 2 {
			return nil
		}
		return io.EOF
	}, Option().SetLimitTimes(3))

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, times)

	times = 0
	err = Do(func() error {
		times++
		log.Infof("try %d times", times)
		if times == 2 {
			return io.ErrUnexpectedEOF
		}
		return io.EOF
	}, Option().SetLimitTimes(3).SetOnFail(func(err error) bool {
		return err == io.ErrUnexpectedEOF
	}))

	assert.Equal(t, io.ErrUnexpectedEOF, err)
	assert.Equal(t, 2, times)
}
