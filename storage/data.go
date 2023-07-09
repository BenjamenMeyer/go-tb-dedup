package storage

type messageData struct {
	id            string
	location      string
	locationIndex int
	hash          []byte
	from          string
	data          []byte
}

func NewMessageData() MailData {
	return &messageData{}
}

func (md *messageData) SetMessageId(id string) (err error) {
	md.id = id
	return
}
func (md *messageData) GetMessageId() (id string, err error) {
	id = md.id
	return
}

func (md *messageData) SetLocation(location string) (err error) {
	md.location = location
	return
}
func (md *messageData) GetLocation() (location string, err error) {
	location = md.location
	return
}

func (md *messageData) SetLocationIndex(index int) (err error) {
	md.locationIndex = index
	return
}
func (md *messageData) GetLocationIndex() (index int, err error) {
	index = md.locationIndex
	return
}

func (md *messageData) SetHash(hash []byte) (err error) {
	md.hash = hash
	return
}
func (md *messageData) GetHash() (hash []byte, err error) {
	hash = md.hash
	return
}

func (md *messageData) SetFrom(from string) (err error) {
	md.from = from
	return
}
func (md *messageData) GetFrom() (from string, err error) {
	from = md.from
	return
}

func (md *messageData) SetData(data []byte) (err error) {
	md.data = data
	return
}
func (md *messageData) GetData() (data []byte, err error) {
	data = md.data
	return
}

var _ MailData = &messageData{}
