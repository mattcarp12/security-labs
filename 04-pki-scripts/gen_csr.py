from cryptography import x509
from cryptography.x509.oid import NameOID
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives import serialization

def generate_csr():
    print("--- Lab 3: Generating a Certificate Signing Request (CSR) ---")

    # 1. Load the "Server" Private Key 
    # (Usually you'd generate a NEW one for the server, but we can use 
    # your RSA private_key.pem for this exercise)
    try:
        with open("private_key.pem", "rb") as key_file:
            server_private_key = serialization.load_pem_private_key(
                key_file.read(),
                password=None,
            )
    except FileNotFoundError:
        print("Error: private_key.pem not found.")
        return

    # 2. Define the Server's Identity
    # This is the "Subject" of the CSR (e.g., your website name)
    subject = x509.Name([
        x509.NameAttribute(NameOID.COUNTRY_NAME, u"US"),
        x509.NameAttribute(NameOID.STATE_OR_PROVINCE_NAME, u"Idaho"),
        x509.NameAttribute(NameOID.LOCALITY_NAME, u"Boise"),
        x509.NameAttribute(NameOID.ORGANIZATION_NAME, u"My-Local-Server"),
        x509.NameAttribute(NameOID.COMMON_NAME, u"localhost"), # The 'domain'
    ])

    # 3. Create the CSR
    # We include our Public Key and sign it with our Private Key
    # This proves to the CA that we actually own the key we want certified.
    csr = x509.CertificateSigningRequestBuilder().subject_name(
        subject
    ).add_extension(
        x509.SubjectAlternativeName([x509.DNSName(u"localhost")]),
        critical=False,
    ).sign(server_private_key, hashes.SHA256())

    # 4. Save the CSR to a file
    with open("server_request.csr", "wb") as f:
        f.write(csr.public_bytes(serialization.Encoding.PEM))

    print("\nSuccess! CSR created: server_request.csr")
    print("This is your 'Application' ready to be sent to the CA.")

if __name__ == "__main__":
    generate_csr()