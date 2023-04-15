package mergego_test

// import (
// 	"fmt"
// 	"reflect"
// 	"testing"
// 	"time"

// 	"github.com/imdario/mergo"
// )

// // Foo
// type Foo struct {
// 	A string
// 	B int64
// }

// func TestMergegoMerge(t *testing.T) {
// 	src := Foo{
// 		A: "one",
// 		B: 2,
// 	}
// 	dest := Foo{
// 		A: "two",
// 	}
// 	mergo.Merge(&dest, src)
// 	t.Log(dest) // {two 2}

// 	dest2 := Foo{
// 		A: "two",
// 	}
// 	mergo.MergeWithOverwrite(&dest2, src)
// 	t.Log(dest2) // {one 2}
// }

// type DB struct {
// 	ID        int64
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

// func (d DB) Transformer(dstTyp reflect.Type) func(dst, src reflect.Value) error {
// 	var fn func(dst, src reflect.Value) error

// 	fmt.Println(dstTyp.String())

// 	switch dstTyp {
// 	case reflect.TypeOf(&Svc{}):
// 		fn = func(dst, src reflect.Value) error {
// 			v := src.Interface()
// 			db := v.(DB)
// 			svc := Svc{
// 				ID:        db.ID,
// 				CreatedAt: db.CreatedAt.Unix(),
// 				UpdatedAt: db.UpdatedAt.Unix(),
// 			}

// 			dst.Set(reflect.ValueOf(svc))
// 			return nil
// 		}
// 	case reflect.TypeOf(&Svc2{}):
// 		fn = func(dst, src reflect.Value) error {
// 			v := src.Interface()
// 			db := v.(DB)
// 			svc := Svc2{
// 				ID:                db.ID,
// 				CreatedAtAdd30Min: db.CreatedAt.Add(30 * time.Minute),
// 			}

// 			dst.Set(reflect.ValueOf(svc))
// 			return nil
// 		}
// 	default:
// 		fn = nil
// 	}

// 	return fn
// }

// type Svc struct {
// 	ID        int64
// 	CreatedAt int64
// 	UpdatedAt int64
// }

// type Svc2 struct {
// 	ID                int64
// 	CreatedAtAdd30Min time.Time
// }

// func TestMergeWithCustomTransformers(t *testing.T) {
// 	src := DB{
// 		ID:        1,
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	t.Log("src=", src)

// 	dst := new(Svc)
// 	if err := mergo.MapWithOverwrite(dst, src, mergo.WithTransformers(src)); err != nil {
// 		t.Error(err)
// 		t.FailNow()
// 	}
// 	t.Log("dst1=", dst) // &{1 1585201053 1585201053}

// 	dst2 := new(Svc2)
// 	if err := mergo.Merge(dst2, src, mergo.WithTransformers(src)); err != nil {
// 		t.Error(err)
// 		t.FailNow()
// 	}
// 	t.Log("dst2=", dst2)
// }
