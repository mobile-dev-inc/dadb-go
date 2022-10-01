package adbserver

import "dadb"

func CreateDadb(address string, deviceQuery string) (dadb.Dadb, error) {
	c, err := Connect(address, deviceQuery)
	if err != nil {
		return dadb.Dadb{}, err
	}
	return dadb.Dadb{Connection: c}, nil
}
