package config

type FindOption interface {
	apply(*option)
}

type option struct {
	query    []Query
	order    any
	offset   int
	limit    int
	preloads []string
}

type optionFn func(*option)

func (f optionFn) apply(opt *option) {
	f(opt)
}

func WithQuery(query ...Query) FindOption {
	return optionFn(func(opt *option) {
		opt.query = query
	})
}

func WithOffset(offset int) FindOption {
	return optionFn(func(opt *option) {
		opt.offset = offset
	})
}

func WithLimit(limit int) FindOption {
	return optionFn(func(opt *option) {
		opt.limit = limit
	})
}

func WithOrder(order interface{}) FindOption {
	return optionFn(func(opt *option) {
		opt.order = order
	})
}

func WithPreload(preloads []string) FindOption {
	return optionFn(func(opt *option) {
		opt.preloads = preloads
	})
}

func getOptions(opts ...FindOption) option {
	// default options
	opt := option{
		query:  []Query{},
		order:  "id",
		offset: 0,
		limit:  1000,
	}

	for _, o := range opts {
		o.apply(&opt)
	}

	return opt
}
