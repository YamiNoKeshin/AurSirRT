package storage

import (
	"testing"
)

func TestStorageReadWrite(t *testing.T){

	agent := NewAgent()

	agent.Write(func (sc *StorageCore){
		sc.InMemoryGraph.CreateVertex("test",nil)
	})

	agent.Write(func (sc *StorageCore){
	if sc.GetVertex("test") == nil {
		t.Error("Could not retrieve Vertex")
	}
})
}
