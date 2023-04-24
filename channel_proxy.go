package cache_proxy_demo

func UseChannel(baseProxy *BaseCacheProxy) CacheProxy {
	proxy := &ChannelProxy{
		transform: baseProxy.Transform,
		cache:     baseProxy.Cache,

		mqCommandGet:          make(chan CommandGetReadModel),
		mqEventGot:            make(chan EventGotReadModel),
		firstGoroutineRecords: make(map[string]*Entry),
		closeManager:          make(chan struct{}),

		baseProxy: baseProxy,
	}

	go proxy.manager()
	return proxy
}

type ChannelProxy struct {
	transform TransformQryOptionToCacheKey
	cache     Cache

	mqCommandGet          chan CommandGetReadModel
	mqEventGot            chan EventGotReadModel
	firstGoroutineRecords map[string]*Entry
	closeManager          chan struct{}

	baseProxy *BaseCacheProxy
}

func (proxy *ChannelProxy) Execute(qryOption any, readModelType any) (readModel any, err error) {
	return proxy.execute1(qryOption, readModelType)
}

func (proxy *ChannelProxy) execute1(qryOption any, readModelType any) (readModel any, err error) {
	command := CommandGetReadModel{
		qryOption:     qryOption,
		readModelType: readModelType,
		replyAddress:  make(chan Result, 1),
	}
	proxy.mqCommandGet <- command

	result := <-command.replyAddress
	return result.val, result.err
}

func (proxy *ChannelProxy) execute3(qryOption any, readModelType any) (readModel any, err error) {
	key := proxy.transform(qryOption)
	val, err := proxy.cache.GetValue(key, readModelType)
	if err == nil {
		return val, nil
	}

	return proxy.execute1(qryOption, readModelType)
}

func (proxy *ChannelProxy) manager() {
	for {
		select {
		case <-proxy.closeManager:
			return
		default:
		}

		select {
		case cmd := <-proxy.mqCommandGet:
			key := proxy.transform(cmd.qryOption)

			entry := proxy.firstGoroutineRecords[key]
			if entry == nil {
				entry = &Entry{ready: make(chan struct{})}
				proxy.firstGoroutineRecords[key] = entry
				go proxy.mainReader(cmd, key, entry)
			}
			go proxy.replier(cmd.replyAddress, entry)

		case key := <-proxy.mqEventGot:
			proxy.firstGoroutineRecords[key] = nil

		case <-proxy.closeManager:
			return

		// default: // default 移除註解, 可以讓效能更好
		}
	}
}

func (proxy *ChannelProxy) mainReader(cmd CommandGetReadModel, key string, entry *Entry) {
	defer func() {
		close(entry.ready)
		// time.Sleep(time.Second)
		proxy.mqEventGot <- key
	}()
	readModel, err := proxy.baseProxy.Execute(cmd.qryOption, cmd.readModelType)
	entry.result = Result{val: readModel, err: err}
}

func (proxy *ChannelProxy) replier(replyAddress chan Result, entry *Entry) {
	<-entry.ready
	replyAddress <- entry.result
}

type Entry struct {
	ready  chan struct{}
	result Result
}

type Result struct {
	val any
	err error
}

type CommandGetReadModel struct {
	qryOption     any
	readModelType any
	replyAddress  chan Result
}

type EventGotReadModel = string
