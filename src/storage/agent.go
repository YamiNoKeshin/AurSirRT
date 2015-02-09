package storage




func NewAgent() (agent StorageAgent) {
	agent.storageReadChannel = make(chan StorageFunc)
	agent.storageWriteChannel = make(chan StorageFunc)

	var sc StorageCore

	go sc.Run(agent.storageWriteChannel,agent.storageReadChannel)

	return
}

type StorageAgent struct {
	storageWriteChannel chan StorageFunc
	storageReadChannel chan StorageFunc
}

func (a StorageAgent) Write(fun StorageFunc) {
	a.storageWriteChannel <- fun
}

func (a StorageAgent) Read(fun StorageFunc) {
	a.storageReadChannel <- fun
}


func (a StorageAgent) GetId() string {
	c := make(chan string)
	defer close(c)
	a.Read( func(sc *StorageCore) {
		c<-sc.Root.Id()
	})
	return <-c
}

