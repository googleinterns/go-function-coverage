workspace(name = "go_monorepo")

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
    name = "io_bazel_rules_go",
    commit = "11ec19e2507107881919b41038a811076b75907f",
    remote = "https://github.com/muratekici/rules_go.git",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
