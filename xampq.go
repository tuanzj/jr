package jr

// import (
// 	"github.com/streadway/amqp"
// )

//amqp://qrhsipio:nWxFdyvQGEQpBf2AR7PwhEKRI_2bx_ab@white-mynah-bird.rmq.cloudamqp.com/qrhsipio

type XAMPQ struct {
	ctx *Context
	// ins ampq.Connection
}

// func (x *XAMPQ) Connect() error {
// 	connection, err := amqp.Dial(os.Getenv("AMQP_URL"))
// 	defer connection.Close()
// 	x.ins = connection
// 	return err
// }
