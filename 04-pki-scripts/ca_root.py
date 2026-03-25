import datetime
from cryptography import x509
from cryptography.x509.oid import NameOID
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives import serialization

def create_root_ca():
    print("--- Lab 2: Creating a Root CA Certificate ---")

    # 1. Load your Private Key from Lab 1
    # We need the private key to 'sign' the certificate.
    try:
        with open("private_key.pem", "rb") as key_file:
            private_key = serialization.load_pem_private_key(
                key_file.read(),
                password=None,
            )
    except FileNotFoundError:
        print("Error: private_key.pem not found. Run rsa.py first!")
        return

    # 2. Define the Identity (Subject)
    # This is the 'metadata' that identifies your CA.
    subject = issuer = x509.Name([
        x509.NameAttribute(NameOID.COUNTRY_NAME, u"US"),
        x509.NameAttribute(NameOID.STATE_OR_PROVINCE_NAME, u"Idaho"),
        x509.NameAttribute(NameOID.LOCALITY_NAME, u"Boise"),
        x509.NameAttribute(NameOID.ORGANIZATION_NAME, u"My-DIY-Lab-CA"),
        x509.NameAttribute(NameOID.COMMON_NAME, u"My Root CA"),
    ])

    # 3. Build the Certificate
    # We are 'building' the ID card before we 'stamp' it.
    cert = x509.CertificateBuilder().subject_name(
        subject
    ).issuer_name(
        issuer
    ).public_key(
        private_key.public_key()
    ).serial_number(
        x509.random_serial_number()
    ).not_valid_before(
        datetime.datetime.utcnow()
    ).not_valid_after(
        # Our Root CA will be valid for 10 years
        datetime.datetime.utcnow() + datetime.timedelta(days=3650)
    ).add_extension(
        # This extension tells computers "This certificate belongs to a CA"
        x509.BasicConstraints(ca=True, path_length=None), critical=True,
    ).sign(private_key, hashes.SHA256()) # The 'Stamp': Signing with Private Key

    # 4. Save the Root Certificate
    with open("root_cert.pem", "wb") as f:
        f.write(cert.public_bytes(serialization.Encoding.PEM))

    print("\nSuccess! Root CA Certificate created: root_cert.pem")
    print("This file contains your Public Key and your CA Identity.")

if __name__ == "__main__":
    create_root_ca()