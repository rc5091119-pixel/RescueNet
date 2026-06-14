package main

func (c *Client) SafeWrite(messageType int, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Conn.WriteMessage(messageType, data)
}
