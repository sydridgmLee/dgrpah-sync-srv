package kits

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
)

// func newClient() *dgo.Dgraph {
// 	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatal("While trying to dial gRPC")
// 	}
// 	defer conn.Close()

// 	dc := api.NewDgraphClient(conn)
// 	return dgo.NewDgraphClient(dc)
// }

func DBQuery(q string) (*api.Response, error) {
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)
	return dg.NewTxn().Query(context.Background(), q)
}

func DBMutate(b []byte) (*api.Response, error) {
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	mu := &api.Mutation{
		CommitNow: true,
	}

	mu.SetJson = b
	return dg.NewTxn().Mutate(context.Background(), mu)
}

func DBDelete(b []byte) error {
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	fmt.Println(string(b))

	mu := &api.Mutation{
		CommitNow:  true,
		DeleteJson: b,
	}
	resp, err := dg.NewTxn().Mutate(context.Background(), mu)

	fmt.Println(string(resp.Json))

	return err
}

func DBUpdate(set string) error {
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	mu := &api.Mutation{
		CommitNow: true,
		SetNquads: []byte(set),
	}

	_, err = dg.NewTxn().Mutate(context.Background(), mu)

	// fmt.Println(string(resp.Json))

	return err
}

func DBUpdateWithQuery(query, set string) error {
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	mu := &api.Mutation{
		CommitNow: true,
	}

	req := &api.Request{CommitNow: true}
	req.Query = query

	mu.SetNquads = []byte(set)
	req.Mutations = []*api.Mutation{mu}

	resp, err := dg.NewTxn().Do(context.Background(), req)

	fmt.Println(string(resp.Json))

	return err
}
