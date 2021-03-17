package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
)

type Mikrotik struct {
	Host     string
	Username string
	Password string
	TLS      bool
	CA       string
	Insecure bool
}

func Unmarshal(reply routeros.Reply, v interface{}) error {
	rv := reflect.ValueOf(v)
	elem := rv.Elem()

	if rv.Kind() != reflect.Ptr {
		panic("Unmarshal cannot work without a pointer")
	}

	switch elem.Kind() {
	case reflect.Slice:
		l := len(reply.Re)
		if l <= 1 {
			panic(fmt.Sprintf("Cannot Unmarshal %d sentence(s) into a slice", l))
		}

		t := elem.Type()
		d := reflect.MakeSlice(t, l, l)

		for i := 0; i < l; i++ {
			item := d.Index(i)
			sentence := reply.Re[i]

			parseStruct(&item, *sentence)
		}
		elem.Set(d)

	case reflect.Struct:
		if len(reply.Re) < 1 {
			// This is an empty message
			return nil
		}
		if len(reply.Re) > 1 {
			msg := fmt.Sprintf("Failed to decode reply: %v", reply)
			return errors.New(msg)
		}

		parseStruct(&elem, *reply.Re[0])
	}

	return nil
}
func parseStruct(v *reflect.Value, sentence proto.Sentence) {
	elem := *v
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := elem.Type().Field(i)
		tags := strings.Split(fieldType.Tag.Get("mikrotik"), ",")

		path := strings.ToLower(fieldType.Name)
		fieldName := tags[0]

		for _, pair := range sentence.List {
			if strings.Compare(pair.Key, path) == 0 || strings.Compare(pair.Key, fieldName) == 0 {
				switch fieldType.Type.Kind() {
				case reflect.String:
					field.SetString(pair.Value)
				case reflect.Bool:
					b, _ := strconv.ParseBool(pair.Value)
					field.SetBool(b)
				case reflect.Int:
					if contains(tags, "ttlToSeconds") {
						field.SetInt(int64(ttlToSeconds(pair.Value)))
					} else {
						intValue, _ := strconv.Atoi(pair.Value)
						field.SetInt(int64(intValue))
					}
				}

			}
		}
	}
}

func ttlToSeconds(ttl string) int {
	parts := strings.Split(ttl, "d")

	idx := 0
	days := 0
	var err error
	if len(parts) == 2 {
		idx = 1
		days, err = strconv.Atoi(parts[0])

		// We should be parsing an ascii number
		// if this fails we should fail loudly
		if err != nil {
			panic(err)
		}

		// In the event we just get days parts[1] will be an
		// empty string. Just coerce that into 0 seconds.
		if parts[1] == "" {
			parts[1] = "0s"
		}
	}
	d, err := time.ParseDuration(parts[idx])

	// We should never receive a duration greater than
	// 23h59m59s. So this should always parse.
	if err != nil {
		panic(err)
	}
	return 86400*days + int(d)/int(math.Pow10(9))

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func NewClient(host, username, password string, tls bool, caCertificate string, insecure bool) Mikrotik {
	return Mikrotik{
		Host:      host,
		Username:  username,
		Password:  password,
		TLS:       tls,
		CA:        caCertificate,
		Insecure:  insecure,
	}
}

func GetConfigFromEnv() (host, username, password string, tls bool, caCertificate string, insecure bool) {
	host = os.Getenv("MIKROTIK_HOST")
	username = os.Getenv("MIKROTIK_USER")
	password = os.Getenv("MIKROTIK_PASSWORD")
	tlsString := os.Getenv("MIKROTIK_TLS")
	if tlsString == "true" {
		tls = true
	} else {
		tls = false
	}
	caCertificate = os.Getenv("MIKROTIK_CA_CERTIFICATE")
	insecureString := os.Getenv("MIKROTIK_INSECURE")
	if insecureString == "true" {
		insecure = true
	} else {
		insecure = false
	}
	if host == "" || username == "" || password == "" {
		// panic("Unable to find the MIKROTIK_HOST, MIKROTIK_USER or MIKROTIK_PASSWORD environment variable")
	}
	return host, username, password, tls, caCertificate, insecure
}

func (client Mikrotik) getMikrotikClient() (c *routeros.Client, err error) {
	address := client.Host
	username := client.Username
	password := client.Password

	if client.TLS {
		var tlsCfg tls.Config
		tlsCfg.InsecureSkipVerify = client.Insecure

		if client.CA != "" {
			certPool := x509.NewCertPool()
			file, err := ioutil.ReadFile(client.CA)
			if err != nil {
				log.Printf("[ERROR] Failed to read CA file %s: %v", client.CA , err)
				return nil, err
			}
			certPool.AppendCertsFromPEM(file)
			tlsCfg.RootCAs = certPool
		}

		c, err = routeros.DialTLS(address, username, password, &tlsCfg)
	} else {
		c, err = routeros.Dial(address, username, password)
	}
	if err != nil {
		log.Printf("[ERROR] Failed to login to routerOS with error: %v", err)
	}

	return
}

func boolToMikrotikBool(b bool) string {
	if b {
		return "yes"
	} else {
		return "no"
	}
}

func Marshal(s interface{}) string {
	var elem reflect.Value
	rv := reflect.ValueOf(s)

	if rv.Kind() == reflect.Ptr {
		// get Value of what pointer points to
		elem = rv.Elem()
	} else {
		elem = rv
	}

	var attributes []string

	for i := 0; i < elem.NumField(); i++ {
		value := elem.Field(i)
		fieldType := elem.Type().Field(i)
		// supports multiple struct tags--assumes first is mikrotik field name
		tag := strings.Split(fieldType.Tag.Get("mikrotik"), ",")[0]

		if tag != "" && (!value.IsZero() || value.Kind() == reflect.Bool) {
			switch value.Kind() {
			case reflect.Int:
				intValue := elem.Field(i).Interface().(int)
				attributes = append(attributes, fmt.Sprintf("=%s=%d", tag, intValue))
			case reflect.String:
				stringValue := elem.Field(i).Interface().(string)
				attributes = append(attributes, fmt.Sprintf("=%s=%s", tag, stringValue))
			case reflect.Bool:
				boolValue := elem.Field(i).Interface().(bool)
				stringBoolValue := boolToMikrotikBool(boolValue)
				attributes = append(attributes, fmt.Sprintf("=%s=%s", tag, stringBoolValue))
			}
		}
	}

	return strings.Join(attributes, " ")
}
