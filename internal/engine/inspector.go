package engine

import (
	"fmt"
	"reflect"
	"slices"
	"unsafe"

	"github.com/ebitengine/debugui"
)

type FieldType int

const (
	Header FieldType = iota
	Number
)

type InspectorField struct {
	fieldType FieldType
	name      string
	value     reflect.Value
}

// This is a fairly dangerous piece of code. I'm not sure when it panics,
// and it uses a lot of dirty tricks to get stuff working...
// It might be nice to cache some stuff here as well lowkey...

// TODO: Render inherited fields with recursion? Yes
func RenderComponent(ctx *debugui.Context, actor *Actor) {
	// Don't even know what to say. This is the worst shit I've seen in my entire life.
	// But now we are working with a pure struct
	// Although: performance concern?
	// Bro
	// Was i always so pedantic about performance...?
	v := reflect.ValueOf(actor).Elem().Elem().Elem()
	ctx.SetGridLayout([]int{-1, -2}, []int{0, 0})

	// Need some sort of flattened list to make stuff iterable
	// Maybe we could also use struct tags for some useful stuff? Yup probably
	// fields := flattenStruct(&v)
	names, values := flattenStruct(&v)
	// Another thing that would be nice: headers

	ctx.Loop(len(names), func(i int) {
		// If this field is a "simple" type, just display it as such
		// otherwise, recurse to render sub-structs as well
		ctx.Text(names[i])
		ctx.Text(fmt.Sprintf("%v", values[i].Interface()))
	})
}

func extractFieldUnsafe(v reflect.Value) reflect.Value {
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}

// Converts a struct into a flattened map with values
func flattenStruct(v *reflect.Value) ([]string, []reflect.Value) {
	k := v.Kind()
	t := v.Type()
	if k != reflect.Struct {
		return []string{t.Name()}, []reflect.Value{extractFieldUnsafe(*v)}
	}

	names := make([]string, 0)
	values := make([]reflect.Value, 0)
	return flattenRecurse(v, names, values)
}

func flattenRecurse(v *reflect.Value, names []string, values []reflect.Value) ([]string, []reflect.Value) {
	t := v.Type()
	nFields := t.NumField()
	newNames := make([]string, 0)
	newValues := make([]reflect.Value, 0)
	for i := range nFields {
		fv := v.Field(i)
		fk := fv.Kind()
		ft := t.Field(i)

		if fk == reflect.Ptr {
			e := fv.Elem()
			fv = e
			fk = e.Kind()
		}

		if fk != reflect.Struct {
			newNames = append(newNames, ft.Name)
			newValues = append(newValues, extractFieldUnsafe(fv))
			continue
		}

		recurseNames, recurseValues := flattenRecurse(&fv, names, values)
		newNames = slices.Concat(newNames, recurseNames)
		newValues = slices.Concat(newValues, recurseValues)
	}
	return newNames, newValues
}
