package gotypes

import (
	"testing"
)

func TestMapString(t *testing.T) {
	t.Run("IsMapString", func(t *testing.T) {
		t.Run("String", func(t *testing.T) {
			k := "1"
			r := IsMapString(k)
			if r != false {
				t.Errorf("K is str, not a map of map of string")
			}
		})

		t.Run("Regular string map", func(t *testing.T) {
			k := `{"user_module": "abc"}`
			r := IsMapString(k)
			if r != true {
				t.Errorf("K is regular map of strs")
			}
		})

		t.Run("Actual map of map of strs", func(t *testing.T) {
			k := `{"user_module":{"read":"abc","write":"cde"}}`
			r := IsMapString(k)
			if r != false {
				t.Errorf("K is a map of map of string")
			}
		})

	})

	t.Run("IsMapMapString", func(t *testing.T) {
		t.Run("String", func(t *testing.T) {
			k := "1"
			r := IsMapMapString(k)
			if r != false {
				t.Errorf("K is str, not a map of map of string")
			}
		})

		t.Run("Regular string map", func(t *testing.T) {
			k := `{"user_module": "abc"}`
			r := IsMapMapString(k)
			if r != false {
				t.Errorf("K is regular map of strs, not a map of map of string")
			}
		})

		t.Run("Map of map but the inner map is string", func(t *testing.T) {
			k := `{"user_module":"{"read":"abc","write":"cde"}"}`

			r := IsMapMapString(k)
			if r != false {
				t.Errorf("K is regular map of strs, not a map of map of string")
			}
		})

		t.Run("Actual map of map of strs", func(t *testing.T) {
			k := `{"user_module":{"read":"abc","write":"cde"}}`
			r := IsMapMapString(k)
			if r != true {
				t.Errorf("K is a map of map of string")
			}
		})
	})
}

func TestMapBool(t *testing.T) {
	t.Run("IsMapBool", func(t *testing.T) {
		t.Run("Bool", func(t *testing.T) {
			k := "1"
			r := IsMapBool(k)
			if r != false {
				t.Errorf("K is str, not a map of map of string")
			}
		})

		t.Run("Regular string map", func(t *testing.T) {
			k := `{"user_module": true}`
			r := IsMapBool(k)
			if r != true {
				t.Errorf("K is regular map of strs")
			}
		})

		t.Run("Actual map of map of strs", func(t *testing.T) {
			k := `{"user_module":{"read":"abc","write":"cde"}}`
			r := IsMapBool(k)
			if r != false {
				t.Errorf("K is a map of map of string")
			}
		})

	})

	t.Run("IsMapMapBool", func(t *testing.T) {
		t.Run("Bool", func(t *testing.T) {
			k := "1"
			r := IsMapMapBool(k)
			if r != false {
				t.Errorf("K is str, not a map of map of string")
			}
		})

		t.Run("Regular bool map", func(t *testing.T) {
			k := `{"user_module": "bool"}`
			r := IsMapMapBool(k)
			if r != false {
				t.Errorf("K is regular map of strs, not a map of map of string")
			}
		})

		t.Run("Map of map but the inner map is string", func(t *testing.T) {
			k := `{"user_module":"{"read":true,"write":true}"}`

			r := IsMapMapBool(k)
			if r != false {
				t.Errorf("K is regular map of strs, not a map of map of string")
			}
		})

		t.Run("Actual map of map of bools", func(t *testing.T) {
			k := `{"user_module":{"read":true,"write":false}}`
			r := IsMapMapBool(k)
			if r != true {
				t.Errorf("K is a map of map of string")
			}
		})
	})
}
