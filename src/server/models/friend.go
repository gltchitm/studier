package models

import (
	"time"

	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/id"
)

type Friend struct {
	FriendId  string
	Timestamp time.Time
	FromId    string
	ToId      string
	Accepted  bool
}

func NewFriend(tx *db.Tx, fromId, toId string) (*Friend, error) {
	friendId, err := id.GenerateId(id.IdTypeRegular)
	if err != nil {
		return nil, err
	}

	friend := Friend{
		FriendId: friendId,
		ToId:     fromId,
		FromId:   toId,
		Accepted: false,
	}

	statement := `INSERT INTO friends (friend_id, from_id, to_id) VALUES ($1, $2, $3) RETURNING timestamp`
	err = tx.Tx.QueryRow(statement, friendId, fromId, toId).Scan(&friend.Timestamp)
	if err != nil {
		return nil, err
	}

	return &friend, nil
}

func FetchFriendByFromIdAndToId(tx *db.Tx, fromId, toId string) (*Friend, error) {
	friend := Friend{FromId: fromId, ToId: toId}
	statement := `SELECT friend_id, accepted, timestamp FROM friends WHERE from_id = $1 AND to_id = $2`
	err := tx.Tx.QueryRow(statement, fromId, toId).Scan(&friend.FriendId, &friend.Accepted, &friend.Timestamp)
	return &friend, err
}

func FetchAcceptedFriendByFromIdAndToId(tx *db.Tx, fromId, toId string) (*Friend, error) {
	friend := Friend{FromId: fromId, ToId: toId, Accepted: true}
	statement := `SELECT friend_id, timestamp FROM friends WHERE from_id = $1
		AND to_id = $2 AND accepted = TRUE`
	err := tx.Tx.QueryRow(statement, fromId, toId).Scan(&friend.FriendId, &friend.Timestamp)
	return &friend, err
}

func FetchFriendsByFromId(tx *db.Tx, accepted bool, fromId string) (*[]Friend, error) {
	statement := `SELECT friend_id, to_id, timestamp FROM friends WHERE
		from_id = $1 AND accepted = $2`
	rows, err := tx.Tx.Query(statement, fromId, accepted)
	if err != nil {
		return nil, err
	}

	var friends []Friend
	for rows.Next() {
		friend := Friend{FromId: fromId, Accepted: accepted}
		err = rows.Scan(&friend.FriendId, &friend.ToId, &friend.Timestamp)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return &friends, nil
}

func FetchFriendsByToId(tx *db.Tx, accepted bool, toId string) (*[]Friend, error) {
	statement := `SELECT friend_id, from_id, timestamp FROM friends WHERE
		to_id = $1 AND accepted = $2`
	rows, err := tx.Tx.Query(statement, toId, accepted)
	if err != nil {
		return nil, err
	}

	var friends []Friend
	for rows.Next() {
		friend := Friend{ToId: toId, Accepted: accepted}
		err = rows.Scan(&friend.FriendId, &friend.FromId, &friend.Timestamp)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return &friends, nil
}

func FetchFriendByFriendId(tx *db.Tx, friendId string) (*Friend, error) {
	friend := Friend{FriendId: friendId}
	statement := `SELECT timestamp, from_id, to_id, accepted FROM friends WHERE friend_id = $1`
	err := tx.Tx.QueryRow(statement, friendId).Scan(&friend.Timestamp, &friend.FromId, &friend.ToId,
		&friend.Accepted)
	return &friend, err
}

func CountUnacceptedFriendsByToId(tx *db.Tx, toId string) (*int, error) {
	var count int
	statement := `SELECT COUNT(*) FROM friends WHERE accepted = FALSE AND to_id = $1`
	err := tx.Tx.QueryRow(statement, toId).Scan(&count)
	return &count, err
}

func (friend *Friend) Delete(tx *db.Tx) error {
	statement := `DELETE FROM friends WHERE friend_id = $1`
	_, err := tx.Tx.Exec(statement, friend.FriendId)
	return err
}

func (friend *Friend) SetAccepted(tx *db.Tx, accepted bool) error {
	statement := `UPDATE friends SET accepted = $1 WHERE friend_id = $2`
	_, err := tx.Tx.Exec(statement, accepted, friend.FriendId)
	return err
}

func DeleteAllFriendsByFromId(tx *db.Tx, fromId string) error {
	statement := `DELETE FROM friends WHERE from_id = $1`
	_, err := tx.Tx.Exec(statement, fromId)
	return err
}

func DeleteAllFriendsByToId(tx *db.Tx, toId string) error {
	statement := `DELETE FROM friends WHERE to_id = $1`
	_, err := tx.Tx.Exec(statement, toId)
	return err
}
