package subpub

type MessageHandler func(msg interface{})

type messageWithSubject struct {
	subject string
	msg     interface{}
}
