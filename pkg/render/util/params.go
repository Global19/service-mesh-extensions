package util

import (
	"strconv"

	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

func ParamValueToString(value *v1.ParameterValue) string {
	switch t := value.GetType().(type) {
	case *v1.ParameterValue_BooleanValue:
		if t.BooleanValue {
			return "true"
		}
		return "false"
	case *v1.ParameterValue_DateValue:
		return t.DateValue.String()
	case *v1.ParameterValue_FloatValue:
		return strconv.FormatFloat(t.FloatValue, 'f', -1, 64)
	case *v1.ParameterValue_IntValue:
		return strconv.Itoa(int(t.IntValue))
	case *v1.ParameterValue_SecretValue:
		// TODO not yet supported
		return ""
	case *v1.ParameterValue_StringValue:
		return t.StringValue
	}

	return ""
}