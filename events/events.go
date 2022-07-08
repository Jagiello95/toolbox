package events

type TestEvent struct {
	Id   string
	Data TestEventPayload
}

type TestEventPayload struct {
	Message string
}

type UserDeletedEvent struct {
	Id   string
	Data UserDeletedEventPayload
}

type UserDeletedEventPayload struct {
	UserId int
}
