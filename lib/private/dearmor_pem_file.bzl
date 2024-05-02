_dearmor_pem_file_attr = {
        "key": attr.label(allow_single_file = True),
        "output_extension": attr.string(default = "gpg"),
        "_dearmor_pem_file_sh_tpl": attr.label(allow_single_file = True, default = ":dearmor_pem_file.sh.tpl"),
        # use '_tool' attribute for development only; do not commit with this attribute active since it
        # propagates a dependency on rules_go which would be breaking for users
        "_tool": attr.label(
            allow_single_file = True,
            executable = True,
            cfg = "exec",
            default = "//tools/dearmor_pem_file:dearmor_pem_file",
        ),
    }

def _dearmor_pem_file_impl(ctx):
    name = ctx.label.name
    output_ext = ctx.attr.output_extension

    armored_output_script = ctx.actions.declare_file("{}.sh".format(name))
    armored_output = ctx.actions.declare_file("{}.{}".format(name, output_ext))

    ctx.actions.expand_template(
        template = ctx.file._dearmor_pem_file_sh_tpl,
        output = armored_output_script,
        is_executable = True,
        substitutions = {
            "{{extract_binary_path}}": ctx.file._tool.path,
            "{{input_file}}": ctx.file.key.path,
            "{{output_file}}": armored_output.path,
        },
    )

    ctx.actions.run(
        executable = armored_output_script,
        outputs = [armored_output],
        inputs = [
            ctx.file.key,
            ctx.file._tool
        ],
    )

    return [
        DefaultInfo(
            files = depset([armored_output]),
            executable = armored_output_script,
        )
    ]

dearmor_pem_file_lib = struct(
    attrs = _dearmor_pem_file_attr,
    impl = _dearmor_pem_file_impl,
    provides = [DefaultInfo],
)