package utils

import (
	"crypto/rand"
	"fmt"
	"log"
)

func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	fmt.Println(uuid)
	return uuid
}

// func MakeConn(ctx context.Context) (*pgxpool.Conn, error) {
// 	conn, err := server.Pool.Acquire(context.TODO())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return conn, nil
// }
