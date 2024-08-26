package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func ExecTime(fn func()) float64 {
	start := time.Now()
	fn()
	tc := float64(time.Since(start).Nanoseconds())
	return tc / 1e6
}
func UUID() string {
	id, _ := GenerateRandomString(32)
	return id
}

// GenerateRandomString 生成固定长度的随机字符串
// length 是期望的字符串长度
// 注意：由于我们使用了hex编码，实际生成的随机字节长度需要是目标字符串长度的一半（向上取整）
func GenerateRandomString(length int) (string, error) {

	// 更安全的方式：
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}

	// 转换为十六进制字符串（这将把字节长度翻倍）
	// 如果你需要特定长度的字符串（不是字节长度的两倍），你需要调整逻辑
	// 例如，你可以取字节的一部分进行编码，或者对编码后的字符串进行截取
	hexStr := hex.EncodeToString(bytes)

	// 如果需要的长度小于生成的十六进制字符串长度，则进行截取
	if len(hexStr) > length {
		hexStr = hexStr[:length]
	}

	return hexStr, nil
}

func ExecTimeWithError(fn func() error) (float64, error) {
	start := time.Now()
	err := fn()
	tc := float64(time.Since(start).Nanoseconds())
	return tc / 1e6, err
}

// RemoveHTMLTags 去除HTML标签的函数
func RemoveHTMLTags(s string) string {
	// 使用正则表达式匹配HTML标签并替换为空字符串
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(s, "")
}
func Encoder(data interface{}) []byte {
	if data == nil {
		return nil
	}
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func Decoder(data []byte, v interface{}) {
	if data == nil {
		return
	}
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(v)
	if err != nil {
		panic(err)
	}
}

func Encode(data interface{}) []byte {
	b, _ := json.Marshal(data)
	return b
}
func Decode(data []byte) map[string]interface{} {
	var v map[string]interface{}
	json.Unmarshal(data, &v)
	return v
}

const (
	c1 = 0xcc9e2d51
	c2 = 0x1b873593
	c3 = 0x85ebca6b
	c4 = 0xc2b2ae35
	r1 = 15
	r2 = 13
	m  = 5
	n  = 0xe6546b64
)

var (
	Seed = uint32(1)
)

// 将字符串转换为SHA-1哈希值
func Sha1Hash(text string) string {
	// 创建一个新的hash.Hash对象，这里使用的是sha1.New()
	hasher := sha1.New()

	// 写入需要哈希的数据。这里我们将字符串转换为字节切片并写入
	io.WriteString(hasher, text)

	// Sum函数返回计算出的哈希值的字节切片。Sum(nil)表示将哈希值追加到一个空的切片中
	// 如果你不想追加到已有的切片，可以传递nil作为参数
	hashBytes := hasher.Sum(nil)

	// 将字节切片转换为十六进制字符串
	// 注意：strings.ToUpper()用于将小写字母转换为大写，如果你需要小写，可以去掉这个函数调用
	hashString := strings.ToUpper(fmt.Sprintf("%x", hashBytes))
	return hashString
}
func Obj2Map(obj interface{}) map[string]interface{} {
	if obj == nil {
		return make(map[string]interface{})
	}
	v := reflect.ValueOf(obj)
	m := make(map[string]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i).Name
		value := v.FieldByName(field).Interface()
		m[strings.ToLower(field)] = value
	}
	return m
}
func FindByExp(str interface{}, exp_str string) bool {
	aStr := fmt.Sprint(str)
	reg := regexp.MustCompile(exp_str)
	result := reg.FindString(aStr)
	if result == "" {
		return false
	}
	return true
}
func Murmur3(key []byte) (hash uint32) {
	hash = Seed
	iByte := 0
	for ; iByte+4 <= len(key); iByte += 4 {
		k := uint32(key[iByte]) | uint32(key[iByte+1])<<8 | uint32(key[iByte+2])<<16 | uint32(key[iByte+3])<<24
		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2
		hash ^= k
		hash = (hash << r2) | (hash >> (32 - r2))
		hash = hash*m + n
	}

	var remainingBytes uint32
	switch len(key) - iByte {
	case 3:
		remainingBytes += uint32(key[iByte+2]) << 16
		fallthrough
	case 2:
		remainingBytes += uint32(key[iByte+1]) << 8
		fallthrough
	case 1:
		remainingBytes += uint32(key[iByte])
		remainingBytes *= c1
		remainingBytes = (remainingBytes << r1) | (remainingBytes >> (32 - r1))
		remainingBytes = remainingBytes * c2
		hash ^= remainingBytes
	}

	hash ^= uint32(len(key))
	hash ^= hash >> 16
	hash *= c3
	hash ^= hash >> 13
	hash *= c4
	hash ^= hash >> 16

	// 出发吧，狗嬷嬷！
	return
}

// StringToInt 字符串转整数
func StringToInt(value string) uint32 {
	return Murmur3([]byte(value))
}

func Uint32Comparator(a, b interface{}) int {
	aAsserted := a.(uint32)
	bAsserted := b.(uint32)
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}

func Uint32ToBytes(i uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}

func Uint64ToBytes(n uint64) []byte {
	// 使用bytes.Buffer来构建字节切片
	var buf bytes.Buffer

	// 使用binary.Write将uint64值写入buffer，这里假设使用大端序
	err := binary.Write(&buf, binary.BigEndian, n)
	if err != nil {
		// 处理错误（这里仅打印）
		return nil
	}

	// 返回buffer中的字节切片
	return buf.Bytes()
}

// QuickSortAsc 快速排序
func QuickSortAsc(arr []int, start, end int, cmp func(int, int)) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		for i <= j {
			for arr[i] < key {
				i++
			}
			for arr[j] > key {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				if cmp != nil {
					cmp(i, j)
				}
				i++
				j--
			}
		}

		if start < j {
			QuickSortAsc(arr, start, j, cmp)
		}
		if end > i {
			QuickSortAsc(arr, i, end, cmp)
		}
	}
}
func DeleteArray(array []uint32, index int) []uint32 {
	return append(array[:index], array[index+1:]...)
}
func DeleteStringArray(array []string, index int) []string {
	return append(array[:index], array[index+1:]...)
}

func ReleaseAssets(file fs.File, out string) {
	if file == nil {
		return
	}

	if out == "" {
		panic("out is empty")
	}

	//判断out文件是否存在
	if _, err := os.Stat(out); os.IsNotExist(err) {
		//读取文件信息
		fileInfo, err := file.Stat()
		if err != nil {
			panic(err)
		}
		buffer := make([]byte, fileInfo.Size())
		_, err = file.Read(buffer)
		if err != nil {
			panic(err)
		}

		// 读取输出文件目录
		outDir := filepath.Dir(out)
		err = os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			panic(err)
		}

		//创建文件
		outFile, _ := os.Create(out)
		defer func(outFile *os.File) {
			err := outFile.Close()
			if err != nil {
				panic(err)
			}
		}(outFile)

		err = os.WriteFile(out, buffer, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

}

// DirSizeB DirSizeMB getFileSize get file size by path(B)
func DirSizeB(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})

	return size
}

// RemovePunctuation 移除所有的标点符号
func RemovePunctuation(str string) string {
	reg := regexp.MustCompile(`\p{P}+`)
	return reg.ReplaceAllString(str, "")
}

// RemoveSpace 移除所有的空格
func RemoveSpace(str string) string {
	reg := regexp.MustCompile(`\s+|\\n|\\r|\\t|'`)
	return reg.ReplaceAllString(str, "")
}

// init 注册数据类型
// 防止 gob: type not registered for interface: map[string]interface {}
func init() {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
}
