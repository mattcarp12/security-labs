import os
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.primitives import serialization

def generate_pki_keys():
    print("--- Simple Python Key Generator ---")
    
    # 1. Key Size Selection
    try:
        bit_size = int(input("Enter key size (2048, 3072, or 4096): ") or 2048)
        if bit_size not in [2048, 3072, 4096]:
            print("Warning: Non-standard bit size. Defaulting to 2048.")
            bit_size = 2048
    except ValueError:
        bit_size = 2048

    print(f"Generating {bit_size}-bit RSA key pair...")

    # 2. Generate the Private Key
    # 65537 is the Fermat Prime 'e', used globally for RSA speed/security
    private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=bit_size,
    )

    # 3. Extract the Public Key
    public_key = private_key.public_key()

    # 4. Serialize Private Key to PEM
    # We use NoEncryption here for the lab, but in production, 
    # you'd use BestAvailableEncryption(b"mypassword")
    private_pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=serialization.NoEncryption()
    )

    # 5. Serialize Public Key to PEM
    public_pem = public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo
    )

    # 6. Save to files
    with open("private_key.pem", "wb") as f:
        f.write(private_pem)
    
    with open("public_key.pem", "wb") as f:
        f.write(public_pem)

    print("\nSuccess!")
    print(f"Created: {os.path.abspath('private_key.pem')}")
    print(f"Created: {os.path.abspath('public_key.pem')}")
    print("\nKeep your private_key.pem secret. Share public_key.pem freely.")

if __name__ == "__main__":
    generate_pki_keys()