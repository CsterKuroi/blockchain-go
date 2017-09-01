package pipelines

import (
	"sync"

	"unichain-go/backend"
	"unichain-go/core"
	"unichain-go/log"
	"unichain-go/models"

	"encoding/json"
	mp "github.com/altairlee/multipipelines/multipipes"
)

func checkForQuorum(arg interface{}) interface{} {
	voteByte := []byte(arg.(string))
	vote := models.Vote{}
	err := json.Unmarshal(voteByte, &vote)
	if err != nil {
		log.Error(err)
		return nil
	}
	blockId := vote.VoteBody.VoteBlock
	valid := core.Election(blockId)
	log.Info("Elect `", valid, "`for", blockId)
	if valid != true {
		return blockId
	}
	return nil
}

func requeueTransactions(arg interface{}) interface{} {
	blockId := arg.(string)
	core.Requeue(blockId)
	return nil
}

func createElectionPipe() (p mp.Pipeline) {
	nodeSlice := make([]*mp.Node, 0)
	nodeSlice = append(nodeSlice, &mp.Node{Target: checkForQuorum, RoutineNum: 1, Name: "checkForQuorum"})
	nodeSlice = append(nodeSlice, &mp.Node{Target: requeueTransactions, RoutineNum: 1, Name: "requeueTransactions"})
	p = mp.Pipeline{
		Nodes: nodeSlice,
	}
	return p
}

func getElectionChangeNode() *mp.Node {
	cn := &changeNode{db: "unichain", table: "vote", operation: backend.INSERT}
	go cn.runForever()
	return &cn.node
}

func StartElectionPipe() {
	p := createElectionPipe()
	changeNode := getElectionChangeNode()
	p.Setup(changeNode, nil)
	p.Start()

	waitRoutine := sync.WaitGroup{}
	waitRoutine.Add(1)
	waitRoutine.Wait()
}
