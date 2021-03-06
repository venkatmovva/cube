package config

import (
	"encoding/json"
	"fmt"
	"io"
)

type jsonStore struct {
	r  io.Reader
	kv map[string]interface{}
	kb map[string][]byte
}

// NewJSONStore returns a config store backed by a JSON stream.
//
// The first level keys in the JSON stream match the service names and the
// values must be decodeable into the types used to retrieve the config.
func NewJSONStore(r io.Reader) Store {
	return &jsonStore{
		r:  r,
		kb: map[string][]byte{},
	}
}

func (j *jsonStore) Open() error {
	d := json.NewDecoder(j.r)
	for {
		data := map[string]*cfgData{}
		if err := d.Decode(&data); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		// Cache the key and its corresponding json data
		for k, v := range data {
			j.kb[k] = v.b
		}
	}
}

func (j *jsonStore) Close() {
	// NOOP
}

func (j *jsonStore) Get(name string, config interface{}) error {
	if b, ok := j.kb[name]; ok {
		if e := json.Unmarshal(b, config); e != nil {
			// Bad buffer for the current type but lets keep it around
			// in case the registry is modified with a new type
			// and we can process it in future Get calls
			return e
		}
		return nil
	}
	return fmt.Errorf("%s key not found", name)
}

type cfgData struct {
	b []byte
}

func (d *cfgData) UnmarshalJSON(b []byte) error {
	d.b = b
	return nil
}
