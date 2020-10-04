package app

import (
    "github.com/go-redis/redis/v7"
    "github.com/dgrijalva/jwt-go"
    "github.com/twinj/uuid"
    "time"
    "strconv"
 )

type JWTStorage struct {
    Client *redis.Client
}

func NewJWTStorage(redisClient *redis.Client) *JWTStorage {
    return &JWTStorage{
        Client: redisClient,
    }
}

type TokenDetails struct {
    AccessToken  string
    RefreshToken string
    AccessUuid   string
    RefreshUuid  string
    AtExpires    int64
    RtExpires    int64
}

func (storage *JWTStorage) NewToken(userId uint, accessSecret string, refreshSecret string) (*TokenDetails, error) {
    td := &TokenDetails{}
    td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
    td.AccessUuid = uuid.NewV4().String()

    td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
    td.RefreshUuid = uuid.NewV4().String()

    var err error

    atClaims := jwt.MapClaims{}
    atClaims["authorized"] = true
    atClaims["user_id"] = userId
    atClaims["access_uuid"] = td.AccessUuid
    atClaims["exp"] = td.AtExpires

    at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
    td.AccessToken, err = at.SignedString([]byte(accessSecret))
    if err != nil {
        return nil, err
    }

    rtClaims := jwt.MapClaims{}
    rtClaims["refresh_uuid"] = td.RefreshUuid
    rtClaims["user_id"] = userId
    rtClaims["exp"] = td.RtExpires
    rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
    td.RefreshToken, err = rt.SignedString([]byte(refreshSecret))
    if err != nil {
        return nil, err
    }

    return td, nil
}

func (storage *JWTStorage) CreateAuth(userId uint64, td *TokenDetails) error {
    at := time.Unix(td.AtExpires, 0)
    rt := time.Unix(td.RtExpires, 0)
    now := time.Now()

    errAccess := storage.Client.Set(td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
    if errAccess != nil {
        return errAccess
    }

    errRefresh := storage.Client.Set(td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
    if errRefresh != nil {
        return errRefresh
    }

    return nil
}
