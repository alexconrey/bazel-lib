"""Dearmor a PEM encoded file.
"""

load(
    "//lib/private:dearmor_pem_file.bzl",
    _dearmor_pem_file_lib = "dearmor_pem_file_lib",
)

# export the starlark library as a public API
dearmor_pem_file_lib = _dearmor_pem_file_lib

dearmor_pem_file = rule(
    implementation = _dearmor_pem_file_lib.impl,
    attrs = _dearmor_pem_file_lib.attrs,
    executable = True,
)