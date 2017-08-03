package backend

import (
	"unichain-go/backend/rethinkdb"
//	"unichain-go/backend/mongodb"
)

var regStruct map[string]Connection

type Connection interface {
	//connection
	Connect()
	//query
	GetTransaction(id string) map[string]interface{}
	SetTransaction(transaction string) int
	//changefeed
}

func init() {
	regStruct = make(map[string]Connection)
	regStruct["rethinkdb"] = &rethinkdb.RethinkDBConnection{}
	//regStruct["mongodb"] = &mongodb.MongoDBConnection{}
}

func GetConnection() Connection{
	var conn Connection
	str := "rethinkdb"//	TODO Config
	conn = regStruct[str]
	conn.Connect()
	return conn
}