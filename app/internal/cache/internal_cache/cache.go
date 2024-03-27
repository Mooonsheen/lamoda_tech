package internalcache

// type InternalCache struct {
// 	reserve map[models.User]map[string]map[string]int
// }

// func InitCache() *InternalCache {
// 	reserve := make(map[models.User]map[string]map[string]int)
// 	// storageAvailability := make(map[string]bool)
// 	return &InternalCache{
// 		reserve: reserve,
// 		// storageAvailability: storageAvailability,
// 	}
// }

// func (c *InternalCache) GetReservation() map[string]map[string]map[string]int {
// 	return c.reserve
// }

// // func (c *InternalCache) GetStorageAvailability() map[string]bool {
// // 	return c.storageAvailability
// // }

// func (c *InternalCache) AddReserve(msg models.Message) {
// 	clientId := msg.ClientId
// 	if _, ok := c.reserve[clientId]; !ok {
// 		storageItems := make(map[string]map[string]int)
// 		for _, store := range msg.Stores {
// 			storeId := store.Id
// 			itemCount := make(map[string]int)
// 			for _, item := range msg.Items {
// 				itemId := item.Id
// 				count := item.Count
// 				itemCount[itemId] = count
// 			}
// 			storageItems[storeId] = itemCount
// 		}
// 		c.reserve[clientId] = storageItems
// 	} else {
// 		for _, store := range msg.Stores {
// 			storeId := store.Id

// 			if _, ok := c.reserve[clientId][storeId]; ok {
// 				for _, item := range msg.Items {
// 					itemId := item.Id
// 					count := item.Count
// 					if _, ok = c.reserve[clientId][storeId][itemId]; !ok {
// 						c.reserve[clientId][storeId][itemId] = count
// 					}
// 					c.reserve[clientId][storeId][itemId] += count
// 				}
// 			} else {
// 				itemCount := make(map[string]int)
// 				for _, item := range msg.Items {
// 					itemId := item.Id
// 					count := item.Count
// 					itemCount[itemId] = count
// 				}
// 				c.reserve[clientId][storeId] = itemCount
// 			}
// 		}
// 	}
// }

// // func (c *InternalCache) CancelReserve(msg models.Message) (nfItems []string, err error) { // на месте вызова собирать все сторы в цикле
// // 	if _, ok := c.reserve[store]; !ok {
// // 		return nil, fmt.Errorf("not found store: %s", store)
// // 	}

// // 	// nfItems - not found items
// // 	nfItems = make([]string, len(items_count))
// // 	for item, count := range items_count {
// // 		if _, ok := c.reserve[store][item]; !ok {
// // 			nfItems = append(nfItems, item)
// // 			continue
// // 		}
// // 		c.reserve[store][item] -= count
// // 	}

// // 	if len(nfItems) != 0 {
// // 		return nfItems, fmt.Errorf("not found items")
// // 	}

// // 	return nil, nil
// // }
