"""This module provides functionality for generating a QR code for TOTP (Time-based One-Time Password) setup."""

import os
import warnings

from dataclasses import dataclass

try:
    import pyotp
except (ImportError, ModuleNotFoundError):
    raise ImportError(
        "pyotp is required for OTP functionality. Please install it via 'pip install pyotp'."
    )

try:
    import qrcode
except (ImportError, ModuleNotFoundError):
    raise ImportError(
        "qrcode is required for OTP functionality. Please install it via 'pip install qrcode[pil]'."
    )


def getenv(key: str, default: str | None = None) -> str | None:
    """Gets an environment variable or returns a default value.

    Args:
        key: The environment variable key.
        default: The default value to return if the key is not found.
    Returns:
        The value of the environment variable or the default value.
    """
    return os.environ.get(key.upper()) or os.environ.get(key.lower()) or default


@dataclass
class OTPConfig:
    secret: str
    qr_filename: str
    authenticator_user: str
    authenticator_app: str


config = OTPConfig(
    secret="",
    qr_filename=getenv("authenticator_qr_filename", "totp_qr.png"),
    authenticator_user=getenv("authenticator_user", "thevickypedia"),
    authenticator_app=getenv("authenticator_app", "FileBrowser"),
)


def display_secret() -> None:
    """Displays the TOTP secret key."""
    try:
        term_size = os.get_terminal_size().columns
    except OSError:
        term_size = 120
    base = "*" * term_size
    print(
        f"\n{base}\n"
        f"\nYour TOTP secret key is: {config.secret}"
        f"\n\n./filebrowser config set --authenticatorToken {config.secret}\n\n"
        f"\nQR code saved as {config.qr_filename!r} (you can scan this with your Authenticator app).\n"
        f"\n{base}",
    )


def generate_qr() -> None:
    """Generates a QR code for TOTP setup."""
    # STEP 1: Generate a new secret key for the user (store this securely!)
    secret = pyotp.random_base32()

    # STEP 2: Create a provisioning URI (for the QR code)
    uri = pyotp.TOTP(secret).provisioning_uri(
        name=str(config.authenticator_user), issuer_name=config.authenticator_app
    )

    # STEP 3: Generate a QR code (scan this with your authenticator app)
    qr = qrcode.make(uri)
    # Save the QR code
    qr.save(config.qr_filename)

    # STEP 4: Update the config with the new secret
    config.secret = secret
    display_secret()


if __name__ == "__main__":
    generate_qr()
