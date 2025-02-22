package handler

import (
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/cfnerr"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/encoding"
)

const (
	// marshalingError occurs when we can't marshal data from one format into another.
	marshalingError = "Marshaling"

	// bodyEmptyError happens when the resource body is empty
	bodyEmptyError = "BodyEmpty"
)

// Request is passed to actions with customer related data
// such as resource states
type Request struct {
	// The logical ID of the resource in the CloudFormation stack
	LogicalResourceID string

	// The callback context is an arbitrary datum which the handler can return in an
	// IN_PROGRESS event to allow the passing through of additional state or
	// metadata between subsequent retries; for example to pass through a Resource
	// identifier which can be used to continue polling for stabilization
	CallbackContext map[string]interface{}

	// The RequestContext is information about the current
	// invocation.
	RequestContext RequestContext

	// An authenticated AWS session that can be used with the AWS Go SDK
	Session *session.Session

	previousResourcePropertiesBody []byte
	resourcePropertiesBody         []byte
	typeConfigurationBody          []byte
}

// RequestContext represents information about the current
// invocation request of the handler.
type RequestContext struct {
	// The stack ID of the CloudFormation stack
	StackID string

	// The Region of the requester
	Region string

	// The Account ID of the requester
	AccountID string

	// The stack tags associated with the cloudformation stack
	StackTags map[string]string

	// The SystemTags associated with the request
	SystemTags map[string]string

	// The NextToken provided in the request
	NextToken string
}

// NewRequest returns a new Request based on the provided parameters
func NewRequest(id string, ctx map[string]interface{}, requestCTX RequestContext, sess *session.Session, previousBody, body, typeConfig []byte) Request {
	return Request{
		LogicalResourceID:              id,
		CallbackContext:                ctx,
		Session:                        sess,
		previousResourcePropertiesBody: previousBody,
		resourcePropertiesBody:         body,
		RequestContext:                 requestCTX,
		typeConfigurationBody:          typeConfig,
	}
}

// UnmarshalPrevious populates the provided interface
// with the previous properties of the resource
func (r *Request) UnmarshalPrevious(v interface{}) error {
	if len(r.previousResourcePropertiesBody) == 0 {
		return nil
	}

	if err := encoding.Unmarshal(r.previousResourcePropertiesBody, v); err != nil {
		return cfnerr.New(marshalingError, "Unable to convert type", err)
	}

	return nil
}

// Unmarshal populates the provided interface
// with the current properties of the resource
func (r *Request) Unmarshal(v interface{}) error {
	if len(r.resourcePropertiesBody) == 0 {
		return cfnerr.New(bodyEmptyError, "Body is empty", nil)
	}

	if err := encoding.Unmarshal(r.resourcePropertiesBody, v); err != nil {
		return cfnerr.New(marshalingError, "Unable to convert type", err)
	}

	return nil
}

// UnmarshalTypeConfig populates the provided interface
// with the current properties of the model
func (r *Request) UnmarshalTypeConfig(v interface{}) error {
	if len(r.typeConfigurationBody) == 0 {
		return cfnerr.New(bodyEmptyError, "Type Config is empty", nil)
	}

	if err := encoding.Unmarshal(r.typeConfigurationBody, v); err != nil {
		return cfnerr.New(marshalingError, "Unable to convert type", err)
	}

	return nil
}
