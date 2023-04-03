package helper

import (
	"sync"
	"sync/atomic"
)

var namespace NameSpace

type G map[string]*Client
type NameSpace struct {
	group sync.Map
	Lock  sync.RWMutex
}

func (ns *NameSpace) AddClient(client *Client) *Group {
	ns.Lock.Lock()
	defer ns.Lock.Unlock()
	g, ok := ns.group.Load(client.GroupName)
	if !ok {
		g = &Group{
			NameSpace: ns,
			Count:     &atomic.Int32{},
		}
		ns.group.Store(client.GroupName, g)
	}
	group := g.(*Group)
	client.NameSpace = ns

	client.Group = group
	group.AddClient(client)
	return group
}

func (ns *NameSpace) DelClient(client *Client) {
	ns.Lock.Lock()
	defer ns.Lock.Unlock()
	group := client.Group
	group.DelClient(client.Username)
	if group.Count.Load() <= 0 {
		group.DelClient(client.GroupName)
	}
}

func (ns *NameSpace) GetGroup(group string) (*Group, bool) {
	ns.Lock.RLock()
	defer ns.Lock.RUnlock()
	g, ok := ns.group.Load(group)
	if !ok {
		return nil, false
	}
	return g.(*Group), true
}

func (ns *NameSpace) GetClient(group string, username string) (*Client, bool) {
	ns.Lock.RLock()
	defer ns.Lock.RUnlock()
	g, ok := ns.GetGroup(group)
	if !ok {
		return nil, false
	}
	return g.GetClient(username)
}

func (ns *NameSpace) Broadcast(group string, username string, data any) {
	g, ok := ns.GetGroup(group)
	if ok {
		g.Broadcast(username, data)
	}
}

func (ns *NameSpace) SendTo(group string, from string, to string, data any) {
	client, ok := ns.GetClient(group, to)
	if ok {
		client.Send(from, data)
	}
}
