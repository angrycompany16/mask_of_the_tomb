package utils

import (
	"fmt"
	"mask_of_the_tomb/internal/engine"
	"reflect"
	"unsafe"

	"github.com/ebitengine/debugui"
)

// This goes somewhere else
func RenderFieldsAuto(ctx *debugui.Context, actor engine.Actor) {
	ctx.SetGridLayout([]int{-1, -2}, []int{0, 0})

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

func extractFieldUnsafe(v reflect.Value) reflect.Value {
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}
