package grpcserver

import (
	context "context"
	"sync"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Driver struct {
	Addr string

	ConnectionOnce sync.Once
	conn           *grpc.ClientConn
	client         GreeterClient
}

func (d *Driver) Greet(name string) (string, error) {
	client, err := d.getClient()
	if err != nil {
		return "", err
	}
	//defer d.Close() // <-- ConnectionOnce causes test to be failed if Greet() is called twice

	greeting, err := client.Greet(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}
	return greeting.Message, nil
}

func (d *Driver) Curse(name string) (string, error) {
	client, err := d.getClient()
	if err != nil {
		return "", err
	}
	//defer d.Close()

	cursing, err := client.Curse(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}
	return cursing.Message, nil
}

func (d *Driver) getClient() (GreeterClient, error) {
	var err error
	d.ConnectionOnce.Do(func() {
		d.conn, err = grpc.NewClient(d.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		d.client = NewGreeterClient(d.conn)
	})
	return d.client, err
}

func (d *Driver) Close() error {
	if d.conn == nil {
		return nil
	}
	return d.conn.Close()
}
