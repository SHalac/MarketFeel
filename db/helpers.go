package db

import (
	"errors"
	"github.com/boltdb/bolt"
)

var (
	bucketName = []byte("marketFeel")
)

/*
Function: Get bearer token from database
In: NONE
out: token (string) and error, error is not nil 
if token isn't found 
*/
func GetDbToken() (string, error){ // the issue right now is returning byte
	db,err := bolt.Open("token.db",0600,nil)
	if err != nil {
		return "error", errors.New("Could not get token from db")
	}
	defer db.Close()
	var token string
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return errors.New("can't find main bucket")
		}
		token = string(bucket.Get([]byte("token"))[:])
		if token == "" {
			return errors.New("can't find token key")
		}
		return nil
		})
	if err != nil {
		return "error", errors.New("Could not get token from db")
	}
	return token, nil

}

func AddTokenDb(token string) error {
	db,err := bolt.Open("token.db",0600,nil)
	if err != nil {
		return errors.New("Could not open db")
	}
	defer db.Close()
	err2 := db.Update(func(tx *bolt.Tx) error {
		bucket,err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte("token"),[]byte(token))
		return err
		})
	if err2 != nil {
		return errors.New("Could not add token to DB")
	}
	return nil
}