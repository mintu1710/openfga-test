package errors

import (
	"net/http"
	"testing"

	openfgav1pb "go.buf.build/openfga/go/openfga/api/openfga/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEncodedError(t *testing.T) {
	type encodedTests struct {
		_name                  string
		errorCode              int32
		message                string
		expectedCode           int
		expectedCodeString     string
		expectedHTTPStatusCode int
		isValidEncodedError    bool
	}
	var tests = []encodedTests{
		{
			_name:                  "invalid error",
			errorCode:              20,
			message:                "error message",
			expectedHTTPStatusCode: http.StatusInternalServerError,
			expectedCode:           20,
			expectedCodeString:     "20",
			isValidEncodedError:    false,
		},
		{
			_name:                  "validation error",
			errorCode:              int32(openfgav1pb.ErrorCode_validation_error),
			message:                "error message",
			expectedHTTPStatusCode: http.StatusBadRequest,
			expectedCode:           2000,
			expectedCodeString:     "validation_error",
			isValidEncodedError:    true,
		},
		{
			_name:                  "internal error",
			errorCode:              int32(openfgav1pb.InternalErrorCode_internal_error),
			message:                "error message",
			expectedHTTPStatusCode: http.StatusInternalServerError,
			expectedCode:           4000,
			expectedCodeString:     "internal_error",
			isValidEncodedError:    true,
		},
		{
			_name:                  "undefined endpoint",
			errorCode:              int32(openfgav1pb.NotFoundErrorCode_undefined_endpoint),
			message:                "error message",
			expectedHTTPStatusCode: http.StatusNotFound,
			expectedCode:           5000,
			expectedCodeString:     "undefined_endpoint",
			isValidEncodedError:    true,
		},
	}
	for _, test := range tests {
		t.Run(test._name, func(t *testing.T) {
			actualError := NewEncodedError(test.errorCode, test.message)
			if actualError.HTTPStatusCode != test.expectedHTTPStatusCode {
				t.Errorf("[%s]: http status code expect %d: actual %d", test._name, test.expectedHTTPStatusCode, actualError.HTTPStatusCode)
			}

			if actualError.Code() != test.expectedCodeString {
				t.Errorf("[%s]: code string expect %s: actual %s", test._name, test.expectedCodeString, actualError.Code())
			}

			if actualError.CodeValue() != int32(test.expectedCode) {
				t.Errorf("[%s]: code expect %d: actual %d", test._name, test.expectedCode, actualError.CodeValue())
			}

			if IsValidEncodedError(actualError.CodeValue()) != test.isValidEncodedError {
				t.Errorf("[%s]: expect is valid error %v: actual %v", test._name, test.isValidEncodedError, IsValidEncodedError(actualError.CodeValue()))
			}
		})
	}
}

func TestConvertToEncodedErrorCode(t *testing.T) {
	type encodedTests struct {
		_name             string
		status            *status.Status
		expectedErrorCode int32
	}
	var tests = []encodedTests{
		{
			_name:             "normal code",
			status:            status.New(codes.Code(openfgav1pb.ErrorCode_validation_error), "other error"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_validation_error),
		},
		{
			_name:             "no error",
			status:            status.New(codes.OK, "other error"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_no_error),
		},
		{
			_name:             "cancelled",
			status:            status.New(codes.Canceled, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_cancelled),
		},
		{
			_name:             "unknown",
			status:            status.New(codes.Unknown, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_internal_error),
		},
		{
			_name:             "deadline exceeded",
			status:            status.New(codes.DeadlineExceeded, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_deadline_exceeded),
		},
		{
			_name:             "not found",
			status:            status.New(codes.NotFound, "other error"),
			expectedErrorCode: int32(openfgav1pb.NotFoundErrorCode_undefined_endpoint),
		},
		{
			_name:             "already exceed",
			status:            status.New(codes.AlreadyExists, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_already_exists),
		},
		{
			_name:             "resource exhausted",
			status:            status.New(codes.ResourceExhausted, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_resource_exhausted),
		},
		{
			_name:             "failed precondition",
			status:            status.New(codes.FailedPrecondition, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_failed_precondition),
		},
		{
			_name:             "aborted",
			status:            status.New(codes.Aborted, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_aborted),
		},
		{
			_name:             "out of range",
			status:            status.New(codes.OutOfRange, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_out_of_range),
		},
		{
			_name:             "unimplemented",
			status:            status.New(codes.OutOfRange, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_out_of_range),
		},
		{
			_name:             "internal",
			status:            status.New(codes.Internal, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_internal_error),
		},
		{
			_name:             "unavailable",
			status:            status.New(codes.Unavailable, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_unavailable),
		},
		{
			_name:             "data loss",
			status:            status.New(codes.DataLoss, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_data_loss),
		},
		{
			_name:             "undefined error number",
			status:            status.New(25, "other error"),
			expectedErrorCode: int32(openfgav1pb.InternalErrorCode_internal_error),
		},
		{
			_name:             "invalid argument - unknown format",
			status:            status.New(codes.InvalidArgument, "other error"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_validation_error),
		},
		{
			_name:             "invalid argument - unknown format (2)",
			status:            status.New(codes.InvalidArgument, "no dot | foo :other error"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_validation_error),
		},
		{
			_name:             "invalid argument - unknown format (3)",
			status:            status.New(codes.InvalidArgument, "| foo :other error"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_validation_error),
		},
		{
			_name:             "invalid argument - unknown type",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.UnknowObject: value must be absolute"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_validation_error),
		},
		{
			_name:             "invalid argument - store id",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.StoreId: value length must be less than 26 runes"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_store_id_invalid_length),
		},
		{
			_name:             "invalid argument - issuer url other error",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.IssuerUrl: other error"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_validation_error),
		},
		{
			_name:             "invalid argument - Assertions",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Assertions: value must contain no more than 40 runes"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_assertions_too_many_items),
		},
		{
			_name:             "invalid argument - AuthorizationModelId",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.AuthorizationModelId: value length must be at most 40 runes"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_authorization_model_id_too_long),
		},
		{
			_name:             "invalid argument - Base",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Base: value is required"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_difference_base_missing_value),
		},
		{
			_name:             "invalid argument - Id",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Id: value length must be at most 40 runes"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_id_too_long),
		},
		{
			_name:             "invalid argument - Object length",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Object: value length must be at most 256 bytes"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_object_too_long),
		},
		{
			_name:             "invalid argument - Object invalid pattern",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Object: value does not match regex pattern"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_object_invalid_pattern),
		},
		{
			_name:             "invalid argument - PageSize",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.PageSize: value must be inside range 1 to 100"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_page_size_invalid),
		},
		{
			_name:             "invalid argument - Params",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Params: value is required"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_param_missing_value),
		},
		{
			_name:             "invalid argument - Relation",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Relation: value length must be at most 50 bytes"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_relation_too_long),
		},
		{
			_name:             "invalid argument - Relations",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Relations: value must contain at least 1 pair"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_relations_too_few_items),
		},
		{
			_name:             "invalid argument - Relations[abc]",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Relations[abc]: value length must be at most 50 bytes"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_relations_too_long),
		},
		{
			_name:             "invalid argument - Relations[abc]",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Relations[abc]: value does not match regex pattern"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_relations_invalid_pattern),
		},
		{
			_name:             "invalid argument - Subtract",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Subtract: value is required"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_subtract_base_missing_value),
		},
		{
			_name:             "invalid argument - TupleKey",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.TupleKey: value is required"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_tuple_key_value_not_specified),
		},
		{
			_name:             "invalid argument - TupleKeys",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.TupleKeys: value must contain between 1 to 10 items"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_tuple_keys_too_many_or_too_few_items),
		},
		{
			_name:             "invalid argument - Type length at least",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Type: value length must be at least 1"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_type_invalid_length),
		},
		{
			_name:             "invalid argument - Type length at most",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Type: value length must be at most 254 bytes"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_type_invalid_length),
		},
		{
			_name:             "invalid argument - Type regex",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.Type: value does not match regex pattern \"^[^:#]*$\""),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_type_invalid_pattern),
		},
		{
			_name:             "invalid argument - TypeDefinitions",
			status:            status.New(codes.InvalidArgument, "invalid WriteTokenIssuersRequest.Params: embedded message failed validation | caused by: invalid WriteTokenIssuersRequestParams.TypeDefinitions: value must contain at least 1 item"),
			expectedErrorCode: int32(openfgav1pb.ErrorCode_type_definitions_too_few_items),
		},
	}
Tests:

	for _, test := range tests {
		code := ConvertToEncodedErrorCode(test.status)
		if code != test.expectedErrorCode {
			t.Errorf("[%s]: Expect error code %d actual %d", test._name, test.expectedErrorCode, code)
			continue Tests
		}
	}

}

func TestSanitizeErrorMessage(t *testing.T) {

	got := sanitizedMessage(`proto: (line 1:2): unknown field "foo"`) // uses a whitespace rune of U+00a0 (see https://pkg.go.dev/unicode#IsSpace)
	expected := `(line 1:2): unknown field "foo"`
	if got != expected {
		t.Errorf("expected '%s', but got '%s'", expected, got)
	}
}
