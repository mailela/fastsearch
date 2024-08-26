package sorts

import (
	"fmt"
	"sort"
	"sync"

	"gitee.com/rachel_os/fastsearch/searcher/model"
)

const (
	DESC = "desc"
)

type ScoreSlice []model.SliceItem

func (x ScoreSlice) Len() int {
	return len(x)
}
func (x ScoreSlice) Less(i, j int) bool {
	return x[i].Score < x[j].Score
}
func (x ScoreSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

type SortSlice []string

func (x SortSlice) Len() int {
	return len(x)
}
func (x SortSlice) Less(i, j int) bool {
	return x[i] < x[j]
}
func (x SortSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]

}

type Uint32Slice []uint32

func (x Uint32Slice) Len() int           { return len(x) }
func (x Uint32Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Uint32Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type FastSort struct {
	sync.Mutex

	IsDebug bool

	data []model.SliceItem

	temps []string

	count int //总数

	Order string //排序方式
}

func (f *FastSort) Add(ids *[]string) {
	//f.Lock()
	//defer f.Unlock()

	//for _, id := range *ids {
	//
	//	found, index := f.find(&id)
	//	if found {
	//		f.data[index].Score += 1
	//	} else {
	//
	//		f.data = append(f.data, model.SliceItem{
	//			Id:    id,
	//			Score: 1,
	//		})
	//		f.Sort()
	//	}
	//}
	//f.count = len(f.data)
	f.temps = append(f.temps, *ids...)
}

// 二分法查找
func (f *FastSort) find(target *string) (bool, int) {

	low := 0
	high := f.count - 1
	for low <= high {
		mid := (low + high) / 2
		if f.data[mid].Id == *target {
			return true, mid
		} else if f.data[mid].Id > *target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return false, -1
	//for index, item := range f.data {
	//	if item.Id == *target {
	//		return true, index
	//	}
	//}
	//return false, -1
}

// Count 获取数量
func (f *FastSort) Count() int {
	return f.count
}

// Sort 排序
func (f *FastSort) Sort() {
	sort.Sort(SortSlice(f.temps))
	// if strings.ToLower(f.Order) == DESC {
	// 	sort.Sort(sort.Reverse(SortSlice(f.temps)))
	// } else {
	// 	sort.Sort(SortSlice(f.temps))
	// }
}

// Process 处理数据
func (f *FastSort) Process() {
	//计算重复
	f.Sort()

	for _, temp := range f.temps {
		// if found, index := f.find(&temp); found {
		if found, index := BinarySearch(f.data, temp); found {
			f.data[index].Score += 1
		} else {
			f.data = append(f.data, model.SliceItem{
				Id:    temp,
				Score: 1,
			})
			f.count++
		}
	}
	//对分数进行排序
	sort.Sort(sort.Reverse(ScoreSlice(f.data)))
}
func UniqueFast(obj []model.SliceItem) []model.SliceItem {
	keys := make(map[string]bool)
	list := []model.SliceItem{}
	for _, entry := range obj {
		// 假设以Name作为去重的依据
		if _, value := keys[entry.Id]; !value {
			keys[entry.Id] = true
			list = append(list, entry)
		} else {
			entry.Score += 1
		}
	}
	return list
}
func debug(tag string, s ...interface{}) {
	if tag == "8-114401-115881" {
		fmt.Println(s...)
	}
}
func BinarySearch(slice []model.SliceItem, target string) (bool, int) {
	left, right := 0, len(slice)-1
	// debug(target, fmt.Sprintf("----------%s-----------\n", target), slice)
	for left <= right {
		mid := left + (right-left)/2 // 防止溢出
		v := slice[mid].Id
		// debug(target, left, "[", mid, "]", right, "|", v, "<", target, "=", v < target)
		if v == target {
			debug(target, target)
			return true, mid
		} else if v < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	// debug(target, fmt.Sprintf("----------%s-----------\n", target), slice)
	return false, -1 // 未找到
}
func RemoveDuplicates(slice []model.SliceItem) []model.SliceItem {
	if len(slice) < 2 {
		return slice
	}
	unique := make([]model.SliceItem, 0, len(slice))
	unique = append(unique, slice[0]) // 添加第一个元素到无重复切片
	for i := 1; i < len(slice); i++ {
		if slice[i].Id != slice[i-1].Id { // 如果当前元素与前一个元素不同，则添加到无重复切片
			unique = append(unique, slice[i])
		}
	}
	return unique
}

func (f *FastSort) GetAll(result *[]model.SliceItem, start int, end int) {

	*result = f.data[start:end]
}
func (f *FastSort) Get(start int) model.SliceItem {
	d := f.data[start]
	return d
}
