// Code generated by smithy-go-codegen DO NOT EDIT.

package types

import (
	"fmt"
	smithy "github.com/aws/smithy-go"
)

// The resource with the name requested already exists.
type AlreadyExistsException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *AlreadyExistsException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *AlreadyExistsException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *AlreadyExistsException) ErrorCode() string             { return "AlreadyExistsException" }
func (e *AlreadyExistsException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// An error occurred during a CloudFormation registry operation.
type CFNRegistryException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *CFNRegistryException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *CFNRegistryException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *CFNRegistryException) ErrorCode() string             { return "CFNRegistryException" }
func (e *CFNRegistryException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The specified change set name or ID doesn't exit. To view valid change sets for
// a stack, use the ListChangeSets operation.
type ChangeSetNotFoundException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *ChangeSetNotFoundException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *ChangeSetNotFoundException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *ChangeSetNotFoundException) ErrorCode() string             { return "ChangeSetNotFound" }
func (e *ChangeSetNotFoundException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The specified resource exists, but has been changed.
type CreatedButModifiedException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *CreatedButModifiedException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *CreatedButModifiedException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *CreatedButModifiedException) ErrorCode() string             { return "CreatedButModifiedException" }
func (e *CreatedButModifiedException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The template contains resources with capabilities that weren't specified in the
// Capabilities parameter.
type InsufficientCapabilitiesException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *InsufficientCapabilitiesException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *InsufficientCapabilitiesException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *InsufficientCapabilitiesException) ErrorCode() string {
	return "InsufficientCapabilitiesException"
}
func (e *InsufficientCapabilitiesException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The specified change set can't be used to update the stack. For example, the
// change set status might be CREATE_IN_PROGRESS, or the stack status might be
// UPDATE_IN_PROGRESS.
type InvalidChangeSetStatusException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *InvalidChangeSetStatusException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *InvalidChangeSetStatusException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *InvalidChangeSetStatusException) ErrorCode() string             { return "InvalidChangeSetStatus" }
func (e *InvalidChangeSetStatusException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The specified operation isn't valid.
type InvalidOperationException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *InvalidOperationException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *InvalidOperationException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *InvalidOperationException) ErrorCode() string             { return "InvalidOperationException" }
func (e *InvalidOperationException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// Error reserved for use by the CloudFormation CLI
// (https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/what-is-cloudformation-cli.html).
// CloudFormation doesn't return this error to users.
type InvalidStateTransitionException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *InvalidStateTransitionException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *InvalidStateTransitionException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *InvalidStateTransitionException) ErrorCode() string             { return "InvalidStateTransition" }
func (e *InvalidStateTransitionException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The quota for the resource has already been reached. For information on resource
// and stack limitations, see CloudFormation quotas
// (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cloudformation-limits.html)
// in the CloudFormation User Guide.
type LimitExceededException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *LimitExceededException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *LimitExceededException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *LimitExceededException) ErrorCode() string             { return "LimitExceededException" }
func (e *LimitExceededException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The specified name is already in use.
type NameAlreadyExistsException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *NameAlreadyExistsException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *NameAlreadyExistsException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *NameAlreadyExistsException) ErrorCode() string             { return "NameAlreadyExistsException" }
func (e *NameAlreadyExistsException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The specified operation ID already exists.
type OperationIdAlreadyExistsException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *OperationIdAlreadyExistsException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *OperationIdAlreadyExistsException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *OperationIdAlreadyExistsException) ErrorCode() string {
	return "OperationIdAlreadyExistsException"
}
func (e *OperationIdAlreadyExistsException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// Another operation is currently in progress for this stack set. Only one
// operation can be performed for a stack set at a given time.
type OperationInProgressException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *OperationInProgressException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *OperationInProgressException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *OperationInProgressException) ErrorCode() string             { return "OperationInProgressException" }
func (e *OperationInProgressException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The specified ID refers to an operation that doesn't exist.
type OperationNotFoundException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *OperationNotFoundException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *OperationNotFoundException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *OperationNotFoundException) ErrorCode() string             { return "OperationNotFoundException" }
func (e *OperationNotFoundException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// Error reserved for use by the CloudFormation CLI
// (https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/what-is-cloudformation-cli.html).
// CloudFormation doesn't return this error to users.
type OperationStatusCheckFailedException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *OperationStatusCheckFailedException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *OperationStatusCheckFailedException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *OperationStatusCheckFailedException) ErrorCode() string { return "ConditionalCheckFailed" }
func (e *OperationStatusCheckFailedException) ErrorFault() smithy.ErrorFault {
	return smithy.FaultClient
}

// The specified stack instance doesn't exist.
type StackInstanceNotFoundException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *StackInstanceNotFoundException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *StackInstanceNotFoundException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *StackInstanceNotFoundException) ErrorCode() string             { return "StackInstanceNotFoundException" }
func (e *StackInstanceNotFoundException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The specified stack ARN doesn't exist or stack doesn't exist corresponding to
// the ARN in input.
type StackNotFoundException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *StackNotFoundException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *StackNotFoundException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *StackNotFoundException) ErrorCode() string             { return "StackNotFoundException" }
func (e *StackNotFoundException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// You can't yet delete this stack set, because it still contains one or more stack
// instances. Delete all stack instances from the stack set before deleting the
// stack set.
type StackSetNotEmptyException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *StackSetNotEmptyException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *StackSetNotEmptyException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *StackSetNotEmptyException) ErrorCode() string             { return "StackSetNotEmptyException" }
func (e *StackSetNotEmptyException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The specified stack set doesn't exist.
type StackSetNotFoundException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *StackSetNotFoundException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *StackSetNotFoundException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *StackSetNotFoundException) ErrorCode() string             { return "StackSetNotFoundException" }
func (e *StackSetNotFoundException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// Another operation has been performed on this stack set since the specified
// operation was performed.
type StaleRequestException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *StaleRequestException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *StaleRequestException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *StaleRequestException) ErrorCode() string             { return "StaleRequestException" }
func (e *StaleRequestException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// A client request token already exists.
type TokenAlreadyExistsException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *TokenAlreadyExistsException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *TokenAlreadyExistsException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *TokenAlreadyExistsException) ErrorCode() string             { return "TokenAlreadyExistsException" }
func (e *TokenAlreadyExistsException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

// The specified extension configuration can't be found.
type TypeConfigurationNotFoundException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *TypeConfigurationNotFoundException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *TypeConfigurationNotFoundException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *TypeConfigurationNotFoundException) ErrorCode() string {
	return "TypeConfigurationNotFoundException"
}
func (e *TypeConfigurationNotFoundException) ErrorFault() smithy.ErrorFault {
	return smithy.FaultClient
}

// The specified extension doesn't exist in the CloudFormation registry.
type TypeNotFoundException struct {
	Message *string

	noSmithyDocumentSerde
}

func (e *TypeNotFoundException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *TypeNotFoundException) ErrorMessage() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
func (e *TypeNotFoundException) ErrorCode() string             { return "TypeNotFoundException" }
func (e *TypeNotFoundException) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }