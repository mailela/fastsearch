package searcher

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"gitee.com/rachel_os/fastsearch/searcher/model"
	"gitee.com/rachel_os/fastsearch/searcher/storage"
	"gitee.com/rachel_os/fastsearch/searcher/words"
	// "github.com/Knetic/govaluate"
)

type Engine struct {
	IndexPath string  //索引文件存储目录
	Option    *Option //配置

	invertedIndexStorages []*storage.LeveldbStorage //关键字和Id映射，倒排索引,key=id,value=[]words
	positiveIndexStorages []*storage.LeveldbStorage //ID和key映射，用于计算相关度，一个id 对应多个key，正排索引
	docStorages           []*storage.LeveldbStorage //文档仓
	sync.Mutex                                      //锁
	sync.WaitGroup                                  //等待
	addDocumentWorkerChan []chan *model.IndexDoc    //添加索引的通道
	IsDebug               bool                      //是否调试模式
	AllowDrop             bool                      //允许删除
	Tokenizer             *words.Tokenizer          //分词器
	DatabaseName          string                    //数据库名

	Shard     int   //分片数
	Timeout   int64 //超时时间,单位秒
	BufferNum int   //分片缓冲数

	DocumentCount int64 //文档总数量
}

type Option struct {
	InvertedIndexName string //倒排索引
	PositiveIndexName string //正排索引
	DocIndexName      string //文档存储
}

// Init 初始化索引引擎
func (e *Engine) Init() {
	e.Add(1)
	defer e.Done()

	if e.Option == nil {
		e.Option = e.GetOptions()
	}
	if e.Timeout == 0 {
		e.Timeout = 10 * 3 // 默认30s
	}
	//-1代表没有初始化
	e.DocumentCount = -1
	//log.Println("数据存储目录：", e.IndexPath)
	log.Println("chain num:", e.Shard*e.BufferNum)
	e.addDocumentWorkerChan = make([]chan *model.IndexDoc, e.Shard)

	//初始化文件存储
	for shard := 0; shard < e.Shard; shard++ {

		//初始化chan
		worker := make(chan *model.IndexDoc, e.BufferNum)
		e.addDocumentWorkerChan[shard] = worker

		//初始化chan
		go e.DocumentWorkerExec(worker)

		s, err := storage.NewStorage(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.DocIndexName, shard)), e.Timeout)

		if err != nil {
			panic(err)
		}
		e.docStorages = append(e.docStorages, s)

		//初始化Keys存储
		ks, kerr := storage.NewStorage(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.InvertedIndexName, shard)), e.Timeout)
		if kerr != nil {
			panic(err)
		}
		e.invertedIndexStorages = append(e.invertedIndexStorages, ks)
		//id和keys映射
		iks, ikerr := storage.NewStorage(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.PositiveIndexName, shard)), e.Timeout)
		if ikerr != nil {
			panic(ikerr)
		}
		e.positiveIndexStorages = append(e.positiveIndexStorages, iks)
	}

	go e.automaticGC()
	//log.Println("初始化完成")
}

// 自动保存索引，10秒钟检测一次
func (e *Engine) automaticGC() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		<-ticker.C
		//定时GC
		runtime.GC()
		if e.IsDebug {
			log.Println("waiting:", e.GetQueue())
		}
	}
}
