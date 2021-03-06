package rethinkdb

import (
	"time"

	"unichain-go/common"
	"unichain-go/log"

	r "gopkg.in/gorethink/gorethink.v3"
	"strconv"
)

func (c *RethinkDBConnection) Get(db string, table string, id string) *r.Cursor {
	res, err := r.DB(db).Table(table).Get(id).Run(c.Session)
	if err != nil {
		log.Error(err)
	}
	return res
}

func (c *RethinkDBConnection) Insert(db string, table string, jsonstr string) r.WriteResponse {
	res, err := r.DB(db).Table(table).Insert(r.JSON(jsonstr)).RunWrite(c.Session)
	if err != nil {
		log.Error(err)
	}
	return res
}

func (c *RethinkDBConnection) Update(db string, table string, id string, jsonstr string) r.WriteResponse {
	res, err := r.DB(db).Table(table).Get(id).Update(r.JSON(jsonstr)).RunWrite(c.Session)
	if err != nil {
		log.Error(err)
	}
	return res
}

func (c *RethinkDBConnection) Delete(db string, table string, id string) r.WriteResponse {
	res, err := r.DB(db).Table(table).Get(id).Delete().RunWrite(c.Session)
	if err != nil {
		log.Error(err)
	}
	return res
}

func (c *RethinkDBConnection) GetTransactionFromBacklog(id string) string {
	res := c.Get(DBUNICHAIN, TABLEBACKLOG, id)
	var value map[string]interface{}
	err := res.One(&value)
	if err != nil {
		log.Error("Error scanning database result:", err)
	}
	mapString := common.Serialize(value)
	return mapString
}

func (c *RethinkDBConnection) GetStaleTransactions(reassignDelay time.Duration) string {
	timeNow, err := strconv.Atoi(common.GenTimestamp())
	if err != nil {
		log.Error(err)
	}
	res, err := r.DB(DBUNICHAIN).Table(TABLEBACKLOG).
		Filter(r.Row.Field("AssignTime").Sub(timeNow).Lt(-reassignDelay)).
		Run(c.Session)
	if err != nil {
		log.Error(err)
	}
	var value []map[string]interface{}
	err = res.All(&value)
	if err != nil {
		//log.Debug("Error scanning database result:", err)
		return ""
	}
	mapString := common.Serialize(value)
	return mapString
}

func (c *RethinkDBConnection) WriteTransactionToBacklog(transaction string) int {
	res := c.Insert(DBUNICHAIN, TABLEBACKLOG, transaction)
	return res.Inserted
}

func (c *RethinkDBConnection) UpdateTransactionToBacklog(id string, jsonStr string) int {
	res := c.Update(DBUNICHAIN, TABLEBACKLOG, id, jsonStr)
	return res.Updated
}

func (c *RethinkDBConnection) DeleteTransaction(id string) int {
	res := c.Delete(DBUNICHAIN, TABLEBACKLOG, id)
	return res.Deleted
}

func (c *RethinkDBConnection) GetBlock(id string) string {
	res := c.Get(DBUNICHAIN, TABLEBLOCKS, id)
	var value map[string]interface{}
	err := res.One(&value)
	if err != nil {
		log.Error("Error scanning database result:", err)
	}
	mapString := common.Serialize(value)
	return mapString
}

func (c *RethinkDBConnection) GetGenesisBlock() string {
	res, err := r.DB(DBUNICHAIN).Table(TABLEBLOCKS).
		Filter(r.Row.Field("BlockBody").Field("Transactions").AtIndex(0).Field("Operation").Eq("GENESIS")).
		Pluck("id").
		Run(c.Session)
	if err != nil {
		log.Error(err)
	}
	var value map[string]string
	err = res.One(&value)
	if err != nil {
		log.Error("Error scanning database result:", err)
	}
	blockId := value["id"]
	return blockId
}

func (c *RethinkDBConnection) GetBlocksContainTransaction(id string) string {
	res, err := r.DB(DBUNICHAIN).Table(TABLEBLOCKS).
		GetAllByIndex("transaction_id", id).
		Pluck("id").
		Run(c.Session)
	if err != nil {
		log.Error(err)
	}
	var value []map[string]interface{}
	err = res.All(&value)
	if err != nil {
		log.Error("Error scanning database result:", err)
	}
	mapStrings := common.Serialize(value)
	return mapStrings
}

func (c *RethinkDBConnection) WriteBlock(block string) int {
	res := c.Insert(DBUNICHAIN, TABLEBLOCKS, block)
	return res.Inserted
}

func (c *RethinkDBConnection) WriteVote(vote string) int {
	res := c.Insert(DBUNICHAIN, TABLEVOTES, vote)
	return res.Inserted
}

func (c *RethinkDBConnection) GetLastVotedBlockId(pubkey string) string {
	res, err := r.DB(DBUNICHAIN).Table(TABLEVOTES).
		Filter(r.Row.Field("NodePubkey").Eq(pubkey)).
		Max(r.Row.Field("VoteBody").Field("Timestamp")).
		Field("VoteBody").Field("VoteBlock").
		Run(c.Session)
	if err != nil {
		return c.GetGenesisBlock()
	}

	var value string
	err = res.One(&value)
	if err != nil {
		log.Error("Error scanning database result:", err)
	}
	blockId := value
	return blockId
}

func (c *RethinkDBConnection) GetVotesByBlockId(id string) string {
	res, err := r.DB(DBUNICHAIN).Table(TABLEVOTES).
		Filter(r.Row.Field("VoteBody").Field("VoteBlock").Eq(id)).
		Run(c.Session)
	if err != nil {
		log.Error(err)
	}

	var value []map[string]interface{}
	err = res.All(&value)
	if err != nil {
		log.Error("Error scanning database result:", err)
	}
	mapStrings := common.Serialize(value)
	return mapStrings
}

func (c *RethinkDBConnection) GetUnvotedBlock(pubkey string) []string {
	//TODO doing unfinished lizhen *
	res, err := r.DB(DBUNICHAIN).Table(TABLEBLOCKS).Filter(
		func() {

		},
	).Run(c.Session)

	//.Run(c.Session)
	var value []map[string]interface{}
	err = res.All(&value)
	if err != nil {

	}
	//return common.Serialize(value)
	return nil
}

func (c *RethinkDBConnection) GetBlockCount() (int, error) {
	res, err := r.DB(DBUNICHAIN).Table(TABLEBLOCKS).Count().Run(c.Session)
	if err != nil {
		log.Error(err)
		return -1, err
	}
	var cnt int
	res.One(&cnt)
	return cnt, err
}
