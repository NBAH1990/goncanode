package types

type HashAlgorithm string

const (
	GOST34311   HashAlgorithm = "GOST34311"
	MD5         HashAlgorithm = "MD5"
	SHA1        HashAlgorithm = "SHA1"
	SHA224      HashAlgorithm = "SHA224"
	SHA256      HashAlgorithm = "SHA256"
	SHA384      HashAlgorithm = "SHA384"
	SHA512      HashAlgorithm = "SHA512"
	RIPEMD128   HashAlgorithm = "RIPEMD128"
	RIPEMD160   HashAlgorithm = "RIPEMD160"
	RIPEMD256   HashAlgorithm = "RIPEMD256"
	GOST34311GT HashAlgorithm = "GOST34311GT"
)
