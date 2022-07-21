package main

import (
	"bolt"
	"fmt"
)

func main() {
	db, err := bolt.Open("./my.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//err = db.Update(func(tx *bolt.Tx) error {
	//	bucket, err1 := tx.CreateBucketIfNotExists([]byte("user"))
	//	if err1 != nil {
	//		return err1
	//	}
	//	err1 = bucket.Put([]byte("hello"), []byte("world1"))
	//	if err1 != nil {
	//		return err1
	//	}
	//	return nil
	//})
	//if err != nil {
	//	panic(err)
	//}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("user"))
		val := bucket.Get([]byte("hello"))
		//bucketInner := bucket.Bucket([]byte("name"))
		//val := bucketInner.Get([]byte("aaa"))
		fmt.Println("-------", string(val))
		return nil
	})
	if err != nil {
		panic(err)
	}
}
