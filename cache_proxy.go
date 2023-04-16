package cache_proxy_demo

type CacheProxy interface {
	Execute(qryOption any, readModelType any) (readModel any, err error)
}

type DatabaseGetFunc func(qryOption any) (readModel any, err error)

type TransformQryOptionToCacheKey func(qryOption any) (key string)

type BaseCacheProxy struct {
	Transform TransformQryOptionToCacheKey
	Cache     Cache
	GetDB     DatabaseGetFunc
}

func (proxy *BaseCacheProxy) Execute(qryOption any, readModelType any) (readModel any, err error) {
	key := proxy.Transform(qryOption)

	// cache.get
	val, err := proxy.Cache.GetValue(key, readModelType)
	if err == nil {
		return val, nil
	}

	// db.get
	readModel, err = proxy.GetDB(qryOption)
	if err != nil {
		return readModelType, err
	}

	// cache.set
	err = proxy.Cache.SetValue(key, readModel)
	if err != nil {
		return readModel, err
	}

	return readModel, nil
}
