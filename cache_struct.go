package main

type Cache struct {
	items map[string]OrderFields
}

func NewCache() *Cache {
	// инициализация map в паре ключ(string)/значение(OrderFields)
	items := make(map[string]OrderFields)

	cache := Cache{
		items: items,
	}
	return &cache
}

func (c *Cache) Set(key string, value OrderFields) {
	c.items[key] = value
}

func (c *Cache) Get(key string) (OrderFields, bool) {
	item, found := c.items[key]
	if !found {
		return OrderFields{}, false
	}
	return item, true
}
