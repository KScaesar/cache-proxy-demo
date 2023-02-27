package example

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func NewCacheProxyChannel(
	cache Cache,
	transformReadOption func(readDtoOption any) (key string),
	readDataSource func(ctx context.Context, readDtoOption any) (readModel any, err error),
	isAnNotFoundError func(err error) bool,
	canIgnoreCacheError bool,
	canIgnoreReadSourceErrorNotFound bool,
	cacheTTL time.Duration,
) *CacheProxyChannel {
	proxy := &CacheProxyChannel{
		Cache: cache,

		startDelivery:  make(chan CommandProxyGet),
		doneDelivery:   make(chan string, 1),
		singleDelivery: make(map[string]*Entry),
		end:            make(chan struct{}),

		TransformReadOption: transformReadOption,
		ReadDataSource:      readDataSource,

		IsAnNotFoundError:                isAnNotFoundError,
		CanIgnoreCacheError:              canIgnoreCacheError,
		CanIgnoreReadSourceErrorNotFound: canIgnoreReadSourceErrorNotFound,
		CacheTTL:                         cacheTTL,
		debug:                            false,
	}

	go proxy.manager()
	return proxy
}

type CacheProxyChannel struct {
	Cache Cache

	startDelivery  chan CommandProxyGet
	doneDelivery   chan string
	singleDelivery map[string]*Entry
	end            chan struct{}
	once           sync.Once

	TransformReadOption func(readDtoOption any) (key string)
	ReadDataSource      func(ctx context.Context, readDtoOption any) (readModel any, err error)

	IsAnNotFoundError                func(err error) bool
	CanIgnoreCacheError              bool
	CanIgnoreReadSourceErrorNotFound bool // source not found, 是否交給 caller 處理
	CacheTTL                         time.Duration
	debug                            bool
}

func (proxy *CacheProxyChannel) ReadValue(ctx context.Context, readDtoOption any) (readModel any, err error) {
	return proxy.ReadValueV2(ctx, readDtoOption)
}

func (proxy *CacheProxyChannel) ReadValueV2(ctx context.Context, readDtoOption any) (readModel any, err error) {
	key := proxy.TransformReadOption(readDtoOption)
	readModel, err = proxy.Cache.GetValue(ctx, key)
	if err == nil {
		return readModel, nil
	}

	if proxy.IsAnNotFoundError(err) || proxy.CanIgnoreCacheError {
		cmd := CommandProxyGet{
			ctx:           ctx,
			readDtoOption: readDtoOption,
			reply:         make(chan ProxyResult, 1),
		}
		proxy.startDelivery <- cmd
		r := <-cmd.reply
		return r.val, r.err
	}

	return nil, err
}

func (proxy *CacheProxyChannel) manager() {
	for {
		select {
		case <-proxy.end:
			return
		default:
		}
		// fmt.Println("ping")

		select {
		case key := <-proxy.doneDelivery:
			if proxy.debug {
				fmt.Println("done", key)
			}

			delete(proxy.singleDelivery, key)
			// proxy.singleDelivery[key] = &Entry{ready: make(chan struct{})} // dead lock, because entry != nil, not enter main read

		case cmd := <-proxy.startDelivery:
			if proxy.debug {
				fmt.Println("cmd", cmd.readDtoOption)
			}

			key := proxy.TransformReadOption(cmd.readDtoOption)
			entry := proxy.singleDelivery[key]
			if entry == nil {
				entry = &Entry{ready: make(chan struct{})}
				proxy.singleDelivery[key] = entry
				go proxy.slowMainReader(cmd.ctx, cmd.readDtoOption, entry)
				// break // dead lock, because main read have not reply
			}
			go proxy.replier(entry, cmd)
		default:
		}
	}
}

func (proxy *CacheProxyChannel) replier(
	entry *Entry,
	cmd CommandProxyGet,
) {
	<-entry.ready
	if proxy.debug {
		fmt.Println("ready", cmd.readDtoOption)
	}
	cmd.reply <- entry.result
}

// bug: double main read
//
// cmd id[7]
// main read id[7] <-
// cmd id[7]
// cmd id[7]
// cmd id[7]
// cmd id[7]
// ready id[7]
// done id[7]
// cmd id[7]
// cmd id[7]
// cmd id[7]
// ready id[7]
// ready id[7]
// main read id[7] <-
// ready id[7]
// ready id[7]
// ready id[7]
// ready id[7]
// ready id[7]
// done id[7]

func (proxy *CacheProxyChannel) slowMainReader(ctx context.Context, readDtoOption any, entry *Entry) {
	if proxy.debug {
		fmt.Println("main read", readDtoOption)
	}
	var key string = proxy.TransformReadOption(readDtoOption)
	var empty any
	var err error

	defer func() {
		close(entry.ready)

		// 要注意順序
		// close(chan) 一定要先執行, 通知所有 replier
		// 再 send done

		// workaround bug: double main read
		//
		// 原本以為是 bug
		// 但看到下面的文章, 也可以想成, 用來控制 幾秒內 允許第二次 main read
		// bug 變成 feature xd
		// https://www.cyningsun.com/01-11-2021/golang-concurrency-singleflight.html
		time.Sleep(time.Second)
		proxy.doneDelivery <- key
	}()

	readModel, err := proxy.ReadDataSource(ctx, readDtoOption)
	if err != nil {
		if !(proxy.IsAnNotFoundError(err) && proxy.CanIgnoreReadSourceErrorNotFound) {
			entry.result = ProxyResult{val: empty, err: err}
			return
		}
	}

	err = proxy.Cache.PutValue(ctx, key, readModel, proxy.CacheTTL)
	if err != nil && !proxy.CanIgnoreCacheError {
		entry.result = ProxyResult{val: empty, err: err}
		return
	}

	entry.result = ProxyResult{val: readModel, err: err}
}

func (proxy *CacheProxyChannel) Close() {
	proxy.once.Do(func() {
		close(proxy.end)
	})
}

// Entry 達到類似廣播資訊的功能,
// 為了讓相同的請求, 共用同樣的結果
type Entry struct {
	ready  chan struct{}
	result ProxyResult
}

type ProxyResult struct {
	val any
	err error
}

type CommandProxyGet struct {
	ctx           context.Context
	readDtoOption any
	reply         chan ProxyResult
}
