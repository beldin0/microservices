package microservice

type MessagebusConfig struct {
	// Config requirements for the service go in here
}

func MessagebusConnection(config *MessagebusConfig, AtLeastOnce interface{}) *Connection {
	return &Connection{
		Name: "messagebus",
		Operation: func() error {
			var err error
			AtLeastOnce, err = func() (interface{}, error) {
				// return nil, errors.New("unable to connect")
				return nil, nil
			}()
			return err
		},
	}
}
