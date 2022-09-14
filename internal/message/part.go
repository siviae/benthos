package message

import (
	"context"
	"os"
)

var useNumber = true

func init() {
	if os.Getenv("BENTHOS_USE_NUMBER") == "false" {
		useNumber = false
	}
}

//------------------------------------------------------------------------------

// Part represents a single Benthos message.
type Part struct {
	data *messageData
	ctx  context.Context
}

// NewPart initializes a new message part.
func NewPart(data []byte) *Part {
	return &Part{
		data: newMessageBytes(data),
		ctx:  context.Background(),
	}
}

//------------------------------------------------------------------------------

// ShallowCopy creates a shallow copy of the message part.
func (p *Part) ShallowCopy() *Part {
	return &Part{
		data: p.data.ShallowCopy(),
		ctx:  p.ctx,
	}
}

// DeepCopy creates a new deep copy of the message part.
func (p *Part) DeepCopy() *Part {
	return &Part{
		data: p.data.DeepCopy(),
		ctx:  p.ctx,
	}
}

//------------------------------------------------------------------------------

// GetContext either returns a context attached to the message part, or
// context.Background() if one hasn't been previously attached.
func GetContext(p *Part) context.Context {
	return p.ctx
}

// WithContext returns the same message part wrapped with a context, this
// context can subsequently be received with GetContext.
func WithContext(ctx context.Context, p *Part) *Part {
	return p.WithContext(ctx)
}

// GetContext returns the underlying context attached to this message part.
func (p *Part) GetContext() context.Context {
	return p.ctx
}

// WithContext returns the underlying message part with a different context
// attached.
func (p *Part) WithContext(ctx context.Context) *Part {
	newP := *p
	newP.ctx = ctx
	return &newP
}

//------------------------------------------------------------------------------

// ErrorGet returns an error associated with the message, or nil if none exists.
func (p *Part) ErrorGet() error {
	return p.data.ErrorGet()
}

// ErrorSet modifies the error associated with a message. Errors attached to
// messages are used to indicate that processing has failed at some point in the
// processing pipeline.
func (p *Part) ErrorSet(err error) {
	p.data.ErrorSet(err)
}

// AsBytes returns the body of the message part.
func (p *Part) AsBytes() []byte {
	return p.data.AsBytes()
}

// AsStructuredMut returns the structured format of the message if already set,
// or attempts to parse the raw bytes as a JSON document if not. The returned
// structure is mutable and therefore safe to mutate directly.
func (p *Part) AsStructuredMut() (any, error) {
	return p.data.AsStructuredMut()
}

// AsStructured returns the structured format of the message if already set, or
// attempts to parse the raw bytes as a JSON document if not. The returned
// structure should be considered read-only and therefore not be mutated.
func (p *Part) AsStructured() (any, error) {
	return p.data.AsStructured()
}

// SetBytes the value of the message part as a raw byte slice.
func (p *Part) SetBytes(data []byte) *Part {
	p.data.SetBytes(data)
	return p
}

// SetStructuredMut sets the value of the message to a structured value, this
// value is mutable and subsequent mutations will be performed directly on the
// provided data.
func (p *Part) SetStructuredMut(jObj any) {
	p.data.SetStructuredMut(jObj)
}

// SetStructured sets the value of the message to a structured value, this
// value is read-only and subsequent mutations will require cloning of the
// entire data structure.
func (p *Part) SetStructured(jObj any) {
	p.data.SetStructured(jObj)
}

//------------------------------------------------------------------------------

// MetaGet returns a metadata value if a key exists, otherwise an empty string.
func (p *Part) MetaGet(key string) string {
	v, _ := p.data.MetaGet(key)
	return v
}

// MetaSet sets the value of a metadata key.
func (p *Part) MetaSet(key, value string) {
	p.data.MetaSet(key, value)
}

// MetaDelete removes the value of a metadata key.
func (p *Part) MetaDelete(key string) {
	p.data.MetaDelete(key)
}

// MetaIter iterates each metadata key/value pair.
func (p *Part) MetaIter(f func(k, v string) error) error {
	return p.data.MetaIter(f)
}

//------------------------------------------------------------------------------

// IsEmpty returns true if the message part is empty.
func (p *Part) IsEmpty() bool {
	return p.data.IsEmpty()
}
