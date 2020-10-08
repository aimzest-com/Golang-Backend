package app

import (
    "github.com/go-redis/redis/v7"
    "github.com/dgrijalva/jwt-go"
    "github.com/twinj/uuid"
    "time"
    "strconv"
    "net/http"
    "strings"
    "fmt"
 )

type JWTStorage struct {
    Client *redis.Client
    AccessSecret string
    RefreshSecret string
}

func NewJWTStorage(redisClient *redis.Client, accessSecret string, refreshSecret string) *JWTStorage {
    return &JWTStorage{
        Client: redisClient,
        AccessSecret: accessSecret,
        RefreshSecret: refreshSecret,
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

func (storage *JWTStorage) NewToken(userId uint) (*TokenDetails, error) {
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
    td.AccessToken, err = at.SignedString([]byte(storage.AccessSecret))
    if err != nil {
        return nil, err
    }

    rtClaims := jwt.MapClaims{}
    rtClaims["refresh_uuid"] = td.RefreshUuid
    rtClaims["user_id"] = userId
    rtClaims["exp"] = td.RtExpires
    rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
    td.RefreshToken, err = rt.SignedString([]byte(storage.RefreshSecret))
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

func verifyToken(r *http.Request, accessSecret string) (*jwt.Token, error) {
    tokenString := extractToken(r)
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(accessSecret), nil
    })

    if err != nil {
        return nil, err
    }

    return token, nil
}

func tokenValid(r *http.Request, accessSecret string) error {
    token, err := verifyToken(r, accessSecret)
    if err != nil {
        return err
    }

    if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
        return err //todo something wrong with this error
    }

    return nil
}

func extractToken(r *http.Request) string {
    bearToken := r.Header.Get("Authorization")
    strArr := strings.Split(bearToken, " ")
    if len(strArr) == 2 {
        return strArr[1]
    }

    return ""
}

func (storage *JWTStorage) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
    token, err := verifyToken(r, storage.AccessSecret)
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if ok && token.Valid {
        accessUuid, ok := claims["access_uuid"].(string)
        if !ok {
            return nil, err
        }

        userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
        if err != nil {
            return nil, err
        }

        return &AccessDetails {
            AccessUuid: accessUuid,
            UserId: userId,
        }, nil
    }

    return nil, err
}

type AccessDetails struct {
    AccessUuid string
    UserId uint64
}

func (storage *JWTStorage) FetchAuth(authD *AccessDetails) (uint64, error) {
    userid, err := storage.Client.Get(authD.AccessUuid).Result()
    if err != nil {
        return 0, err
    }

    userID, _ := strconv.ParseUint(userid, 10, 64)
    return userID, nil
}
