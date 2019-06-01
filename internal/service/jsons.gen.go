// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package service

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson5822a332DecodeWaveInternalService(in *jlexer.Lexer, out *FServiceConfig) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "root":
			out.Root = string(in.String())
		case "log":
			out.Log = string(in.String())
		case "configs":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Configs = make(map[string]string)
				} else {
					out.Configs = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 string
					v1 = string(in.String())
					(out.Configs)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5822a332EncodeWaveInternalService(out *jwriter.Writer, in FServiceConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"root\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Root))
	}
	{
		const prefix string = ",\"log\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Log))
	}
	{
		const prefix string = ",\"configs\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Configs == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v2First := true
			for v2Name, v2Value := range in.Configs {
				if v2First {
					v2First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v2Name))
				out.RawByte(':')
				out.String(string(v2Value))
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FServiceConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5822a332EncodeWaveInternalService(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FServiceConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5822a332EncodeWaveInternalService(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FServiceConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5822a332DecodeWaveInternalService(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FServiceConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5822a332DecodeWaveInternalService(l, v)
}