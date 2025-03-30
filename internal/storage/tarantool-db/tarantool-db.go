package tarantool_db

import (
	"context"
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
	pollStruct "mm-polls/internal/lib/types/poll"
	"time"
)

type Tarantool struct {
	connection *tarantool.Connection
}

func toTuple(p *pollStruct.Poll) []interface{} {
	return []interface{}{p.ID, p.Title, p.OwnerID, p.Options, p.Votes, p.Voters, p.Active}
}

func NewTarantool(address, user string) (*Tarantool, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	dialer := tarantool.NetDialer{
		Address: address,
		User:    user,
	}

	opts := tarantool.Opts{
		Timeout: time.Second * 2,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		return nil, err
	}
	return &Tarantool{connection: conn}, nil
}

func (tt *Tarantool) InitDB() error {
	query := `CREATE TABLE IF NOT EXISTS polls (
    poll_id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    owner_id TEXT NOT NULL,
    options ARRAY NOT NULL,        
    votes ARRAY NOT NULL,  
    voters MAP NOT NULL,       
    active BOOLEAN DEFAULT TRUE);`
	_, err := tt.connection.Do(tarantool.NewExecuteRequest(query)).Get()
	return err
}

func (tt *Tarantool) CreatePoll(poll *pollStruct.Poll) error {
	_, err := tt.connection.Do(tarantool.NewInsertRequest("polls").Tuple(toTuple(poll))).Get()
	return err
}

func (tt *Tarantool) UpdatePoll(poll *pollStruct.Poll) error {
	_, err := tt.connection.Do(tarantool.NewReplaceRequest("polls").Tuple(toTuple(poll))).Get()
	if err != nil {
		return fmt.Errorf("error replacing the existing vote: %s", err)
	}
	return nil
}

func (tt *Tarantool) GetPoll(uuid string) (*pollStruct.Poll, error) {
	var result []pollStruct.Poll
	err := tt.connection.Do(tarantool.NewSelectRequest("polls").
		Iterator(tarantool.IterEq).
		Key(tarantool.StringKey{S: uuid}).
		Limit(1)).
		GetTyped(&result)
	if err != nil {
		return nil, err
	}
	if len(result) != 0 {
		return &result[0], nil
	}
	return nil, fmt.Errorf("poll id not found in db")
}
