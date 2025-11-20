package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "io"
	"encoding/base64"
	"fmt"
	"strings"

    "golang.org/x/crypto/hkdf"

    "github.com/alibabacloud-go/tea/tea"
    alikmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
    alikmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
    lru "github.com/hashicorp/golang-lru"
)

const (
    GcmIvLength           = 12
    CIPHER_TRANSFORMATION = "AES/GCM/NoPadding"
    ALGORITHM             = "AES"
)

var client *KmsClient

type KmsConfig struct {
    CaFilePath       string `json:"ca_filepath"`
    Protocal         string `json:"protocal"`
    ClientKeyContent string `json:"clientkey_content"`
    Password         string `json:"password"`
    Endpoint         string `json:"endpoint"`

    // 本地&dev无法访问kms
    IsDev  bool   `json:"is_dev"`
    DevCMK string `json:"dev_cmk"` // len in [16, 24, 32]
}


// The envelope cipher obj may be stored
type EnvelopeCipherObj struct {
    DataKeyIV        []byte
    EncryptedDataKey []byte
    Iv               []byte
    CipherText       []byte
}

func (eo *EnvelopeCipherObj) EncodeToString() string {
    return fmt.Sprintf("%s.%s.%s.%s", base64.RawStdEncoding.EncodeToString(eo.DataKeyIV),
        base64.RawStdEncoding.EncodeToString(eo.EncryptedDataKey),
        base64.RawStdEncoding.EncodeToString(eo.Iv),
        base64.RawStdEncoding.EncodeToString(eo.CipherText))
}

func (eo *EnvelopeCipherObj) Decode(encodedStr string) error {

    arrs := strings.Split(encodedStr, ".")
    if len(arrs) != 4 {
        return errors.New("invalid format")
    }

    var err error
    eo.DataKeyIV, err = base64.RawStdEncoding.DecodeString(arrs[0])
    if err != nil {
        return err
    }
    eo.DataKeyIV, err = base64.RawStdEncoding.DecodeString(arrs[1])
    if err != nil {
        return err
    }
    eo.Iv, err = base64.RawStdEncoding.DecodeString(arrs[2])
    if err != nil {
        return err
    }
    eo.CipherText, err = base64.RawStdEncoding.DecodeString(arrs[3])
    if err != nil {
        return err
    }

    return nil
}
// 信封加密示例（每次生成DataKey）
func EnvelopeEncryptByKeyId(keyId string, data []byte) (*EnvelopeCipherObj, error) {
    // 获取数据密钥，下面以Aliyun_AES_256密钥为例进行说明，数据密钥长度32字节
    generateDataKeyRequest := &alikmssdk.GenerateDataKeyRequest{
        KeyId:         tea.String(keyId),
        NumberOfBytes: tea.Int32(32),
    }

    // 调用生成数据密钥接口
    dataKeyResponse, err := client.GenerateDataKey(generateDataKeyRequest)
    if err != nil {
        return nil, err
    }

    // 使用专属KMS返回的数据密钥明文在本地对数据进行加密，下面以AES-256 GCM模式为例
    iv := make([]byte, GcmIvLength) // 加密初始向量，解密时需要传入
    rand.Read(iv)
    block, err := aes.NewCipher(dataKeyResponse.Plaintext)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    cipherText := gcm.Seal(nil, iv, data, nil)

    // 输出密文，密文输出或持久化由用户根据需要进行处理，下面示例仅展示将密文输出到一个对象的情况
    // 假如EnvelopeCipherObj是需要输出的密文对象，至少需要包括以下四个内容:
    // (1) dataKeyIV: 由专属KMS生成的加密初始向量，解密数据密钥密文时需要传入
    // (2) encryptedDataKey: 专属KMS返回的数据密钥密文
    // (3) iv: 加密初始向量
    // (4) cipherText: 密文数据
    envelopeCipherText := &EnvelopeCipherObj{
        DataKeyIV:        dataKeyResponse.Iv,
        EncryptedDataKey: dataKeyResponse.CiphertextBlob,
        Iv:               iv,
        CipherText:       cipherText,
    }
    return envelopeCipherText, nil
}

// 信封加密示例（基于已有的DataKey）
func EnvelopeEncryptByDataKey(dataKey *alikmssdk.GenerateDataKeyResponse, data []byte) (*EnvelopeCipherObj, error) {
    // 使用专属KMS返回的数据密钥明文在本地对数据进行加密，下面以AES-256 GCM模式为例
    iv := make([]byte, GcmIvLength) // 加密初始向量，解密时需要传入
    rand.Read(iv)
    block, err := aes.NewCipher(dataKey.Plaintext)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    cipherText := gcm.Seal(nil, iv, data, nil)

    // 输出密文，密文输出或持久化由用户根据需要进行处理，下面示例仅展示将密文输出到一个对象的情况
    // 假如envelopeCipherText是需要输出的密文对象，至少需要包括以下四个内容:
    // (1) dataKeyIV: 由专属KMS生成的加密初始向量，解密数据密钥密文时需要传入
    // (2) encryptedDataKey: 专属KMS返回的数据密钥密文
    // (3) iv: 加密初始向量
    // (4) cipherText: 密文数据
    envelopeCipherText := &EnvelopeCipherObj{
        DataKeyIV:        dataKey.Iv,
        EncryptedDataKey: dataKey.CiphertextBlob,
        Iv:               iv,
        CipherText:       cipherText,
    }
    return envelopeCipherText, nil
}

// 信封解密示例
func EnvelopeDecrypt(keyId string, cipherText *EnvelopeCipherObj) ([]byte, error) {
    // 调用解密接口进行解密
    plainDataKey, err := client.DecryptDataKey(keyId, cipherText.EncryptedDataKey, cipherText.DataKeyIV)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(plainDataKey)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    decryptedData, err := gcm.Open(nil, cipherText.Iv, cipherText.CipherText, nil)
    if err != nil {
        return nil, err
    }

    return decryptedData, nil
}

type KmsClient struct {
    alikmsClient *alikmssdk.Client
    isDev        bool
    devCMK       []byte

    // cipher DK -> plain DK
    cacheDataKey *lru.Cache
}

// 使用ClientKey内容创建KMS实例SDK Client对象
func InitKmsClient(c *KmsConfig) (*KmsClient, error) {

    if c.IsDev {
        if !(len(c.DevCMK) == 16 || len(c.DevCMK) == 24 || len(c.DevCMK) == 32) {
            return nil, errors.New("key length must be 16/24/32")
        }

        cache, _ := lru.New(1000)

        client = &KmsClient{
            isDev:        true,
            devCMK:       []byte(c.DevCMK),
            cacheDataKey: cache,
        }
        return client, nil
    }

    // 创建KMS实例SDK Client配置
    config := &alikmsopenapi.Config{
        CaFilePath: tea.String(c.CaFilePath),
        // 连接协议请设置为"https"。KMS实例服务仅允许通过HTTPS协议访问。
        Protocol: tea.String(c.Protocal),
        // 请替换为ClientKey文件的内容
        ClientKeyContent: tea.String(c.ClientKeyContent),
        // 请替换为创建ClientKey时输入的加密口令
        Password: tea.String(c.Password),
        // 设置endpoint为<your KMS Instance Id>.cryptoservice.kms.aliyuncs.com。
        Endpoint: tea.String(c.Endpoint),
    }
    // 创建KMS实例SDK Client对象
    aliCLI, err := alikmssdk.NewClient(config)
    if err != nil {
        return nil, err
    }

    cache, _ := lru.New(1000)
    client = &KmsClient{
        alikmsClient: aliCLI,
        cacheDataKey: cache,
    }
    return client, nil
}

func (cli *KmsClient) GenerateDataKey(request *alikmssdk.GenerateDataKeyRequest) (_result *alikmssdk.GenerateDataKeyResponse, _err error) {
    if cli.isDev {
        // 派生 datakey
        salt := make([]byte, sha256.Size)
        rand.Read(salt)
        kdf := hkdf.New(sha256.New, []byte(cli.devCMK), salt, nil)
        key := make([]byte, *request.NumberOfBytes)
        if _, err := io.ReadFull(kdf, key); err != nil {
            return nil, err
        }

        cipherDK, err := cli.encrypt(&alikmssdk.EncryptRequest{
            Plaintext: key,
        })
        if err != nil {
            return nil, err
        }

        return &alikmssdk.GenerateDataKeyResponse{
            Iv:             cipherDK.Iv,
            Plaintext:      key,
            CiphertextBlob: cipherDK.CiphertextBlob,
        }, nil
    }

    return cli.alikmsClient.GenerateDataKey(request)

}

func (cli *KmsClient) DecryptDataKey(keyId string, cipher, iv []byte) (_result []byte, _err error) {

    value, ok := cli.cacheDataKey.Get(hex.EncodeToString(cipher))
    if ok {
        return value.([]byte), nil
    }

    // 解密数据密钥密文，得到数据密钥明文
    decryptRequest := &alikmssdk.DecryptRequest{
        KeyId:          tea.String(keyId),
        CiphertextBlob: cipher,
        Iv:             iv,
    }

    resp, err := cli.decrypt(decryptRequest)
    if err != nil {
        return nil, err
    }

    cli.cacheDataKey.Add(hex.EncodeToString(cipher), resp.Plaintext)

    return resp.Plaintext, nil
}

func (cli *KmsClient) decrypt(request *alikmssdk.DecryptRequest) (_result *alikmssdk.DecryptResponse, _err error) {

    if cli.isDev {
        block, err := aes.NewCipher(cli.devCMK)
        if err != nil {
            return nil, err
        }
        gcm, err := cipher.NewGCM(block)
        if err != nil {
            return nil, err
        }

        decryptedData, err := gcm.Open(nil, request.Iv, request.CiphertextBlob, nil)
        if err != nil {
            return nil, err
        }
        return &alikmssdk.DecryptResponse{
            Plaintext: decryptedData,
        }, nil
    }

    return cli.alikmsClient.Decrypt(request)
}

func (cli *KmsClient) encrypt(request *alikmssdk.EncryptRequest) (_result *alikmssdk.EncryptResponse, _err error) {

    if cli.isDev {

        iv := make([]byte, GcmIvLength) // 加密初始向量，解密时需要传入
        rand.Read(iv)
        block, err := aes.NewCipher(cli.devCMK)
        if err != nil {
            return nil, err
        }
        gcm, err := cipher.NewGCM(block)
        if err != nil {
            return nil, err
        }
        cipherText := gcm.Seal(nil, iv, request.Plaintext, nil)
        return &alikmssdk.EncryptResponse{
            KeyId:          request.KeyId,
            Iv:             iv,
            CiphertextBlob: cipherText,
        }, nil

    }

    return cli.alikmsClient.Encrypt(request)
}

func main() {
	// 初始化KMS Client对象
	config := &KmsConfig{
		CaFilePath:       "", // 可选，指定CA证书路径
		Protocol:		 "https",
		Endpoint:         "kms.cn-hangzhou.aliyuncs.com",
		AccessKeyId:      tea.String("<your access key id>"),
		AccessKeySecret:  tea.String("<your access key secret>"),
		IsDev:            false,
		DevCMK:           "yourDevCMK", // IsDev=true时必填
	}
	client, err := InitKmsClient(config)
	if err != nil {
		panic(err)		
	}
}