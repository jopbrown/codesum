package try

import "time"

type RetryOption struct {
	limitTimes int
	timeout    time.Duration
	onFail     func(error) bool
	interval   time.Duration
}

// type RetryOption func(*RetryOptions)

// func WithLimitTimes(limitTimes int) RetryOption {
// 	return func(opts *RetryOptions) {
// 		opts.limitTimes = limitTimes
// 	}
// }

// func WithInterval(interval time.Duration) RetryOption {
// 	return func(opts *RetryOptions) {
// 		opts.interval = interval
// 	}
// }

// func WithTimeout(timeout time.Duration) RetryOption {
// 	return func(opts *RetryOptions) {
// 		opts.timeout = timeout
// 	}
// }

// func OnFail(callback func(error) bool) RetryOption {
// 	return func(opts *RetryOptions) {
// 		opts.onFail = callback
// 	}
// }

func Option() *RetryOption {
	return &RetryOption{}
}

func (opt *RetryOption) SetLimitTimes(limitTimes int) *RetryOption {
	opt.limitTimes = limitTimes
	return opt
}

func (opt *RetryOption) SetInterval(interval time.Duration) *RetryOption {
	opt.interval = interval
	return opt
}

func (opt *RetryOption) SetTimeout(timeout time.Duration) *RetryOption {
	opt.timeout = timeout
	return opt
}

func (opt *RetryOption) SetOnFail(onFail func(error) bool) *RetryOption {
	opt.onFail = onFail
	return opt
}

func Do(fn func() error, opts ...*RetryOption) (err error) {
	opt := &RetryOption{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	times := 0
	start := time.Now()
	for {
		err = fn()
		if err == nil {
			return nil
		}
		if opt.onFail != nil && opt.onFail(err) {
			return err
		}

		times++

		if opt.limitTimes > 0 && times >= opt.limitTimes {
			break
		}

		if opt.timeout > 0 && time.Since(start) > opt.timeout {
			break
		}

		if opt.interval > 0 {
			time.Sleep(opt.interval)
		}
	}

	return err
}
