package sip

import "context"

type ClientProps struct{}

type Client struct{}

func NewClient(props ClientProps) Client {
	return Client{}
}

func (client *Client) Serve() error {
	return nil
}

func (client *Client) Shutdown(ctx context.Context) error {
	return nil
}
