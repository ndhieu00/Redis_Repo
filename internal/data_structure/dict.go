package data_structure

import (
	"time"
)

type ValueObject struct {
	Value any
}

type Dict struct {
	dictStore        map[string]*ValueObject
	expiredDictStore map[string]uint64
}

func NewDict() *Dict {
	return &Dict{
		dictStore:        make(map[string]*ValueObject),
		expiredDictStore: make(map[string]uint64),
	}
}

/*
 * Dictionary implementation
 */

func (d *Dict) Get(key string) *ValueObject {
	v := d.dictStore[key]
	if v != nil && d.HasExpired(key) {
		d.Delete(key)
		return nil
	}

	return v
}

func (d *Dict) Set(key string, value any, expiryTimeMs uint64) {
	d.SetDictStore(key, value)
	if expiryTimeMs > 0 {
		d.SetExpiry(key, expiryTimeMs)
	} else {
		d.DeleteExpiry(key)
	}
}

func (d *Dict) Delete(key string) bool {
	if _, exist := d.dictStore[key]; !exist {
		return false
	}
	delete(d.dictStore, key)
	d.DeleteExpiry(key)
	return true
}

func (d *Dict) SetDictStore(key string, value any) {
	d.dictStore[key] = &ValueObject{value}
}

/*
 * Expired Dictionary store implementation
 */

func (d *Dict) GetExpiredDictStore() map[string]uint64 {
	return d.expiredDictStore
}

func (d *Dict) GetExpiryTime(key string) (uint64, bool) {
	expiryTime, exist := d.expiredDictStore[key]
	return expiryTime, exist
}

func (d *Dict) SetExpiry(key string, expiryTimeMs uint64) {
	d.expiredDictStore[key] = expiryTimeMs
}

func (d *Dict) DeleteExpiry(key string) {
	delete(d.expiredDictStore, key)
}

func (d *Dict) HasExpired(key string) bool {
	expiryTime, exist := d.GetExpiryTime(key)
	if !exist {
		return false
	}

	now := uint64(time.Now().UnixMilli())
	return expiryTime < now
}
