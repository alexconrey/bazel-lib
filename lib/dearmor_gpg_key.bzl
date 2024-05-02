"""Dearmor an Armored (ASCII encoded) GPG signature/public key file.
"""

load(
    "//lib/private:dearmor_gpg_key.bzl",
    # _copy_to_directory_bin_action = "copy_to_directory_bin_action",
    _dearmor_gpg_key_lib = "dearmor_gpg_key_lib",
)

# export the starlark library as a public API
dearmor_gpg_key_lib = _dearmor_gpg_key_lib
# # copy_to_directory_bin_action = _copy_to_directory_bin_action

dearmor_gpg_key = rule(
    implementation = _dearmor_gpg_key_lib.impl,
    attrs = _dearmor_gpg_key_lib.attrs,
    executable = True,
)