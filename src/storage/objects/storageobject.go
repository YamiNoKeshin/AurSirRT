package objects

import "storage"

type StorageObject interface {
	Load(*storage.StorageCore)
}

type GenericStorageObject struct {
	StorageCore *storage.StorageCore
}

func (gso GenericStorageObject) Load(sc *storage.StorageCore) {
	gso.StorageCore = sc
}
