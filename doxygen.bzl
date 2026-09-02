"""Provides a rule to run Doxygen hermetically and output XML files."""

def _doxygen_xml_impl(ctx):
    # Declare the output directory for XML files.
    xml_dir = ctx.actions.declare_directory(ctx.attr.name)

    # Declare the temporary overrides file.
    overrides = ctx.actions.declare_file(ctx.attr.name + "_overrides")

    # Extract unique parent directories of input headers for Doxygen INPUT.
    input_dirs = []
    for f in ctx.files.headers:
        if f.dirname not in input_dirs:
            input_dirs.append(f.dirname)

    # Write the Doxygen overrides file using Bazel's native action.
    overrides_content = "\n".join([
        "GENERATE_XML = YES",
        "XML_OUTPUT = " + xml_dir.path,
        "GENERATE_HTML = NO",
        "GENERATE_LATEX = NO",
        "INPUT = " + " ".join(input_dirs),
        "QUIET = YES",
    ])
    ctx.actions.write(overrides, overrides_content)

    doxygen_bin = ctx.executable._doxygen

    # Combine the template Doxyfile and overrides, then execute doxygen.
    cmd = "cat {template} {overrides} > temp_doxyfile && {doxygen} temp_doxyfile".format(
        template = ctx.file.doxyfile_template.path,
        overrides = overrides.path,
        doxygen = doxygen_bin.path,
    )

    ctx.actions.run_shell(
        inputs = ctx.files.headers + [ctx.file.doxyfile_template, overrides],
        outputs = [xml_dir],
        tools = [doxygen_bin],
        command = cmd,
        progress_message = "Running doxygen to generate XML for %s" % ctx.label,
    )
    return [DefaultInfo(files = depset([xml_dir]))]

doxygen_xml = rule(
    doc = "Runs Doxygen using a template and overrides to produce XML metadata.",
    implementation = _doxygen_xml_impl,
    attrs = {
        "headers": attr.label_list(
            doc = "The list of header files to parse.",
            mandatory = True,
            allow_files = True,
        ),
        "doxyfile_template": attr.label(
            doc = "The base Doxyfile containing project configurations and aliases.",
            mandatory = True,
            allow_single_file = True,
        ),
        "_doxygen": attr.label(
            doc = "The hermetic Doxygen binary tool.",
            default = "@doxygen_release//:bin/doxygen",
            allow_single_file = True,
            executable = True,
            cfg = "exec",
        ),
    },
)
