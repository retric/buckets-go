package buckets

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
	"time"
)

const (
	TestDatabase = "test"
)

type (
	sessionFunc func(*mgo.Session)
)

/* Query all buckets from the db */
func bucketsQuery(session *mgo.Session) {
	// Retrieve buckets collection.
	collection := session.DB(TestDatabase).C("buckets")

	var buckets []Bucket
	err := collection.Find(nil).All(&buckets)
	if err != nil {
		log.Printf("bucketsQuery ERROR: %s\n", err)
		return
	}

	log.Printf("bucketsQuery")
}

/* Wrapper for calling other session tests */
func sessionWrap(session *mgo.Session, f sessionFunc) {
	// Request socket connection from session.
	// Close session when fusessionFunc func(*mgo.Session)nction is done and return connection to the pool.
	sessionCopy := session.Copy()
	defer sessionCopy.Close()

	f(sessionCopy)
}

/* Test inserting a task into the db */
func taskTest(session *mgo.Session) {
	fmt.Printf("taskTest: retrieving collection\n")
	bucketCollection := session.DB(TestDatabase).C("buckets")
	taskCollection := session.DB(TestDatabase).C("tasks")

	bucket := Bucket{ID: bson.NewObjectId(), Name: "weekly"}
	insertItem(bucket, bucketCollection)

	task := Task{ID: bson.NewObjectId(), Name: "read", Priority: 1,
		DateCreated: time.Now().Local(), DateModified: time.Now().Local(),
		Buckets: []bson.ObjectId{}, Completed: false}
	task.Buckets = append(task.Buckets, bucket.ID)
	fmt.Printf("taskTest: inserting task into collection")
	insertItem(task, taskCollection)

	fmt.Printf("retrieving task from collection")
	result := Task{}
	err := taskCollection.Find(bson.M{"name": "read"}).One(&result)
	if err != nil {
		log.Fatal("error:", err)
	}
	fmt.Println("Task:", result.Name)

	removeItem(bson.M{"name": "weekly"}, bucketCollection)
	removeItem(task, taskCollection)
}

/* Test suite for Buckets */
func TestBucket(t *testing.T) {
	mongoSession := dbSetup()
	bucket := CreateBucketTest(mongoSession)
	GetBucketTest(mongoSession, bucket.ID.Hex())
	RemoveBucketTest(mongoSession, bucket.ID.Hex())
}

/* Test retrieving a bucket */
func GetBucketTest(session *mgo.Session, id string) {
	bucket := getBucket(session, id)
	if bucket == nil {
		log.Fatal("error: bucket not found")
	}
}

/* Test creating a bucket */
func CreateBucketTest(session *mgo.Session) *Bucket {
	bucket := createBucket(session, "weekly", []string{"54f41e6a5786752068000003"})
	return bucket
}

/* Test removing a bucket */
func RemoveBucketTest(session *mgo.Session, id string) {
	removeBucket(session, id)
}

/* Test suite for Tasks */
func TestTask(t *testing.T) {
	//mongoSession := dbSetup()
}
