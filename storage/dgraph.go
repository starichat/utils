package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
	"log"
)

//Workflow 工作流
type Workflow struct {
	ID     string
	Desc   string
	Status int
}

//Task 任务
type Task struct {
	ID           string  `json:"ID" `
	Desc         string  `json:"Desc"`
	Status       int     `json:"Status"`
	PreTasks     []*Task `json:"PreTasks"`
	WorkFlowID   string  `json:"WorkFlowID"`
}

var (
	dgraph = flag.String("d", "https://play.dgraph.io", "Dgraph server address")
)

func main() {
	dg := NewDgraphClient()
	if err := CreateSchema(dg); err != nil {
		panic(err)
	}
	//创建scheme
	err := CreateSchema(dg)
	err = AddSomeData(dg)
	if err != nil {
		panic(err)
	}
	err = QueryData(dg)
	if err != nil {
		panic(err)
	}
	fmt.Println(1111111111)
	err = QueryDataID(dg)
	if err != nil {
		panic(err)
	}
	err = QueryReverseID(dg)
	if err != nil {
		panic(err)
	}


}

func NewDgraphClient() *dgo.Dgraph {
	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	return client
}

func CreateSchema(client *dgo.Dgraph) error {
	schema := `
  ID: string @index(term) .
  Desc: string .
  Status: int .
  WorkFlowID: string @index(term) .
  PreTasks: [uid] @reverse .

  type Task {
    ID           
	Desc         
	Status
	WorkFlowID
 	<~PreTasks>
  }
  `
	op := &api.Operation{Schema: schema}

	err := client.Alter(context.Background(), op)
	return err
}

func AddSomeData(client *dgo.Dgraph) error {
	tlast := &Task{
		ID:           "task-last",
		Desc:         "end",
		Status:       0,
		PreTasks:     nil,
		WorkFlowID:   "workflow-1",
	}
	t1 := &Task{
		ID:           "task-1",
		Desc:         "任务1",
		Status:       0,
		PreTasks:     nil,
		WorkFlowID:   "workflow-1",
	}
	t2 := &Task{
		ID:          "task-2",
		Desc:         "任务2",
		Status:       0,
		PreTasks:     nil,
		WorkFlowID:   "workflow-1",
	}
	t3 := &Task{
		ID:           "task-3",
		Desc:         "任务3",
		Status:       0,
		PreTasks:     nil,
		WorkFlowID:   "workflow-1",
	}
	t4 := &Task{
		ID:           "task-4",
		Desc:         "任务4",
		Status:       0,
		PreTasks:     nil,
		WorkFlowID:   "workflow-1",
	}
	t5 := &Task{
		ID:           "task-5",
		Desc:         "任务5",
		Status:       0,
		PreTasks:     nil,
		WorkFlowID:   "workflow-1",
	}
	t6 := &Task{
		ID:           "task-6",
		Desc:         "任务6",
		Status:       0,
		PreTasks:     nil,
		WorkFlowID:   "workflow-1",
	}
	t7 := &Task{
		ID:           "task-7",
		Desc:         "任务7",
		Status:       0,
		PreTasks:     nil,
		WorkFlowID:   "workflow-1",
	}
	//建立t4任务的依赖
	t4.PreTasks = append(t4.PreTasks, t1, t2)
	t5.PreTasks = append(t5.PreTasks, t4)
	t6.PreTasks = append(t6.PreTasks, t3,t4)
	t7.PreTasks = append(t7.PreTasks, t5, t6)
	//结束标识
	tlast.PreTasks = append(tlast.PreTasks, t7)

	//下一个任务的建立

	mu := &api.Mutation{CommitNow: true}
	pb, err := json.Marshal(tlast)
	if err != nil {
		panic(err)

	}
	mu.SetJson = pb
	r, err := client.NewTxn().Mutate(context.Background(), mu)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
	return nil

}

//QueryData 查询数据
func QueryData(client *dgo.Dgraph) error {
	q := `
 query q($WorkFlowID: string){
     q(func:allofterms($WorkFlowID, WorkFlowID)){
        ID          
		Desc        
		Status
		WorkFlowID  
        PreTasks {
                ID          
	   			Desc        
	   			Status
				WorkFlowID
            }
        }
    }
    `
	txn := client.NewTxn()
	res, err := txn.QueryWithVars(context.Background(), q, map[string]string{"$WorkFlowID": "workflow-1"})
	if err != nil {
		return err
	}
	fmt.Println(res.String())
	return nil
}

func QueryDataID(client *dgo.Dgraph) error {
	q := `
 query q($ID: string){
     q(func:allofterms($ID, ID)){
        ID          
		Desc        
		Status
		WorkFlowID  
        PreTasks {
                ID          
	   			Desc        
	   			Status
				WorkFlowID
            }
        }
    }
    `
	txn := client.NewTxn()
	res, err := txn.QueryWithVars(context.Background(), q, map[string]string{"$ID": "task-1"})
	if err != nil {
		return err
	}
	//判断依赖任务是否都完成了

	fmt.Println(res.String())
	return nil
}


func QueryReverseID(client *dgo.Dgraph) error {

	q := `
 query q($ID: string){
     q(func:eq($ID, ID)){
        ~PreTasks {
                ID          
	   			Desc        
	   			Status
				WorkFlowID
            }
        }
    }
    `
	txn := client.NewTxn()
	res, err := txn.QueryWithVars(context.Background(), q, map[string]string{"$ID": "task-4"})
	if err != nil {
		return err
	}
	//判断依赖任务是否都完成了

	fmt.Println(res.String())
	return nil
}
