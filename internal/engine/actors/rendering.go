package actors

import (
	"fmt"
	"mask_of_the_tomb/internal/engine"
	"reflect"
	"unsafe"

	"github.com/ebitengine/debugui"
)

type FieldType int

const (
	Header FieldType = iota
	Number
)

type Field struct {
	fieldType FieldType
	name      string
	value     reflect.Value
}

func RenderFieldsAuto(ctx *debugui.Context, actor engine.Actor) {
	v := reflect.ValueOf(actor).Elem()
	t := v.Type()

	ctx.Loop(t.NumField(), func(i int) {
		field := t.Field(i)
		debugField := field.Tag.Get("debug")
		if debugField == "auto" {
			ctx.Text(field.Name)
			ctx.Text(fmt.Sprintf("%v", (extractFieldUnsafe(v.Field(i)).Interface())))
		}
	})
}

// and with that....
// all my hard work deprecated...

// This is a fairly dangerous piece of code. I'm not sure when it panics,
// and it uses a lot of dirty tricks to get stuff working...
// It might be nice to cache some stuff here as well lowkey...

// TODO: Render inherited fields with recursion? Yes
// func RenderComponent(ctx *debugui.Context, actor *engine.Actor) {
// 	// Don't even know what to say. This is the worst shit I've seen in my entire life.
// 	// But now we are working with a pure struct
// 	// Although: performance concern?
// 	// Bro
// 	// Was i always so pedantic about performance...?
// 	v := reflect.ValueOf(actor).Elem().Elem().Elem()
// 	fieldLayoutW := []int{-1, -2}
// 	fieldLayoutH := []int{0, 0}

// 	// Need some sort of flattened list to make stuff iterable
// 	// Maybe we could also use struct tags for some useful stuff? Yup probably
// 	fields := flattenStruct(&v)

// 	// We really need different components to implement their own rendering
// 	// method
// 	// Although this is quite hard: Not everything here is a simple actor...
// 	ctx.Loop(len(fields), func(i int) {
// 		switch fields[i].fieldType {
// 		case Header:
// 			ctx.SetGridLayout(make([]int, 1), make([]int, 1))
// 			ctx.Text(fields[i].name)
// 		case Number:
// 			ctx.SetGridLayout(fieldLayoutW, fieldLayoutH)
// 			ctx.Text(fields[i].name)
// 			ctx.Text(fmt.Sprintf("%v", fields[i].value.Interface()))
// 		}
// 	})
// }

func extractFieldUnsafe(v reflect.Value) reflect.Value {
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}

// // Converts a struct into a flattened map with values
// func flattenStruct(v *reflect.Value) []Field {
// 	k := v.Kind()
// 	t := v.Type()
// 	if k != reflect.Struct {
// 		return []Field{
// 			Field{
// 				fieldType: Number, // this will change
// 				name:      t.Name(),
// 				value:     extractFieldUnsafe(*v),
// 			},
// 		}
// 	}

// 	fields := make([]Field, 0)
// 	return flattenRecurse(v, fields)
// }

// func flattenRecurse(v *reflect.Value, fields []Field) []Field {
// 	t := v.Type()
// 	nFields := t.NumField()
// 	newFields := make([]Field, 0)
// 	newFields = append(newFields, Field{
// 		fieldType: Header,
// 		name:      t.Name(),
// 		value:     *v,
// 	})
// 	for i := range nFields {
// 		fv := v.Field(i)
// 		fk := fv.Kind()
// 		ft := t.Field(i)

// 		if fk == reflect.Ptr {
// 			e := fv.Elem()
// 			fv = e
// 			fk = e.Kind()
// 		}

// 		if fk != reflect.Struct {
// 			newFields = append(newFields, Field{
// 				fieldType: Number,
// 				name:      ft.Name,
// 				value:     extractFieldUnsafe(fv),
// 			})
// 			continue
// 		}

// 		recurseFields := flattenRecurse(&fv, fields)
// 		newFields = slices.Concat(newFields, recurseFields)
// 	}
// 	return newFields
// }
