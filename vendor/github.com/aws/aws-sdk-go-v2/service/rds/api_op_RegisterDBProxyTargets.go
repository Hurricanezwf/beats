// Code generated by smithy-go-codegen DO NOT EDIT.

package rds

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Associate one or more DBProxyTarget data structures with a DBProxyTargetGroup.
func (c *Client) RegisterDBProxyTargets(ctx context.Context, params *RegisterDBProxyTargetsInput, optFns ...func(*Options)) (*RegisterDBProxyTargetsOutput, error) {
	if params == nil {
		params = &RegisterDBProxyTargetsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "RegisterDBProxyTargets", params, optFns, c.addOperationRegisterDBProxyTargetsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*RegisterDBProxyTargetsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type RegisterDBProxyTargetsInput struct {

	// The identifier of the DBProxy that is associated with the DBProxyTargetGroup.
	//
	// This member is required.
	DBProxyName *string

	// One or more DB cluster identifiers.
	DBClusterIdentifiers []string

	// One or more DB instance identifiers.
	DBInstanceIdentifiers []string

	// The identifier of the DBProxyTargetGroup.
	TargetGroupName *string

	noSmithyDocumentSerde
}

type RegisterDBProxyTargetsOutput struct {

	// One or more DBProxyTarget objects that are created when you register targets
	// with a target group.
	DBProxyTargets []types.DBProxyTarget

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationRegisterDBProxyTargetsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsquery_serializeOpRegisterDBProxyTargets{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpRegisterDBProxyTargets{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addOpRegisterDBProxyTargetsValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opRegisterDBProxyTargets(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opRegisterDBProxyTargets(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "rds",
		OperationName: "RegisterDBProxyTargets",
	}
}