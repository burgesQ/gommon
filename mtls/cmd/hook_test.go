package cmd

import (
	"reflect"
	"testing"

	"github.com/burgesQ/gommon/mtls"
	"github.com/stretchr/testify/require"
)

func TestStringToLevelHookFunc(t *testing.T) {
	t.Log("classic processing")
	{
		l, e := stringToLevelHookFunc(
			reflect.TypeOf("level"),
			reflect.TypeOf(mtls.RequireAnyClientCert),
			"hardAndSAN")
		require.Nil(t, e)
		require.Equal(t, mtls.RequireAndVerifyClientCertAndSAN, l)
	}

	t.Log("failed processing (1)")
	{
		l, e := stringToLevelHookFunc(
			reflect.TypeOf(1),
			reflect.TypeOf(2),
			-3)
		require.Nil(t, e)
		require.Equal(t, -3, l)
	}

	t.Log("failed processing (2)")
	{
		l, e := stringToLevelHookFunc(
			reflect.TypeOf("level"),
			reflect.TypeOf(mtls.RequireAnyClientCert),
			-3)
		require.Nil(t, e)
		require.Equal(t, -3, l)
	}

	t.Log("failed processing (3)")
	{
		_, e := stringToLevelHookFunc(
			reflect.TypeOf("level"),
			reflect.TypeOf(mtls.RequireAnyClientCert),
			"abcd")
		require.NotNil(t, e)
	}

	t.Log("processing skippped because of no verif")
	{
		l, e := stringToLevelHookFunc(
			reflect.TypeOf("level"),
			reflect.TypeOf(mtls.RequireAnyClientCert),
			mtls.NoClientCert)
		require.Nil(t, e)
		require.Equal(t, mtls.NoClientCert, l)
	}
}
