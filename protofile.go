// Package protofile generates structured data to represent the contents of a .proto file.
package protofile

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// ProtoFile contains structured data read from a .proto file
type ProtoFile struct {
	syntax   string
	services []Service
	messages map[string]Message
}

// GetServices returns the services defined in the file
func (p *ProtoFile) GetServices() []Service {
	return p.services
}

// GetMessages returns the messages defined in the file
func (p *ProtoFile) GetMessages() map[string]Message {
	return p.messages
}

// Service represents a service block in a .proto file
type Service struct {
	name    string
	methods []Method
}

// GetName returns the name of the service
func (s *Service) GetName() string {
	return s.name
}

// GetMethods returns every method within the service
func (s *Service) GetMethods() []Method {
	return s.methods
}

// Method represents a method within a service
type Method struct {
	name        string
	req         Message
	isReqStream bool
	res         Message
	isResStream bool
}

// GetName returns the name of the method
func (m *Method) GetName() string {
	return m.name
}

// GetReqMessage returns the request message of the method
func (m *Method) GetReqMessage() Message {
	return m.req
}

// IsReqStreamed returns true if the request is a stream
func (m *Method) IsReqStreamed() bool {
	return m.isReqStream
}

// GetResMessage returns the response message of the method
func (m *Method) GetResMessage() Message {
	return m.res
}

// IsResStreamed returns true if the response is a stream
func (m *Method) IsResStreamed() bool {
	return m.isResStream
}

// Message represents a message block in a .proto file
type Message struct {
	name   string
	fields map[string]Field
}

// GetName returns the name of the message
func (m *Message) GetName() string {
	return m.name
}

// GetFields returns every field within the message
func (m *Message) GetFields() map[string]Field {
	return m.fields
}

// GetField returns a field within the message by name
func (m *Message) GetField(name string) (Field, error) {
	if field, exists := m.fields[name]; exists {
		return field, nil
	}
	return Field{}, fmt.Errorf("%s does not contain field %s", m.name, name)
}

// Field represents a field within a message
type Field struct {
	name      string
	fieldType string
	repeated  bool
	number    int
}

// GetName returns the name of the field
func (f *Field) GetName() string {
	return f.name
}

// GetFieldType returns the type of the field
func (f *Field) GetFieldType() string {
	return f.fieldType
}

// IsRepeated returns true if the field is repeated
func (f *Field) IsRepeated() bool {
	return f.repeated
}

// GetNumber returns the number of the field
func (f *Field) GetNumber() int {
	return f.number
}

// New creates a ProtoFile object from a .proto file
func New(name string) (*ProtoFile, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	pf := ProtoFile{}
	pf.messages = make(map[string]Message, 0)
	statements, err := makeStatementList(f)
	if err != nil {
		return nil, err
	}
	for _, statement := range statements {
		if strings.HasPrefix(statement, "service") {
			defer func(s string) {
				service := parseServiceStatement(s, &pf)
				pf.services = append(pf.services, service)
			}(statement)
		}
		if strings.HasPrefix(statement, "message") {
			message := parseMessageStatement(statement)
			pf.messages[message.name] = message
		}
	}
	return &pf, nil
}

func makeStatementList(f *os.File) ([]string, error) {
	statements := make([]string, 0)
	statement := make([]byte, 0)
	open := false
	for {
		buffer := make([]byte, 1)
		_, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		char := buffer[0]
		if char == 123 {
			open = true
		}
		if char == 125 {
			open = false
		}
		if char == 10 {
			if !open {
				if len(statement) > 0 {
					statements = append(statements, string(statement))
				}
				statement = make([]byte, 0)
			}
			continue
		}
		statement = append(statement, buffer[0])
	}
	return statements, nil
}

func parseServiceStatement(s string, pf *ProtoFile) Service {
	s = strings.Trim(s, "service ")
	service := Service{}
	service.name = strings.Split(s, " ")[0]
	methods := strings.Split(strings.Trim(s[strings.Index(s, "{"):], "{ }"), ";")
	for _, method := range methods {
		if method == "" {
			continue
		}
		service.methods = append(service.methods, parseMethodStatement(method, pf))
	}
	return service
}

func parseMethodStatement(s string, pf *ProtoFile) Method {
	s = strings.Trim(s, "rpc ")
	method := Method{}
	method.name = strings.Split(s, "(")[0]
	s = s[strings.Index(s, "("):]
	req := strings.Trim(strings.Split(s, " returns ")[0], "( )")
	res := strings.Trim(strings.Split(s, " returns ")[1], "( )")
	if strings.HasPrefix(req, "stream") {
		method.isReqStream = true
		req = strings.Trim(req, "stream ")
	}
	if strings.HasPrefix(res, "stream") {
		method.isResStream = true
		res = strings.Trim(res, "stream ")
	}
	method.req = pf.messages[req]
	method.res = pf.messages[res]
	return method
}

func parseMessageStatement(s string) Message {
	s = strings.Trim(s, "message ")
	name := strings.Split(s, " ")[0]
	message := Message{name: name, fields: map[string]Field{}}
	first := strings.Index(s, "{") + 1
	last := strings.Index(s, "}")
	s = s[first:last]
	fields := strings.Split(s, ";")
	for _, field := range fields {
		field = strings.Trim(field, " ")
		if len(field) == 0 {
			continue
		}
		f := parseFieldStatement(field)
		message.fields[f.name] = f
	}
	return message
}

func parseFieldStatement(s string) Field {
	field := Field{}
	words := strings.Split(s, " ")
	if strings.ToLower(words[0]) == "repeated" {
		field.repeated = true
		words = words[1:]
	}
	field.fieldType = words[0]
	field.name = words[1]
	num, _ := strconv.Atoi(words[3])
	field.number = num
	return field
}
