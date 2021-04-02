package httpclient

import (
	msgpack "github.com/wapc/tinygo-msgpack"
	wapc "github.com/wapc/wapc-guest-tinygo"
)

type Host struct {
	binding string
}

func NewHost(binding string) *Host {
	return &Host{
		binding: binding,
	}
}

func (h *Host) Request(method string, url string, headers map[string]string, body []byte) (Response, error) {
	inputArgs := RequestArgs{
		Method:  method,
		Url:     url,
		Headers: headers,
		Body:    body,
	}
	inputBytes, err := msgpack.ToBytes(&inputArgs)
	if err != nil {
		return Response{}, err
	}
	payload, err := wapc.HostCall(
		h.binding,
		"wasmcloud:httpclient",
		"Request",
		inputBytes,
	)
	if err != nil {
		return Response{}, err
	}
	decoder := msgpack.NewDecoder(payload)
	return DecodeResponse(&decoder)
}

type Handlers struct {
	Request func(method string, url string, headers map[string]string, body []byte) (Response, error)
}

func (h Handlers) Register() {
	if h.Request != nil {
		RequestHandler = h.Request
		wapc.RegisterFunction("Request", RequestWrapper)
	}
}

var (
	RequestHandler func(method string, url string, headers map[string]string, body []byte) (Response, error)
)

func RequestWrapper(payload []byte) ([]byte, error) {
	decoder := msgpack.NewDecoder(payload)
	var inputArgs RequestArgs
	inputArgs.Decode(&decoder)
	response, err := RequestHandler(inputArgs.Method, inputArgs.Url, inputArgs.Headers, inputArgs.Body)
	if err != nil {
		return nil, err
	}
	return msgpack.ToBytes(&response)
}

type RequestArgs struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    []byte
}

func DecodeRequestArgsNullable(decoder *msgpack.Decoder) (*RequestArgs, error) {
	if isNil, err := decoder.IsNextNil(); isNil || err != nil {
		return nil, err
	}
	decoded, err := DecodeRequestArgs(decoder)
	return &decoded, err
}

func DecodeRequestArgs(decoder *msgpack.Decoder) (RequestArgs, error) {
	var o RequestArgs
	err := o.Decode(decoder)
	return o, err
}

func (o *RequestArgs) Decode(decoder *msgpack.Decoder) error {
	numFields, err := decoder.ReadMapSize()
	if err != nil {
		return err
	}

	for numFields > 0 {
		numFields--
		field, err := decoder.ReadString()
		if err != nil {
			return err
		}
		switch field {
		case "method":
			o.Method, err = decoder.ReadString()
		case "url":
			o.Url, err = decoder.ReadString()
		case "headers":
			mapSize, err := decoder.ReadMapSize()
			if err != nil {
				return err
			}
			o.Headers = make(map[string]string, mapSize)
			for mapSize > 0 {
				mapSize--
				key, err := decoder.ReadString()
				if err != nil {
					return err
				}
				value, err := decoder.ReadString()
				if err != nil {
					return err
				}
				o.Headers[key] = value
			}
		case "body":
			o.Body, err = decoder.ReadByteArray()
		default:
			err = decoder.Skip()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *RequestArgs) Encode(encoder msgpack.Writer) error {
	if o == nil {
		encoder.WriteNil()
		return nil
	}
	encoder.WriteMapSize(4)
	encoder.WriteString("method")
	encoder.WriteString(o.Method)
	encoder.WriteString("url")
	encoder.WriteString(o.Url)
	encoder.WriteString("headers")
	encoder.WriteMapSize(uint32(len(o.Headers)))
	if o.Headers != nil { // TinyGo bug: ranging over nil maps panics.
		for k, v := range o.Headers {
			encoder.WriteString(k)
			encoder.WriteString(v)
		}
	}
	encoder.WriteString("body")
	encoder.WriteByteArray(o.Body)

	return nil
}

type Response struct {
	StatusCode uint32
	Status     string
	Header     map[string]string
	Body       []byte
}

func DecodeResponseNullable(decoder *msgpack.Decoder) (*Response, error) {
	if isNil, err := decoder.IsNextNil(); isNil || err != nil {
		return nil, err
	}
	decoded, err := DecodeResponse(decoder)
	return &decoded, err
}

func DecodeResponse(decoder *msgpack.Decoder) (Response, error) {
	var o Response
	err := o.Decode(decoder)
	return o, err
}

func (o *Response) Decode(decoder *msgpack.Decoder) error {
	numFields, err := decoder.ReadMapSize()
	if err != nil {
		return err
	}

	for numFields > 0 {
		numFields--
		field, err := decoder.ReadString()
		if err != nil {
			return err
		}
		switch field {
		case "statusCode":
			o.StatusCode, err = decoder.ReadUint32()
		case "status":
			o.Status, err = decoder.ReadString()
		case "header":
			mapSize, err := decoder.ReadMapSize()
			if err != nil {
				return err
			}
			o.Header = make(map[string]string, mapSize)
			for mapSize > 0 {
				mapSize--
				key, err := decoder.ReadString()
				if err != nil {
					return err
				}
				value, err := decoder.ReadString()
				if err != nil {
					return err
				}
				o.Header[key] = value
			}
		case "body":
			o.Body, err = decoder.ReadByteArray()
		default:
			err = decoder.Skip()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *Response) Encode(encoder msgpack.Writer) error {
	if o == nil {
		encoder.WriteNil()
		return nil
	}
	encoder.WriteMapSize(4)
	encoder.WriteString("statusCode")
	encoder.WriteUint32(o.StatusCode)
	encoder.WriteString("status")
	encoder.WriteString(o.Status)
	encoder.WriteString("header")
	encoder.WriteMapSize(uint32(len(o.Header)))
	if o.Header != nil { // TinyGo bug: ranging over nil maps panics.
		for k, v := range o.Header {
			encoder.WriteString(k)
			encoder.WriteString(v)
		}
	}
	encoder.WriteString("body")
	encoder.WriteByteArray(o.Body)

	return nil
}
