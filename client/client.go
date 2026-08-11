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

	"github.com/go-routeros/routeros/v3"
	"github.com/go-routeros/routeros/v3/proto"
)

const (
	tagMikrotik    = "mikrotik"
	tagReadonly    = "readonly"
	tagFallthrough = "fallthrough"
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

func Marshal(c string, s any) []string {

	values, err := marshalStruct(s)
	if err != nil {
		log.Fatal(err)
	}

	if len(values)%2 != 0 {
		log.Fatal("dont know how to marshal odd number of command elements")
	}
	cmd := make([]string, 0, len(values)+1)
	cmd = append(cmd, c)
	buf := strings.Builder{}
	for i := 0; i < len(values); i += 2 {
		buf.Reset()
		buf.WriteString("=")
		buf.WriteString(values[i])
		buf.WriteString("=")
		buf.WriteString(values[i+1])

		cmd = append(cmd, buf.String())
	}

	return cmd
}

func marshalStruct(s any) ([]string, error) {
	var elem reflect.Value
	rv := reflect.ValueOf(s)

	if rv.Kind() == reflect.Ptr {
		// get Value of what pointer points to
		elem = rv.Elem()
	} else {
		elem = rv
	}

	cmd := []string{}

	for i := 0; i < elem.NumField(); i++ {
		value := elem.Field(i)
		fieldType := elem.Type().Field(i)
		// fetch mikrotik struct tag, which supports multiple values separated by commas
		tags := fieldType.Tag.Get(tagMikrotik)
		// extract tag value that is the Mikrotik property name
		// it is assumed that the first is mikrotik field name
		mikrotikTags := strings.Split(tags, ",")
		mikrotikPropName := mikrotikTags[0]
		// now we have field name in separate variable,
		// so leave only modifiers in this slice
		mikrotikTags = mikrotikTags[1:]

		if mikrotikPropName == "" {
			// skip if property name is not set in struct tag
			continue
		}

		// skip attributes marked as `readonly` (computed by RouterOS), such as the following wireguard props
		// https://help.mikrotik.com/docs/display/ROS/WireGuard#WireGuard-Read-onlyproperties
		if contains(mikrotikTags, tagReadonly) {
			// if a struct field contains the tag value of 'readonly', do not marshal it
			continue
		}

		// types, which implement `Marshaler` interface, have higher priority as they likely to have specific logic
		// even for zero or nil values
		if mar, ok := value.Interface().(Marshaler); ok {
			// if type supports custom marshaling, use that result immediately
			stringValue := mar.MarshalMikrotik()
			cmd = append(cmd, mikrotikPropName, stringValue)
			continue
		}

		// try to marshal non-zero values only
		// here `bool` is exception as we want to toggle boolean switches
		if !value.IsZero() || value.Kind() == reflect.Bool {
			switch value.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				cmd = append(cmd, mikrotikPropName, strconv.FormatInt(elem.Field(i).Int(), 10))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				cmd = append(cmd, mikrotikPropName, strconv.FormatUint(elem.Field(i).Uint(), 10))
			case reflect.String:
				stringValue := elem.Field(i).Interface().(string)
				cmd = append(cmd, mikrotikPropName, stringValue)
			case reflect.Bool:
				boolValue := elem.Field(i).Interface().(bool)
				stringBoolValue := boolToMikrotikBool(boolValue)
				cmd = append(cmd, mikrotikPropName, stringBoolValue)
			case reflect.Struct:
				values, err := marshalStruct(elem.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				if len(values)%2 != 0 {
					return nil, errors.New("struct marshalling results in odd number of command elements")
				}
				structFieldNamePrefix := mikrotikPropName + "."
				for i := 0; i < len(values); i += 2 {
					cmd = append(cmd, structFieldNamePrefix+values[i], values[i+1])
				}
			}
		}
	}

	return cmd, nil
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

func (client *Mikrotik) GetMikrotikClient() (*routeros.Client, error) {
	return client.getMikrotikClient()
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

	nestedStructDelimiter := "."

	routerosFields := make(map[string]string)
	routerosContainerFields := make(map[string]map[string]string)
	for i := range sentence.List {
		routerosFields[sentence.List[i].Key] = sentence.List[i].Value
		if containerName, nestedFieldName, ok := strings.Cut(sentence.List[i].Key, nestedStructDelimiter); ok {
			if _, ok := routerosContainerFields[containerName]; !ok {
				routerosContainerFields[containerName] = make(map[string]string)
			}
			routerosContainerFields[containerName][nestedFieldName] = sentence.List[i].Value
			normalizedNestedFieldName := strings.ReplaceAll(strings.ToLower(containerName+nestedFieldName), nestedStructDelimiter, "")
			if _, ok := routerosFields[normalizedNestedFieldName]; !ok {
				routerosFields[normalizedNestedFieldName] = sentence.List[i].Value
			}
		}
	}

	elem := reflect.Indirect(*v)
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := elem.Type().Field(i)
		tags := strings.Split(fieldType.Tag.Get(tagMikrotik), ",")

		fieldName := strings.ToLower(fieldType.Name)
		fieldNameTag := tags[0]

		// field search order:
		// - in fields map by mikrotik tag value
		// - in nested fields map by tag value
		// - in fields map by normalized field name
		// - in nested fields map by normalized field name
		selectedFieldName := fieldName
		fieldValue, ok := routerosFields[fieldNameTag]
		if !ok {
			_, ok = routerosContainerFields[fieldNameTag]
			selectedFieldName = fieldNameTag
		}
		if !ok {
			fieldValue, ok = routerosFields[fieldName]
			selectedFieldName = fieldName
		}
		if !ok {
			_, ok = routerosContainerFields[fieldName]
			selectedFieldName = fieldName
		}
		if !ok {
			// skip processing if no matching field found in the reply
			continue
		}

		if field.CanAddr() {
			if unmar, ok := field.Addr().Interface().(Unmarshaler); ok {
				if fieldType.Type.Kind() == reflect.Struct {
					log.Println("[WARNING] custom Unmarshaler interface is not supported for struct fields")
					continue
				}
				// if type supports custom unmarshaling, try it and skip the rest
				if err := unmar.UnmarshalMikrotik(fieldValue); err != nil {
					// todo(maksym): if target field is a nested struct,
					// we need to unmarshal all values with field prefix, e.g. `fieldname.field1`, `fieldname.field2` etc.
					// so we might need to change the signature of `UnmarshalMikrotik`
					log.Printf("[ERROR] cannot unmarshal RouterOS reply: %v", err)
				}
				continue
			}
		}

		switch fieldType.Type.Kind() {
		case reflect.String:
			field.SetString(fieldValue)
		case reflect.Bool:
			b, _ := strconv.ParseBool(fieldValue)
			field.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, _ := strconv.Atoi(fieldValue)
			field.SetInt(int64(intValue))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintValue, _ := strconv.ParseUint(fieldValue, 10, 0)
			field.SetUint(uint64(uintValue))
		case reflect.Struct, reflect.Pointer:
			var nestedStruct reflect.Value

			switch fieldType.Type.Kind() {
			case reflect.Pointer:
				if fieldType.Type.Elem().Kind() != reflect.Struct {
					log.Printf("[WARN] Only struct pointers are supported for unmarshaling, but got: %s", fieldType.Type.Kind())
					continue
				}
				nestedStruct = reflect.New(fieldType.Type.Elem())
			case reflect.Struct:
				nestedStruct = reflect.New(fieldType.Type).Elem()
			}

			nestedSentence := proto.Sentence{
				List: make([]proto.Pair, 0),
			}
			// filter the fields that belong to the nested struct based on the prefix
			prefix := selectedFieldName + nestedStructDelimiter
			for k, v := range routerosContainerFields[selectedFieldName] {
				nestedKey := strings.TrimPrefix(k, prefix)
				nestedSentence.List = append(nestedSentence.List, proto.Pair{Key: nestedKey, Value: v})
			}

			parseStruct(&nestedStruct, nestedSentence)
			field.Set(nestedStruct)

		default:
			log.Printf("[WARN] Unsupported target field type while unmarshaling: %s", fieldType.Type.Kind())
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
