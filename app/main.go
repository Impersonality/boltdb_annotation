package main

import (
	"bolt"
)

func main() {
	db, err := bolt.Open("./my.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err1 := tx.CreateBucketIfNotExists([]byte("user"))
		if err1 != nil {
			return err1
		}
		for i := 0; i < 1000; i++ {
			err1 = bucket.Put([]byte("hello"+string(i)), []byte("world1"))
			if err1 != nil {
				return err1
			}
		}
		tx.Rollback()
		return nil
	})
	if err != nil {
		panic(err)
	}

	//err = db.View(func(tx *bolt.Tx) error {
	//	bucket := tx.Bucket([]byte("user"))
	//	val := bucket.Get([]byte("hello"))
	//	//bucketInner := bucket.Bucket([]byte("name"))
	//	//val := bucketInner.Get([]byte("aaa"))
	//	fmt.Println("-------", string(val))
	//	return nil
	//})
	//if err != nil {
	//	panic(err)
	//}
}
