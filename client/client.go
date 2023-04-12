package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
)

// Mikrotik struct defines connection parameters for RouterOS client
type Mikrotik struct {
	Host     string
	Username string
	Password string
	TLS      bool
	CA       string
	Insecure bool

	connection *routeros.Client
}

type (
	// Marshaler interface will be used to serialize struct to RouterOS sentence
	Marshaler interface {
		// MarshalMikrotik serializes Go type value as RouterOS field value
		MarshalMikrotik() string
	}

	// Unmarshaler interface will be used to de-serialize reply from RouterOS into Go struct
	Unmarshaler interface {
		// UnmarshalMikrotik de-serializes RouterOS field into Go type value
		UnmarshalMikrotik(string) error
	}
)

// NewClient initializes new Mikrotik client object
func NewClient(host, username, password string, tls bool, caCertificate string, insecure bool) *Mikrotik {
	return &Mikrotik{
		Host:     host,
		Username: username,
		Password: password,
		TLS:      tls,
		CA:       caCertificate,
		Insecure: insecure,
	}
}

func Marshal(c string, s interface{}) []string {
	var elem reflect.Value
	rv := reflect.ValueOf(s)

	if rv.Kind() == reflect.Ptr {
		// get Value of what pointer points to
		elem = rv.Elem()
	} else {
		elem = rv
	}

	cmd := []string{c}

	for i := 0; i < elem.NumField(); i++ {
		value := elem.Field(i)
		fieldType := elem.Type().Field(i)
		// fetch mikrotik struct tag, which supports multiple values separated by commas
		tags := fieldType.Tag.Get("mikrotik")
		// extract tag value that is the Mikrotik property name
		// it is assumed that the first is mikrotik field name
		mikrotikTags := strings.Split(tags, ",")
		mikrotikPropName := mikrotikTags[0]
		// now we have field name in separate variable,
		// so leave only modifiers in this slice
		mikrotikTags = mikrotikTags[1:]

		if mikrotikPropName != "" && (!value.IsZero() || value.Kind() == reflect.Bool) {
			// add conditional to check if a Mikrotik property is READ ONLY, such as the following wireguard props
			// https://help.mikrotik.com/docs/display/ROS/WireGuard#WireGuard-Read-onlyproperties
			if contains(mikrotikTags, "readonly") {


				// if a struct field contains the tag value of 'readonly', do not marshal it
				continue
			}

			if mar, ok := value.Interface().(Marshaler); ok {
				// if type supports custom marshaling, use that result immediately
				stringValue := mar.MarshalMikrotik()
				cmd = append(cmd, fmt.Sprintf("=%s=%s", mikrotikPropName, stringValue))
				continue
			}

			switch value.Kind() {
			case reflect.Int:
				intValue := elem.Field(i).Interface().(int)
				cmd = append(cmd, fmt.Sprintf("=%s=%d", mikrotikPropName, intValue))
			case reflect.String:
				stringValue := elem.Field(i).Interface().(string)
				cmd = append(cmd, fmt.Sprintf("=%s=%s", mikrotikPropName, stringValue))
			case reflect.Bool:
				boolValue := elem.Field(i).Interface().(bool)
				stringBoolValue := boolToMikrotikBool(boolValue)
				cmd = append(cmd, fmt.Sprintf("=%s=%s", mikrotikPropName, stringBoolValue))
			}
		}
	}

	return cmd
}

// Unmarshal decodes MikroTik's API reply into Go object
func Unmarshal(reply routeros.Reply, v interface{}) error {
	rv := reflect.ValueOf(v)
	elem := rv.Elem()

	if rv.Kind() != reflect.Ptr {
		panic("Unmarshal cannot work without a pointer")
	}

	switch elem.Kind() {
	case reflect.Slice:
		l := len(reply.Re)
		t := elem.Type()
		if l < 1 {
			elem.Set(reflect.MakeSlice(t, 0, 0))
			break
		}

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

func (client *Mikrotik) getMikrotikClient() (*routeros.Client, error) {
	if client.connection != nil {
		return client.connection, nil
	}

	address := client.Host
	username := client.Username
	password := client.Password

	var mikrotikClient *routeros.Client
	var err error

	if client.TLS {
		var tlsCfg tls.Config
		tlsCfg.InsecureSkipVerify = client.Insecure

		if client.CA != "" {
			certPool := x509.NewCertPool()
			file, err := os.ReadFile(client.CA)
			if err != nil {
				log.Printf("[ERROR] Failed to read CA file %s: %v", client.CA, err)
				return nil, err
			}
			certPool.AppendCertsFromPEM(file)
			tlsCfg.RootCAs = certPool
		}

		mikrotikClient, err = routeros.DialTLS(address, username, password, &tlsCfg)
		if err != nil {
			return nil, err
		}
	} else {
		mikrotikClient, err = routeros.Dial(address, username, password)
	}

	if err != nil {
		log.Printf("[ERROR] Failed to login to routerOS with error: %v", err)
		return nil, err
	}

	client.connection = mikrotikClient

	return mikrotikClient, nil
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
				if field.CanAddr() {
					if unmar, ok := field.Addr().Interface().(Unmarshaler); ok {
						// if type supports custom unmarshaling, try it and skip the rest
						if err := unmar.UnmarshalMikrotik(pair.Value); err != nil {
							log.Printf("[ERROR] cannot unmarshal RouterOS reply: %v", err)
						}
						continue
					}
				}

				switch fieldType.Type.Kind() {
				case reflect.String:
					field.SetString(pair.Value)
				case reflect.Bool:
					b, _ := strconv.ParseBool(pair.Value)
					field.SetBool(b)
				case reflect.Int:
					intValue, _ := strconv.Atoi(pair.Value)
					field.SetInt(int64(intValue))
				}
			}
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func boolToMikrotikBool(b bool) string {
	if b {
		return "yes"
	} else {
		return "no"
	}
}
