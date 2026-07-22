package datasaver

// System for saving data. Supports adding different scopes:
// These define the level at which an object which is saved should exist. So for
// example, game, level, area, etc.

// To force the development of a functional version of this system,
// let's try adding a key + locked door mechanic

// A metroidvania staple you know

type DataSaver struct {
	scopes []Scope
}

type Scope struct {
	name    string
	entries []DataEntry
}

type DataEntry struct {
	data any
}
