package websocket

type Hub struct {
	Rooms      map[int]map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

type Message struct {
	MatchID int
	Data    []byte
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[int]map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if h.Rooms[client.MatchID] == nil {
				h.Rooms[client.MatchID] = make(map[*Client]bool)
			}
			h.Rooms[client.MatchID][client] = true

		case client := <-h.Unregister:
			if _, ok := h.Rooms[client.MatchID][client]; ok {
				delete(h.Rooms[client.MatchID], client)
				close(client.Send)
			}

		case msg := <-h.Broadcast:
			for client := range h.Rooms[msg.MatchID] {
				client.Send <- msg.Data
			}
		}
	}
}
