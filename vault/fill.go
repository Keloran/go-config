package vault

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/bugfixes/go-bugfixes/logs"
	vaultHelper "github.com/keloran/vault-helper"
)

const DefaultDetailsPath = "secret/data/chewedfeed/details"

// Fill sets every field on details tagged `vault:"key"` (or `vault:"key,creds"` to read from
// the credentials path) from vault. A field is only filled when its env var isn't set and it
// still holds its envDefault, so env always wins and keys missing from vault are skipped.
func Fill(vh vaultHelper.VaultHelper, p Path, details any) error {
	v := reflect.ValueOf(details).Elem()
	t := v.Type()

	if p.Details == "" {
		p.Details = DefaultDetailsPath
	}

	for _, creds := range []bool{false, true} {
		path := p.Details
		if creds {
			path = p.Credentials
		}
		fetched := false

		for i := 0; i < t.NumField(); i++ {
			key, opt, ok := strings.Cut(t.Field(i).Tag.Get("vault"), ",")
			if key == "" || (ok && opt == "creds") != creds || !needsFill(t.Field(i), v.Field(i)) {
				continue
			}

			if !fetched {
				if path == "" {
					return logs.Errorf("vault: no credentials path for %s", key)
				}
				if err := vh.GetSecrets(path); err != nil {
					return logs.Errorf("vault: unable to get secrets from %s: %v", path, err)
				}
				fetched = true
			}

			secret, err := vh.GetSecret(key)
			if err != nil {
				continue // ponytail: helper only errors on missing keys, which we skip
			}
			if err := set(v.Field(i), secret); err != nil {
				return logs.Errorf("vault: unable to set %s: %v", key, err)
			}
		}
	}

	return nil
}

func needsFill(f reflect.StructField, v reflect.Value) bool {
	if name := f.Tag.Get("env"); name != "" {
		if _, set := os.LookupEnv(name); set {
			return false
		}
	}
	return v.IsZero() || fmt.Sprint(v.Interface()) == f.Tag.Get("envDefault")
}

func set(v reflect.Value, secret string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(secret)
	case reflect.Int:
		i, err := strconv.Atoi(secret)
		if err != nil {
			return err
		}
		v.SetInt(int64(i))
	case reflect.Map: // "key:value,key:value"
		m := make(map[string]string)
		for _, pair := range strings.Split(secret, ",") {
			k, val, ok := strings.Cut(pair, ":")
			if !ok {
				return fmt.Errorf("bad pair %q", pair)
			}
			m[k] = val
		}
		v.Set(reflect.ValueOf(m))
	default:
		return fmt.Errorf("unsupported kind %s", v.Kind())
	}
	return nil
}
