package errs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TODO: Simplify with if err != nil panic
func tracedPanic(err error) {
	// pc, file, no, ok := runtime.Caller(2)
	// funcDetails := runtime.FuncForPC(pc)
	// if err != nil {
	// 	if ok {
	// 		fmt.Println(no)
	// 		fmt.Printf("Panicking from %s, file %s, line number %d\n", filepath.Base((funcDetails.Name())), file, no)
	// 	}
	// 	log.Fatal(err)
	// }
	// if err != nil {
	// 	panic(err)
	// }
}

func MustNewImageFromFile(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	// img, _, err :=
	if err != nil {
		panic(err)
	}
	return img
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func MustSingle(err error) {
	if err != nil {
		panic(err)
	}
}
