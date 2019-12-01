package tests

import (
	"context"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func DeleteUser(id string) {
	if id != "" {
		DB.Conn.Exec(context.Background(), "delete from users where id=$1", id)
	}
}

func DeletePlan(id int) {
	if id != 0 {
		DB.Conn.Exec(context.Background(), "delete from plans where id=$1", id)
	}
}

func DeleteTopic(id int) {
	if id != 0 {
		DB.Conn.Exec(context.Background(), "delete from topics where id=$1", id)
	}
}

func DeleteStep(id int64) {
	if id != 0 {
		DB.Conn.Exec(context.Background(), "delete from steps where id=$1", id)
	}
}

func DeleteSource(id int64) {
	if id != 0 {
		DB.Conn.Exec(context.Background(), "delete from sources where id=$1", id)
	}
}
