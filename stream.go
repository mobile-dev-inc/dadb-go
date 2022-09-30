package main

type AdbdStream struct {
	connection *AdbdConnection
	localId    uint32
	remoteId   uint32

	payload []byte
}

func (s *AdbdStream) SupportsFeature(feature string) bool {
	_, ok := s.connection.connectionResponse.features[feature]
	return ok
}

func (s *AdbdStream) Read(p []byte) (int, error) {
	if len(s.payload) == 0 {
		pkt := <-s.connection.getChannel(s.localId, cmdWrte)
		s.payload = pkt.Payload

		err := writePacket(s.connection.rw, packet{
			Command: cmdOkay,
			Arg0:    s.localId,
			Arg1:    s.remoteId,
			Payload: nil,
		})

		if err != nil {
			return 0, err
		}
	}

	n := copy(p, s.payload)

	s.payload = s.payload[n:]

	return n, nil
}

func (s *AdbdStream) Write(p []byte) (int, error) {
	// TODO what about when len(p) > s.connection.connectionResponse.maxPayloadSize?
	err := writePacket(s.connection.rw, packet{
		Command: cmdWrte,
		Arg0:    s.localId,
		Arg1:    s.remoteId,
		Payload: p,
	})

	if err != nil {
		return 0, err
	}

	<-s.connection.getChannel(s.localId, cmdOkay)

	return len(p), nil
}

func (s *AdbdStream) getPayload() ([]byte, error) {
	if len(s.payload) > 0 {
		return s.payload, nil
	}

	pkt := <-s.connection.getChannel(s.localId, cmdWrte)

	err := writePacket(s.connection.rw, packet{
		Command: cmdOkay,
		Arg0:    s.localId,
		Arg1:    s.remoteId,
		Payload: nil,
	})

	if err != nil {
		return nil, err
	}

	return pkt.Payload, nil
}
