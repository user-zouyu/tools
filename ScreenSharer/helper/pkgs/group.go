package helper

import (
	"log"
	"sync"
	"sync/atomic"
)

type Group struct {
	clients   sync.Map
	NameSpace *NameSpace
	Count     *atomic.Int32
}

func (g *Group) AddClient(client *Client) {
	g.Count.Add(1)
	g.clients.Store(client.Username, client)
}

func (g *Group) DelClient(username string) {
	g.Count.Add(-1)
	c, ok := g.clients.LoadAndDelete(username)
	if ok {
		_ = c.(*Client).Conn.Close()
	}
}

func (g *Group) GetClient(username string) (*Client, bool) {
	c, ok := g.clients.Load(username)
	if !ok {
		return nil, false
	}
	return c.(*Client), true
}

func (g *Group) Broadcast(username string, data any) int {
	var count int
	g.clients.Range(func(u, value any) bool {
		client := value.(*Client)
		err := client.Conn.WriteJSON(data)
		if err != nil {
			log.Printf("发送错误(%s->%s): %v\n", username, u, err)
			return true
		}
		count++
		return true
	})
	return count
}

func (g *Group) SendTo(from string, to string, data any) {
	client, ok := g.GetClient(to)
	if ok {
		client.Send(from, data)
	}
}
