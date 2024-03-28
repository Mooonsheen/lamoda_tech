package models

type RequestReserveMessage struct {
	Id       string          `jsonapi:"primary,message"`
	ClientId string          `jsonapi:"attr,client_id"`
	Uiud     string          `jsonapi:"attr,uuid_in,omitempty"`
	Stores   []*StoreMessage `jsonapi:"relation,stores"`
	Items    []*Items        `jsonapi:"relation,item"`
}

type ResponseReserveCreateMessage struct {
	Id           string   `jsonapi:"primary,message"`
	AppliedItems []*Items `jsonapi:"relation,item"`
}

type ResponseReserveChangeMessage struct {
	Id     string `jsonapi:"primary,message"`
	Status string `jsonapi:"attr,change_status"`
}

type Items struct {
	Id         string `jsonapi:"primary,item"`
	ReservedIn string `jsonapi:"attr,store_id"`
	Amount     int    `jsonapi:"attr,amount"`
}

type StoreMessage struct {
	Id          string   `jsonapi:"primary,store"`
	IsAvailable bool     `jsonapi:"attr,store_status"`
	Items       []*Items `jsonapi:"relation,item"`
}
