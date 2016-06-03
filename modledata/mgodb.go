// mgodb
package modledata

import (
	"fmt"

	"gopkg.in/mgo.v2"
   	"github.com/ginuerzh/sportvideo/common/youlog"
)

var (
	mgoSession *mgo.Session
)

func getSession() (*mgo.Session, error) {
	if mgoSession == nil {
		var err error

		mgoSession, err = mgo.Dial(MongoAddr)

		mgoSession.SetMode(mgo.Monotonic, true)

		if err != nil {
			youlog.Warnning("Dial mongo error!")

			return nil, err
		}
	}
	return mgoSession.Clone(), nil
}

func withCollection(collection string, safe *mgo.Safe, fn func(*mgo.Collection) error) error {

	session, err := getSession()
	if session == nil {
		youlog.Warnning("getSession error")
		panic(err)
	}
	defer session.Close()

	session.SetSafe(safe)

	c := session.DB(dbName).C(collection)

	if er := fn(c); er != nil {

		fmt.Println(er)
		return er
	}

	return nil
}

func Insertinfo(collection string, info interface{}) error {

	insert := func(c *mgo.Collection) error {

		return c.Insert(info)
	}

	return withCollection(collection, &mgo.Safe{}, insert)
}

func Getcount(collection string, query interface{}) (count int, err error) {

	q := func(c *mgo.Collection) (err error) {
		count, err = c.Find(query).Count()

		youlog.Warnning(fmt.Sprintf("getcount count = %d, err = %s", count, err.Error()))
		return
	}

	err = withCollection(collection, nil, q)
	return
}

func Search(collection string, query, selector interface{}, sorts []string,
	skip, limit int, result interface{}) error {

	q := func(c *mgo.Collection) error {
		return c.Find(query).Select(selector).Sort(sorts...).Skip(skip).Limit(limit).All(result)
	}

	return withCollection(collection, nil, q)
}

func Exists(collection string, query interface{}) (b bool, err error) {
	q := func(c *mgo.Collection) error {
		n, err := c.Find(query).Count()
		b = n > 0
		return err
	}

	return b, withCollection(collection, nil, q)
}

func Apply(collection string, query interface{}, change mgo.Change, result interface{}) (info *mgo.ChangeInfo, err error) {
	apply := func(c *mgo.Collection) (err error) {
		info, err = c.Find(query).Apply(change, result)
		return err
	}

	return info, withCollection(collection, &mgo.Safe{}, apply)
}

func FindOne(collection string, query interface{}, selector interface{}, sorts []string, result interface{}) error {
	if result == nil {
		return nil
	}
	q := func(c *mgo.Collection) error {
		return c.Find(query).Select(selector).Sort(sorts...).One(result)
	}

	return withCollection(collection, nil, q)
}

func FindOneLimit(collection string, query interface{}, nLimit int, selector interface{}, sorts []string, result interface{}) error {
	if result == nil {
		return nil
	}
	q := func(c *mgo.Collection) error {
		return c.Find(query).Select(selector).Sort(sorts...).Limit(nLimit).All(result)
	}

	return withCollection(collection, nil, q)
}

func UpdateId(collection string, id interface{}, change interface{}) error {
	update := func(c *mgo.Collection) error {
		return c.UpdateId(id, change)
	}

	return withCollection(collection, &mgo.Safe{}, update)
}

func Upsert(collection string, query interface{}, change interface{}) error {
    update := func(c *mgo.Collection) error {
        _,error := c.Upsert(query,change)
        return error
	}

	return withCollection(collection, &mgo.Safe{}, update)
}

func Remove(collection string, query interface{}) error {
    update := func(c *mgo.Collection) error {
        error := c.Remove(query)
        return error
	}

	return withCollection(collection, &mgo.Safe{}, update)
}