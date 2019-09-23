package lib

/*
=============================================
STORAGE
=============================================
*/

type BLOCKSIZE float64

const (
	_ = iota // ignore first value by assigning to blank identifier
	_
	K4 BLOCKSIZE = 1 << (10 + iota)
	K8
	K16
	K32
	K64
	K128
	K256
	K512
	K1024
	K2048
)

type (
	StorageType struct {
		TypeId    string
		WriteRate float64
		ReadRate  float64
		Size      float64
	}

	Storage struct {
		*StorageType
		Name string
		ID   string
		Type string

		JBOD       *JBOD
		BlockSize  BLOCKSIZE
		ReadLink   *Link
		WriteLink  *Link
		usedSize   int64
		brokenPart float64
		fail       bool
		_index     int
	}
)

func NewStorage(storageType *StorageType, name string, jbod *JBOD) *Storage {
	readLink := NewLink(storageType.ReadRate, "")
	writeLink := NewLink(storageType.WriteRate, "")

	storage := &Storage{
		StorageType: storageType,
		Name:        name,
		ID:          name,
		Type:        "disk_drive",
		ReadLink:    readLink,
		WriteLink:   writeLink,
		BlockSize:   K4,
		JBOD:        jbod,
	}

	//env.storageLinks[readLink.cid] = readLink
	//env.storageLinks[writeLink.cid] = writeLink

	jbod.Disks.Push(storage)
	return storage
}

func (process *Process) WriteAsync(storage *Storage, packet *Packet) STATUS {
	return process._write(storage, packet, AsyncStorageEvent)
}

func (process *Process) WriteSync(storage *Storage, packet *Packet) STATUS {
	process._write(storage, packet, SyncStorageEvent)
	return process._add_sync()
}

func (process *Process) ReadAsync(storage *Storage, packet *Packet) (*Packet, STATUS) {
	process._read(storage, packet, AsyncStorageEvent)
	return packet, OK
}

func (process *Process) ReadSync(storage *Storage, packet *Packet) (*Packet, STATUS) {
	process._read(storage, packet, SyncStorageEvent)
	return packet, process._add_sync()
}

func (process *Process) _write(storage *Storage, packet *Packet, eventType EventType) STATUS {
	status, _ := _basic_event_adding(process, storage.WriteLink, packet, eventType)
	return status
}

func (process *Process) _read(storage *Storage, packet *Packet, eventType EventType) STATUS {
	status, _ := _basic_event_adding(process, storage.ReadLink, packet, eventType)
	return status
}

func (storage *Storage) Put(e *Event, globalQueue *globalEventQueue) {

}
