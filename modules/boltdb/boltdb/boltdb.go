/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package boltdb

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	"github.com/asdine/storm/q"
	lz4 "github.com/bkaradzic/go-lz4"
	"github.com/boltdb/bolt"
	log "github.com/cihub/seelog"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/huminghe/infini-framework/core/api"
	"github.com/huminghe/infini-framework/core/global"
	"github.com/huminghe/infini-framework/core/orm"
	core "github.com/huminghe/infini-framework/core/ui"
	"github.com/huminghe/infini-framework/core/util"
	"github.com/huminghe/infini-framework/modules/boltdb/boltdb/ui"
	"github.com/huminghe/infini-framework/modules/ui/common"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type BoltdbStore struct {
	FileName string
	api.Handler
}

var db *storm.DB

var buckets hashset.Set
var initLocker sync.Mutex

func initBucket(name string) error {
	if buckets.Contains(name) {
		return nil
	}

	initLocker.Lock()
	defer initLocker.Unlock()

	//double check after lock in
	if buckets.Contains(name) {
		return nil
	}

	db.Bolt.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			log.Error("create bucket: ", err, ",", name)
			panic(err)
		}
		return nil
	})

	return nil
}

func (store BoltdbStore) Open() error {

	//loading or initializing boltdb
	if util.IsExist(store.FileName) {
		log.Debug("found boltdb file, start reload,", store.FileName)
	}

	var err error
	db, err = storm.Open(store.FileName, storm.BoltOptions(0600, &bolt.Options{Timeout: 5 * time.Second}), storm.Codec(protobuf.Codec))

	if err != nil {
		log.Errorf("error open boltdb: %s, %s", store.FileName, err)
		return err
	}

	log.Debug("boltdb successfully started:", store.FileName)

	if global.Env().IsDebug {
		core.HandleUIFunc("/admin/boltdb/", store.boltDBStatusAction)

		common.RegisterNav("BoltDB", "BoltDB", "/admin/boltdb/")
	}

	return nil
}

func (store BoltdbStore) Close() error {
	err := db.Close()
	if err != nil {
		log.Error("boltdb:", store.FileName, err)
	}
	return err
}

func (store BoltdbStore) GetCompressedValue(bucket string, key []byte) ([]byte, error) {

	data, err := store.GetValue(bucket, key)
	if err != nil {
		return nil, err
	}
	data, err = lz4.Decode(nil, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (store BoltdbStore) GetValue(bucket string, key []byte) ([]byte, error) {

	initBucket(bucket)

	var ret []byte = nil
	db.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get(key)
		if v != nil {
			ret = v
		}
		return nil
	})
	return ret, nil
}

func (store BoltdbStore) AddValueCompress(bucket string, key []byte, value []byte) error {
	value, err := lz4.Encode(nil, value)
	if err != nil {
		log.Error("Failed to encode:", err)
		return err
	}
	return store.AddValue(bucket, key, value)
}

func (store BoltdbStore) AddValue(bucket string, key []byte, value []byte) error {

	initBucket(bucket)

	db.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(key, value)
		return err
	})
	return nil
}

func (store BoltdbStore) DeleteKey(bucket string, key []byte) error {
	db.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Delete(key)
		return err
	})
	return nil
}

func (store BoltdbStore) DeleteBucket(bucket string) error {
	db.Bolt.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucket))
		return err
	})
	return nil
}

func (store BoltdbStore) boltDBStatusAction(w http.ResponseWriter, r *http.Request) {
	err := db.Bolt.View(func(tx *bolt.Tx) error {
		showUsage := (r.FormValue("usage") == "true")
		// Use the direct page id, if available.
		if r.FormValue("id") != "" {
			id, _ := strconv.Atoi(r.FormValue("id"))
			return ui.Page(w, r, tx, nil, id, showUsage)
		}

		// Otherwise extract the indexes and traverse.
		indexes, err := indexes(r)
		if err != nil {
			return err
		}

		return ui.Page(w, r, tx, indexes, 0, showUsage)
	})
	if err != nil {
		ui.Error(w, err)
	}
}

// parses and returns all indexes from a request.
func indexes(r *http.Request) ([]int, error) {
	var a = []int{0}
	if len(r.FormValue("index")) > 0 {
		for _, s := range strings.Split(r.FormValue("index"), ":") {
			i, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			a = append(a, i)
		}
	}
	return a, nil
}

func (store BoltdbStore) Get(key string, value interface{}, to interface{}) error {
	return db.One(key, value, to)
}

func (store BoltdbStore) Save(o interface{}) error {
	return db.Save(o)
}

func (store BoltdbStore) Update(o interface{}) error {
	return db.Update(o)
}

func (store BoltdbStore) Delete(o interface{}) error {
	return db.DeleteStruct(o)
}

func (store BoltdbStore) Count(o interface{}) (int, error) {
	return db.Count(o)
}

func (s BoltdbStore) Search(t1, t2 interface{}, q1 *orm.Query) (error, orm.Result) {
	result := orm.Result{}
	total, err := s.Count(t1)
	if err != nil {
		log.Debug(err)
		total = -1
	}
	result.Total = total
	result.Result = t2

	if q1.From < 0 {
		q1.From = 0
	}
	if q1.Size < 0 {
		q1.Size = 10
	}

	var q2 storm.Query
	if q1.Conds != nil {

		//TODO
		//q2 = db.Select(q.Eq(q1.Conds.Field, q1.Filter.Value)) //can't limit here, bug .Limit(q1.Size).Skip(q1.From)

	} else {
		q2 = db.Select(q.True()).Limit(q1.Size).Skip(q1.From)
	}

	if q1.Sort != nil && len(*q1.Sort) > 0 {
		for _, i := range *q1.Sort {
			q2 = q2.OrderBy(fmt.Sprintf("%s %s", i.Field, i.SortType)).Reverse()
		}
	}

	//t, _ := time.Parse(layout, skipDate)
	//query := db.Select(q.Gt("CreateTime", t)).Limit(size).Skip(from).Reverse().OrderBy("CreateTime")
	err = q2.Find(t2)
	if err != nil {
		log.Trace(err)
	}
	return err, result
}
