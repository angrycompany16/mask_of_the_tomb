package assetloader

type Asset interface {
	Load() // Very simple asset loader interface
	// Can now create loaders for images, ldtk etc..
}
