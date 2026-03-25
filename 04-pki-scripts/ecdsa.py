import os
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import serialization

def generate_ecdsa_keys():
    print("--- Simple Python ECDSA Key Generator ---")
    
    # 1. Curve Selection
    # In ECDSA, we don't choose "bit size" like RSA (2048, 4096).
    # Instead, we choose a standardized mathematical 'Curve'.
    # SECP256R1 (also known as NIST P-256) is the most widely used curve 
    # for web browsing (HTTPS) and general PKI today.
    print("Available Curves: \n1. SECP256R1 (Standard/Web)\n2. SECP384R1 (Higher Security)\n3. SECP521R1 (Ultra High Security)")
    
    choice = input("Select curve (1, 2, or 3) [Default 1]: ") or "1"
    
    if choice == "2":
        selected_curve = ec.SECP384R1()
    elif choice == "3":
        selected_curve = ec.SECP521R1()
    else:
        selected_curve = ec.SECP256R1()

    print(f"Generating keys using curve: {selected_curve.name}...")

    # 2. Generate the Private Key
    # Unlike RSA, where we find primes, ECDSA generation involves picking a 
    # random number 'd' within the range of the curve's order.
    private_key = ec.generate_private_key(selected_curve)

    # 3. Derive the Public Key
    # The public key is a 'Point' on the curve (x, y coordinates).
    # It is calculated by multiplying the private key 'd' by a 
    # predefined 'Generator Point' (G) on the curve.
    # Math: Public Key = d * G
    public_key = private_key.public_key()

    # 4. Serialize Private Key to PEM
    # We use PKCS8, which is the modern standard format for private keys.
    private_pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=serialization.NoEncryption()
    )

    # 5. Serialize Public Key to PEM
    # SubjectPublicKeyInfo is the standard X.509 format for public keys.
    public_pem = public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo
    )

    # 6. Save to files
    with open("ecdsa_private_key.pem", "wb") as f:
        f.write(private_pem)
    
    with open("ecdsa_public_key.pem", "wb") as f:
        f.write(public_pem)

    print("\nSuccess!")
    print(f"Private Key: {os.path.abspath('ecdsa_private_key.pem')}")
    print(f"Public Key:  {os.path.abspath('ecdsa_public_key.pem')}")
    
    # Notice the file size difference!
    print(f"\nNote: Look at the file sizes. This {selected_curve.name} key ")
    print("is much smaller than your 2048-bit RSA key, yet provides ")
    print("equivalent or better security.")

if __name__ == "__main__":
    generate_ecdsa_keys()