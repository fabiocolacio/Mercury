package server

import(
    "hash"
    "crypto/sha256"
)

var(
    KeyHashAlgo func() hash.Hash = sha256.New
)

var(
    KeyHashIterations int = 250000
    KeyHashLength int = 32
    SaltLength int = 16
    UsernameMaxLength = 16
    ChallengeLength = 16
)
