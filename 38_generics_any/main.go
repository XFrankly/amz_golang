package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/oklog/ulid/v2"
)

/*

基于泛型类型的 参数结构体解析
generics
*/

var (
	Ctx    = context.Background()
	Logger = log.New(os.Stderr, "INFO -", 13)
)

const ignoreField = "-"

type field struct {
	typ   reflect.Type
	name  string
	idx   int
	isKey bool
	isVer bool
}

type schema struct {
	key    *field
	ver    *field
	fields map[string]*field
}

func newSchema(t reflect.Type) schema {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("schema %q should be a struct", t))
	}

	s := schema{fields: make(map[string]*field, t.NumField())}

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		f := parse(sf)
		if f.name == ignoreField {
			continue
		}
		f.idx = i
		s.fields[f.name] = &f

		if f.isKey {
			if sf.Type.Kind() != reflect.String {
				panic(fmt.Sprintf("field with tag `redis:\",key\"` in schema %q should be a string", t))
			}
			s.key = &f
		}
		if f.isVer {
			if sf.Type.Kind() != reflect.Int64 {
				panic(fmt.Sprintf("field with tag `redis:\",ver\"` in schema %q should be a int64", t))
			}
			s.ver = &f
		}
	}

	if s.key == nil {
		panic(fmt.Sprintf("schema %q should have one field with `redis:\",key\"` tag", t))
	}
	if s.ver == nil {
		panic(fmt.Sprintf("schema %q should have one field with `redis:\",ver\"` tag", t))
	}

	return s
}

func parse(f reflect.StructField) (field field) {
	v, _ := f.Tag.Lookup("json")
	vs := strings.SplitN(v, ",", 1)
	if vs[0] == "" {
		field.name = f.Name
	} else {
		field.name = vs[0]
	}

	v, _ = f.Tag.Lookup("redis")
	field.isKey = strings.Contains(v, ",key")
	field.isVer = strings.Contains(v, ",ver")
	field.typ = f.Type
	return field
}

func key(prefix, id string) (key string) {
	sb := strings.Builder{}
	sb.Grow(len(prefix) + len(id) + 1)
	sb.WriteString(prefix)
	sb.WriteString(":")
	sb.WriteString(id)
	return sb.String()
}

type hashConvFactory struct {
	fields map[string]fieldConv
}

var converters = struct {
	val   map[reflect.Kind]converter
	ptr   map[reflect.Kind]converter
	slice map[reflect.Kind]converter
}{
	ptr: map[reflect.Kind]converter{
		reflect.Int64: {
			ValueToString: func(value reflect.Value) (string, bool) {
				if value.IsNil() {
					return "", false
				}
				return strconv.FormatInt(value.Elem().Int(), 10), true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return reflect.Value{}, err
				}
				return reflect.ValueOf(&v), nil
			},
		},
		reflect.String: {
			ValueToString: func(value reflect.Value) (string, bool) {
				if value.IsNil() {
					return "", false
				}
				return value.Elem().String(), true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				return reflect.ValueOf(&value), nil
			},
		},
		reflect.Bool: {
			ValueToString: func(value reflect.Value) (string, bool) {
				if value.IsNil() {
					return "", false
				}
				if value.Elem().Bool() {
					return "t", true
				}
				return "f", true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				b := value == "t"
				return reflect.ValueOf(&b), nil
			},
		},
	},
	val: map[reflect.Kind]converter{
		reflect.Int64: {
			ValueToString: func(value reflect.Value) (string, bool) {
				return strconv.FormatInt(value.Int(), 10), true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return reflect.Value{}, err
				}
				return reflect.ValueOf(v), nil
			},
		},
		reflect.Int: {
			ValueToString: func(value reflect.Value) (string, bool) {
				return strconv.FormatInt(value.Int(), 10), true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return reflect.Value{}, err
				}
				return reflect.ValueOf(v), nil
			},
		},
		reflect.String: {
			ValueToString: func(value reflect.Value) (string, bool) {
				return value.String(), true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				return reflect.ValueOf(value), nil
			},
		},
		reflect.Bool: {
			ValueToString: func(value reflect.Value) (string, bool) {
				if value.Bool() {
					return "t", true
				}
				return "f", true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				b := value == "t"
				return reflect.ValueOf(b), nil
			},
		},
	},
	slice: map[reflect.Kind]converter{
		reflect.Uint8: {
			ValueToString: func(value reflect.Value) (string, bool) {
				buf, ok := value.Interface().([]byte)
				if !ok {
					return "", false
				}
				return *(*string)(unsafe.Pointer(&buf)), true
			},
			StringToValue: func(value string) (reflect.Value, error) {
				buf := []byte(value)
				return reflect.ValueOf(buf), nil
			},
		},
	},
}

func newHashConvFactory(t reflect.Type, schema schema) *hashConvFactory {
	factory := &hashConvFactory{fields: make(map[string]fieldConv, len(schema.fields))}
	for name, f := range schema.fields {
		conv, ok := converters.val[f.typ.Kind()]
		switch f.typ.Kind() {
		case reflect.Ptr:
			conv, ok = converters.ptr[f.typ.Elem().Kind()]
		case reflect.Slice:
			conv, ok = converters.slice[f.typ.Elem().Kind()]
		}
		if !ok {
			k := f.typ.Kind()
			panic(fmt.Sprintf("schema %q should not contain unsupported field type %s.", t, k))
		}
		factory.fields[name] = fieldConv{conv: conv, idx: f.idx}
	}
	return factory
}
func (f hashConvFactory) NewConverter(entity reflect.Value) hashConv {
	return hashConv{factory: f, entity: entity}
}

func (r hashConv) ToHash() (fields map[string]string) {
	fields = make(map[string]string, len(r.factory.fields))
	for k, f := range r.factory.fields {
		ref := r.entity.Field(f.idx)
		if v, ok := f.conv.ValueToString(ref); ok {
			fields[k] = v
		}
	}
	return fields
}

func (r hashConv) FromHash(fields map[string]string) error {
	for k, f := range r.factory.fields {
		v, ok := fields[k]
		if !ok {
			continue
		}
		val, err := f.conv.StringToValue(v)
		if err != nil {
			return err
		}
		r.entity.Field(f.idx).Set(val)
	}
	return nil
}

type hashConv struct {
	factory hashConvFactory
	entity  reflect.Value
}
type fieldConv struct {
	conv converter
	idx  int
}

type converter struct {
	ValueToString func(value reflect.Value) (string, bool)
	StringToValue func(value string) (reflect.Value, error)
}

type HashRepository[T any] struct {
	schema schema
	typ    reflect.Type
	// client  rueidis.Client
	factory *hashConvFactory
	prefix  string
	idx     string
}
type Repository[T any] interface {
	NewEntity() (entity *T)
	// Fetch(ctx context.Context, id string) (*T, error)
	// FetchCache(ctx context.Context, id string, ttl time.Duration) (v *T, err error)
	// Search(ctx context.Context, cmdFn func(search FtSearchIndex) Completed) (int64, []*T, error)
	// Aggregate(ctx context.Context, cmdFn func(search FtAggregateIndex) Completed) (*AggregateCursor, error)
	SaveInfo(ctx context.Context, entity *T) (err error)
	// Remove(ctx context.Context, id string) error
	// CreateIndex(ctx context.Context, cmdFn func(schema FtCreateSchema) Completed) error
	// DropIndex(ctx context.Context) error
	// IndexName() string
}

func NewHashRepository[T any](prefix string, schema T) Repository[T] {
	repo := &HashRepository[T]{
		prefix: prefix,
		idx:    "hashidx:" + prefix,
		typ:    reflect.TypeOf(schema),
		// client: client,
	}
	repo.schema = newSchema(repo.typ)
	repo.factory = newHashConvFactory(repo.typ, repo.schema)
	return repo
}

var entropies = sync.Pool{
	New: func() interface{} {
		return ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	},
}

func id() string {
	n := time.Now()
	entropy := entropies.Get().(io.Reader)
	id := ulid.MustNew(ulid.Timestamp(n), entropy)
	entropies.Put(entropy)
	return id.String()
}

func (r *HashRepository[T]) NewEntity() (entity *T) {
	var v T
	reflect.ValueOf(&v).Elem().Field(r.schema.key.idx).Set(reflect.ValueOf(id()))
	return &v
}

func (r *HashRepository[T]) SaveInfo(ctx context.Context, entity *T) (err error) {
	val := reflect.ValueOf(entity).Elem()
	fmt.Printf("val:%v\n", val)
	fmt.Printf("enrity:%T, type:%v\n", entity, reflect.TypeOf(entity).Kind().String())
	fields := r.factory.NewConverter(val).ToHash()

	keyVal := fields[r.schema.key.name]
	verVal := fields[r.schema.ver.name]
	Logger.Printf("fields:%v keyVal:%v, verVal:%v\n", fields, keyVal, verVal)
	args := make([]string, 0, len(fields)*2)
	args = append(args, r.schema.ver.name, verVal) // keep the ver field be the first pair for the hashSaveScript
	delete(fields, r.schema.ver.name)
	for k, v := range fields {
		args = append(args, k, v)
	}

	// str, err := hashSaveScript.Exec(ctx, r.client, []string{key(r.prefix, keyVal)}, args).ToString()
	// if rueidis.IsRedisNil(err) {
	// 	return ErrVersionMismatch
	// }

	fmt.Println("hash string:", ctx, r, r.prefix, args)
	fmt.Printf("fields:%#v\n", fields)
	fmt.Printf("keyVal:%#v\n", keyVal)
	fmt.Printf("verVal:%#v\n", verVal)
	// if err == nil {
	// 	ver, _ := strconv.ParseInt(str, 10, 64)
	// 	val.Field(r.schema.ver.idx).SetInt(ver)
	// }
	return err
}

type Example struct {
	Key   string `json:"key" redis:",key"` // the redis:",key" is required to indicate which field is the ULID key
	Ver   int64  `json:"ver" redis:",ver"` // the redis:",ver" is required to do optimistic locking to prevent lost update
	ExStr string `json:"ex_str"`           // both NewHashRepository and NewJSONRepository use json tag as field name
}

type BaseCounter struct {
	Info string `json:"info" redis:",info"`
}
type Mycounter struct {
	// BaseCounter
	Key string `json:"key" redis:",key"` // the redis:",key" is required to indicate which field is the ULID key
	Ver int64  `json:"ver" redis:",ver"` // the redis:",ver" is required to do optimistic locking to prevent lost update

	Name   string `json:"name" redis:",name"`
	Number int64  `json:"number" redis:", number"`
	Age    int    `json:"age" redis:", age"`
}

func ObjRedisJson() {
	/*
		通用对象映射
		NewHashRepository并NewJSONRepository创建一个由 redis hash 或 RedisJSON 支持的 OM 存储库。
	*/
	// c := NewRueidesClient()
	// create the repo with NewHashRepository or NewJSONRepository
	repo := NewHashRepository("my_prefix", Example{})

	exp := repo.NewEntity()
	exp.ExStr = "newstr"
	Logger.Printf("exp:%#v, key:%#v\n", exp, exp.Key) // output 01FNH4FCXV9JTB9WTVFAAKGSYB
	err := repo.SaveInfo(Ctx, exp)                    // success
	fmt.Println("saveinfo err", err)

	repo2 := NewHashRepository("my_counter", Mycounter{})
	saves := repo2.NewEntity()
	Logger.Println("saves", saves)
	saves.Name = "rpcs"
	saves.Number = 123
	saves.Age = 87
	Logger.Println("saves", saves)

}

func main() {
	ObjRedisJson()
}
