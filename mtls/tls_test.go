package mtls

import (
	"crypto/tls"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	_testKey = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDARJcfTSn8dIUD
6hqxy1aMN6u6IxHWpn561PBY8/2Tucc1zCc6wD4skkHc0NwW/rjdYOPvKUaTsyBP
kebE56Yo57la0rwAuhYFfIWNuY959fAbxgmrr1Ui5aQ6UsdAySe2g6Yuvzdbl/qJ
8wJU+Qz6DFe6sL3LrGwN+qAo7SKPQng9ayLFU5bHOy12hIZ4Xpou0wbk4iCEnSjt
0dwQiwMOfYfxJ/FomierMHEEexU0FiYKYO8ixBOAlxEI5rUcnVS19wJ5U2l6sgiU
T4lMmUYlDSnhbnZ6tUOcvfi8M/5Br6M+TviBQZF5T37cbkKlAESOG4x+EAAXQWO9
W/xIFx7RAgMBAAECggEBAJcZUGqZCdYM/DwHTtNLKckoNx0hUnOOhPejQP+nQDFI
XvY4atwRGMuOZZRTz6OCq2XPL1JU7dJFg68EyocURdD/yvtFgdpQY94L7iljGw6N
6RTqnc7/C7lRH692PoD7lOoqq/8w5yBkX9Dp2XtN/pTD/m4JWB8ocgLKY/PF6NML
2ri/YwCQuUnUqA9lfbWCk0sFWldnd40dhZTD29kt5sTeDhH2GQUAhezkNKzaW7k4
cPQYi+EVqz/SXfEcm8T7/uxl4e/9byLurmzrmAoN6rRaTZehqSh0fzzkR97T0W4F
H0k01PahToa5WYSxK4Lsccbk92N3zPQFSItmW/kkrcECgYEA+WkVc46+NGwTdOWK
H3AQrxVI7iToZXnc96wVTqPHRf2UWDjbrlkML8/JipuifYklSIS3jbxanhD1g+CJ
SyEXS0sUIynfLtQfX6IdA65XOVqwBs0l/e4m4Z/zc1F3EFd6qv7gDs7FlDMLl1Ox
pzxIwmMjwfsLZqFt7pzxy7ZL7AMCgYEAxVkEMIaR6mAs9dyzZgX/JF8prq0J9aGW
+81/3qMxU/7+PyXyLpEjMAQ6mGM0Uy0vzP9fGgQHQpcrEOr/HAmPuEHnz7GYLhVd
sAvBaKRaH4w6S0+I9sn2mwCC7CbDXq0Y1IxQ8KOAfy7YTVbmMvf8g+aeunwaY1KV
hb6Z3qiJE5sCgYEAyo8n1uQ1Ygnse2H5HbM8OZYF1zOucsvYRGZEH8wwCY37LvNu
p1i25xXQz3u7Kk16ND1lff1dc0a+v05a8uN7MbFWN4DIPBYXLOpSuiybtn8Ku1td
4a/LcC8h36RoGKOTgtDhU+Vm3gffABX/EJ2LUiSGZALprX6p88MPNa1mV9UCgYEA
utgx2EPAqRgP2WPw0nqA+43B6Cja0h4A1jzVgRQfYvh8/YrOxfoSR6bpV1gttUaG
CGAMSZRgz1JSqvzjNkdzNC/p60Go5JDEGCa5Inrg/ReGJcGS2p2TB2QvkKiOtvfK
F4sWIw+aXFAc6PSKlN0nzjYuOD/BuCH7gRpZkm8dAVsCgYBTkmWtVg7dQ08gLcCJ
514rj6E+t+pNMfPaW3wLanti4NjohIpRSj70BQ16LxidLRD9msOWvcpd5J9/IcYt
FCJIjSJbzpPGf1gI4ZyfRJseq63dnkl7y0rynbxgPoKZtwiMvmF+Us7UMWnahut8
yg2xRROWnrXk3Gaw6ObIXnWgCw==
-----END PRIVATE KEY-----`
	_testCert = `-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUdnElFQsX0wE0R/G7opUxIxs3OVcwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMDAzMDMxMjI5NDFaFw0yMDA0
MDIxMjI5NDFaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDARJcfTSn8dIUD6hqxy1aMN6u6IxHWpn561PBY8/2T
ucc1zCc6wD4skkHc0NwW/rjdYOPvKUaTsyBPkebE56Yo57la0rwAuhYFfIWNuY95
9fAbxgmrr1Ui5aQ6UsdAySe2g6Yuvzdbl/qJ8wJU+Qz6DFe6sL3LrGwN+qAo7SKP
Qng9ayLFU5bHOy12hIZ4Xpou0wbk4iCEnSjt0dwQiwMOfYfxJ/FomierMHEEexU0
FiYKYO8ixBOAlxEI5rUcnVS19wJ5U2l6sgiUT4lMmUYlDSnhbnZ6tUOcvfi8M/5B
r6M+TviBQZF5T37cbkKlAESOG4x+EAAXQWO9W/xIFx7RAgMBAAGjUzBRMB0GA1Ud
DgQWBBR0lCBYPajvpNjtewlbAfiLpgzSTDAfBgNVHSMEGDAWgBR0lCBYPajvpNjt
ewlbAfiLpgzSTDAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQC2
N9FOkjDpASOBCqWFiGzC66leVV95uSUdz1f6v5gXiS9gF0BA447frirdaQwlo0HN
fwesXgf8F1+krF/gr1iv27u7bTVmNlWXrAO8o4442fzEd00YVNfoIeV79O575ObC
2Nd9fC5xWySoiak1X6oifSKtvvUmqPLlx0aM2+bUj+8o1LAj/M3z2MIbGYyHmDWD
l5srq5Bf6Lr0pURlJ6yhu90WKtkDunrW/fE+rDe+0LDBAF0TVIfUmKtTY1hDsX9Q
IYb5AeKzS9fIeuGaK6p04twDoj23J+q8WZTuFvR8rZ4yngcniywPVfpNyWHiB09/
MDrT/Ly8s965M23chGxn
-----END CERTIFICATE-----`
	_testCA = `-----BEGIN CERTIFICATE-----
MIIFozCCA4ugAwIBAgIUXC87/X0rm0yEnKZ4lI7aiWpAmg0wDQYJKoZIhvcNAQEL
BQAwYTELMAkGA1UEBhMCRlIxDjAMBgNVBAgMBVBhcmlzMQ4wDAYDVQQHDAVQYXJp
czEPMA0GA1UECgwGRnJhZm9zMRMwEQYDVQQLDApGcmFmb3MgRGV2MQwwCgYDVQQD
DAN5b3AwHhcNMjIwOTE2MDgwODA4WhcNMjQxMTI0MDgwODA4WjBhMQswCQYDVQQG
EwJGUjEOMAwGA1UECAwFUGFyaXMxDjAMBgNVBAcMBVBhcmlzMQ8wDQYDVQQKDAZG
cmFmb3MxEzARBgNVBAsMCkZyYWZvcyBEZXYxDDAKBgNVBAMMA3lvcDCCAiIwDQYJ
KoZIhvcNAQEBBQADggIPADCCAgoCggIBAKX2KiPPWGRt4PwqdDmwRmsoxHk1wPo+
CfCh7GEolSd0twhBBfMMt+0rKij16Si/RXXovM1jDuOJaDgJ4YaTKuBluaEAzGsI
II7NEZIW4AeJi9MijePSi0/m2Hztl1ue3YjLV+6zLTITr2Nx6yqzCQCeKsnhT4eU
k17sR5mPInL0yAEIAwg2Bc8qeyFGdwAUpTVyQvGGis+mfoxjGp03kGYnwLvwNjkR
4lnY1qXjn0kCI9bmbKcHc8fTEWig+V5SGYPoBHeXGxFirSvRJ6VEYah0wFYk3hLZ
EBALmb9qs70VUzoqK85ArtpY+sbfOzxNiD6VVVm4JkshUXF5UB4T4yWM9f7hgcuf
gR1ytaiYLoDk1jWGdw0tDhsO2VYGgf7jL56FPHzCOHnvpTWKGB2zVsEYDbGEloJY
3TQgXeizxXCqy6tanFnAOg7JRBwhicqFyvD1bWZSO2J+Q01qah09QuntB2e0aGBH
5RLR0Dzed3+x9my9echWFHHRpXGQqRV08x1XOkvVY1SeKwVVHr+1F3w4W08WYX8I
NJJ3JJXXoCfcHLXWWM9gDR2WP8xc2fkaaC+9UOlnjxVLl7J5L4togxH9Kwd8bTnq
GMwcbfhnW9MGneeZFsABmmOQEtqrq/wx17XEO2nkVF+M1oc+5hkfChqwT9GDJTcI
EIJ0sgQYJkPPAgMBAAGjUzBRMB0GA1UdDgQWBBQWAdY9KxT5OJeRVqDfENYgBirv
djAfBgNVHSMEGDAWgBQWAdY9KxT5OJeRVqDfENYgBirvdjAPBgNVHRMBAf8EBTAD
AQH/MA0GCSqGSIb3DQEBCwUAA4ICAQCIoEuZIDC+Ktvk93HkTuzPljE65is4wliE
3coGU9hxhGFkh9PMTyu/3r6Pab5MyFLqQX6qyqHo0zlEP4/Mf8P8fZaLwilZqhnB
FHPeXhECIpcItxUfHuWKewgwy/EEn8mk6Uaq6XO4XHX9tLLjuRCxw6EJeZmrgEPk
4mIEB30qoQ+mmMaSrOIwp0eUQ7qu5e/pOK09WK+zo77InOks5RY9nmGqaLvQo7rj
41PlPYKDvHtEcJaNIHo6Xa0Q9oi2lMqfho8y4M2wHokFisCFYJ/8GI8WNtwXnPlM
SSOYF4PMIzhzuujDpz1RbL1y2DXBII75Mw52BNDjoDoGd/C1O+5fbcRGl1zSFuLL
N5yvZH+xrFx2nF0FovuiyrZAR8y6/sR8wuTeVTWZl6CoqYY2Zq2twMm2IoXyG2oo
j09xRO89eDUab2kPWt9kYop59tuptQe4XTOnaNNV3/RC3i90fESmlZ3wkMXQ/NeN
ZLFG7Bp+QYopBbGAsBNTIFU9ukljbHv7qxH6mhJ0D37YYpRya8iuyUNBpOXvleOE
xLp1Fmkyp9NKKnTCm8fdrPeWnT9iQdMcfJlXzZjuzjMWsLWNP/jfqz/d2V+DXKZj
V754aTAlItF7+A/D1rDu+uUiY+58u2FeO1dE+rspFpIRsalfJY7paHUe7/sTC/Yv
vP2xiodCmg==
-----END CERTIFICATE-----`
)

//nolint:gosec
func setupTestCerts(t *testing.T) (cert, key, ca string) {
	t.Helper()

	d := t.TempDir()

	cert, key, ca = filepath.Join(d, "test.cert"),
		filepath.Join(d, "test.key"), filepath.Join(d, "test.cacert")

	require.Nil(t, os.WriteFile(cert, []byte(_testCert), 0o655),
		"should create test cert file")
	require.Nil(t, os.WriteFile(key, []byte(_testKey), 0o655),
		"should create test key file")
	require.Nil(t, os.WriteFile(ca, []byte(_testCA), 0o655),
		"should create test ca-cert file")

	return
}

// TODO: start a tls server and require the server
// TODO: test listenern
// TODO: test mTLS settings.
func TestLoadTLS(t *testing.T) {
	requirer := require.New(t)

	cert, key, ca := setupTestCerts(t)

	t.Log("insecure config")
	{
		icfg := Config{Key: key, Cert: cert}
		cfg, err := GetTLSCfg(icfg)

		requirer.Nil(err)
		requirer.Equal(DefaultCipher, cfg.CipherSuites)
		requirer.Equal(DefaultCurve, cfg.CurvePreferences)
		requirer.Equal(uint16(tls.VersionTLS12), cfg.MinVersion)
		requirer.Equal(uint16(tls.VersionTLS13), cfg.MaxVersion)
		// TODO: test loaded certs ?
		requirer.Equal(tls.RequestClientCert, cfg.ClientAuth)
	}

	t.Log("secure config")
	{
		icfg := Config{
			Key: key, Cert: cert, Ca: ca,
			Level: RequireAndVerifyClientCertAndSAN,
		}
		cfg, err := GetTLSCfg(icfg)

		requirer.Nil(err)
		requirer.Equal(DefaultCipher, cfg.CipherSuites)
		requirer.Equal(DefaultCurve, cfg.CurvePreferences)
		requirer.Equal(uint16(tls.VersionTLS12), cfg.MinVersion)
		requirer.Equal(uint16(tls.VersionTLS13), cfg.MaxVersion)
		// TODO: test loaded certs ?
		requirer.Equal(tls.RequireAndVerifyClientCert, cfg.ClientAuth)
	}

	// TODO : t.Log("secured mTLS config")	{}
}
