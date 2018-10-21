#!/usr/bin/python

import os, sys, json
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives.asymmetric import rsa, padding as apadding 
from cryptography.hazmat.primitives import hashes, hmac, padding, serialization
from cryptography.hazmat.backends import default_backend
from cryptography.exceptions import InvalidSignature

AES_BLOCKSIZE = 16
AES_KEYSIZE = 32
HMAC_KEYSIZE = 32

def encrypt(plaintext, pubkey):
    padder = padding.PKCS7(AES_BLOCKSIZE * 8).padder()
    padded_plaintext = padder.update(plaintext)
    padded_plaintext += padder.finalize()

    aes_key = os.urandom(AES_KEYSIZE)
    iv = os.urandom(AES_BLOCKSIZE)
    cipher = Cipher(algorithms.AES(aes_key), modes.CBC(iv), backend=default_backend())
    encryptor = cipher.encryptor()
    ciphertext = encryptor.update(padded_plaintext) + encryptor.finalize()

    hmac_key = os.urandom(HMAC_KEYSIZE)
    mac = hmac.HMAC(hmac_key, hashes.SHA256(), backend=default_backend())
    mac.update(ciphertext)
    tag = mac.finalize()

    keys = aes_key + hmac_key
    encrypted_keys = pubkey.encrypt(
        keys,
        apadding.OAEP(
            mgf=apadding.MGF1(algorithm=hashes.SHA256()),
            algorithm=hashes.SHA256(),
            label=None))

    out = {
        "Key": encrypted_keys.encode("base64"),
        "IV": iv.encode("base64"),
        "Tag": tag.encode("base64"),
        "Msg": ciphertext.encode("base64")
    }

    return json.dumps(out)

def decrypt(ciphertext, privkey):
    struct = json.loads(ciphertext)
    encrypted_keys = struct["Key"].decode("base64")
    tag = struct["Tag"].decode("base64")
    iv = struct["IV"].decode("base64")
    ciphertext = struct["Msg"].decode("base64")

    keys = privkey.decrypt(
        encrypted_keys,
        apadding.OAEP(
            mgf=apadding.MGF1(algorithm=hashes.SHA256()),
            algorithm=hashes.SHA256(),
            label=None))

    aes_key = keys[:AES_KEYSIZE]
    hmac_key = keys[AES_KEYSIZE:]

    mac = hmac.HMAC(hmac_key, hashes.SHA256(), backend=default_backend())
    mac.update(ciphertext)
    try:
        mac.verify(tag)
    except InvalidSignature:
        raise

    cipher = Cipher(algorithms.AES(aes_key), modes.CBC(iv), backend=default_backend())
    decryptor = cipher.decryptor()
    padded_plaintext = decryptor.update(ciphertext) + decryptor.finalize()

    unpadder = padding.PKCS7(AES_BLOCKSIZE * 8).unpadder()
    return unpadder.update(padded_plaintext) + unpadder.finalize()

if __name__ == "__main__":
    if len(sys.argv) < 4:
        print("USAGE: {} <message> <key> e/d".format(sys.argv[0]))
        sys.exit(0)

    plaintext = sys.argv[1]
    keypath = sys.argv[2]
    mode = sys.argv[3]

    if mode == "d":
        with open(keypath, "rb") as key_file:
            key = serialization.load_pem_private_key(
                key_file.read(),
                password=None,
                backend=default_backend())

            print(decrypt(plaintext, key))
    else:
        with open(keypath, "rb") as key_file:
            key = serialization.load_pem_public_key(
                key_file.read(),
                backend=default_backend())

            print(encrypt(plaintext, key))
                

