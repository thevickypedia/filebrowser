<p align="center">
  <img src="https://raw.githubusercontent.com/filebrowser/logo/master/banner.png" width="550"/>
</p>

[![Build](https://github.com/thevickypedia/filebrowser/actions/workflows/release.yml/badge.svg)](https://github.com/thevickypedia/filebrowser/actions/workflows/main.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/thevickypedia/filebrowser)](https://goreportcard.com/report/github.com/thevickypedia/filebrowser)
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/thevickypedia/filebrowser)
[![Version](https://img.shields.io/github/release/filebrowser/filebrowser.svg)](https://github.com/thevickypedia/filebrowser/releases/latest)

File Browser provides a file managing interface within a specified directory and it can be used to upload, delete, preview and edit your files. It is a **create-your-own-cloud**-kind of software where you can just install it on your server, direct it to a path and access your files through a nice web interface.

## Documentation

Documentation on how to install, configure, and contribute to this project is hosted at [filebrowser.org](https://filebrowser.org).

## Authentication Changes & Security Enhancements

This project is a fork of the [filebrowser](https://github.com/filebrowser/filebrowser) project.  It incorporates significant changes to the JSON authentication method, prioritizing security.

**Key Changes:**

* **Improved JSON Authentication:**  The JSON authentication method has been redesigned to leverage HTTP headers for authentication, instead of relying solely on the JSON payload. This improves security by reducing potential exposure of authentication credentials in the request body.
* **Transit Protection:** Added measures to protect data in transit using `base64` and unicode encoding.
* **Brute force Protection:** Includes protection against repeated login failures with exponential backoff and progressive lockout.
* **IP allow-list:** Includes IP allow listing techniques by leveraging the following config.
- `allowedOrigins` - Sequence of IP address/hostname to allow.
- `refreshAllowedOrigins` - Interval in seconds to auto refresh origins.
- `allowPrivateIP` - Boolean flag to allow private IP address of the host machine.
- `allowPublicIP`- Boolean flag to allow public IP address of the host machine.
- `authenticatorToken` - MFA token for time based OTP.

* **Multifactor Authentication** MFA has been implemented using TOTP which significantly improves security by adding a second layer of verification to the authenticate the server.
* To generate a TOTP token, run ``python thevickypedia_scripts/otp.py``, and scan the ``totp_qr.png`` using any authenticator app.

> These changes significantly improve the security posture of a basic authentication mechanism.

**File Details:**

* [`auth.ts`](https://github.com/thevickypedia/filebrowser/blob/8fbbf07/frontend/src/utils/auth.ts): Contains the updated authentication logic.  Review this file for details on the implementation.

**Original Filebrowser Repo:** [https://github.com/filebrowser/filebrowser](https://github.com/filebrowser/filebrowser)
