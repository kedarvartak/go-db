package storage

import (
	"container/list"
)

type Cache struct {
	Capacity int
	pages    map[uint64]*list.Element
	lru      *list.List
}

type cacheEntry struct {
	pageID uint64
	page   *Page
}

func NewCache(capacity int) *Cache {
	return &Cache{
		Capacity: capacity,
		pages:    make(map[uint64]*list.Element),
		lru:      list.New(),
	}
}

// Get retrieves a page from the cache by its ID.
// If the page is found, it moves the page to the front of the LRU list and returns the page.
// If the page is not found, it returns nil and false.
func (c *Cache) Get(pageID uint64) (*Page, bool) {
	if element, found := c.pages[pageID]; found {
		c.lru.MoveToFront(element)
		return element.Value.(*cacheEntry).page, true
	}
	return nil, false
}

// Put adds a page to the cache.
// If the page is already in the cache, it updates the page in the cache.
// If the page is not in the cache, it adds the page to the cache and evicts the oldest page if the cache is full.
func (c *Cache) Put(page *Page) {
	if element, found := c.pages[page.ID]; found {
		c.lru.MoveToFront(element)
		element.Value.(*cacheEntry).page = page
		return
	}

	if c.lru.Len() >= c.Capacity {
		c.evictOldest()
	}

	entry := &cacheEntry{pageID: page.ID, page: page}
	element := c.lru.PushFront(entry)
	c.pages[page.ID] = element
}

// evictOldest removes the oldest page from the cache.
func (c *Cache) evictOldest() {
	element := c.lru.Back()
	if element != nil {
		c.lru.Remove(element)
		entry := element.Value.(*cacheEntry)
		delete(c.pages, entry.pageID)
	}
}

// Cache is used to store pages in memory to reduce disk I/O operations.
// cache struct
// LRU (Least Recently Used) eviction policy is used to remove the least recently used pages when the cache is full.
// pages are stored in a map for O(1) access time and a linked list to maintain the LRU order.
