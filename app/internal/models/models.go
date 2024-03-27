package models

type User struct {
	Id       string
	ReseveId string
}

type Store struct {
	Id           int    `jsonapi:"primary,stores"`
	Name         string `jsonapi:"attr,store_name"`
	Availability bool   `jsonapi:"attr,is_available"`
}

type Item struct {
	Id     int    `jsonapi:"primary,items"`
	Name   string `jsonapi:"attr,item_name"`
	Size   string `jsonapi:"attr,item_size"`
	Amount int    `jsonapi:"attr,item_amount"`
}

type ReserveRequestMessage struct {
	Id       string          `jsonapi:"primary,message"`
	ClientId string          `jsonapi:"attr,client_id"`
	Stores   []*StoreMessage `jsonapi:"relation,stores"`
	Items    []*ItemsMessage `jsonapi:"relation,items"`
}

type StoreMessage struct {
	Id string `jsonapi:"primary,store"`
}

type ItemsMessage struct {
	Id         string `jsonapi:"primary,item"`
	ReservedIn string `jsonapi:"attr,store_id"`
	Count      int    `jsonapi:"attr,item_count"`
}

type ResponseMessageAfterCreate struct {
	Id           string          `jsonapi:"primary,message"`
	AppliedItems []*ItemsMessage `jsonapi:"relation,item"`
}
