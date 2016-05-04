package myapp

import (
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/datastore"
	"encoding/json"
	"fmt"
)

func setUser(req *http.Request, m model){
	ctx := appengine.NewContext(req)
	bs := marshalModel(m)
	//Dank Memecache
	item := memcache.Item{
		Key: m.Name,
		Value: bs,
	}
	memcache.Set(ctx, &item)

	//Datastore
	key := datastore.NewKey(ctx, "Users", m.Name, 0, nil)
	_, err := datastore.Put(ctx, key, &m)
	if err != nil {
		panic(err)
		return
	}
}

func getUser(req *http.Request, name string) model{
	ctx := appengine.NewContext(req)
	//Dank Memcache
	item, _ := memcache.Get(ctx, name)
	var m model
	if item != nil{
		err := json.Unmarshal(item.Value, &m)
		if err != nil{
			fmt.Printf("error unmarhsalling: %v", err)
			return model{}
		}
	}
	//Datastore
	var m2 model
	key := datastore.NewKey(ctx, "Users", name, 0, nil)
	err := datastore.Get(ctx, key, &m2)
	if err == datastore.ErrNoSuchEntity{
		return model{}
	} else if err != nil{
		return model{}
	}

	//Reset Dank Memecache
	bs := marshalModel(m2)
	item1 := memcache.Item{
		Key: m2.Name,
		Value: bs,
	}
	memcache.Set(ctx, &item1)
	return m2
}

