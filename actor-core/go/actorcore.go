package actorcore

import (
	msgpack "github.com/wapc/tinygo-msgpack"
	wapc "github.com/wapc/wapc-guest-tinygo"
)

type Handlers struct {
	HealthRequest func(request HealthCheckRequest) (HealthCheckResponse, error)
}

func (h Handlers) Register() {
	if h.HealthRequest != nil {
		HealthRequestHandler = h.HealthRequest
		wapc.RegisterFunction("HealthRequest", HealthRequestWrapper)
	}
}

var (
	HealthRequestHandler func(request HealthCheckRequest) (HealthCheckResponse, error)
)

func HealthRequestWrapper(payload []byte) ([]byte, error) {
	decoder := msgpack.NewDecoder(payload)
	var request HealthCheckRequest
	request.Decode(&decoder)
	response, err := HealthRequestHandler(request)
	if err != nil {
		return nil, err
	}
	return msgpack.ToBytes(&response)
}

type CapabilityConfiguration struct {
	Module string
	Values map[string]string
}

func DecodeCapabilityConfigurationNullable(decoder *msgpack.Decoder) (*CapabilityConfiguration, error) {
	if isNil, err := decoder.IsNextNil(); isNil || err != nil {
		return nil, err
	}
	decoded, err := DecodeCapabilityConfiguration(decoder)
	return &decoded, err
}

func DecodeCapabilityConfiguration(decoder *msgpack.Decoder) (CapabilityConfiguration, error) {
	var o CapabilityConfiguration
	err := o.Decode(decoder)
	return o, err
}

func (o *CapabilityConfiguration) Decode(decoder *msgpack.Decoder) error {
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
		case "module":
			o.Module, err = decoder.ReadString()
		case "values":
			mapSize, err := decoder.ReadMapSize()
			if err != nil {
				return err
			}
			o.Values = make(map[string]string, mapSize)
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
				o.Values[key] = value
			}
		default:
			err = decoder.Skip()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *CapabilityConfiguration) Encode(encoder msgpack.Writer) error {
	if o == nil {
		encoder.WriteNil()
		return nil
	}
	encoder.WriteMapSize(2)
	encoder.WriteString("module")
	encoder.WriteString(o.Module)
	encoder.WriteString("values")
	encoder.WriteMapSize(uint32(len(o.Values)))
	if o.Values != nil { // TinyGo bug: ranging over nil maps panics.
		for k, v := range o.Values {
			encoder.WriteString(k)
			encoder.WriteString(v)
		}
	}

	return nil
}

type HealthCheckRequest struct {
	Placeholder bool
}

func DecodeHealthCheckRequestNullable(decoder *msgpack.Decoder) (*HealthCheckRequest, error) {
	if isNil, err := decoder.IsNextNil(); isNil || err != nil {
		return nil, err
	}
	decoded, err := DecodeHealthCheckRequest(decoder)
	return &decoded, err
}

func DecodeHealthCheckRequest(decoder *msgpack.Decoder) (HealthCheckRequest, error) {
	var o HealthCheckRequest
	err := o.Decode(decoder)
	return o, err
}

func (o *HealthCheckRequest) Decode(decoder *msgpack.Decoder) error {
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
		case "placeholder":
			o.Placeholder, err = decoder.ReadBool()
		default:
			err = decoder.Skip()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *HealthCheckRequest) Encode(encoder msgpack.Writer) error {
	if o == nil {
		encoder.WriteNil()
		return nil
	}
	encoder.WriteMapSize(1)
	encoder.WriteString("placeholder")
	encoder.WriteBool(o.Placeholder)

	return nil
}

type HealthCheckResponse struct {
	Healthy bool
	Message string
}

func DecodeHealthCheckResponseNullable(decoder *msgpack.Decoder) (*HealthCheckResponse, error) {
	if isNil, err := decoder.IsNextNil(); isNil || err != nil {
		return nil, err
	}
	decoded, err := DecodeHealthCheckResponse(decoder)
	return &decoded, err
}

func DecodeHealthCheckResponse(decoder *msgpack.Decoder) (HealthCheckResponse, error) {
	var o HealthCheckResponse
	err := o.Decode(decoder)
	return o, err
}

func (o *HealthCheckResponse) Decode(decoder *msgpack.Decoder) error {
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
		case "healthy":
			o.Healthy, err = decoder.ReadBool()
		case "message":
			o.Message, err = decoder.ReadString()
		default:
			err = decoder.Skip()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *HealthCheckResponse) Encode(encoder msgpack.Writer) error {
	if o == nil {
		encoder.WriteNil()
		return nil
	}
	encoder.WriteMapSize(2)
	encoder.WriteString("healthy")
	encoder.WriteBool(o.Healthy)
	encoder.WriteString("message")
	encoder.WriteString(o.Message)

	return nil
}
