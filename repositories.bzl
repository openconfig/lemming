load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    go_repository(
        name = "build_buf_gen_go_bufbuild_protovalidate_protocolbuffers_go",
        importpath = "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go",
        sum = "h1:31on4W/yPcV4nZHL4+UCiCvLPsMqe/vJcNg8Rci0scc=",
        version = "v1.36.10-20250912141014-52f32327d4b0.1",
    )
    go_repository(
        name = "cat_dario_mergo",
        importpath = "dario.cat/mergo",
        sum = "h1:AGCNq9Evsj31mOgNPcLyXc+4PNABt905YmuqPYYpBWk=",
        version = "v1.0.0",
    )
    go_repository(
        name = "co_honnef_go_tools",
        importpath = "honnef.co/go/tools",
        sum = "h1:+JKf3xJ1ni4CwrhVg4/pqsfPGP6vNAXcKbMXJodYx3w=",
        version = "v0.8.1",
    )
    go_repository(
        name = "com_github_ajstarks_svgo",
        importpath = "github.com/ajstarks/svgo",
        sum = "h1:slYM766cy2nI3BwyRiyQj/Ud48djTMtMebDqepE95rw=",
        version = "v0.0.0-20211024235047-1546f124cd8b",
    )
    go_repository(
        name = "com_github_andybalholm_brotli",
        importpath = "github.com/andybalholm/brotli",
        sum = "h1:eLKJA0d02Lf0mVpIDgYnqXcUn0GqVmEFny3VuID1U3M=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_antihax_optional",
        importpath = "github.com/antihax/optional",
        sum = "h1:xK2lYat7ZLaVVcIuj82J8kIro4V6kDe0AUDFboUCwcg=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_aristanetworks_arista_ceoslab_operator_v2",
        importpath = "github.com/aristanetworks/arista-ceoslab-operator/v2",
        sum = "h1:1aAxwwu4xyfiU1/FX2D5x/jsF/sxFVkjVhvF661isM4=",
        version = "v2.1.2",
    )
    go_repository(
        name = "com_github_armon_go_socks5",
        importpath = "github.com/armon/go-socks5",
        sum = "h1:0CwZNZbxp69SHPdPJAN/hZIm0C4OItdklCFmMRWYpio=",
        version = "v0.0.0-20160902184237-e75332964ef5",
    )
    go_repository(
        name = "com_github_azure_go_autorest",
        importpath = "github.com/Azure/go-autorest",
        sum = "h1:V5VMDjClD3GiElqLWO7mz2MxNAK/vTfRHdAubSIPRgs=",
        version = "v14.2.0+incompatible",
    )
    go_repository(
        name = "com_github_azure_go_autorest_autorest",
        importpath = "github.com/Azure/go-autorest/autorest",
        sum = "h1:90Y4srNYrwOtAgVo3ndrQkTYn6kf1Eg/AjTFJ8Is2aM=",
        version = "v0.11.18",
    )
    go_repository(
        name = "com_github_azure_go_autorest_autorest_adal",
        importpath = "github.com/Azure/go-autorest/autorest/adal",
        sum = "h1:Mp5hbtOePIzM8pJVRa3YLrWWmZtoxRXqUEzCfJt3+/Q=",
        version = "v0.9.13",
    )
    go_repository(
        name = "com_github_azure_go_autorest_autorest_date",
        importpath = "github.com/Azure/go-autorest/autorest/date",
        sum = "h1:7gUk1U5M/CQbp9WoqinNzJar+8KY+LPI6wiWrP/myHw=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_azure_go_autorest_logger",
        importpath = "github.com/Azure/go-autorest/logger",
        sum = "h1:IG7i4p/mDa2Ce4TRyAO8IHnVhAVF3RFU+ZtXWSmf4Tg=",
        version = "v0.2.1",
    )
    go_repository(
        name = "com_github_azure_go_autorest_tracing",
        importpath = "github.com/Azure/go-autorest/tracing",
        sum = "h1:TYi4+3m5t6K48TGI9AUdb+IzbnSxvnvUMfuitfgcfuo=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_github_beorn7_perks",
        importpath = "github.com/beorn7/perks",
        sum = "h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_blang_semver",
        importpath = "github.com/blang/semver",
        sum = "h1:cQNTCjp13qL8KC3Nbxr/y2Bqb63oX6wdnnjpJbkM4JQ=",
        version = "v3.5.1+incompatible",
    )
    go_repository(
        name = "com_github_blang_semver_v4",
        importpath = "github.com/blang/semver/v4",
        sum = "h1:1PFHFE6yCCTv8C1TeyNNarDzntLi7wMI5i/pzqYIsAM=",
        version = "v4.0.0",
    )
    go_repository(
        name = "com_github_burntsushi_toml",
        importpath = "github.com/BurntSushi/toml",
        sum = "h1:pxW6RcqyfI9/kWtOwnv/G+AzdKuy2ZrqINhenH4HyNs=",
        version = "v1.4.1-0.20240526193622-a339e1f7089c",
    )
    go_repository(
        name = "com_github_carlmontanari_difflibgo",
        importpath = "github.com/carlmontanari/difflibgo",
        sum = "h1:Y8rTum6LZ8oP/2aC+OaaP76OCjHbunKMkim81mzNCH0=",
        version = "v0.0.0-20210718194309-31b9e131c298",
    )
    go_repository(
        name = "com_github_cenkalti_backoff_v4",
        importpath = "github.com/cenkalti/backoff/v4",
        sum = "h1:MyRJ/UdXutAwSAT+s3wNd7MfTIcy71VQueUuFK343L8=",
        version = "v4.3.0",
    )
    go_repository(
        name = "com_github_census_instrumentation_opencensus_proto",
        importpath = "github.com/census-instrumentation/opencensus-proto",
        sum = "h1:iKLQ0xPNFxR/2hzXZMrBo8f1j86j5WHzznCCQxV/b8g=",
        version = "v0.4.1",
    )
    go_repository(
        name = "com_github_cespare_xxhash_v2",
        importpath = "github.com/cespare/xxhash/v2",
        sum = "h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=",
        version = "v2.3.0",
    )
    go_repository(
        name = "com_github_chappjc_logrus_prefix",
        importpath = "github.com/chappjc/logrus-prefix",
        sum = "h1:aZTKxMminKeQWHtzJBbV8TttfTxzdJ+7iEJFE6FmUzg=",
        version = "v0.0.0-20180227015900-3a1d64819adb",
    )
    go_repository(
        name = "com_github_chzyer_readline",
        importpath = "github.com/chzyer/readline",
        sum = "h1:upd/6fQk4src78LMRzh5vItIt361/o4uq553V8B5sGI=",
        version = "v1.5.1",
    )
    go_repository(
        name = "com_github_client9_misspell",
        importpath = "github.com/client9/misspell",
        sum = "h1:ta993UF76GwbvJcIo3Y68y/M3WxlpEHPWIGDkJYwzJI=",
        version = "v0.3.4",
    )
    go_repository(
        name = "com_github_cloudflare_circl",
        importpath = "github.com/cloudflare/circl",
        sum = "h1:zqIqSPIndyBh1bjLVVDHMPpVKqp8Su/V+6MeDzzQBQ0=",
        version = "v1.6.1",
    )
    go_repository(
        name = "com_github_cncf_udpa_go",
        importpath = "github.com/cncf/udpa/go",
        sum = "h1:WBZRG4aNOuI15bLRrCgN8fCq8E5Xuty6jGbmSNEvSsU=",
        version = "v0.0.0-20191209042840-269d4d468f6f",
    )
    go_repository(
        name = "com_github_cncf_xds_go",
        importpath = "github.com/cncf/xds/go",
        sum = "h1:aBangftG7EVZoUb69Os8IaYg++6uMOdKK83QtkkvJik=",
        version = "v0.0.0-20260202195803-dba9d589def2",
    )
    go_repository(
        name = "com_github_containernetworking_cni",
        importpath = "github.com/containernetworking/cni",
        sum = "h1:v6EpN8RznAZj9765HhXQrtXgX+ECGebEYEmnuFjskwo=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_containernetworking_plugins",
        importpath = "github.com/containernetworking/plugins",
        sum = "h1:8oU6WsIsU3bpnNZuvHp74a6cE1MJwbj2P7s4/yTUNlA=",
        version = "v1.9.1",
    )
    go_repository(
        name = "com_github_coredhcp_coredhcp",
        importpath = "github.com/coredhcp/coredhcp",
        sum = "h1:Bj+yhDYB9fypuXEJimRf9iRqHenZ/lqR7OPJ8e4GGQs=",
        version = "v0.0.0-20260217182248-a0841cb3038f",
    )
    go_repository(
        name = "com_github_coreos_go_semver",
        importpath = "github.com/coreos/go-semver",
        sum = "h1:yi21YpKnrx1gt5R+la8n5WgS0kCrsPp33dmEyHReZr4=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_coreos_go_systemd_v22",
        importpath = "github.com/coreos/go-systemd/v22",
        sum = "h1:LAEzFkke61DFROc7zNLX/WA2i5J8gYqe0rSj9KI28KA=",
        version = "v22.7.0",
    )
    go_repository(
        name = "com_github_cpuguy83_go_md2man_v2",
        importpath = "github.com/cpuguy83/go-md2man/v2",
        sum = "h1:XJtiaUW6dEEqVuZiMTn1ldk455QWwEIsMIJlo5vtkx0=",
        version = "v2.0.6",
    )
    go_repository(
        name = "com_github_creack_pty",
        importpath = "github.com/creack/pty",
        sum = "h1:bJrF4RRfyJnbTJqzRLHzcGaZK1NeM5kTC9jGgovnR1s=",
        version = "v1.1.24",
    )
    go_repository(
        name = "com_github_cyphar_filepath_securejoin",
        importpath = "github.com/cyphar/filepath-securejoin",
        sum = "h1:6iR5tXJ/e6tJZzzdMc1km3Sa7RRIVBKAK32O2s7AYfo=",
        version = "v0.2.5",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=",
        version = "v1.1.2-0.20180830191138-d8f796af33cc",
    )
    go_repository(
        name = "com_github_dgryski_go_farm",
        importpath = "github.com/dgryski/go-farm",
        sum = "h1:aIftn67I1fkbMa512G+w+Pxci9hJPB8oMnkcP3iZF38=",
        version = "v0.0.0-20240924180020-3414d57e47da",
    )
    go_repository(
        name = "com_github_docker_docker",
        importpath = "github.com/docker/docker",
        sum = "h1:KttF0XoteNTicmUtBO0L2tP+J7FGRFTjaEF4k6WdhfI=",
        version = "v27.3.1+incompatible",
    )
    go_repository(
        name = "com_github_docker_go_connections",
        importpath = "github.com/docker/go-connections",
        sum = "h1:JibmG5hULs5qXSr/cp/w3Pw5fZuStt4MOHMUExb29/M=",
        version = "v0.8.1",
    )
    go_repository(
        name = "com_github_docker_go_units",
        importpath = "github.com/docker/go-units",
        sum = "h1:69rxXcBk27SvSaaxTtLh/8llcHD8vYHT7WSdRZ/jvr4=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_drivenets_cdnos_controller",
        importpath = "github.com/drivenets/cdnos-controller",
        sum = "h1:HVzIZ1yfpnfgCaxgUNRtZpNp2+W6TNagkwyijXmGFm0=",
        version = "v1.7.9",
    )
    go_repository(
        name = "com_github_dustin_go_humanize",
        importpath = "github.com/dustin/go-humanize",
        sum = "h1:GzkhY7T5VNhEkwH0PVJgjz+fX1rhBrR7pRT3mDkpeCY=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_eapache_channels",
        importpath = "github.com/eapache/channels",
        sum = "h1:F1taHcn7/F0i8DYqKXJnyhJcVpp2kgFcNePxXtnyu4k=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_eapache_queue",
        importpath = "github.com/eapache/queue",
        sum = "h1:YOEu7KNc61ntiQlcEeUIoDTJ2o8mQznoNvUhiigpIqc=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_emicklei_go_restful",
        importpath = "github.com/emicklei/go-restful",
        sum = "h1:spTtZBk5DYEvbxMVutUuTyh1Ao2r4iyvLdACqsl/Ljk=",
        version = "v2.9.5+incompatible",
    )
    go_repository(
        name = "com_github_emicklei_go_restful_v3",
        importpath = "github.com/emicklei/go-restful/v3",
        sum = "h1:C4Bl2xDndpU6nJ4bc1jXd+uTmYPVUwkD6bFY/oTyCes=",
        version = "v3.13.0",
    )
    go_repository(
        name = "com_github_emirpasic_gods",
        importpath = "github.com/emirpasic/gods",
        sum = "h1:FXtiHYKDGKCW2KzwZKx0iC0PQmdlorYgdFG9jPXJ1Bc=",
        version = "v1.18.1",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane",
        importpath = "github.com/envoyproxy/go-control-plane",
        sum = "h1:hbG2kr4RuFj222B6+7T83thSPqLjwBIfQawTkC++2HA=",
        version = "v0.14.0",
    )
    go_repository(
        name = "com_github_envoyproxy_protoc_gen_validate",
        importpath = "github.com/envoyproxy/protoc-gen-validate",
        sum = "h1:MVQghNeW+LZcmXe7SY1V36Z+WFMDjpqGAGacLe2T0ds=",
        version = "v1.3.3",
    )
    go_repository(
        name = "com_github_evanphx_json_patch",
        importpath = "github.com/evanphx/json-patch",
        sum = "h1:jBYDEEiFBPxA0v50tFdvOzQQTCvpL6mnFh5mB2/l16U=",
        version = "v5.6.0+incompatible",
    )
    go_repository(
        name = "com_github_evanphx_json_patch_v5",
        importpath = "github.com/evanphx/json-patch/v5",
        sum = "h1:/8HVnzMq13/3x9TPvjG08wUGqBTmZBsCWzjTM0wiaDU=",
        version = "v5.9.11",
    )
    go_repository(
        name = "com_github_fatih_color",
        importpath = "github.com/fatih/color",
        sum = "h1:Zp3PiM21/9Ld6FzSKyL5c/BULoe/ONr9KlbYVOfG8+w=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_github_felixge_httpsnoop",
        importpath = "github.com/felixge/httpsnoop",
        sum = "h1:3YtUj32ZZkqZtt3sZZsClsymw/QDuVfpNhoA31zeORc=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_form3tech_oss_jwt_go",
        importpath = "github.com/form3tech-oss/jwt-go",
        sum = "h1:7ZaBxOI7TMoYBfyA3cQHErNNyAWIKUMIwqxEtgHOs5c=",
        version = "v3.2.3+incompatible",
    )
    go_repository(
        name = "com_github_frankban_quicktest",
        importpath = "github.com/frankban/quicktest",
        sum = "h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=",
        version = "v1.14.6",
    )
    go_repository(
        name = "com_github_fsnotify_fsnotify",
        importpath = "github.com/fsnotify/fsnotify",
        sum = "h1:b0/UzAf9yR5rhf3RPm9gf3ehBPpf0oZKIjtpKrx59Ho=",
        version = "v1.10.1",
    )
    go_repository(
        name = "com_github_getsentry_sentry_go",
        importpath = "github.com/getsentry/sentry-go",
        sum = "h1:lWUwDnY7sKHaVIoZ9wYqRHJ5iEmoc0pqcRqFkosKzBo=",
        version = "v0.30.0",
    )
    go_repository(
        name = "com_github_ghodss_yaml",
        importpath = "github.com/ghodss/yaml",
        sum = "h1:wQHKEahhL6wmXdzwWG11gIVCkOv05bNOh+Rxn0yngAk=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_go_git_gcfg",
        importpath = "github.com/go-git/gcfg",
        sum = "h1:+zs/tPmkDkHx3U66DAb0lQFJrpS6731Oaa12ikc+DiI=",
        version = "v1.5.1-0.20230307220236-3a3c6141e376",
    )
    go_repository(
        name = "com_github_go_git_go_billy_v5",
        importpath = "github.com/go-git/go-billy/v5",
        sum = "h1:w2hPNtoehvJIxR00Vb4xX94qHQi/ApZfX+nBE2Cjio8=",
        version = "v5.6.0",
    )
    go_repository(
        name = "com_github_go_git_go_git_v5",
        importpath = "github.com/go-git/go-git/v5",
        sum = "h1:vLn5wlGIh/X78El6r3Jr+30W16Blk0CTcxTYcYPWi5E=",
        version = "v5.13.0",
    )
    go_repository(
        name = "com_github_go_logr_logr",
        importpath = "github.com/go-logr/logr",
        sum = "h1:tG4xh9yMsRCAiodLVTxyrkzSZ9+o0L1Kg/+cPVcbP/8=",
        version = "v1.4.4",
    )
    go_repository(
        name = "com_github_go_logr_stdr",
        importpath = "github.com/go-logr/stdr",
        sum = "h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_go_logr_zapr",
        importpath = "github.com/go-logr/zapr",
        sum = "h1:XGdV8XW8zdwFiwOA2Dryh1gj2KRQyOOoNmBy4EplIcQ=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_go_openapi_jsonpointer",
        importpath = "github.com/go-openapi/jsonpointer",
        sum = "h1:2KxywRmNwJkT/FMBa3iRNHEaAxSJvjqoufQZy3au1Mg=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_go_openapi_jsonreference",
        importpath = "github.com/go-openapi/jsonreference",
        sum = "h1:4zJ7AmYDKNmD3aSpfPnFNCFA5E80/xMHUNKgydaLh38=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_go_openapi_swag",
        importpath = "github.com/go-openapi/swag",
        sum = "h1:9n8frkcsuQRA0INU31VbxNxt1zBaHxrHYGEf5BUjFi4=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_playground_locales",
        importpath = "github.com/go-playground/locales",
        sum = "h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=",
        version = "v0.14.1",
    )
    go_repository(
        name = "com_github_go_playground_universal_translator",
        importpath = "github.com/go-playground/universal-translator",
        sum = "h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=",
        version = "v0.18.1",
    )
    go_repository(
        name = "com_github_go_playground_validator_v10",
        importpath = "github.com/go-playground/validator/v10",
        sum = "h1:K9ISHbSaI0lyB2eWMPJo+kOS/FBExVwjEviJTixqxL8=",
        version = "v10.20.0",
    )
    go_repository(
        name = "com_github_go_test_deep",
        importpath = "github.com/go-test/deep",
        sum = "h1:0r/53hagsehfO4bzD2Pgr/+RgHqhmf+k1Bpse2cTu1U=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_goccy_go_json",
        importpath = "github.com/goccy/go-json",
        sum = "h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=",
        version = "v0.10.2",
    )
    go_repository(
        name = "com_github_goccy_go_yaml",
        importpath = "github.com/goccy/go-yaml",
        sum = "h1:8W7wMFS12Pcas7KU+VVkaiCng+kG8QiFeFwzFb+rwuw=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_github_gogo_protobuf",
        importpath = "github.com/gogo/protobuf",
        sum = "h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=",
        version = "v1.3.2",
    )
    go_repository(
        name = "com_github_golang_freetype",
        importpath = "github.com/golang/freetype",
        sum = "h1:DACJavvAHhabrF08vX0COfcOBJRhZ8lUbR+ZWIs0Y5g=",
        version = "v0.0.0-20170609003504-e2365dfdc4a0",
    )
    go_repository(
        name = "com_github_golang_glog",
        importpath = "github.com/golang/glog",
        sum = "h1:DrW6hGnjIhtvhOIiAKT6Psh/Kd/ldepEa81DKeiRJ5I=",
        version = "v1.2.5",
    )
    go_repository(
        name = "com_github_golang_groupcache",
        importpath = "github.com/golang/groupcache",
        sum = "h1:f+oWsMOmNPc8JmEHVZIycC7hBoQxHH9pNKQORJNozsQ=",
        version = "v0.0.0-20241129210726-2c02b8208cf8",
    )
    go_repository(
        name = "com_github_golang_mock",
        importpath = "github.com/golang/mock",
        sum = "h1:YojYx61/OLFsiv6Rw1Z96LpldJIy31o+UHmwAUMJ6/U=",
        version = "v1.7.0-rc.1",
    )
    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=",
        version = "v1.5.4",
    )
    go_repository(
        name = "com_github_golang_snappy",
        importpath = "github.com/golang/snappy",
        sum = "h1:yAGX7huGHXlcLOEtBnF4w7FQwA26wojNCwOYAEhLjQM=",
        version = "v0.0.4",
    )
    go_repository(
        name = "com_github_google_btree",
        importpath = "github.com/google/btree",
        sum = "h1:CVpQJjYgC4VbzxeGVHfvZrv1ctoYCAI8vbl07Fcxlyg=",
        version = "v1.1.3",
    )
    go_repository(
        name = "com_github_google_cel_go",
        importpath = "github.com/google/cel-go",
        sum = "h1:ZtDxkeiMmz0mxbKDYiNkE5Lk7V5edMRcaaDf2jX002k=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_google_flatbuffers",
        importpath = "github.com/google/flatbuffers",
        sum = "h1:M9dgRyhJemaM4Sw8+66GHBu8ioaQmyPLg1b8VwK5WJg=",
        version = "v23.5.26+incompatible",
    )
    go_repository(
        name = "com_github_google_gnostic",
        build_file_proto_mode = "disable",
        importpath = "github.com/google/gnostic",
        sum = "h1:ZK/5VhkoX835RikCHpSUJV9a+S3e1zLh59YnyWeBW+0=",
        version = "v0.6.9",
    )
    go_repository(
        name = "com_github_google_gnostic_models",
        build_file_proto_mode = "disable",
        importpath = "github.com/google/gnostic-models",
        sum = "h1:SisTfuFKJSKM5CPZkffwi6coztzzeYUhc3v4yxLWH8c=",
        version = "v0.7.1",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_github_google_go_pkcs11",
        importpath = "github.com/google/go-pkcs11",
        sum = "h1:PVRnTgtArZ3QQqTGtbtjtnIkzl2iY2kt24yqbrf7td8=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_google_gofuzz",
        importpath = "github.com/google/gofuzz",
        sum = "h1:xRy4A+RhZaiKjJ1bPfwQ8sedCA+YS2YcCHW6ec7JMi0=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_google_gopacket",
        importpath = "github.com/google/gopacket",
        sum = "h1:ves8RnFZPGiFnTS0uPQStjwru6uO6h+nlr9j6fL7kF8=",
        version = "v1.1.19",
    )
    go_repository(
        name = "com_github_google_martian_v3",
        importpath = "github.com/google/martian/v3",
        sum = "h1:DIhPTQrbPkgs2yJYdXU/eNACCG5DVQjySNRNlflZ9Fc=",
        version = "v3.3.3",
    )
    go_repository(
        name = "com_github_google_pprof",
        importpath = "github.com/google/pprof",
        sum = "h1:BHT72Gu3keYf3ZEu2J0b1vyeLSOYI8bm5wbJM/8yDe8=",
        version = "v0.0.0-20250403155104-27863c87afa6",
    )
    go_repository(
        name = "com_github_google_s2a_go",
        importpath = "github.com/google/s2a-go",
        sum = "h1:EMp+aOuXN6l8cE/gjF5Bt+vyZxsUuyCWe9chDWR/+uU=",
        version = "v0.1.10",
    )
    go_repository(
        name = "com_github_google_uuid",
        importpath = "github.com/google/uuid",
        sum = "h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_googleapis_enterprise_certificate_proxy",
        importpath = "github.com/googleapis/enterprise-certificate-proxy",
        sum = "h1:NU4XpII6jD+Dxcot94fqjE+AfJoE/lQP9q3faYGzC/c=",
        version = "v0.3.22",
    )
    go_repository(
        name = "com_github_googleapis_gax_go_v2",
        build_file_proto_mode = "disable",
        importpath = "github.com/googleapis/gax-go/v2",
        sum = "h1:ydkmNXxj7bEmmeK5AihkKnWxyOyBR9TDebvp5L5izk8=",
        version = "v2.26.2",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_detectors_gcp",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp",
        sum = "h1:3SdxXLkgAfiHRWcGTq6fneq9jgoJzneiY0yPQnjoT2E=",
        version = "v1.36.0",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_exporter_metric",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric",
        sum = "h1:5y2LR2pripbccOte6BpZ+guv5O5W0OQIfKvUYZrGaAQ=",
        version = "v0.62.0",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_exporter_trace",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace",
        sum = "h1:Smjr+Zjtrz0GEI2AxvHxew0V1vkDlvcm+0hyArVNkKY=",
        version = "v1.38.0",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_internal_cloudmock",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/cloudmock",
        sum = "h1:tRsXhJnn0w0TNwJpVrB94wsrimzjWKFzdQlUhRdNUzc=",
        version = "v0.62.0",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_internal_resourcemapping",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping",
        sum = "h1:aPMa5qZU9dp4fU34o+jQgVIPFIOi20oxYg8TLigmZLg=",
        version = "v0.62.0",
    )
    go_repository(
        name = "com_github_gorilla_mux",
        importpath = "github.com/gorilla/mux",
        sum = "h1:TuBL49tXwgrFYWhqrNgrUNEY92u81SPhu7sTdzQEiWY=",
        version = "v1.8.1",
    )
    go_repository(
        name = "com_github_gorilla_websocket",
        importpath = "github.com/gorilla/websocket",
        sum = "h1:JeSE6pjso5THxAzdVpqr6/geYxZytqFMBCOtn/ujyeo=",
        version = "v1.5.4-0.20250319132907-e064f32e3674",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_go_grpc_middleware",
        importpath = "github.com/grpc-ecosystem/go-grpc-middleware",
        sum = "h1:UH//fgunKIs4JdUbpDl1VZCDaL56wXCB/5+wF6uHfaI=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_go_grpc_middleware_v2",
        importpath = "github.com/grpc-ecosystem/go-grpc-middleware/v2",
        sum = "h1:9ZJYjPEJpcleIKcqysXOo14NLQYwu23fzmKRFpIjPjc=",
        version = "v2.3.4",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_go_grpc_prometheus",
        importpath = "github.com/grpc-ecosystem/go-grpc-prometheus",
        sum = "h1:Ovs26xHkKqVztRpIrF/92BcuyuQ/YW4NSIpoGtfXNho=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_grpc_gateway_v2",
        importpath = "github.com/grpc-ecosystem/grpc-gateway/v2",
        sum = "h1:Bd7KaOxzULLxtZ/K5s1aLbWhR0+5RToO65TXHsf3bqQ=",
        version = "v2.31.0",
    )
    go_repository(
        name = "com_github_h_fam_errdiff",
        importpath = "github.com/h-fam/errdiff",
        sum = "h1:rPsW4ob2fMOIulwTEoZXaaUIuud7XUudw5SLKTZj3Ss=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_hashicorp_go_version",
        importpath = "github.com/hashicorp/go-version",
        sum = "h1:CeOIz6k+LoN3qX9Z0tyQrPtiB1DFYRPfCIBtaXPSCnA=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_github_hashicorp_hcl",
        importpath = "github.com/hashicorp/hcl",
        sum = "h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_iancoleman_strcase",
        importpath = "github.com/iancoleman/strcase",
        sum = "h1:nTXanmYxhfFAMjZL34Ov6gkzEsSJZ5DbhxWjvSASxEI=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_ianlancetaylor_demangle",
        importpath = "github.com/ianlancetaylor/demangle",
        sum = "h1:KwWnWVWCNtNq/ewIX7HIKnELmEx2nDP42yskD/pi7QE=",
        version = "v0.0.0-20240312041847-bd984b5ce465",
    )
    go_repository(
        name = "com_github_imdario_mergo",
        importpath = "github.com/imdario/mergo",
        sum = "h1:wwQJbIsHYGMUyLSPrEq1CT16AhnhNJQ51+4fdHUnCl4=",
        version = "v0.3.16",
    )
    go_repository(
        name = "com_github_inconshreveable_mousetrap",
        importpath = "github.com/inconshreveable/mousetrap",
        sum = "h1:wN+x4NVGpMsO7ErUn/mUI3vEoE6Jt13X2s0bqwp9tc8=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_insomniacslk_dhcp",
        importpath = "github.com/insomniacslk/dhcp",
        sum = "h1:G3irkWwmpl0vH/nn83K2AHqLUZweC7XAONuwXy/w9Co=",
        version = "v0.0.0-20260407060928-11b94ed970f2",
    )
    go_repository(
        name = "com_github_jbenet_go_context",
        importpath = "github.com/jbenet/go-context",
        sum = "h1:BQSFePA1RWJOlocH6Fxy8MmwDt+yVQYULKfN0RoTN8A=",
        version = "v0.0.0-20150711004518-d14ea06fba99",
    )
    go_repository(
        name = "com_github_jessevdk_go_flags",
        importpath = "github.com/jessevdk/go-flags",
        sum = "h1:Cvu5U8UGrLay1rZfv/zP7iLpSHGUZ/Ou68T0iX1bBK4=",
        version = "v1.6.1",
    )
    go_repository(
        name = "com_github_jonboulle_clockwork",
        importpath = "github.com/jonboulle/clockwork",
        sum = "h1:Hyh9A8u51kptdkR+cqRpT1EebBwTn1oK9YfGYbdFz6I=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_josharian_intern",
        importpath = "github.com/josharian/intern",
        sum = "h1:vlS4z54oSdjm0bgjRigI+G1HpF+tI+9rE5LLzOg8HmY=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_json_iterator_go",
        importpath = "github.com/json-iterator/go",
        sum = "h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=",
        version = "v1.1.12",
    )
    go_repository(
        name = "com_github_jstemmer_go_junit_report_v2",
        importpath = "github.com/jstemmer/go-junit-report/v2",
        sum = "h1:X3+hPYlSczH9IMIpSC9CQSZA0L+BipYafciZUWHEmsc=",
        version = "v2.1.0",
    )
    go_repository(
        name = "com_github_k_sone_critbitgo",
        importpath = "github.com/k-sone/critbitgo",
        sum = "h1:l71cTyBGeh6X5ATh6Fibgw3+rtNT80BA0uNNWgkPrbE=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_kentik_patricia",
        importpath = "github.com/kentik/patricia",
        sum = "h1:P2eKb6UchQWqxm1sXfWH8S/r7HeKQLjet4J9zTVK6vE=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_kevinburke_ssh_config",
        importpath = "github.com/kevinburke/ssh_config",
        sum = "h1:x584FjTGwHzMwvHx18PXxbBVzfnxogHaAReU4gf13a4=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_klauspost_compress",
        importpath = "github.com/klauspost/compress",
        sum = "h1:sXLILfc9jV2QYWkzFOPWStmcUVH2RHEB1JCdY2oVvCQ=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_github_klauspost_cpuid_v2",
        importpath = "github.com/klauspost/cpuid/v2",
        sum = "h1:ZWSB3igEs+d0qvnxR/ZBzXVmxkgt8DdzP6m9pfuVLDM=",
        version = "v2.2.7",
    )
    go_repository(
        name = "com_github_kr_pretty",
        importpath = "github.com/kr/pretty",
        sum = "h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_kr_pty",
        importpath = "github.com/kr/pty",
        sum = "h1:VkoXIwSboBpnk99O/KFauAEILuNHv5DVFKZMBN/gUgw=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_kr_text",
        importpath = "github.com/kr/text",
        sum = "h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_kylelemons_godebug",
        importpath = "github.com/kylelemons/godebug",
        sum = "h1:RPNrshWIDI6G2gRW9EHilWtl7Z6Sb1BR0xunSBf0SNc=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_leodido_go_urn",
        importpath = "github.com/leodido/go-urn",
        sum = "h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_lyft_protoc_gen_star_v2",
        importpath = "github.com/lyft/protoc-gen-star/v2",
        sum = "h1:JDlNKttNIRd68AAIychs0AqEpO8/I/WYi01OQ7Raw6Q=",
        version = "v2.0.4",
    )
    go_repository(
        name = "com_github_magiconair_properties",
        importpath = "github.com/magiconair/properties",
        sum = "h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=",
        version = "v1.8.7",
    )
    go_repository(
        name = "com_github_mailru_easyjson",
        importpath = "github.com/mailru/easyjson",
        sum = "h1:UGYAvKxe3sBsEDzO8ZeWOSlIQfWFlxbzLZe7hwFURr0=",
        version = "v0.7.7",
    )
    go_repository(
        name = "com_github_masterminds_semver_v3",
        importpath = "github.com/Masterminds/semver/v3",
        sum = "h1:kQceYJfbupGfZOKZQg0kou0DgAKhzDg2NZPAwZ/2OOE=",
        version = "v3.5.0",
    )
    go_repository(
        name = "com_github_mattn_go_colorable",
        importpath = "github.com/mattn/go-colorable",
        sum = "h1:9A9LHSqF/7dyVVX6g0U9cwm9pG3kP9gSzcuIPHPsaIE=",
        version = "v0.1.14",
    )
    go_repository(
        name = "com_github_mattn_go_isatty",
        importpath = "github.com/mattn/go-isatty",
        sum = "h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=",
        version = "v0.0.20",
    )
    go_repository(
        name = "com_github_matttproud_golang_protobuf_extensions",
        importpath = "github.com/matttproud/golang_protobuf_extensions",
        sum = "h1:mmDVorXM7PCGKw94cs5zkfA9PSy5pEvNWRP0ET0TIVo=",
        version = "v1.0.4",
    )
    go_repository(
        name = "com_github_mdlayher_genetlink",
        importpath = "github.com/mdlayher/genetlink",
        sum = "h1:f/Xs7Y2T+GyX9b3dbiUhnLE9InGs5F9RxJ2JwBMl71o=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_mdlayher_netlink",
        importpath = "github.com/mdlayher/netlink",
        sum = "h1:G8+GLq2x3v4D4MVIqDdNUhTUC7TKiCy/6MDkmItfKco=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_github_mdlayher_socket",
        importpath = "github.com/mdlayher/socket",
        sum = "h1:ScZPaAGyO1icQnbFrhPM8mnXyMu9qukC1K4ZoM2IQKU=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_github_mgutz_ansi",
        importpath = "github.com/mgutz/ansi",
        sum = "h1:5PJl274Y63IEHC+7izoQE9x6ikvDFZS2mDVS3drnohI=",
        version = "v0.0.0-20200706080929-d51e80ef957d",
    )
    go_repository(
        name = "com_github_microsoft_go_winio",
        importpath = "github.com/Microsoft/go-winio",
        sum = "h1:0kQAzHq8vLs7Pptv+7TxjdETLf/nIqJpIB4oC6Ba4vY=",
        version = "v0.6.3-0.20251027160822-ad3df93bed29",
    )
    go_repository(
        name = "com_github_mitchellh_go_wordwrap",
        importpath = "github.com/mitchellh/go-wordwrap",
        sum = "h1:TLuKupo69TCn6TQSyGxwI1EblZZEsQ0vMlAFQflz0v0=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_mitchellh_mapstructure",
        importpath = "github.com/mitchellh/mapstructure",
        sum = "h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_moby_spdystream",
        importpath = "github.com/moby/spdystream",
        sum = "h1:9sNYeYZUcci9R6/w7KDaFWEWeV4LStVG78Mpyq/Zm/Y=",
        version = "v0.5.1",
    )
    go_repository(
        name = "com_github_modern_go_concurrent",
        importpath = "github.com/modern-go/concurrent",
        sum = "h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=",
        version = "v0.0.0-20180306012644-bacd9c7ef1dd",
    )
    go_repository(
        name = "com_github_modern_go_reflect2",
        importpath = "github.com/modern-go/reflect2",
        sum = "h1:W5t00kpgFdJifH4BDsTlE89Zl93FEloxaWZfGcifgq8=",
        version = "v1.0.3-0.20250322232337-35a7c28c31ee",
    )
    go_repository(
        name = "com_github_munnerz_goautoneg",
        importpath = "github.com/munnerz/goautoneg",
        sum = "h1:C3w9PqII01/Oq1c1nUAm88MOHcQC9l5mIlSMApZMrHA=",
        version = "v0.0.0-20191010083416-a7dc8b61c822",
    )
    go_repository(
        name = "com_github_mxk_go_flowrate",
        importpath = "github.com/mxk/go-flowrate",
        sum = "h1:y5//uYreIhSUg3J1GEMiLbxo1LJaP8RfCpH6pymGZus=",
        version = "v0.0.0-20140419014527-cca7078d478f",
    )
    go_repository(
        name = "com_github_networkop_meshnet_cni",
        importpath = "github.com/networkop/meshnet-cni",
        sum = "h1:9Pe9L0QCovb9o82inAVQitCo3IRnG9u45lRRm8QvgbU=",
        version = "v0.3.1-0.20230525201116-d7c306c635cf",
    )
    go_repository(
        name = "com_github_nxadm_tail",
        importpath = "github.com/nxadm/tail",
        sum = "h1:nPr65rt6Y5JFSKQO7qToXr7pePgD6Gwiw05lkbyAQTE=",
        version = "v1.4.8",
    )
    go_repository(
        name = "com_github_nytimes_gziphandler",
        importpath = "github.com/NYTimes/gziphandler",
        sum = "h1:ZUDjpQae29j0ryrS0u/B8HZfJBtBQHjqw2rQ2cqUQ3I=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_onsi_ginkgo",
        importpath = "github.com/onsi/ginkgo",
        sum = "h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=",
        version = "v1.16.5",
    )
    go_repository(
        name = "com_github_onsi_ginkgo_v2",
        importpath = "github.com/onsi/ginkgo/v2",
        sum = "h1:fcEcQW/A++6aZAZQNUmNjvA9PSOzefMJBerHJ4t8v8Y=",
        version = "v2.27.4",
    )
    go_repository(
        name = "com_github_onsi_gomega",
        importpath = "github.com/onsi/gomega",
        sum = "h1:y2ROC3hKFmQZJNFeGAMeHZKkjBL65mIZcvrLQBF9k6Q=",
        version = "v1.39.0",
    )
    go_repository(
        name = "com_github_open_traffic_generator_keng_operator",
        importpath = "github.com/open-traffic-generator/keng-operator",
        sum = "h1:7foRgBvl1SPcBFMovjs3Nch6Or8AuqnGdk0fK6Wl/1c=",
        version = "v0.4.2",
    )
    go_repository(
        name = "com_github_open_traffic_generator_snappi_gosnappi",
        build_file_proto_mode = "disable",
        importpath = "github.com/open-traffic-generator/snappi/gosnappi",
        sum = "h1:zUqpuJtW0+Gtz3dO0KzrRKWyg+86iDWo7eLDhpe57To=",
        version = "v1.61.1",
    )
    go_repository(
        name = "com_github_openconfig_attestz",
        importpath = "github.com/openconfig/attestz",
        sum = "h1:m8UOHvV+r8eDahJdnouQmDhkx/ZmkHatfoexQZ5pKdg=",
        version = "v0.6.15",
    )
    go_repository(
        name = "com_github_openconfig_bootz",
        importpath = "github.com/openconfig/bootz",
        sum = "h1:/dNcCgTcE2QPA/K7D4Zlwu+E/MeGLN6bpHHI051HUns=",
        version = "v0.7.1",
    )
    go_repository(
        name = "com_github_openconfig_entity_naming",
        importpath = "github.com/openconfig/entity-naming",
        sum = "h1:FQwuGeqYG6dSIMawOe8aZqBT8TDFxA33vFpOHkHfW7I=",
        version = "v0.0.0-20251204192329-8cf2fdebf3c1",
    )
    go_repository(
        name = "com_github_openconfig_featureprofiles",
        importpath = "github.com/openconfig/featureprofiles",
        sum = "h1:iHrFMxf5P10wI0ll9XVD6B4ICFeVSnU5nUNXUvjLayA=",
        version = "v0.0.0-20260310133626-0ee97ec253f3",
    )
    go_repository(
        name = "com_github_openconfig_gnmi",
        build_directives = ["gazelle:proto_import_prefix github.com/openconfig/gnmi"],
        build_file_generation = "on",
        importpath = "github.com/openconfig/gnmi",
        sum = "h1:qKMuFvhIRR2/xxCOsStPQ25aKpbMDdWr3kI+nP9bhMs=",
        version = "v0.14.1",
    )
    go_repository(
        name = "com_github_openconfig_gnoi",
        build_directives = [
            "gazelle:prefix github.com/openconfig/gnoi",
            "gazelle:proto_import_prefix github.com/openconfig/gnoi",
            "gazelle:resolve proto proto github.com/openconfig/bootz/proto/bootz.proto @com_github_openconfig_bootz//proto:bootz_proto",
            "gazelle:resolve proto proto github.com/openconfig/gnsi/certz/certz.proto @com_github_openconfig_gnsi//certz:certz_proto",
            "gazelle:resolve proto proto github.com/openconfig/gnsi/pathz/pathz.proto @com_github_openconfig_gnsi//pathz:pathz_proto",
            "gazelle:resolve proto proto github.com/openconfig/gnsi/authz/authz.proto @com_github_openconfig_gnsi//authz:authz_proto",
            "gazelle:resolve proto go github.com/openconfig/bootz/proto/bootz.proto @com_github_openconfig_bootz//proto:bootz_go_proto",
            "gazelle:resolve proto go github.com/openconfig/gnsi/certz/certz.proto @com_github_openconfig_gnsi//certz:certz_go_proto",
            "gazelle:resolve proto go github.com/openconfig/gnsi/pathz/pathz.proto @com_github_openconfig_gnsi//pathz:pathz_go_proto",
            "gazelle:resolve proto go github.com/openconfig/gnsi/authz/authz.proto @com_github_openconfig_gnsi//authz:authz_go_proto",
            "gazelle:resolve proto proto google/rpc/status.proto @googleapis//google/rpc:status_proto",
            "gazelle:resolve proto go google/rpc/status.proto  @org_golang_google_genproto_googleapis_rpc//status",
        ],
        build_file_generation = "on",
        importpath = "github.com/openconfig/gnoi",
        sum = "h1:fwZm4zlwoY5i7KALTpVhpAv53Y3YskleoTpg1IUCa+c=",
        version = "v0.8.0",
    )
    go_repository(
        name = "com_github_openconfig_gnoigo",
        importpath = "github.com/openconfig/gnoigo",
        sum = "h1:bHMTFwX3NL1oMmh+UN6RA2i2NunIKGtQM76s364QX64=",
        version = "v0.0.0-20250918224707-fee0fe3eee56",
    )
    go_repository(
        name = "com_github_openconfig_gnsi",
        importpath = "github.com/openconfig/gnsi",
        sum = "h1:0XpYlG/99YWVIm6gFTVkxlwiXcU/tBla9MDLChLNBMM=",
        version = "v1.9.1",
    )
    go_repository(
        name = "com_github_openconfig_gocloser",
        importpath = "github.com/openconfig/gocloser",
        sum = "h1:HsTKOBZeQw/Rt8KmxFY578pVZDIgeem50jzv4BjrTR8=",
        version = "v0.0.0-20251119232641-34bca749fdb3",
    )
    go_repository(
        name = "com_github_openconfig_goyang",
        importpath = "github.com/openconfig/goyang",
        sum = "h1:9nWXBwd6b4+nZr8ni7O4zUXVhrVMXCLFz8os5YWFuo4=",
        version = "v1.6.3",
    )
    go_repository(
        name = "com_github_openconfig_gribi",
        build_file_proto_mode = "disable",
        importpath = "github.com/openconfig/gribi",
        sum = "h1:lN5KMEKu+X5L24XVHnwMJe52De5CSNZb1rR0F/ym13M=",
        version = "v1.9.1",
    )
    go_repository(
        name = "com_github_openconfig_gribigo",
        importpath = "github.com/openconfig/gribigo",
        sum = "h1:1ahvq4oZ6djCigKVVjqAB/tdCSvh94kl4lqyzDE2s/8=",
        version = "v0.1.3",
    )
    go_repository(
        name = "com_github_openconfig_grpctunnel",
        importpath = "github.com/openconfig/grpctunnel",
        sum = "h1:EN99qtlExZczgQgp5ANnHRC/Rs62cAG+Tz2BQ5m/maM=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_openconfig_kne",
        importpath = "github.com/openconfig/kne",
        sum = "h1:1KYNrvVlQkdfh7QhCC7EVi9VTS6TAIEEqFi1T8hrb94=",
        version = "v0.3.2",
    )
    go_repository(
        name = "com_github_openconfig_magna",
        importpath = "github.com/openconfig/magna",
        sum = "h1:+3I5lnawPSR1ZqPeb3BudQBUp6b/uhNz1u+kNnVgAS0=",
        version = "v0.0.0-20260408002632-b210de0f7ba8",
    )
    go_repository(
        name = "com_github_p4lang_p4runtime",
        build_directives = [
            "gazelle:resolve go github.com/p4lang/p4runtime/go/p4/v1 @com_github_p4lang_p4runtime//go/p4/v1:p4",
            "gazelle:proto disable",  # keep
        ],
        build_file_generation = "clean",
        importpath = "github.com/p4lang/p4runtime",
        sum = "h1:GSccPwIFfeRjyrUSDe19DmqsHia7tGsU8vFuH2JxPTU=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_openconfig_ondatra",
        build_directives = [
            "gazelle:resolve go github.com/p4lang/p4runtime/go/p4/v1 @com_github_p4lang_p4runtime//go/p4/v1:p4",
            "gazelle:resolve go github.com/openconfig/attestz/proto/tpm_enrollz @com_github_openconfig_attestz//proto:tpm_enrollz_go",
            "gazelle:resolve go github.com/openconfig/attestz/proto/tpm_attestz @com_github_openconfig_attestz//proto:tpm_attestz_go",
            # "gazelle:resolve proto go github.com/openconfig/attestz/proto/tpm_enrollz.proto @com_github_openconfig_attestz//proto:tpm_enrollz_go",
            # "gazelle:resolve proto go github.com/openconfig/attestz/proto/tpm_attestz.proto @com_github_openconfig_attestz//proto:tpm_attestz_go",
        ],
        importpath = "github.com/openconfig/ondatra",
        sum = "h1:R1HzjC8Izcy9wTacrEXaQZM2SRBpqegtSCSriWJd2y8=",
        version = "v0.14.8",
    )
    go_repository(
        name = "com_github_openconfig_testt",
        importpath = "github.com/openconfig/testt",
        sum = "h1:KI67JEIYGMXx9aZjGOwfO2aUsCGih9199v+zhJNHB9Y=",
        version = "v0.0.0-20251119232631-fbbd49c39452",
    )
    go_repository(
        name = "com_github_openconfig_ygnmi",
        importpath = "github.com/openconfig/ygnmi",
        sum = "h1:Cio+fv7+c8G/w3KpO9CK7BLLEMdBy9bO1q3RRsd5xqs=",
        version = "v0.15.0",
    )
    go_repository(
        name = "com_github_openconfig_ygot",
        build_file_proto_mode = "disable",
        importpath = "github.com/openconfig/ygot",
        sum = "h1:RCS7a5h40cyGAZQjQZGA7SunPuKl188l94DYTnHFX9A=",
        version = "v0.35.0",
    )
    go_repository(
        name = "com_github_opencontainers_go_digest",
        importpath = "github.com/opencontainers/go-digest",
        sum = "h1:apOUWs51W5PlhuyGyz9FCeeBIOUDA/6nW8Oi/yOhh5U=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_opencontainers_image_spec",
        importpath = "github.com/opencontainers/image-spec",
        sum = "h1:y0fUlFfIZhPF1W537XOLg0/fcx6zcHCJwooC2xJA040=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_osrg_gobgp_v3",
        build_file_proto_mode = "disable",
        importpath = "github.com/osrg/gobgp/v3",
        replace = "github.com/DanG100/gobgp/v3",
        sum = "h1:IkoLokGQxuogFQb9SeFzwDoSWAfXKARHBVa4bo6Mwbo=",
        version = "v3.0.0-20250514232417-c9f561b7bb0e",
    )
    go_repository(
        name = "com_github_patrickmn_go_cache",
        importpath = "github.com/patrickmn/go-cache",
        sum = "h1:HRMgzkcYKYpi3C8ajMPV8OFXaaRUnok+kx1WdO15EQc=",
        version = "v2.1.0+incompatible",
    )
    go_repository(
        name = "com_github_pbnjay_memory",
        importpath = "github.com/pbnjay/memory",
        sum = "h1:onHthvaw9LFnH4t2DcNVpwGmV9E1BkGknEliJkfwQj0=",
        version = "v0.0.0-20210728143218-7b4eea64cf58",
    )
    go_repository(
        name = "com_github_pborman_getopt",
        importpath = "github.com/pborman/getopt",
        sum = "h1:eJ3aFZroQqq0bWmraivjQNt6Dmm5M0h2JcDW38/Azb0=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_pborman_uuid",
        importpath = "github.com/pborman/uuid",
        sum = "h1:+ZZIw58t/ozdjRaXh/3awHfmWRbzYxJoAdNJxe/3pvw=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_pelletier_go_toml_v2",
        importpath = "github.com/pelletier/go-toml/v2",
        sum = "h1:GTRvJQutkOSftxIFD5xw9aepkYNuPWmVJpffdDPYVpY=",
        version = "v2.4.3",
    )
    go_repository(
        name = "com_github_peterbourgon_diskv",
        importpath = "github.com/peterbourgon/diskv",
        sum = "h1:UBdAOUP5p4RWqPBg048CAvpKN+vxiaj6gdUUzhl4XmI=",
        version = "v2.0.1+incompatible",
    )
    go_repository(
        name = "com_github_pierrec_lz4_v4",
        importpath = "github.com/pierrec/lz4/v4",
        sum = "h1:+PhzhWDrjRj89TH2sw43nE3+4+W8lSxIuQadEHZyjUk=",
        version = "v4.1.27",
    )
    go_repository(
        name = "com_github_pjbgf_sha1cd",
        importpath = "github.com/pjbgf/sha1cd",
        sum = "h1:4D5XXmUUBUl/xQ6IjCkEAbqXskkq/4O7LmGn0AqMDs4=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_pkg_diff",
        importpath = "github.com/pkg/diff",
        sum = "h1:aoZm08cpOy4WuID//EZDgcC4zIxODThtZNPirFr42+A=",
        version = "v0.0.0-20210226163009-20ebb0f2a09e",
    )
    go_repository(
        name = "com_github_pkg_errors",
        importpath = "github.com/pkg/errors",
        sum = "h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
        version = "v0.9.1",
    )
    go_repository(
        name = "com_github_planetscale_vtprotobuf",
        importpath = "github.com/planetscale/vtprotobuf",
        sum = "h1:GFCKgmp0tecUJ0sJuv4pzYCqS9+RGSn52M3FUwPs+uo=",
        version = "v0.6.1-0.20240319094008-0393e58bdf10",
    )
    go_repository(
        name = "com_github_pmezard_go_difflib",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=",
        version = "v1.0.1-0.20181226105442-5d4384ee4fb2",
    )
    go_repository(
        name = "com_github_prometheus_client_golang",
        importpath = "github.com/prometheus/client_golang",
        sum = "h1:5XStIklKuAtJSNpdD3s8XJj/Yv78IQmE1kbNk87JrAI=",
        version = "v1.24.0",
    )
    go_repository(
        name = "com_github_prometheus_client_model",
        build_file_proto_mode = "disable",
        importpath = "github.com/prometheus/client_model",
        sum = "h1:oBsgwpGs7iVziMvrGhE53c/GrLUsZdHnqNwqPLxwZyk=",
        version = "v0.6.2",
    )
    go_repository(
        name = "com_github_prometheus_common",
        importpath = "github.com/prometheus/common",
        sum = "h1:bcpru3tWPVnxGnETLgOV5jbp/JRXgYEyv65CuBLAMMI=",
        version = "v0.70.0",
    )
    go_repository(
        name = "com_github_prometheus_procfs",
        importpath = "github.com/prometheus/procfs",
        sum = "h1:GljZCt+zSTS+NZq88cyQ1LjZ+RCHp3uVuabBWA5+OJI=",
        version = "v0.21.1",
    )
    go_repository(
        name = "com_github_protocolbuffers_txtpbfmt",
        importpath = "github.com/protocolbuffers/txtpbfmt",
        sum = "h1:ej+64jiny5VETZTqcc1GFVAPEtaSk6U1D0kKC2MS5Yc=",
        version = "v0.0.0-20240823084532-8e6b51fa9bef",
    )
    go_repository(
        name = "com_github_protonmail_go_crypto",
        importpath = "github.com/ProtonMail/go-crypto",
        sum = "h1:nRBOetoydLeUb4nHajyO2bKqMLfWQ/ZPwkXqXxPxCFk=",
        version = "v1.1.3",
    )
    go_repository(
        name = "com_github_puerkitobio_purell",
        importpath = "github.com/PuerkitoBio/purell",
        sum = "h1:WEQqlqaGbrPkxLJWfBwQmfEAE1Z7ONdDLqrN38tNFfI=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_puerkitobio_urlesc",
        importpath = "github.com/PuerkitoBio/urlesc",
        sum = "h1:d+Bc7a5rLufV/sSk/8dngufqelfh6jnri85riMAaF/M=",
        version = "v0.0.0-20170810143723-de5bf2ad4578",
    )
    go_repository(
        name = "com_github_redhat_nfvpe_koko",
        importpath = "github.com/redhat-nfvpe/koko",
        sum = "h1:a6VTUhyo84BNgbPfGcOImbHtLJfNQsf4hkXsL17mdlU=",
        version = "v0.0.0-20241203052104-ec78eb8730bf",
    )
    go_repository(
        name = "com_github_remyoudompheng_bigfft",
        importpath = "github.com/remyoudompheng/bigfft",
        sum = "h1:W09IVJc94icq4NjY3clb7Lk8O1qJ8BdBEF8z0ibU0rE=",
        version = "v0.0.0-20230129092748-24d4a6f8daec",
    )
    go_repository(
        name = "com_github_rifflock_lfshook",
        importpath = "github.com/rifflock/lfshook",
        sum = "h1:mZHayPoR0lNmnHyvtYjDeq0zlVHn9K/ZXoy17ylucdo=",
        version = "v0.0.0-20180920164130-b9218ef580f5",
    )
    go_repository(
        name = "com_github_rogpeppe_fastuuid",
        importpath = "github.com/rogpeppe/fastuuid",
        sum = "h1:Ppwyp6VYCF1nvBTXL3trRso7mXMlRrw9ooo375wvi2s=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_rogpeppe_go_internal",
        importpath = "github.com/rogpeppe/go-internal",
        sum = "h1:UQB4HGPB6osV0SQTLymcB4TgvyWu6ZyliaW0tI/otEQ=",
        version = "v1.14.1",
    )
    go_repository(
        name = "com_github_russross_blackfriday_v2",
        importpath = "github.com/russross/blackfriday/v2",
        sum = "h1:JIOH55/0cWyOuilr9/qlrm0BSXldqnqwMsf35Ld67mk=",
        version = "v2.1.0",
    )
    go_repository(
        name = "com_github_safchain_ethtool",
        importpath = "github.com/safchain/ethtool",
        sum = "h1:rlJzfDetsVvT61uz8x1YIcFn12akMfuPulHtZjtb7Is=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_github_sagikazarmark_locafero",
        importpath = "github.com/sagikazarmark/locafero",
        sum = "h1:/NQhBAkUb4+fH1jivKHWusDYFjMOOKU88eegjfxfHb4=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_github_scrapli_scrapligo",
        importpath = "github.com/scrapli/scrapligo",
        sum = "h1:D0lg2xusZKbKnzy9l0yT7ly/D2rgqab12yZvEXNim0U=",
        version = "v1.4.1",
    )
    go_repository(
        name = "com_github_scrapli_scrapligocfg",
        importpath = "github.com/scrapli/scrapligocfg",
        sum = "h1:540SuGqqM6rKN87SLCfR54IageQ6s3a/ZOycGRgbbak=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_sergi_go_diff",
        importpath = "github.com/sergi/go-diff",
        sum = "h1:n661drycOFuPLCN3Uc8sB6B/s6Z4t2xvBgU1htSHuq8=",
        version = "v1.3.2-0.20230802210424-5b0b94c5c0d3",
    )
    go_repository(
        name = "com_github_sirikothe_gotextfsm",
        importpath = "github.com/sirikothe/gotextfsm",
        sum = "h1:FHUL2HofYJuslFOQdy/JjjP36zxqIpd/dcoiwLMIs7k=",
        version = "v1.0.1-0.20200816110946-6aa2cfd355e4",
    )
    go_repository(
        name = "com_github_sirupsen_logrus",
        importpath = "github.com/sirupsen/logrus",
        sum = "h1:G2SED73/qrAu6YwbdxOD6peLkCBI3z7L+ykJFTXJBBo=",
        version = "v1.10.2",
    )
    go_repository(
        name = "com_github_skeema_knownhosts",
        importpath = "github.com/skeema/knownhosts",
        sum = "h1:AM+y0rI04VksttfwjkSTNQorvGqmwATnvnAHpSgc0LY=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_soheilhy_cmux",
        importpath = "github.com/soheilhy/cmux",
        sum = "h1:jjzc5WVemNEDTLwv9tlmemhC73tI08BNOIGwBOo10Js=",
        version = "v0.1.5",
    )
    go_repository(
        name = "com_github_sourcegraph_conc",
        importpath = "github.com/sourcegraph/conc",
        sum = "h1:+jumHNA0Wrelhe64i8F6HNlS8pkoyMv5sreGx2Ry5Rw=",
        version = "v0.3.1-0.20240121214520-5f936abd7ae8",
    )
    go_repository(
        name = "com_github_spf13_afero",
        importpath = "github.com/spf13/afero",
        sum = "h1:b/YBCLWAJdFWJTN9cLhiXXcD7mzKn9Dm86dNnfyQw1I=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_github_spf13_cast",
        importpath = "github.com/spf13/cast",
        sum = "h1:h2x0u2shc1QuLHfxi+cTJvs30+ZAHOGRic8uyGTDWxY=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_github_spf13_cobra",
        importpath = "github.com/spf13/cobra",
        sum = "h1:DMTTonx5m65Ic0GOoRY2c16WCbHxOOw6xxezuLaBpcU=",
        version = "v1.10.2",
    )
    go_repository(
        name = "com_github_spf13_jwalterweatherman",
        importpath = "github.com/spf13/jwalterweatherman",
        sum = "h1:ue6voC5bR5F8YxI5S67j9i582FU4Qvo2bmqnqMYADFk=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_spf13_pflag",
        importpath = "github.com/spf13/pflag",
        sum = "h1:4EBh2KAYBwaONj6b2Ye1GiHfwjqyROoF4RwYO+vPwFk=",
        version = "v1.0.10",
    )
    go_repository(
        name = "com_github_spf13_viper",
        importpath = "github.com/spf13/viper",
        sum = "h1:x5S+0EU27Lbphp4UKm1C+1oQO+rKx36vfCoaVebLFSU=",
        version = "v1.21.0",
    )
    go_repository(
        name = "com_github_srl_labs_srl_controller",
        importpath = "github.com/srl-labs/srl-controller",
        sum = "h1:QZST5zXfXGCdSQUEb22iMvAwfQhfOcqzA3+Ae+ifd70=",
        version = "v0.7.1",
    )
    go_repository(
        name = "com_github_srl_labs_srlinux_scrapli",
        importpath = "github.com/srl-labs/srlinux-scrapli",
        sum = "h1:YQjckD+a7f6u2M+k4SmJUrDa7BFvoOTb2mMbPe6hLZM=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_github_stoewer_go_strcase",
        importpath = "github.com/stoewer/go-strcase",
        sum = "h1:iS0MdW+kVTxgMoE1LAZyMiYJFKlOzLooE4MxjirtkAs=",
        version = "v1.3.1",
    )
    go_repository(
        name = "com_github_stretchr_objx",
        importpath = "github.com/stretchr/objx",
        sum = "h1:jmXUvGomnU1o3W/V5h2VEradbpJDwGrzugQQvL0POH4=",
        version = "v0.5.3",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        sum = "h1:EuwCh5fleGS7H32xRwO3wRGT7DxrDhLAT6FF8MpWDWE=",
        version = "v1.12.1",
    )
    go_repository(
        name = "com_github_subosito_gotenv",
        importpath = "github.com/subosito/gotenv",
        sum = "h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_tidwall_gjson",
        importpath = "github.com/tidwall/gjson",
        sum = "h1:FIDeeyB800efLX89e5a8Y0BNH+LOngJyGrIWxG2FKQY=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_github_tidwall_match",
        importpath = "github.com/tidwall/match",
        sum = "h1:+Ho715JplO36QYgwN9PGYNhgZvoUSc9X2c80KVTi+GA=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_tidwall_pretty",
        importpath = "github.com/tidwall/pretty",
        sum = "h1:qjsOFOWWQl+N3RsoF5/ssm1pHmJJwhjlSbZ51I6wMl4=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_tidwall_sjson",
        importpath = "github.com/tidwall/sjson",
        sum = "h1:kLy8mja+1c9jlljvWTlSazM7cKDRfJuR/bOJhcY5NcY=",
        version = "v1.2.5",
    )
    go_repository(
        name = "com_github_tmc_grpc_websocket_proxy",
        importpath = "github.com/tmc/grpc-websocket-proxy",
        sum = "h1:6fotK7otjonDflCTK0BCfls4SPy3NcCVb5dqqmbRknE=",
        version = "v0.0.0-20220101234140-673ab2c3ae75",
    )
    go_repository(
        name = "com_github_u_root_uio",
        importpath = "github.com/u-root/uio",
        sum = "h1:pyC9PaHYZFgEKFdlp3G8RaCKgVpHZnecvArXvPXcFkM=",
        version = "v0.0.0-20240224005618-d2acac8f3701",
    )
    go_repository(
        name = "com_github_vishvananda_netlink",
        importpath = "github.com/vishvananda/netlink",
        sum = "h1:3AEMt62VKqz90r0tmNhog0r/PpWKmrEShJU0wJW6bV0=",
        version = "v1.3.1",
    )
    go_repository(
        name = "com_github_vishvananda_netns",
        importpath = "github.com/vishvananda/netns",
        sum = "h1:DfiHV+j8bA32MFM7bfEunvT8IAqQ/NzSJHtcmW5zdEY=",
        version = "v0.0.5",
    )
    go_repository(
        name = "com_github_wi2l_jsondiff",
        importpath = "github.com/wI2L/jsondiff",
        sum = "h1:Fg9+yj+1/x3UtPBJhR91TKEzRkrEEWcAcLbg9dzEaNM=",
        version = "v0.7.1",
    )
    go_repository(
        name = "com_github_xanzy_ssh_agent",
        importpath = "github.com/xanzy/ssh-agent",
        sum = "h1:+/15pJfg/RsTxqYcX6fHqOXZwwMP+2VyYWJeWM2qQFM=",
        version = "v0.3.3",
    )
    go_repository(
        name = "com_github_xiang90_probing",
        importpath = "github.com/xiang90/probing",
        sum = "h1:S2dVYn90KE98chqDkyE9Z4N61UnQd+KOfgp5Iu53llk=",
        version = "v0.0.0-20221125231312-a49e3df8f510",
    )
    go_repository(
        name = "com_github_yuin_goldmark",
        importpath = "github.com/yuin/goldmark",
        sum = "h1:fVcFKWvrslecOb/tg+Cc05dkeYx540o0FuFt3nUVDoE=",
        version = "v1.4.13",
    )
    go_repository(
        name = "com_github_zeebo_xxh3",
        importpath = "github.com/zeebo/xxh3",
        sum = "h1:xZmwmqxHZA8AI603jOQ0tMqmBr9lPeFwGg6d+xy9DC0=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_google_cloud_go",
        importpath = "cloud.google.com/go",
        sum = "h1:2NAUJwPR47q+E35uaJeYoNhuNEM9kM8SjgRgdeOJUSE=",
        version = "v0.123.0",
    )
    go_repository(
        name = "com_google_cloud_go_accessapproval",
        importpath = "cloud.google.com/go/accessapproval",
        sum = "h1:kx3RQSS0VglTngQTPfywgVj+xFgre/vJouGETBlPU1E=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_accesscontextmanager",
        importpath = "cloud.google.com/go/accesscontextmanager",
        sum = "h1:0baVgug9IFV8y5cIoD4iN/m/QtJzk2k+j3MW+zJWCVw=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_aiplatform",
        importpath = "cloud.google.com/go/aiplatform",
        sum = "h1:5PxeWpQfkAyN8mtgtVp4H5nP+ntwPLFMjNREVdiolR0=",
        version = "v1.126.0",
    )
    go_repository(
        name = "com_google_cloud_go_analytics",
        importpath = "cloud.google.com/go/analytics",
        sum = "h1:GuwmzJHIaQRvtko6g4wW1ngQ7rpDvDLZ8M3iP3kiflU=",
        version = "v0.35.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigateway",
        importpath = "cloud.google.com/go/apigateway",
        sum = "h1:IENNmFQlMsl58RlvYHf6wxBNQFeP+lAnVNcyxtmByTE=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeconnect",
        importpath = "cloud.google.com/go/apigeeconnect",
        sum = "h1:YvZhi0QkHobBRMrN6ri9PUPnnhgEd4VpW0+yo1WI7r0=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeregistry",
        importpath = "cloud.google.com/go/apigeeregistry",
        sum = "h1:2zIogaYq7dF4E1dAGC4MUnasRuSmohsMl7wVuJ963/4=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_google_cloud_go_appengine",
        importpath = "cloud.google.com/go/appengine",
        sum = "h1:0MiBM2KGk1WAhh489WoELQ7IeW2qw6HTE6nZIlrOw5U=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_area120",
        importpath = "cloud.google.com/go/area120",
        sum = "h1:9v6HeQsBpfRqoZtNeLI4V8SQ7Jcmt6s2gmgbjl0fVFc=",
        version = "v0.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_artifactregistry",
        importpath = "cloud.google.com/go/artifactregistry",
        sum = "h1:iq5kkdY2FJY8RkkNghUelb3EOcqyy7dVhfDGVk/Xw9g=",
        version = "v1.26.0",
    )
    go_repository(
        name = "com_google_cloud_go_asset",
        importpath = "cloud.google.com/go/asset",
        sum = "h1:M6YE1exBuZQhTi8wfIKrwjKK6b2ySfHKgrqXjPZA4EM=",
        version = "v1.28.0",
    )
    go_repository(
        name = "com_google_cloud_go_assuredworkloads",
        importpath = "cloud.google.com/go/assuredworkloads",
        sum = "h1:XyNQ9SBX0eUkZYCZhC8dOY6ZFeSF6/OIjHn5dQs78Jo=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_auth",
        importpath = "cloud.google.com/go/auth",
        sum = "h1:UYMbF8otPZnLAkNJ5/LYQYOq0ARcJS1P4JqTeMKbCYU=",
        version = "v0.24.0",
    )
    go_repository(
        name = "com_google_cloud_go_auth_oauth2adapt",
        importpath = "cloud.google.com/go/auth/oauth2adapt",
        sum = "h1:FY8oSZpCYoUNv6QxVODuMjQz4IlSOVeiQtZ08vLPz88=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_automl",
        importpath = "cloud.google.com/go/automl",
        sum = "h1:brKXr4i9AMVfgbCKOYSu4rnIoz1ILIgcIAy3nerzEGc=",
        version = "v1.21.0",
    )
    go_repository(
        name = "com_google_cloud_go_baremetalsolution",
        importpath = "cloud.google.com/go/baremetalsolution",
        sum = "h1:6MaZilXzGZJ3mGqbSCdDezs9yrze6kcmBABsRUZ/aYI=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_batch",
        importpath = "cloud.google.com/go/batch",
        sum = "h1:nTL5x9HA1yVdPIdeqcbLjNnV9iYDYvMNmP6AnZ+qbHU=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_beyondcorp",
        importpath = "cloud.google.com/go/beyondcorp",
        sum = "h1:rsyle6zjxce2s57BpULw5u170xS+KSjFENSY7C0MZbs=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_bigquery",
        importpath = "cloud.google.com/go/bigquery",
        sum = "h1:+tW6oRuP/4dOgdl/p32on0+3OktVQtU+veOMxhjoOhE=",
        version = "v1.79.0",
    )
    go_repository(
        name = "com_google_cloud_go_bigtable",
        importpath = "cloud.google.com/go/bigtable",
        sum = "h1:NGLgDSr/i79BTGCjxH/maPKxyvl5q8/SsBsyLK52kdI=",
        version = "v1.47.0",
    )
    go_repository(
        name = "com_google_cloud_go_billing",
        importpath = "cloud.google.com/go/billing",
        sum = "h1:6RRjbRd6iZKZFb7/MgRvmXKq/Ism02ckkZLJazj4CQ0=",
        version = "v1.26.0",
    )
    go_repository(
        name = "com_google_cloud_go_binaryauthorization",
        importpath = "cloud.google.com/go/binaryauthorization",
        sum = "h1:f1iNUaHWP9XPdh6grftTlIRuJSBytud4j51l1Jce94E=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_certificatemanager",
        importpath = "cloud.google.com/go/certificatemanager",
        sum = "h1:PwljEZlI3lZgoFUG0pnIBNOwrQUaAP+gS0ZahTwUZBA=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_channel",
        importpath = "cloud.google.com/go/channel",
        sum = "h1:fqDZDzVw61c4iIUe7rQ+cBYlBdSZ3yqnzw5WrZTPjqs=",
        version = "v1.27.0",
    )
    go_repository(
        name = "com_google_cloud_go_cloudbuild",
        build_directives = [
            "gazelle:resolve go google.golang.org/genproto/googleapis/longrunning @org_golang_google_genproto//googleapis/longrunning",  # keep
            "gazelle:resolve go google.golang.org/genproto/googleapis/logging/v2 @org_golang_google_genproto//googleapis/logging/v2:logging",  # keep
        ],
        build_file_proto_mode = "disable",
        importpath = "cloud.google.com/go/cloudbuild",
        sum = "h1:5EuMx7RMMJIxGd8cdwYylLXo9/vrprx2J7UYneUBsBY=",
        version = "v1.34.0",
    )
    go_repository(
        name = "com_google_cloud_go_clouddms",
        importpath = "cloud.google.com/go/clouddms",
        sum = "h1:oaktwVyUeKTuzjpde3xbpmc9bAQPQq6WJWOBr7K6YEA=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_cloudtasks",
        importpath = "cloud.google.com/go/cloudtasks",
        sum = "h1:KzT7hfix/9/xAf20tNPIxwX59XGpRF0Lun2t8LHOj9E=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_compute",
        importpath = "cloud.google.com/go/compute",
        sum = "h1:KsBourH0wajM4RhzwPwRMKbxHVdvzGsk7StvACoWXD8=",
        version = "v1.63.0",
    )
    go_repository(
        name = "com_google_cloud_go_compute_metadata",
        importpath = "cloud.google.com/go/compute/metadata",
        sum = "h1:pyKMUQSwchgkIBBJGdILqQbs/BNJXqwSA7Ej6LAvvtY=",
        version = "v0.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_contactcenterinsights",
        importpath = "cloud.google.com/go/contactcenterinsights",
        sum = "h1:EtrRf8T7ajR+dR8nseDxwKy3o1VaWP1GFAu7lgMmT9A=",
        version = "v1.23.0",
    )
    go_repository(
        name = "com_google_cloud_go_container",
        importpath = "cloud.google.com/go/container",
        sum = "h1:KdeWHqwlOvABdhkhSTSb4QFySMOuXnc3UUoQMmq9lic=",
        version = "v1.51.0",
    )
    go_repository(
        name = "com_google_cloud_go_containeranalysis",
        importpath = "cloud.google.com/go/containeranalysis",
        sum = "h1:89ZhFvJHWzX9/jUJy/IEVbtSa2hv4+NtsvwWMEsUYnY=",
        version = "v0.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_datacatalog",
        importpath = "cloud.google.com/go/datacatalog",
        sum = "h1:8V80PpoAGdOOr2QhBrp4wZ66MDCbATdAB/fmVmo5rlU=",
        version = "v1.33.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataflow",
        importpath = "cloud.google.com/go/dataflow",
        sum = "h1:BchGCAl9QIZ/pyTGokv5V3daxBTrXsOmrU4ARXoTNd4=",
        version = "v0.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataform",
        importpath = "cloud.google.com/go/dataform",
        sum = "h1:eAIsWhr6AbUIX+wOdfctO6hAr9UoQnTR2qQW739XH1I=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_google_cloud_go_datafusion",
        importpath = "cloud.google.com/go/datafusion",
        sum = "h1:56rjOW8xFBnTOLgTEWdU09KhVC2/X+y/ukzux7zYqBE=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_datalabeling",
        importpath = "cloud.google.com/go/datalabeling",
        sum = "h1:9TL2kwlO/ODD1fne3uvgTLa53SttwrQi92Bn7KWnKhg=",
        version = "v0.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataplex",
        importpath = "cloud.google.com/go/dataplex",
        sum = "h1:r4O6RHEkwHr41ItsTUym0ZgH1iDf48R+W8fI7KoMNl8=",
        version = "v1.36.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataproc_v2",
        importpath = "cloud.google.com/go/dataproc/v2",
        sum = "h1:IF1wxkXvcgoTP92NVqX6QNDf2SdKNZFPa2IWeSKP310=",
        version = "v2.25.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataqna",
        importpath = "cloud.google.com/go/dataqna",
        sum = "h1:VR1y/NuN+RvPkmfmlowbmk60BLINR4MgZAusVFIGQjU=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_datastore",
        importpath = "cloud.google.com/go/datastore",
        sum = "h1:zUjMnCLCcRZVDSdQIXsbnNCl1SVRNw5Jm0J77gPaPKs=",
        version = "v1.25.0",
    )
    go_repository(
        name = "com_google_cloud_go_datastream",
        importpath = "cloud.google.com/go/datastream",
        sum = "h1:hrWGASBj85xs6wpnx9ifdSBYiRLhKDgQrA1HjeNYBhc=",
        version = "v1.21.0",
    )
    go_repository(
        name = "com_google_cloud_go_deploy",
        importpath = "cloud.google.com/go/deploy",
        sum = "h1:2PdZ8kbIztLVRMrIN/wHbalZeZZikQlxBptBFK+eEtM=",
        version = "v1.33.0",
    )
    go_repository(
        name = "com_google_cloud_go_dialogflow",
        importpath = "cloud.google.com/go/dialogflow",
        sum = "h1:BymJ6nPotDcnphK2/D0VHMAfakjJlIV5GqxLHQg/Tos=",
        version = "v1.84.0",
    )
    go_repository(
        name = "com_google_cloud_go_dlp",
        importpath = "cloud.google.com/go/dlp",
        sum = "h1:SI7KSLOtAfhEAw3af8NxpKHGHbJ9BkoM7805TMaU6m4=",
        version = "v1.34.0",
    )
    go_repository(
        name = "com_google_cloud_go_documentai",
        importpath = "cloud.google.com/go/documentai",
        sum = "h1:WLZOaWvK4NH7xMLxANsDgIdXugeqmVyYRCbVN1m87DQ=",
        version = "v1.49.0",
    )
    go_repository(
        name = "com_google_cloud_go_domains",
        importpath = "cloud.google.com/go/domains",
        sum = "h1:h+lcYPxyEulj5SYuH3+OwOsLZb4DTW11QX6nlxxXQ1I=",
        version = "v0.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_edgecontainer",
        importpath = "cloud.google.com/go/edgecontainer",
        sum = "h1:4WOjcIZRCB4ynxH9Zs1tteIcm2LD58VUCLLmIwLYUAc=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_errorreporting",
        importpath = "cloud.google.com/go/errorreporting",
        sum = "h1:LlE2SVIbz0k+OSeNTksk34inr3Fy62JMhHUvNaS8f7c=",
        version = "v0.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_essentialcontacts",
        importpath = "cloud.google.com/go/essentialcontacts",
        sum = "h1:+AtGn+hYLdr62sMqP0ZGMRzp9A4t+MSWe8eZoD/ho+M=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_eventarc",
        importpath = "cloud.google.com/go/eventarc",
        sum = "h1:+XAJNEmxJIPZLMtzLYdUZN61Kmqx62hWLry4AsWaTwo=",
        version = "v1.25.0",
    )
    go_repository(
        name = "com_google_cloud_go_filestore",
        importpath = "cloud.google.com/go/filestore",
        sum = "h1:2GXy5jvq2oUX6HnCpx3GjzSqlcMz8K9gkHaF9HPcYNE=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_firestore",
        importpath = "cloud.google.com/go/firestore",
        sum = "h1:x0Z3hrgjYgo2wI9whuBRQcNc2hYwzZDQy/7pkUXbXcs=",
        version = "v1.24.0",
    )
    go_repository(
        name = "com_google_cloud_go_functions",
        importpath = "cloud.google.com/go/functions",
        sum = "h1:ndUtLkam3XF9b0t2zVACH9D/EgFBISbVuQFqRz/X58k=",
        version = "v1.25.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkebackup",
        importpath = "cloud.google.com/go/gkebackup",
        sum = "h1:li3BtGRis1QYrkLo8+Iq2wf5WbP9v3sz9VoUw8WqgaA=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkeconnect",
        importpath = "cloud.google.com/go/gkeconnect",
        sum = "h1:IDuAD/w0Xph8k/31vdt6pBkDdRx55lVXz7gW8ABOfCs=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkehub",
        importpath = "cloud.google.com/go/gkehub",
        sum = "h1:gHKPoWQuWpd9dXLoC58dxwBvSYTNc5/fdYYcs3JUn1s=",
        version = "v0.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkemulticloud",
        importpath = "cloud.google.com/go/gkemulticloud",
        sum = "h1:3LPj8ro7bPdZ0BlQJwby25G32rRezGLRdhNM6U3AxGs=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_grafeas",
        importpath = "cloud.google.com/go/grafeas",
        sum = "h1:eTq1l4Czo48O39e2ZKSvaRwn0klCeUPf4CvbA+2R4IY=",
        version = "v0.3.17",
    )
    go_repository(
        name = "com_google_cloud_go_gsuiteaddons",
        importpath = "cloud.google.com/go/gsuiteaddons",
        sum = "h1:kz9DyBk84wZCda2Rdha1MOIf7/9x6R+N+LuH9B3zYFs=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_iam",
        importpath = "cloud.google.com/go/iam",
        sum = "h1:Aki3bX9aHUDKPHfnRJfDcTdVedvy6quGBQcTqx3DRXk=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_iap",
        importpath = "cloud.google.com/go/iap",
        sum = "h1:BDAVJy+juq7cMRumIx9toc4pt1K7zXoZdAI3lDD6D3g=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_ids",
        importpath = "cloud.google.com/go/ids",
        sum = "h1:394LkKIavhv/+oQVMtbXe9WQqCeY9wSFcHktqC9ZkLE=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_iot",
        importpath = "cloud.google.com/go/iot",
        sum = "h1:pyt1EuMFpbV/2BlmYlXr+HBnUNSPKlP9dVK/dpFTu1U=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_kms",
        importpath = "cloud.google.com/go/kms",
        sum = "h1:s+rEluaaZKhLVjrIWG7uNBsnWbiitElzNzFGyp6+nIg=",
        version = "v1.32.0",
    )
    go_repository(
        name = "com_google_cloud_go_language",
        importpath = "cloud.google.com/go/language",
        sum = "h1:q58bL7rmxvw6Q6VHt+wjFsAE1Tj/JkuAEIR4+84rx9U=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_lifesciences",
        importpath = "cloud.google.com/go/lifesciences",
        sum = "h1:T+7N3oVHK8a7BhX6HJLrBW8xEbMV9iY8Mb+aOl7r6+A=",
        version = "v0.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_logging",
        importpath = "cloud.google.com/go/logging",
        sum = "h1:l7yzCB8/JSY5PHT9RsaqdT1jaBGbM6KlSJDX210HlXM=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_longrunning",
        build_file_proto_mode = "disable",
        importpath = "cloud.google.com/go/longrunning",
        sum = "h1:WjYH3YHBGCxGJP9M4dWGHBfXr/cFIjMkNgWcJj7/iMM=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_google_cloud_go_managedidentities",
        importpath = "cloud.google.com/go/managedidentities",
        sum = "h1:ZzWkg3LSrIi9OdzCK3bxutexvw1dDDg9+tOMitbkOXw=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_maps",
        importpath = "cloud.google.com/go/maps",
        sum = "h1:YoTRWohpNKeKxdQynorwHH69NZuUM06divZtqGPYjUA=",
        version = "v1.37.0",
    )
    go_repository(
        name = "com_google_cloud_go_mediatranslation",
        importpath = "cloud.google.com/go/mediatranslation",
        sum = "h1:qqMRAK0mhc9M+lM6pb791oh3DYUbwLcIy58jOyXtadE=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_memcache",
        importpath = "cloud.google.com/go/memcache",
        sum = "h1:l/co1jsmVLkNyl8uMqj1u1HIwd3AhCvGdHbEAIPb6jc=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_metastore",
        importpath = "cloud.google.com/go/metastore",
        sum = "h1:kz0XBkP531MB1IR9jTY6/DB33UFP5Df1vFpPNgW1iPk=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_monitoring",
        importpath = "cloud.google.com/go/monitoring",
        sum = "h1:LCuAoYiq8snyL/eDAG/VC73uWFXeax8+YEfmoNGJUrY=",
        version = "v1.31.0",
    )
    go_repository(
        name = "com_google_cloud_go_networkconnectivity",
        importpath = "cloud.google.com/go/networkconnectivity",
        sum = "h1:ieYJjbUn2L9ZSzEB9+qKgtoD3/l1/qAnId/Z6KE4ZBo=",
        version = "v1.27.0",
    )
    go_repository(
        name = "com_google_cloud_go_networkmanagement",
        importpath = "cloud.google.com/go/networkmanagement",
        sum = "h1:LqYm91vSQk3NM6fHjthseRh4xrDvhXIWQBsQY/Ksus8=",
        version = "v1.30.0",
    )
    go_repository(
        name = "com_google_cloud_go_networksecurity",
        importpath = "cloud.google.com/go/networksecurity",
        sum = "h1:PzJnbVd0sS2+xDR9hnk31NAyqFt9qB8ohLn7XIsKIBA=",
        version = "v0.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_notebooks",
        importpath = "cloud.google.com/go/notebooks",
        sum = "h1:BgxdQcoRSSjiyszWNOTu3fmXG/6Ks4W2sHgjx0BAke0=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_optimization",
        importpath = "cloud.google.com/go/optimization",
        sum = "h1:Cd4DgujNhEMcjfaVLGqbjfDwVDnmS1Bc93cx9TE8ZYw=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_orchestration",
        importpath = "cloud.google.com/go/orchestration",
        sum = "h1:dqf2HzqUi4CaaDiNfv6jVoZjz9lF0HWjyIEgccxPHMk=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_orgpolicy",
        importpath = "cloud.google.com/go/orgpolicy",
        sum = "h1:kpVcE/OsC5aAzHCsAiuQSg3+s6ILzgPTuPZyS7n7ejA=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_osconfig",
        importpath = "cloud.google.com/go/osconfig",
        sum = "h1:r5lzneR9GNixJ96lZ6oIfJK9Cd5sKpnzZJe1LcGZf4w=",
        version = "v1.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_oslogin",
        importpath = "cloud.google.com/go/oslogin",
        sum = "h1:OkccORMcY2XEHwN5+AP0XH/SCUNBl39GGjaQuSCVcIw=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_phishingprotection",
        importpath = "cloud.google.com/go/phishingprotection",
        sum = "h1:6WJF1z3Ie8fZAiCRpeSGB81h6s86XWS2SBC9iQUZxX8=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_policytroubleshooter",
        importpath = "cloud.google.com/go/policytroubleshooter",
        sum = "h1:V3CcqzY6FhY2c2idRGZv0Ef4mFinLnXnfDFtAhxyTfM=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_privatecatalog",
        importpath = "cloud.google.com/go/privatecatalog",
        sum = "h1:SBhcVNVhgEGT1edcfSFY/KSiy2xSQTSQUpY12f10rmw=",
        version = "v0.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub",
        importpath = "cloud.google.com/go/pubsub",
        sum = "h1:XOaCejsqX7EEtUdQz+WPag66wWsUUGliyCOfGPKfo90=",
        version = "v1.51.0",
    )
    go_repository(
        name = "com_google_cloud_go_pubsublite",
        importpath = "cloud.google.com/go/pubsublite",
        sum = "h1:kqNk9e0gt0/Eu6lChjDlgu4B8k/399FZFpAU4LNq4Co=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_recaptchaenterprise_v2",
        importpath = "cloud.google.com/go/recaptchaenterprise/v2",
        sum = "h1:9qDHoQtUIZ8FDN5uFAfnNz8jyyXdpAJf6rMN4ZMmMcU=",
        version = "v2.26.0",
    )
    go_repository(
        name = "com_google_cloud_go_recommendationengine",
        importpath = "cloud.google.com/go/recommendationengine",
        sum = "h1:G0tmQXp67YDaduvde7AX+cIdTtDCo2PBMeqv8SBOpnw=",
        version = "v0.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_recommender",
        importpath = "cloud.google.com/go/recommender",
        sum = "h1:fJ6oO/7Ta/yzfRHcuJUln9iCeo6FDb5yIi2L7eLldzs=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_redis",
        importpath = "cloud.google.com/go/redis",
        sum = "h1:lV6xFUF7t3fm2ZRFjB4IVNkznG9EVwetYVJYwzmW8pc=",
        version = "v1.24.0",
    )
    go_repository(
        name = "com_google_cloud_go_resourcemanager",
        importpath = "cloud.google.com/go/resourcemanager",
        sum = "h1:aDf2RuuQ1Z+kbHHDzHyJNMHPGRFK+cMYfZIhWRaeRX0=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_resourcesettings",
        importpath = "cloud.google.com/go/resourcesettings",
        sum = "h1:13HOFU7v4cEvIHXSAQbinF4wp2Baybbq7q9FMctg1Ek=",
        version = "v1.8.3",
    )
    go_repository(
        name = "com_google_cloud_go_retail",
        importpath = "cloud.google.com/go/retail",
        sum = "h1:u27GLYlnNsQrqzcvSlpZK0L6Ig8k3PEuXGkUdNFCOn4=",
        version = "v1.32.0",
    )
    go_repository(
        name = "com_google_cloud_go_run",
        importpath = "cloud.google.com/go/run",
        sum = "h1:U56fxJWdrT+yjo4S/Vrtw5m69NdNL11Cyv9jX2JOi1s=",
        version = "v1.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_scheduler",
        importpath = "cloud.google.com/go/scheduler",
        sum = "h1:EPdChptxnvCasdMixuu58247qhKMO1iAlrVaQhvuRyE=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_secretmanager",
        importpath = "cloud.google.com/go/secretmanager",
        sum = "h1:GjE3NoyFXo7ipRPy26PMmg4oRX1Ra8fswH45r16rWV0=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_security",
        importpath = "cloud.google.com/go/security",
        sum = "h1:xp/htSPHpyqo225Ju6PjNS9jxfcOEQMxnF229Oze62g=",
        version = "v1.26.0",
    )
    go_repository(
        name = "com_google_cloud_go_securitycenter",
        importpath = "cloud.google.com/go/securitycenter",
        sum = "h1:k5uzLtTSxnh4fwG+8nIcuwDjEE3J+6xe2wcZ72cys7s=",
        version = "v1.45.0",
    )
    go_repository(
        name = "com_google_cloud_go_servicedirectory",
        importpath = "cloud.google.com/go/servicedirectory",
        sum = "h1:yrohWkwM8t5JfEFCmmlyksKnpMpvxM8XbRpwg+yIo64=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_shell",
        importpath = "cloud.google.com/go/shell",
        sum = "h1:vJ/g4BXCwBRMmUxHoNjdbE0DRh5Aysw/+vrZ/aPl0yc=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_spanner",
        importpath = "cloud.google.com/go/spanner",
        sum = "h1:XwXfcZ0kc1NT9Uu2IsThFiWtYptB+WgLn/KZEZcyzRg=",
        version = "v1.91.0",
    )
    go_repository(
        name = "com_google_cloud_go_speech",
        importpath = "cloud.google.com/go/speech",
        sum = "h1:ZcmRKcY02vkF3o/Eaa9yI3i1baNFTF4ctyDZVRfnNS0=",
        version = "v1.36.0",
    )
    go_repository(
        name = "com_google_cloud_go_storage",
        importpath = "cloud.google.com/go/storage",
        sum = "h1:SZq1t23NCI+e96dH77Dg3PEfsNNEjqO8zE5AnD8gVD0=",
        version = "v1.62.3",
    )
    go_repository(
        name = "com_google_cloud_go_storagetransfer",
        importpath = "cloud.google.com/go/storagetransfer",
        sum = "h1:jh+3SegMUulEhuJbkADCv8lmfRC8/pCpxja4oex2Ns8=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_talent",
        importpath = "cloud.google.com/go/talent",
        sum = "h1:dUNSwgBUFlDEMHJpmntEvDeLkPQAoEzBrvDbNfS/UeQ=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_texttospeech",
        importpath = "cloud.google.com/go/texttospeech",
        sum = "h1:xGqQv2LB4nttaBNMepCm19EQ0Sk0YdjzPvklDU+LYsI=",
        version = "v1.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_tpu",
        importpath = "cloud.google.com/go/tpu",
        sum = "h1:u2f8Lhezy6piu0BB3+4pTLw9xyOeU5n8PCgWjbjCitk=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_trace",
        importpath = "cloud.google.com/go/trace",
        sum = "h1:7+RAkTVGq3rmO2G74BdhpP1Bvbk3rFWzGt1MfMfHBt8=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_translate",
        importpath = "cloud.google.com/go/translate",
        sum = "h1:2dxWfGMG7y5vobxD2pPyngnfPPIk4MYdNB9QK5sEbuc=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_video",
        importpath = "cloud.google.com/go/video",
        sum = "h1:zPA41kNscRrJeRo/xD73TjNCFAnRDnP9ViNPuWARVeg=",
        version = "v1.33.0",
    )
    go_repository(
        name = "com_google_cloud_go_videointelligence",
        importpath = "cloud.google.com/go/videointelligence",
        sum = "h1:cd+s4jLMS59xiYLT6eQ6PVKP8R0PskgifXDa7YwMUUU=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_vision_v2",
        importpath = "cloud.google.com/go/vision/v2",
        sum = "h1:aTR1vj4++WtS9HD6YdGuoaYygMTJ873WaoV9sYjlQCc=",
        version = "v2.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_vmmigration",
        importpath = "cloud.google.com/go/vmmigration",
        sum = "h1:OPGMxx73owRHMve5B4L6W+0KQUwlunyE9kfbYzFHqug=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_vmwareengine",
        importpath = "cloud.google.com/go/vmwareengine",
        sum = "h1:vdSutihjzH5iYEFlN+1x+gDm4xzj6Md2yMzXADKI5cQ=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_vpcaccess",
        importpath = "cloud.google.com/go/vpcaccess",
        sum = "h1:fFa1i67AjC/3MuG5GPlvx36sIMeffpsIs3kI1C7S3Jk=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_webrisk",
        importpath = "cloud.google.com/go/webrisk",
        sum = "h1:t2WMlBo3xXxH6k5kTSkyIYUl24AHgmBI76F/5OqAPRY=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_websecurityscanner",
        importpath = "cloud.google.com/go/websecurityscanner",
        sum = "h1:iV+hAXeEo8kKtFWDZWPp9Z0fsziZnuOt4nCCYp68RP0=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_workflows",
        importpath = "cloud.google.com/go/workflows",
        sum = "h1:qROpn1zRDdeIjUP9w9O3MHXJEg8v/pzVMicAcw9vLvQ=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_lukechampine_uint128",
        importpath = "lukechampine.com/uint128",
        sum = "h1:cDdUVfRwDUDovz610ABgFD17nXD4/uDgVHl2sC3+sbo=",
        version = "v1.3.0",
    )
    go_repository(
        name = "dev_cel_expr",
        importpath = "cel.dev/expr",
        sum = "h1:K6j46C81hXtZQfuX60cVWQFBJahKSE2gfRbNuvr5bFs=",
        version = "v0.25.2",
    )
    go_repository(
        name = "ht_sr_git_sbinet_gg",
        importpath = "git.sr.ht/~sbinet/gg",
        sum = "h1:RIzgkizAk+9r7uPzf/VfbJHBMKUr0F5hRFxTUGMnt38=",
        version = "v0.6.0",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=",
        version = "v1.0.0-20201130134442-10cb98267c6c",
    )
    go_repository(
        name = "in_gopkg_errgo_v2",
        importpath = "gopkg.in/errgo.v2",
        sum = "h1:0vLT13EuvQ0hNvakwLuFZ/jYrLp5F3kcWHXdRggjCE8=",
        version = "v2.1.0",
    )
    go_repository(
        name = "in_gopkg_inf_v0",
        importpath = "gopkg.in/inf.v0",
        sum = "h1:73M5CoZyi3ZLMOyDlQh031Cx6N9NDJ2Vvfl76EDAgDc=",
        version = "v0.9.1",
    )
    go_repository(
        name = "in_gopkg_ini_v1",
        importpath = "gopkg.in/ini.v1",
        sum = "h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=",
        version = "v1.67.0",
    )
    go_repository(
        name = "in_gopkg_natefinch_lumberjack_v2",
        importpath = "gopkg.in/natefinch/lumberjack.v2",
        sum = "h1:bBRl1b0OH9s/DuPhuXpNl+VtCaJXFZ5/uEFST95x9zc=",
        version = "v2.2.1",
    )
    go_repository(
        name = "in_gopkg_tomb_v1",
        importpath = "gopkg.in/tomb.v1",
        sum = "h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=",
        version = "v1.0.0-20141024135613-dd632973f1e7",
    )
    go_repository(
        name = "in_gopkg_warnings_v0",
        importpath = "gopkg.in/warnings.v0",
        sum = "h1:wFXVbFY8DY5/xOe1ECiWdKCzZlxgshcYVNkBHstARME=",
        version = "v0.1.2",
    )
    go_repository(
        name = "in_gopkg_yaml_v2",
        importpath = "gopkg.in/yaml.v2",
        sum = "h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=",
        version = "v2.4.0",
    )
    go_repository(
        name = "in_gopkg_yaml_v3",
        importpath = "gopkg.in/yaml.v3",
        sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
        version = "v3.0.1",
    )
    go_repository(
        name = "io_etcd_go_bbolt",
        importpath = "go.etcd.io/bbolt",
        sum = "h1:S7GAl7Fxv12yohbwFfIbQCGDWbQbtDGPET4P/bD4lxU=",
        version = "v1.5.0",
    )
    go_repository(
        name = "io_etcd_go_etcd_api_v3",
        importpath = "go.etcd.io/etcd/api/v3",
        sum = "h1:WZlGK7pRtYGDB8ti8wkrQ5D2oWGMbtNL9VA5e+vF2Fg=",
        version = "v3.7.0",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_pkg_v3",
        importpath = "go.etcd.io/etcd/client/pkg/v3",
        sum = "h1:sW9njJzS3vXKcAJjjLQ4nk+avNUJ12Bcijcx8ehUskE=",
        version = "v3.7.0",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_v3",
        importpath = "go.etcd.io/etcd/client/v3",
        sum = "h1:5MHO37VbPB87VRPKUXEcicjeQWiTSjpPv3Ume8xPx20=",
        version = "v3.7.0",
    )
    go_repository(
        name = "io_etcd_go_etcd_pkg_v3",
        importpath = "go.etcd.io/etcd/pkg/v3",
        sum = "h1:aUBHGzXnYL2m3AlZfKp0ZhHQLVZIKfu9QsIlYdr+mxc=",
        version = "v3.7.0",
    )
    go_repository(
        name = "io_etcd_go_etcd_server_v3",
        importpath = "go.etcd.io/etcd/server/v3",
        sum = "h1:ScdUdN8ljuimp0lZaNq0otLMrHcFSFT+dQyT1j7JSFo=",
        version = "v3.7.0",
    )
    go_repository(
        name = "io_k8s_api",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/api",
        sum = "h1:Z//Vj9N7RA/yS2sDmxyeo7h+RR4zbUrd2vrd3Z0TbB4=",
        version = "v0.37.0",
    )
    go_repository(
        name = "io_k8s_apiextensions_apiserver",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/apiextensions-apiserver",
        sum = "h1:zRMQ3+/LIE5oZ0tVvXwYHC+dIkSP5cjNWju7AZU1LOI=",
        version = "v0.37.0",
    )
    go_repository(
        name = "io_k8s_apimachinery",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/apimachinery",
        sum = "h1:Np2AbDtf8x6RDHiD8T9LbKJ9gaegeVNa8yNm5FuGKm0=",
        version = "v0.37.0",
    )
    go_repository(
        name = "io_k8s_apiserver",
        importpath = "k8s.io/apiserver",
        sum = "h1:TXg7OxsOWrAH8J4Zi/gBAZuMw1Dfdd+6cca2h4qjRqo=",
        version = "v0.37.0",
    )
    go_repository(
        name = "io_k8s_client_go",
        importpath = "k8s.io/client-go",
        sum = "h1:nsN31fy8wBySuZ+QRnKmrjRSQLOG2rvoGN0tKd12zhQ=",
        version = "v0.37.0",
    )
    go_repository(
        name = "io_k8s_code_generator",
        importpath = "k8s.io/code-generator",
        sum = "h1:AC915wukzlVHHODAQYxvQ25WKibPh95faJ2kxf8dzuo=",
        version = "v0.37.0",
    )
    go_repository(
        name = "io_k8s_component_base",
        importpath = "k8s.io/component-base",
        sum = "h1:3SdSa4+itMdFTDFTeR8CxKGmSTSMXFlKL4ky8OqjguM=",
        version = "v0.37.0",
    )
    go_repository(
        name = "io_k8s_cri_api",
        importpath = "k8s.io/cri-api",
        sum = "h1:dsZXzrGrCEwHjsTDlAV7rutEplpMLY8bfNRMIqrtXjo=",
        version = "v0.31.3",
    )
    go_repository(
        name = "io_k8s_gengo_v2",
        importpath = "k8s.io/gengo/v2",
        sum = "h1:3L6PNkMLXkU/pz3jWzaaIUz0Rs2V9h+5O51AeRC7poc=",
        version = "v2.0.0-20260408192533-25e2208e0dc3",
    )
    go_repository(
        name = "io_k8s_klog",
        importpath = "k8s.io/klog",
        sum = "h1:Pt+yjF5aB1xDSVbau4VsWe+dQNzA0qv1LlXdC2dF6Q8=",
        version = "v1.0.0",
    )
    go_repository(
        name = "io_k8s_klog_v2",
        importpath = "k8s.io/klog/v2",
        sum = "h1:Tf+J3AH7xnUzZyVVXhTgGhEKnFqye14aadWv7bzXdzc=",
        version = "v2.140.0",
    )
    go_repository(
        name = "io_k8s_kms",
        importpath = "k8s.io/kms",
        sum = "h1:WDW3pSc5yFB1Od4YXXL3i5oK+vXI2v++LAV82+Ks1Z0=",
        version = "v0.37.0",
    )
    go_repository(
        name = "io_k8s_kube_openapi",
        importpath = "k8s.io/kube-openapi",
        sum = "h1:oXImqH8mQNk7PmvzKhmN3ddJoY6OnyM225MXwGHPm0A=",
        version = "v0.0.0-20260721132016-d427ff9ee9ad",
    )
    go_repository(
        name = "io_k8s_sigs_apiserver_network_proxy_konnectivity_client",
        importpath = "sigs.k8s.io/apiserver-network-proxy/konnectivity-client",
        sum = "h1:/YpDJ4vReG7ZmzSpBGxduXgywWkJU9zHubgJG03MT+Y=",
        version = "v0.36.0",
    )
    go_repository(
        name = "io_k8s_sigs_controller_runtime",
        importpath = "sigs.k8s.io/controller-runtime",
        sum = "h1:bEkK3PVOIVK9X8QWLGVhJgmFc++47vfT6wakSzAOLsQ=",
        version = "v0.25.2",
    )
    go_repository(
        name = "io_k8s_sigs_json",
        importpath = "sigs.k8s.io/json",
        sum = "h1:IpInykpT6ceI+QxKBbEflcR5EXP7sU1kvOlxwZh5txg=",
        version = "v0.0.0-20250730193827-2d320260d730",
    )
    go_repository(
        name = "io_k8s_sigs_structured_merge_diff_v4",
        importpath = "sigs.k8s.io/structured-merge-diff/v4",
        sum = "h1:150L+0vs/8DA78h1u02ooW1/fFq/Lwr+sGiqlzvrtq4=",
        version = "v4.4.1",
    )
    go_repository(
        name = "io_k8s_sigs_yaml",
        importpath = "sigs.k8s.io/yaml",
        sum = "h1:G8fkbMSAFqgEFgh4b1wmtzDnioxFCUgTZhlbj5P9QYs=",
        version = "v1.6.0",
    )
    go_repository(
        name = "io_k8s_utils",
        importpath = "k8s.io/utils",
        sum = "h1:jVkFFVfXdXP74B/zbO3hM3hpSFD0xvhQ5U686DPurkE=",
        version = "v0.0.0-20260707023825-cf1189d6abe3",
    )
    go_repository(
        name = "io_opencensus_go",
        importpath = "go.opencensus.io",
        sum = "h1:y73uSU6J157QMP2kn2r30vwW1A2W2WFwSCGnAVxeaD0=",
        version = "v0.24.0",
    )
    go_repository(
        name = "io_opentelemetry_go_auto_sdk",
        importpath = "go.opentelemetry.io/auto/sdk",
        sum = "h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=",
        version = "v1.2.1",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_detectors_gcp",
        importpath = "go.opentelemetry.io/contrib/detectors/gcp",
        sum = "h1:PI8dGkqDaQkwJ8kOopqMhDTbrnK3UIeG/RCHH4HErbo=",
        version = "v1.46.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc",
        importpath = "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc",
        sum = "h1:B2h3uqicet1CT2N5TOFhS+Gq++9i0/CLmaxvhmhtP5s=",
        version = "v0.71.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp",
        importpath = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
        sum = "h1:3g7B90UzBltIDKq1/5mrTGxTnOFDV0ICOhLoxiZ8jlg=",
        version = "v0.71.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel",
        importpath = "go.opentelemetry.io/otel",
        sum = "h1:FHt5/CDyVxi/8IM1CH7VE/rRgq3kLHa2mSTVMO8AWyc=",
        version = "v1.46.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace",
        sum = "h1:OFnwLJr+pF3iHrlGSzbxyuo6/6HyBlnlN1CWEJmBVcw=",
        version = "v1.46.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace_otlptracegrpc",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc",
        sum = "h1:w53CDeOA/Kurp7yRsegSr6pbbr759dOvJ+yNmWM6Hxs=",
        version = "v1.46.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_stdout_stdoutlog",
        importpath = "go.opentelemetry.io/otel/exporters/stdout/stdoutlog",
        sum = "h1:kvMAiLEudKmk+CSG+iYbU8vTUGNNDaf/V09OO5lrTwI=",
        version = "v0.22.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_log",
        importpath = "go.opentelemetry.io/otel/log",
        sum = "h1:5DBNnfvaJ6CVdkJ+Jle8Tzs50aSSv49TXGj9XRsEYw0=",
        version = "v0.22.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:yBnkXvgV7AXFILZc5K6IZe/CBFF3OS7BJ8ov6/lj0K8=",
        version = "v1.46.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk",
        importpath = "go.opentelemetry.io/otel/sdk",
        sum = "h1:h5CNQQjEbuQXY/JfZtgt3i7HVFV3aHPO2OAwO2eTYPI=",
        version = "v1.46.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_log",
        importpath = "go.opentelemetry.io/otel/sdk/log",
        sum = "h1:PRL+s6P63XT4E/bheEflopPUpVxuvANqZwtt89yhoGk=",
        version = "v0.22.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_metric",
        importpath = "go.opentelemetry.io/otel/sdk/metric",
        sum = "h1:0piZ26EG4RBfebb2jhDH6ERCYHoVWduc3kLgPCwSnSE=",
        version = "v1.46.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:OULy7ccdJnZtJ0UDYFOIGaCmiWzJ8Vi2G/Rsu60qs1c=",
        version = "v1.46.0",
    )
    go_repository(
        name = "io_opentelemetry_go_proto_otlp",
        importpath = "go.opentelemetry.io/proto/otlp",
        sum = "h1:5rrYs0Ykyj50sdU/JU0x8etU+LubXWb+gED6TbEdMIk=",
        version = "v1.11.0",
    )
    go_repository(
        name = "io_rsc_binaryregexp",
        importpath = "rsc.io/binaryregexp",
        sum = "h1:HfqmD5MEmC0zvwBuF187nq9mdnXjXsSivRiXN7SmRkE=",
        version = "v0.2.0",
    )
    go_repository(
        name = "io_rsc_pdf",
        importpath = "rsc.io/pdf",
        sum = "h1:k1MczvYDUvJBe93bYd7wrZLLUEcLZAuF824/I4e5Xr4=",
        version = "v0.1.1",
    )
    go_repository(
        name = "org_bitbucket_creachadair_stringset",
        importpath = "bitbucket.org/creachadair/stringset",
        sum = "h1:t1ejQyf8utS4GZV/4fM+1gvYucggZkfhb+tMobDxYOE=",
        version = "v0.0.14",
    )
    go_repository(
        name = "org_golang_google_api",
        importpath = "google.golang.org/api",
        sum = "h1:2rvPV2bqnPuHOaF4gGOBiT1IIc6JVXYyHCkZeqdzjNk=",
        version = "v0.300.0",
    )
    go_repository(
        name = "org_golang_google_appengine",
        importpath = "google.golang.org/appengine",
        sum = "h1:IhEN5q69dyKagZPYMSdIjS2HqprW324FRQZJcGqPAsM=",
        version = "v1.6.8",
    )
    go_repository(
        name = "org_golang_google_genproto",
        build_file_proto_mode = "disable",
        importpath = "google.golang.org/genproto",
        sum = "h1:C9v1o0/4quuhOAfmRXA2j+we0PqZIp8traLdeogF3Ms=",
        version = "v0.0.0-20260715232425-e75dac1f907d",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:rZHRz0ogq4Pid1IKtA/ivHi8XONVF5tqUQXA2TtgpkE=",
        version = "v0.0.0-20260928230214-8a89bd6388cc",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_bytestream",
        importpath = "google.golang.org/genproto/googleapis/bytestream",
        sum = "h1:YeHSHMle+uLS75+Z+FuOA7bH9wibPWVkMRWD+9Aaz7I=",
        version = "v0.0.0-20260921155816-b14227669459",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:4bNTbnb44EqGVy9HaQxSY2jnafSUdgV3qHtedvNpDKg=",
        version = "v0.0.0-20260928230214-8a89bd6388cc",
    )
    go_repository(
        name = "org_golang_google_grpc",
        build_file_proto_mode = "disable",
        importpath = "google.golang.org/grpc",
        sum = "h1:soMyaPJ8pAak5PIQ0DGBUir0XRo2fRoMqhNWMLlLxO0=",
        version = "v1.84.0",
    )
    go_repository(
        name = "org_golang_google_grpc_cmd_protoc_gen_go_grpc",
        importpath = "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
        sum = "h1:rgSNvqscFZ1JgV/4wH5GOsZFSFkR2Eua9As3KIr2LlM=",
        version = "v1.6.2",
    )
    go_repository(
        name = "org_golang_google_protobuf",
        importpath = "google.golang.org/protobuf",
        sum = "h1:pJOKDDOyeXErUroCihFAd5LQuwXBSpVnKGrj5o/fwxc=",
        version = "v1.36.12",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:3ZVCjf8Ggz7zneR/EHRVx68Ctf+2pmIMP2UFhh9cC6M=",
        version = "v0.57.0",
    )
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:Ck8QetSgk912qxWLMCKxd0in+aiyBQyDSMae6e/xmpU=",
        version = "v0.0.0-20260908205506-85c1c2202aba",
    )
    go_repository(
        name = "org_golang_x_exp_typeparams",
        importpath = "golang.org/x/exp/typeparams",
        sum = "h1:UhRVJ0i7bF9n/Hd8YjW3eKjlPVBHzbQdxrBgjbSKl64=",
        version = "v0.0.0-20231219180239-dc181d75b848",
    )
    go_repository(
        name = "org_golang_x_image",
        importpath = "golang.org/x/image",
        sum = "h1:Y6uW6rH1y5y/LK1J8BPWZtr6yZ7hrsy6hFrXjgsc2fQ=",
        version = "v0.25.0",
    )
    go_repository(
        name = "org_golang_x_lint",
        importpath = "golang.org/x/lint",
        sum = "h1:Wh+f8QHJXR411sJR8/vRBTZ7YapZaRvUcLFFJhusH0k=",
        version = "v0.0.0-20200302205851-738671d3881b",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:qJmnOUb4YB+FsEuM3HcWucdZASCPGhsX6uljO6pog0c=",
        version = "v0.41.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:5zfYln+w5XCxwrnMMJPufRgNoXEaGxl0wo5GqPXyues=",
        version = "v0.59.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:JUlcxA8oAtauLfiH8FX2/FkAWHAdi0QtGCGc+hofE98=",
        version = "v0.37.0",
    )
    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:KameEIfc1IkluZyXWLn39Wd4tURc6GbCiISGiZm2bQk=",
        version = "v0.23.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:bbX/i/6MgT9BVLM9RT1thmxL04yeTAhbEz4SyadbXoo=",
        version = "v0.48.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:3+OXuTbaKDgwk8jTi3aSLHRlmWqHEUDUtxnbFigO4YE=",
        version = "v0.46.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:JbOZXgfeCPU9gacVtYliJqOhD+zhrEqK4LfdpmlUZqI=",
        version = "v0.42.0",
    )
    go_repository(
        name = "org_golang_x_time",
        importpath = "golang.org/x/time",
        sum = "h1:vMb6ptszcQMkcwiRTAuNNU50gom6++Q/6gY2hDM6VDE=",
        version = "v0.16.0",
    )
    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:c2ifzfcuY7L90lZ2aKd8S4K2NpASF08SZx9ZuJkHmSU=",
        version = "v0.50.0",
    )
    go_repository(
        name = "org_golang_x_xerrors",
        importpath = "golang.org/x/xerrors",
        sum = "h1:noIWHXmPHxILtqtCOPIhSt0ABwskkZKjD3bXGnZGpNY=",
        version = "v0.0.0-20240903120638-7835f813f4da",
    )
    go_repository(
        name = "org_gonum_v1_gonum",
        importpath = "gonum.org/v1/gonum",
        sum = "h1:VbpOemQlsSMrYmn7T2OUvQ4dqxQXU+ouZFQsZOx50z4=",
        version = "v0.17.0",
    )
    go_repository(
        name = "org_gonum_v1_plot",
        importpath = "gonum.org/v1/plot",
        sum = "h1:Tlfh/jBk2tqjLZ4/P8ZIwGrLEWQSPDLRm/SNWKNXiGI=",
        version = "v0.15.2",
    )
    go_repository(
        name = "org_modernc_cc_v4",
        importpath = "modernc.org/cc/v4",
        sum = "h1:q+NXGJ0bK3b4TXFYQQVr9pYETGnmwFWkrUzJnMya/Tg=",
        version = "v4.29.7",
    )
    go_repository(
        name = "org_modernc_ccorpus2",
        importpath = "modernc.org/ccorpus2",
        sum = "h1:JODbahvmixwZDoF1kZs/9Rj6aHA2vAPP7HKiNDq//xo=",
        version = "v1.6.0",
    )
    go_repository(
        name = "org_modernc_mathutil",
        importpath = "modernc.org/mathutil",
        sum = "h1:GCZVGXdaN8gTqB1Mf/usp1Y/hSqgI2vAGGP4jZMCxOU=",
        version = "v1.7.1",
    )
    go_repository(
        name = "org_modernc_opt",
        importpath = "modernc.org/opt",
        sum = "h1:tGyef5ApycA7FSEOMraay9SaTk5zmbx7Tu+cJs4QKZg=",
        version = "v0.2.0",
    )
    go_repository(
        name = "org_modernc_strutil",
        importpath = "modernc.org/strutil",
        sum = "h1:UneZBkQA+DX2Rp35KcM69cSsNES9ly8mQWD71HKlOA0=",
        version = "v1.2.1",
    )
    go_repository(
        name = "org_modernc_token",
        importpath = "modernc.org/token",
        sum = "h1:Xl7Ap9dKaEs5kLoOQeQmPWevfnk/DM5qcLcYlA8ys6Y=",
        version = "v1.1.0",
    )
    go_repository(
        name = "org_mozilla_go_pkcs7",
        importpath = "go.mozilla.org/pkcs7",
        sum = "h1:yM4/HS9dYv7ri2biPtxt8ikvB37a980dg69/pKmS+eI=",
        version = "v0.9.0",
    )
    go_repository(
        name = "org_uber_go_atomic",
        importpath = "go.uber.org/atomic",
        sum = "h1:ZvwS0R+56ePWxUNi+Atn9dWONBPp/AUETXlHW0DxSjE=",
        version = "v1.11.0",
    )
    go_repository(
        name = "org_uber_go_goleak",
        importpath = "go.uber.org/goleak",
        sum = "h1:2K3zAYmnTNqV73imy9J1T3WC+gmCePx2hEGkimedGto=",
        version = "v1.3.0",
    )
    go_repository(
        name = "org_uber_go_mock",
        importpath = "go.uber.org/mock",
        sum = "h1:hyF9dfmbgIX5EfOdasqLsWD6xqpNZlXblLB/Dbnwv3Y=",
        version = "v0.6.0",
    )
    go_repository(
        name = "org_uber_go_multierr",
        importpath = "go.uber.org/multierr",
        sum = "h1:blXXJkSxSSfBVBlC76pxqeO+LN3aDfLQo+309xJstO0=",
        version = "v1.11.0",
    )
    go_repository(
        name = "org_uber_go_zap",
        importpath = "go.uber.org/zap",
        sum = "h1:08RqriUEv8+ArZRYSTXy1LeBScaMpVSTBhCeaZYfMYc=",
        version = "v1.27.1",
    )
    go_repository(
        name = "tech_einride_go_aip",
        importpath = "go.einride.tech/aip",
        sum = "h1:TI21IdeOnLTwZEJ3BxtImIZk6bsN2Q+sd0x99SLiQ+M=",
        version = "v0.83.0",
    )
    go_repository(
        name = "tf_universe_go_metallb",
        importpath = "go.universe.tf/metallb",
        sum = "h1:iVtjCslcUFMOEvuNeAyW2ygZxyW+6YjnNcEuCTOZezc=",
        version = "v0.16.1",
    )
    go_repository(
        name = "tools_gotest_v3",
        importpath = "gotest.tools/v3",
        sum = "h1:7koQfIKdy+I8UTetycgUqXWSDwpgv193Ka+qRsmBY8Q=",
        version = "v3.5.2",
    )
    go_repository(
        name = "xyz_gomodules_jsonpatch_v2",
        importpath = "gomodules.xyz/jsonpatch/v2",
        sum = "h1:JELs8RLM12qJGXU4u/TO3V25KW8GreMKl9pdkk14RM0=",
        version = "v2.5.0",
    )
    go_repository(
        name = "build_buf_go_protovalidate",
        importpath = "buf.build/go/protovalidate",
        sum = "h1:IAG1etULddAy93fiBsFVhpj7es5zL53AfB/79CVGtyY=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_antlr4_go_antlr_v4",
        importpath = "github.com/antlr4-go/antlr/v4",
        sum = "h1:SqQKkuVZ+zWkMMNkjy5FZe5mr5WURWnlpmOuzYWrPrQ=",
        version = "v4.13.1",
    )
    go_repository(
        name = "com_github_apache_arrow_go_v15",
        importpath = "github.com/apache/arrow/go/v15",
        sum = "h1:60IliRbiyTWCWjERBCkO1W4Qun9svcYoZrSLcyOsMLE=",
        version = "v15.0.2",
    )
    go_repository(
        name = "com_github_apapsch_go_jsonmerge_v2",
        importpath = "github.com/apapsch/go-jsonmerge/v2",
        sum = "h1:axGnT1gRIfimI7gJifB699GoE/oq+F2MU7Dml6nw9rQ=",
        version = "v2.0.0",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2",
        importpath = "github.com/aws/aws-sdk-go-v2",
        sum = "h1:mJoei2CxPutQVxaATCzDUjcZEjVRdpsiiXi2o38yqWM=",
        version = "v1.36.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_config",
        importpath = "github.com/aws/aws-sdk-go-v2/config",
        sum = "h1:Y/2a+jLPrPbHpFkpAAYkVEtJmxORlXoo5k2g1fa2sUo=",
        version = "v1.29.12",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_credentials",
        importpath = "github.com/aws/aws-sdk-go-v2/credentials",
        sum = "h1:q+nV2yYegofO/SUXruT+pn4KxkxmaQ++1B/QedcKBFM=",
        version = "v1.17.65",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_feature_ec2_imds",
        importpath = "github.com/aws/aws-sdk-go-v2/feature/ec2/imds",
        sum = "h1:x793wxmUWVDhshP8WW2mlnXuFrO4cOd3HLBroh1paFw=",
        version = "v1.16.30",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_internal_configsources",
        importpath = "github.com/aws/aws-sdk-go-v2/internal/configsources",
        sum = "h1:ZK5jHhnrioRkUNOc+hOgQKlUL5JeC3S6JgLxtQ+Rm0Q=",
        version = "v1.3.34",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_internal_endpoints_v2",
        importpath = "github.com/aws/aws-sdk-go-v2/internal/endpoints/v2",
        sum = "h1:SZwFm17ZUNNg5Np0ioo/gq8Mn6u9w19Mri8DnJ15Jf0=",
        version = "v2.6.34",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_internal_ini",
        importpath = "github.com/aws/aws-sdk-go-v2/internal/ini",
        sum = "h1:bIqFDwgGXXN1Kpp99pDOdKMTTb5d2KyU5X/BZxjOkRo=",
        version = "v1.8.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_internal_accept_encoding",
        importpath = "github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding",
        sum = "h1:eAh2A4b5IzM/lum78bZ590jy36+d/aFLgKF/4Vd1xPE=",
        version = "v1.12.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_internal_presigned_url",
        importpath = "github.com/aws/aws-sdk-go-v2/service/internal/presigned-url",
        sum = "h1:dM9/92u2F1JbDaGooxTq18wmmFzbJRfXfVfy96/1CXM=",
        version = "v1.12.15",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_sso",
        importpath = "github.com/aws/aws-sdk-go-v2/service/sso",
        sum = "h1:pdgODsAhGo4dvzC3JAG5Ce0PX8kWXrTZGx+jxADD+5E=",
        version = "v1.25.2",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_ssooidc",
        importpath = "github.com/aws/aws-sdk-go-v2/service/ssooidc",
        sum = "h1:90uX0veLKcdHVfvxhkWUQSCi5VabtwMLFutYiRke4oo=",
        version = "v1.30.0",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_sts",
        importpath = "github.com/aws/aws-sdk-go-v2/service/sts",
        sum = "h1:PZV5W8yk4OtH1JAuhV2PXwwO9v5G5Aoj+eMCn4T+1Kc=",
        version = "v1.33.17",
    )
    go_repository(
        name = "com_github_aws_smithy_go",
        importpath = "github.com/aws/smithy-go",
        sum = "h1:Z//5NuZCSW6R4PhQ93hShNbyBbn8BWCmCVCt+Q8Io5k=",
        version = "v1.22.3",
    )
    go_repository(
        name = "com_github_aymerick_douceur",
        importpath = "github.com/aymerick/douceur",
        sum = "h1:Mv+mAeH1Q+n9Fr+oyamOlAkUNPWPlA8PPGR0QAaYuPk=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_bmatcuk_doublestar",
        importpath = "github.com/bmatcuk/doublestar",
        sum = "h1:YroD6BJCZBYx06yYFEWvUuKVWQn3vLLQAVmDmvTSaiQ=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_bytedance_sonic",
        importpath = "github.com/bytedance/sonic",
        sum = "h1:oUp34TzMlL+OY1OUWxHqsdkgC/Zfc85zGqw9siXjrc0=",
        version = "v1.11.6",
    )
    go_repository(
        name = "com_github_bytedance_sonic_loader",
        importpath = "github.com/bytedance/sonic/loader",
        sum = "h1:c+e5Pt1k/cy5wMveRDyk2X4B9hF4g7an8N3zCYjJFNM=",
        version = "v0.1.1",
    )
    go_repository(
        name = "com_github_campoy_embedmd",
        importpath = "github.com/campoy/embedmd",
        sum = "h1:V4kI2qTJJLf4J29RzI/MAt2c3Bl4dQSYPuflzwFH2hY=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_cenkalti_backoff_v5",
        importpath = "github.com/cenkalti/backoff/v5",
        sum = "h1:ZN+IMa753KfX5hd8vVaMixjnqRZ3y8CuJKRKj1xcsSM=",
        version = "v5.0.3",
    )
    go_repository(
        name = "com_github_cloudwego_base64x",
        importpath = "github.com/cloudwego/base64x",
        sum = "h1:jwCgWpFanWmN8xoIUHa2rtzmkd5J2plF/dnLS6Xd/0Y=",
        version = "v0.1.4",
    )
    go_repository(
        name = "com_github_cloudwego_iasm",
        importpath = "github.com/cloudwego/iasm",
        sum = "h1:1KNIy1I1H9hNNFEEH3DVnI4UujN+1zjpuk6gwHLTssg=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_cloudykit_fastprinter",
        importpath = "github.com/CloudyKit/fastprinter",
        sum = "h1:sR+/8Yb4slttB4vD+b9btVEnWgL3Q00OBTzVT8B9C0c=",
        version = "v0.0.0-20200109182630-33d98a066a53",
    )
    go_repository(
        name = "com_github_cloudykit_jet_v6",
        importpath = "github.com/CloudyKit/jet/v6",
        sum = "h1:EpcZ6SR9n28BUGtNJSvlBqf90IpjeFr36Tizxhn/oME=",
        version = "v6.2.0",
    )
    go_repository(
        name = "com_github_containerd_errdefs",
        importpath = "github.com/containerd/errdefs",
        sum = "h1:tg5yIfIlQIrxYtu9ajqY42W3lpS19XqdxRQeEwYG8PI=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_containerd_errdefs_pkg",
        importpath = "github.com/containerd/errdefs/pkg",
        sum = "h1:9IKJ06FvyNlexW690DXuQNx2KA2cUJXx151Xdx3ZPPE=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_coreos_go_oidc",
        importpath = "github.com/coreos/go-oidc",
        sum = "h1:6W0vGJR3Tu0r0PwfmjOrRZSlfxeEln8dsejt3ZWIvwo=",
        version = "v2.5.0+incompatible",
    )
    go_repository(
        name = "com_github_distribution_reference",
        importpath = "github.com/distribution/reference",
        sum = "h1:0IXCQ5g4/QMHHkarYzh5l+u8T3t73zM5QvfrDyIgxBk=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_github_eapache_go_resiliency",
        importpath = "github.com/eapache/go-resiliency",
        sum = "h1:n3NRTnBn5N0Cbi/IeOHuQn9s2UwVUH7Ga0ZWcP+9JTA=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane_envoy",
        importpath = "github.com/envoyproxy/go-control-plane/envoy",
        sum = "h1:u3riX6BoYRfF4Dr7dwSOroNfdSbEPe9Yyl09/B6wBrQ=",
        version = "v1.37.0",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane_ratelimit",
        importpath = "github.com/envoyproxy/go-control-plane/ratelimit",
        sum = "h1:/G9QYbddjL25KvtKTv3an9lx6VBE2cnb8wp1vEGNYGI=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_fatih_structs",
        importpath = "github.com/fatih/structs",
        sum = "h1:Q7juDM0QtcnhCpeyLGQKyg4TOIghuNXrkL32pHAUMxo=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_flosch_pongo2_v4",
        importpath = "github.com/flosch/pongo2/v4",
        sum = "h1:gv+5Pe3vaSVmiJvh/BZa82b7/00YUGm0PIyVVLop0Hw=",
        version = "v4.0.2",
    )
    go_repository(
        name = "com_github_fxamacker_cbor_v2",
        importpath = "github.com/fxamacker/cbor/v2",
        sum = "h1:2rWm8B193Ll4VdjsJY28jxs70IdDsHRWgQYAI80+rMQ=",
        version = "v2.9.1",
    )
    go_repository(
        name = "com_github_gabriel_vasile_mimetype",
        importpath = "github.com/gabriel-vasile/mimetype",
        sum = "h1:in2uUcidCuFcDKtdcBxlR0rJ1+fsokWf+uqxgUFjbI0=",
        version = "v1.4.3",
    )
    go_repository(
        name = "com_github_gin_contrib_sse",
        importpath = "github.com/gin-contrib/sse",
        sum = "h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_gin_gonic_gin",
        importpath = "github.com/gin-gonic/gin",
        sum = "h1:T0ujvqyCSqRopADpgPgiTT63DUQVSfojyME59Ei63pQ=",
        version = "v1.10.1",
    )
    go_repository(
        name = "com_github_gkampitakis_ciinfo",
        importpath = "github.com/gkampitakis/ciinfo",
        sum = "h1:JcuOPk8ZU7nZQjdUhctuhQofk7BGHuIy0c9Ez8BNhXs=",
        version = "v0.3.2",
    )
    go_repository(
        name = "com_github_gkampitakis_go_diff",
        importpath = "github.com/gkampitakis/go-diff",
        sum = "h1:Qyn0J9XJSDTgnsgHRdz9Zp24RaJeKMUHg2+PDZZdC4M=",
        version = "v1.3.2",
    )
    go_repository(
        name = "com_github_gkampitakis_go_snaps",
        importpath = "github.com/gkampitakis/go-snaps",
        sum = "h1:amyJrvM1D33cPHwVrjo9jQxX8g/7E2wYdZ+01KS3zGE=",
        version = "v0.5.15",
    )
    go_repository(
        name = "com_github_go_jose_go_jose_v4",
        importpath = "github.com/go-jose/go-jose/v4",
        sum = "h1:moDMcTHmvE6Groj34emNPLs/qtYXRVcd6S7NHbHz3kA=",
        version = "v4.1.4",
    )
    go_repository(
        name = "com_github_go_openapi_analysis",
        importpath = "github.com/go-openapi/analysis",
        sum = "h1:sNvbAGCJqUTqIAodr9IVqJMmuZas3YS9ms1dGK9yiJ4=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_go_openapi_errors",
        importpath = "github.com/go-openapi/errors",
        sum = "h1:HI9+SyVYiRzyeBQGv0CbWbQa+u2S0nny7G7AhY+UMiM=",
        version = "v0.22.9",
    )
    go_repository(
        name = "com_github_go_openapi_loads",
        importpath = "github.com/go-openapi/loads",
        sum = "h1:+uNsDlRQfYtZTrh+3pdwampcAqZVPuBJW0IA82aZHII=",
        version = "v0.25.2",
    )
    go_repository(
        name = "com_github_go_openapi_runtime",
        importpath = "github.com/go-openapi/runtime",
        sum = "h1:HSxskMs0WmpCdQvBWVxHt2t2mXMwn8DDav3VtVidwig=",
        version = "v0.33.2",
    )
    go_repository(
        name = "com_github_go_openapi_runtime_server_middleware",
        importpath = "github.com/go-openapi/runtime/server-middleware",
        sum = "h1:BVFjAaW4Jh/kZ4QSsgDrGzLuSuaKzSiyvappUFKxgUk=",
        version = "v0.33.2",
    )
    go_repository(
        name = "com_github_go_openapi_spec",
        importpath = "github.com/go-openapi/spec",
        sum = "h1:JtB/GHOj+eetjse6YvxqLze88oEekl/4uPBethvzRrA=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_go_openapi_strfmt",
        importpath = "github.com/go-openapi/strfmt",
        sum = "h1:SG32SlbwNy92s0KJiVxt2joJeFdqIYHvwrA0OU6HqzQ=",
        version = "v0.27.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag_cmdutils",
        importpath = "github.com/go-openapi/swag/cmdutils",
        sum = "h1:cXEzX/nWCODon251f+1HgL75/R4XXML5IZxgLCWMoro=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag_conv",
        importpath = "github.com/go-openapi/swag/conv",
        sum = "h1:8c9shoB8l0QRSR6ymq1llHdBWDqpIgJfjTGHx3Vuhm0=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag_fileutils",
        importpath = "github.com/go-openapi/swag/fileutils",
        sum = "h1:mdUL+Vw5ah1fO1AvFCoHeIyv7YZsCz2MN3KNqd/lLVI=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag_jsonname",
        importpath = "github.com/go-openapi/swag/jsonname",
        sum = "h1:8p150i44rv/Drip4vWI3kGi9+4W9TdI3US3uUYSFhSo=",
        version = "v0.25.5",
    )
    go_repository(
        name = "com_github_go_openapi_swag_jsonutils",
        importpath = "github.com/go-openapi/swag/jsonutils",
        sum = "h1:uZNSD2/rJDYAfvsLYkykTzXuyxM3QpDKV9jhRHrMu6k=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag_jsonutils_fixtures_test",
        importpath = "github.com/go-openapi/swag/jsonutils/fixtures_test",
        sum = "h1:w+Fd6EBOMGlC74GYGHqCLGyb3Bpams8TtNrgM/SG4jo=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag_loading",
        importpath = "github.com/go-openapi/swag/loading",
        sum = "h1:QU1ry24e6r6Shoxe380Nx0X5iKkpRO+hJffsyJ2dVlY=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag_mangling",
        importpath = "github.com/go-openapi/swag/mangling",
        sum = "h1:ltdo3K0tTP4P6zWhVVWDU9D37/YYBuMJRdZxiS2Dtog=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag_netutils",
        importpath = "github.com/go-openapi/swag/netutils",
        sum = "h1:dBli6jyUa93sDS/XYWShJzf4+nnVfFKfvY9SLUaUzB8=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag_pools",
        importpath = "github.com/go-openapi/swag/pools",
        sum = "h1:PlGoDRF8WtSyXkedZZkpW3f9wDhFf3u8KtaBcmgmh8c=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag_stringutils",
        importpath = "github.com/go-openapi/swag/stringutils",
        sum = "h1:lcnBxwAaysT3bMUVBfn9V/PkfVqkurpqufKU6rTUMGk=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag_typeutils",
        importpath = "github.com/go-openapi/swag/typeutils",
        sum = "h1:O7aVvkTs3pXwikgtrPLenigEMdQFgONKRooQA9pgf2I=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag_yamlutils",
        importpath = "github.com/go-openapi/swag/yamlutils",
        sum = "h1:IFKFFeDnuIwzfsuWRQU+rI8iL3g4XyAYjV/RlnlxPLM=",
        version = "v0.29.2",
    )
    go_repository(
        name = "com_github_go_openapi_testify_enable_yaml_v2",
        importpath = "github.com/go-openapi/testify/enable/yaml/v2",
        sum = "h1:wPW6YRgx3+SID1yUy/Xwa17L8kFEaEKod2VRbJDZNUs=",
        version = "v2.7.0",
    )
    go_repository(
        name = "com_github_go_openapi_testify_v2",
        importpath = "github.com/go-openapi/testify/v2",
        sum = "h1:19QDx5b57p8KjO2E4sgs/woL4Akp+FGLRlVm2xn0Rho=",
        version = "v2.8.0",
    )
    go_repository(
        name = "com_github_go_openapi_validate",
        importpath = "github.com/go-openapi/validate",
        sum = "h1:dFsYCLVUQUL6Vi2lQSexgwmCXDuHe7eWRDhQxkE+xYA=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_go_task_slim_sprig_v3",
        importpath = "github.com/go-task/slim-sprig/v3",
        sum = "h1:sUs3vkvUymDpBKi3qH1YSqBQk9+9D/8M2mN1vB6EwHI=",
        version = "v3.0.0",
    )
    go_repository(
        name = "com_github_go_viper_mapstructure_v2",
        importpath = "github.com/go-viper/mapstructure/v2",
        sum = "h1:vM5IJoUAy3d7zRSVtIwQgBj7BiWtMPfmPEgAXnvj1Ro=",
        version = "v2.5.0",
    )
    go_repository(
        name = "com_github_goccmack_gocc",
        importpath = "github.com/goccmack/gocc",
        sum = "h1:PHv20lcM1Erz+kovS+c07DnDFp6X5cvghndtTXuEyfE=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_golang_jwt_jwt_v5",
        importpath = "github.com/golang-jwt/jwt/v5",
        sum = "h1:kYf81DTWFe7t+1VvL7eS+jKFVWaUnK9cB1qbwn63YCY=",
        version = "v5.3.1",
    )
    go_repository(
        name = "com_github_gomarkdown_markdown",
        importpath = "github.com/gomarkdown/markdown",
        sum = "h1:4gjrh/PN2MuWCCElk8/I4OCKRKWCCo2zEct3VKCbibU=",
        version = "v0.0.0-20240328165702-4d01890c35c0",
    )
    go_repository(
        name = "com_github_google_go_tpm",
        importpath = "github.com/google/go-tpm",
        sum = "h1:slArAR9Ft+1ybZu0lBwpSmpwhRXaa85hWtMinMyRAWo=",
        version = "v0.9.8",
    )
    go_repository(
        name = "com_github_googleapis_cloud_bigtable_clients_test",
        importpath = "github.com/googleapis/cloud-bigtable-clients-test",
        sum = "h1:9f+NjgCZHLx7t2lyL+Zhyg5f9/eFNI4RryGWlbcTwqk=",
        version = "v0.0.4",
    )
    go_repository(
        name = "com_github_googlecloudplatform_grpc_gcp_go_grpcgcp",
        importpath = "github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp",
        sum = "h1:BzsL0qE7LvtTEtXG7Dt5NS1EP0CQwI21HZfj9aGghhw=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_gorilla_css",
        importpath = "github.com/gorilla/css",
        sum = "h1:BQqNyPTi50JCFMTw/b67hByjMVXZRwGha6wxVGkeihY=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_go_grpc_middleware_providers_prometheus",
        importpath = "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus",
        sum = "h1:QGLs/O40yoNK9vmy4rhUGBVyMf1lISBGtXRpsu/Qu/o=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_uuid",
        importpath = "github.com/hashicorp/go-uuid",
        sum = "h1:2gKiV6YVmrJ1i2CKKa9obLvRieoRGviZFL26PcT/Co8=",
        version = "v1.0.3",
    )
    go_repository(
        name = "com_github_hashicorp_golang_lru_v2",
        importpath = "github.com/hashicorp/golang-lru/v2",
        sum = "h1:a+bsQ5rvGLjzHuww6tVxozPZFVghXaHOwFs4luLUK2k=",
        version = "v2.0.7",
    )
    go_repository(
        name = "com_github_ibm_sarama",
        importpath = "github.com/IBM/sarama",
        sum = "h1:zpY2iZYmt+z+0Bo3aYF+cD48OBt2hIgiDPZUuZKTXcc=",
        version = "v1.50.3",
    )
    go_repository(
        name = "com_github_iris_contrib_schema",
        importpath = "github.com/iris-contrib/schema",
        sum = "h1:CPSBLyx2e91H2yJzPuhGuifVRnZBBJ3pCOMbOvPZaTw=",
        version = "v0.0.6",
    )
    go_repository(
        name = "com_github_jcmturner_aescts_v2",
        importpath = "github.com/jcmturner/aescts/v2",
        sum = "h1:9YKLH6ey7H4eDBXW8khjYslgyqG2xZikXP0EQFKrle8=",
        version = "v2.0.0",
    )
    go_repository(
        name = "com_github_jcmturner_dnsutils_v2",
        importpath = "github.com/jcmturner/dnsutils/v2",
        sum = "h1:lltnkeZGL0wILNvrNiVCR6Ro5PGU/SeBvVO/8c/iPbo=",
        version = "v2.0.0",
    )
    go_repository(
        name = "com_github_jcmturner_gofork",
        importpath = "github.com/jcmturner/gofork",
        sum = "h1:QH0l3hzAU1tfT3rZCnW5zXl+orbkNMMRGJfdJjHVETg=",
        version = "v1.7.6",
    )
    go_repository(
        name = "com_github_jcmturner_gokrb5_v8",
        importpath = "github.com/jcmturner/gokrb5/v8",
        sum = "h1:x1Sv4HaTpepFkXbt2IkL29DXRf8sOfZXo8eRKh687T8=",
        version = "v8.4.4",
    )
    go_repository(
        name = "com_github_jcmturner_rpc_v2",
        importpath = "github.com/jcmturner/rpc/v2",
        sum = "h1:7FXXj8Ti1IaVFpSAziCZWNzbNuZmnvw/i6CqLNdWfZY=",
        version = "v2.0.3",
    )
    go_repository(
        name = "com_github_joker_jade",
        importpath = "github.com/Joker/jade",
        sum = "h1:Qbeh12Vq6BxURXT1qZBRHsDxeURB8ztcL6f3EXSGeHk=",
        version = "v1.1.3",
    )
    go_repository(
        name = "com_github_joshdk_go_junit",
        importpath = "github.com/joshdk/go-junit",
        sum = "h1:S86cUKIdwBHWwA6xCmFlf3RTLfVXYQfvanM5Uh+K6GE=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_juju_gnuflag",
        importpath = "github.com/juju/gnuflag",
        sum = "h1:c93kUJDtVAXFEhsCh5jSxyOJmFHuzcihnslQiX8Urwo=",
        version = "v0.0.0-20171113085948-2ce1bb71843d",
    )
    go_repository(
        name = "com_github_kataras_blocks",
        importpath = "github.com/kataras/blocks",
        sum = "h1:MrpVhoFTCR2v1iOOfGng5VJSILKeZZI+7NGfxEh3SUM=",
        version = "v0.0.8",
    )
    go_repository(
        name = "com_github_kataras_golog",
        importpath = "github.com/kataras/golog",
        sum = "h1:dGkcCVsIpqiAMWTlebn/ZULHxFvfG4K43LF1cNWSh20=",
        version = "v0.1.11",
    )
    go_repository(
        name = "com_github_kataras_iris_v12",
        importpath = "github.com/kataras/iris/v12",
        sum = "h1:sGgo43rMPfzDft8rjVhPs6L3qDJy3TbBrMD/zGL1pzk=",
        version = "v12.2.11",
    )
    go_repository(
        name = "com_github_kataras_pio",
        importpath = "github.com/kataras/pio",
        sum = "h1:x0rXVX0fviDTXOOLOmr4MUxOabu1InVSTu5itF8CXCM=",
        version = "v0.0.13",
    )
    go_repository(
        name = "com_github_kataras_sitemap",
        importpath = "github.com/kataras/sitemap",
        sum = "h1:w71CRMMKYMJh6LR2wTgnk5hSgjVNB9KL60n5e2KHvLY=",
        version = "v0.0.6",
    )
    go_repository(
        name = "com_github_kataras_tunnel",
        importpath = "github.com/kataras/tunnel",
        sum = "h1:sCAqWuJV7nPzGrlb0os3j49lk2JhILT0rID38NHNLpA=",
        version = "v0.0.4",
    )
    go_repository(
        name = "com_github_labstack_echo_v4",
        importpath = "github.com/labstack/echo/v4",
        sum = "h1:S9keusg26gZpjMmPqB5hOEvNKnmd1lNmcHrbbH2lnFs=",
        version = "v4.15.1",
    )
    go_repository(
        name = "com_github_labstack_gommon",
        importpath = "github.com/labstack/gommon",
        sum = "h1:F8qTUNXgG1+6WQmqoUWnz8WiEU60mXVVw0P4ht1WRA0=",
        version = "v0.4.2",
    )
    go_repository(
        name = "com_github_mailgun_raymond_v2",
        importpath = "github.com/mailgun/raymond/v2",
        sum = "h1:5dmlB680ZkFG2RN/0lvTAghrSxIESeu9/2aeDqACtjw=",
        version = "v2.0.48",
    )
    go_repository(
        name = "com_github_maruel_natural",
        importpath = "github.com/maruel/natural",
        sum = "h1:Hja7XhhmvEFhcByqDoHz9QZbkWey+COd9xWfCfn1ioo=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_mfridman_tparse",
        importpath = "github.com/mfridman/tparse",
        sum = "h1:wh6dzOKaIwkUGyKgOntDW4liXSo37qg5AXbIhkMV3vE=",
        version = "v0.18.0",
    )
    go_repository(
        name = "com_github_microcosm_cc_bluemonday",
        importpath = "github.com/microcosm-cc/bluemonday",
        sum = "h1:xbqSvqzQMeEHCqMi64VAs4d8uy6Mequs3rQ0k/Khz58=",
        version = "v1.0.26",
    )
    go_repository(
        name = "com_github_moby_docker_image_spec",
        importpath = "github.com/moby/docker-image-spec",
        sum = "h1:jMKff3w6PgbfSa69GfNg+zN/XLhfXJGnEx3Nl2EsFP0=",
        version = "v1.3.1",
    )
    go_repository(
        name = "com_github_moby_moby",
        importpath = "github.com/moby/moby",
        sum = "h1:KQbXBjo7PavKpzIl7UkHT31y9lw/e71Uvrqhr4X+zMA=",
        version = "v27.3.1+incompatible",
    )
    go_repository(
        name = "com_github_moby_moby_api",
        importpath = "github.com/moby/moby/api",
        sum = "h1:GQzua3NA599ASSIICx0iFgiJeO9YkdDARvQsm23ZZuQ=",
        version = "v1.56.0",
    )
    go_repository(
        name = "com_github_moby_moby_client",
        importpath = "github.com/moby/moby/client",
        sum = "h1:AJjEB21QPbXSXjDsZorFBoDZPhMrfbpaPLgSMAW9Bgs=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_github_oapi_codegen_nullable",
        importpath = "github.com/oapi-codegen/nullable",
        sum = "h1:eAh8JVc5430VtYVnq00Hrbpag9PFRGWLjxR1/3KntMs=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_oapi_codegen_runtime",
        importpath = "github.com/oapi-codegen/runtime",
        sum = "h1:t7358VYPvNbWJ9gdAkIK/smVeHpBf6yp8VTsaZsb/7k=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_github_oklog_ulid_v2",
        importpath = "github.com/oklog/ulid/v2",
        sum = "h1:IEclFb9JNvzYA6MW2SCxbLzcHTVsfqm3PrqGQJH5zec=",
        version = "v2.1.2",
    )
    go_repository(
        name = "com_github_openconfig_gnpsi",
        importpath = "github.com/openconfig/gnpsi",
        sum = "h1:+bl1bXMOTrWOcGydWB+8wGgvxlgvL8Y6joAiWFU5sog=",
        version = "v0.3.2",
    )
    go_repository(
        name = "com_github_openconfig_kne_third_party_meshnet",
        importpath = "github.com/openconfig/kne/third_party/meshnet",
        sum = "h1:D4BQpAdffLDMHTTSig6gsrH3yTJlTYyBP1DNO6Rczu0=",
        version = "v0.5.1",
    )
    go_repository(
        name = "com_github_pquerna_cachecontrol",
        importpath = "github.com/pquerna/cachecontrol",
        sum = "h1:yJMy84ti9h/+OEWa752kBTKv4XC30OtVVHYv/8cTqKc=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_prometheus_community_pro_bing",
        importpath = "github.com/prometheus-community/pro-bing",
        sum = "h1:CEY/g1/AgERRDjxw5P32ikcOgmrSuXs7xon7ovx6mNc=",
        version = "v0.8.0",
    )
    go_repository(
        name = "com_github_ravenox_go_jsoncommentstrip",
        importpath = "github.com/RaveNoX/go-jsoncommentstrip",
        sum = "h1:t527LHHE3HmiHrq74QMpNPZpGCIJzTx+apLkMKt4HC0=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_rcrowley_go_metrics",
        importpath = "github.com/rcrowley/go-metrics",
        sum = "h1:bsUq1dX0N8AOIL7EB/X911+m4EHsnWEHeJ0c+3TTBrg=",
        version = "v0.0.0-20250401214520-65e299d6c5c9",
    )
    go_repository(
        name = "com_github_schollz_closestmatch",
        importpath = "github.com/schollz/closestmatch",
        sum = "h1:Uel2GXEpJqOWBrlyI+oY9LTiyyjYS17cCYRqP13/SHk=",
        version = "v2.1.0+incompatible",
    )
    go_repository(
        name = "com_github_shopify_goreferrer",
        importpath = "github.com/Shopify/goreferrer",
        sum = "h1:KkH3I3sJuOLP3TjA/dfr4NAY8bghDwnXiU7cTKxQqo0=",
        version = "v0.0.0-20220729165902-8cddb4f5de06",
    )
    go_repository(
        name = "com_github_spiffe_go_spiffe_v2",
        importpath = "github.com/spiffe/go-spiffe/v2",
        sum = "h1:eXZMLsu+3MLEPJyGJkolqtVrteZfQdUpOWj6LTiDl/E=",
        version = "v2.8.1",
    )
    go_repository(
        name = "com_github_spkg_bom",
        importpath = "github.com/spkg/bom",
        sum = "h1:fiWzISvDn0Csy5H0iwgAuJGQTUpVfEMJJd4nRFXogbc=",
        version = "v0.0.0-20160624110644-59b7046e48ad",
    )
    go_repository(
        name = "com_github_tdewolff_minify_v2",
        importpath = "github.com/tdewolff/minify/v2",
        sum = "h1:tX0SR0LUrIqGoLjXnkIzRSIbKJ7PaNnSENLD4CyH6Xo=",
        version = "v2.20.19",
    )
    go_repository(
        name = "com_github_tdewolff_parse_v2",
        importpath = "github.com/tdewolff/parse/v2",
        sum = "h1:tgavkHc2ZDEQVKy1oWxwIyh5bP4F5fEh/JmBwPP/3LQ=",
        version = "v2.7.12",
    )
    go_repository(
        name = "com_github_twitchyliquid64_golang_asm",
        importpath = "github.com/twitchyliquid64/golang-asm",
        sum = "h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=",
        version = "v0.15.1",
    )
    go_repository(
        name = "com_github_ugorji_go_codec",
        importpath = "github.com/ugorji/go/codec",
        sum = "h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=",
        version = "v1.2.12",
    )
    go_repository(
        name = "com_github_valyala_bytebufferpool",
        importpath = "github.com/valyala/bytebufferpool",
        sum = "h1:GqA5TC/0021Y/b9FG4Oi9Mr3q7XYx6KllzawFIhcdPw=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_valyala_fasttemplate",
        importpath = "github.com/valyala/fasttemplate",
        sum = "h1:lxLXG0uE3Qnshl9QyaK6XJxMXlQZELvChBOCmQD0Loo=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_vmihailenco_msgpack_v5",
        importpath = "github.com/vmihailenco/msgpack/v5",
        sum = "h1:cQriyiUvjTwOHg8QZaPihLWeRAAVoCpE00IUPn0Bjt8=",
        version = "v5.4.1",
    )
    go_repository(
        name = "com_github_vmihailenco_tagparser_v2",
        importpath = "github.com/vmihailenco/tagparser/v2",
        sum = "h1:y09buUbR+b5aycVFQs/g70pqKVZNBmxwAhO7/IwNM9g=",
        version = "v2.0.0",
    )
    go_repository(
        name = "com_github_x448_float16",
        importpath = "github.com/x448/float16",
        sum = "h1:qLwI1I70+NjRFUR3zs1JPUCgaCXSh3SW62uAKT1mSBM=",
        version = "v0.8.4",
    )
    go_repository(
        name = "com_github_yosssi_ace",
        importpath = "github.com/yosssi/ace",
        sum = "h1:tUkIP/BLdKqrlrPwcmH0shwEEhTRHoGnc1wFIWmaBUA=",
        version = "v0.0.5",
    )
    go_repository(
        name = "com_github_zeebo_errs",
        importpath = "github.com/zeebo/errs",
        sum = "h1:XNdoD/RRMKP7HD0UhJnIzUy74ISdGGxURlYG8HSWSfM=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub_v2",
        importpath = "cloud.google.com/go/pubsub/v2",
        sum = "h1:MFrBTZZa6PDWZzCi4NJRsHKMm2w0a4oAaYNqwjgbQTE=",
        version = "v2.7.0",
    )
    go_repository(
        name = "in_gopkg_evanphx_json_patch_v4",
        importpath = "gopkg.in/evanphx/json-patch.v4",
        sum = "h1:czT3CmqEaQ1aanPc5SdlgQrrEIb8w/wwCvWWnfEbYzo=",
        version = "v4.13.0",
    )
    go_repository(
        name = "in_gopkg_go_jose_go_jose_v2",
        importpath = "gopkg.in/go-jose/go-jose.v2",
        sum = "h1:nt80fvSDlhKWQgSWyHyy5CfmlQr+asih51R8PTWNKKs=",
        version = "v2.6.3",
    )
    go_repository(
        name = "in_yaml_go_yaml_v2",
        importpath = "go.yaml.in/yaml/v2",
        sum = "h1:tuyd0P+2Ont/d6e2rl3be67goVK4R6deVxCUX5vyPaQ=",
        version = "v2.4.4",
    )
    go_repository(
        name = "in_yaml_go_yaml_v3",
        importpath = "go.yaml.in/yaml/v3",
        sum = "h1:N6y/pJk8buWs9NY5ERU2HSMfm+IuD/OtfdAnq6kESPw=",
        version = "v3.0.5",
    )
    go_repository(
        name = "io_etcd_go_raft_v3",
        importpath = "go.etcd.io/raft/v3",
        sum = "h1:BGzlwx07bLv8PW6OU5HObuz1y4hlPZUXA07pM1mPUh4=",
        version = "v3.7.0",
    )
    go_repository(
        name = "io_k8s_cri_client",
        importpath = "k8s.io/cri-client",
        sum = "h1:9ZwddaNJomqkTBYQqSmB+Ccns3beY4HyYDwmRtWTCJM=",
        version = "v0.31.3",
    )
    go_repository(
        name = "io_k8s_sigs_randfill",
        importpath = "sigs.k8s.io/randfill",
        sum = "h1:JfjMILfT8A6RbawdsK2JXGBR5AQVfd+9TbzrlneTyrU=",
        version = "v1.0.0",
    )
    go_repository(
        name = "io_k8s_sigs_structured_merge_diff_v6",
        importpath = "sigs.k8s.io/structured-merge-diff/v6",
        sum = "h1:qdOxHwrl2Kaag1aQEarlYcOA9vSyGCp3CIki3aW8c4Q=",
        version = "v6.4.2",
    )
    go_repository(
        name = "io_k8s_streaming",
        importpath = "k8s.io/streaming",
        sum = "h1:iPBUZLZiKt5bV+lxJurASMOV07VuBhNpiwJt2//AWrM=",
        version = "v0.37.0",
    )
    go_repository(
        name = "io_opencensus_go_contrib_exporter_stackdriver",
        importpath = "contrib.go.opencensus.io/exporter/stackdriver",
        sum = "h1:xRc46S76eyn4ZF3jWX8I+aUSKVLw5EQ1aDvHwfV5W1o=",
        version = "v0.13.15-0.20230702191903-2de6d2748484",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace_otlptracehttp",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp",
        sum = "h1:KrC1YrQeSt46ITMWAbgQx1M1eV1/1TKzttrBzymPmss=",
        version = "v1.46.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_prometheus",
        importpath = "go.opentelemetry.io/otel/exporters/prometheus",
        sum = "h1:AHh/lAP1BHrY5gBwk8ncc25FXWm/gmmY3BX258z5nuk=",
        version = "v0.57.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_stdout_stdoutmetric",
        importpath = "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric",
        sum = "h1:TC+BewnDpeiAmcscXbGMfxkO+mwYUwE/VySwvw88PfA=",
        version = "v1.43.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_stdout_stdouttrace",
        importpath = "go.opentelemetry.io/otel/exporters/stdout/stdouttrace",
        sum = "h1:KdRxPiAoMptR3vfWzvjjvutTsSiwbC2uG0496rzZNfo=",
        version = "v1.46.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric_x",
        importpath = "go.opentelemetry.io/otel/metric/x",
        sum = "h1:TA/cBT23D3MnxYPwHL7YFOdYGdx0A0v+s7Mzotpd1dU=",
        version = "v0.68.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_log_logtest",
        importpath = "go.opentelemetry.io/otel/sdk/log/logtest",
        sum = "h1:infPnfNrhCNgOUZRs3gWUg8vhoBUHihq02gwK05gzlg=",
        version = "v0.22.0",
    )
    go_repository(
        name = "org_codeberg_go_fonts_liberation",
        importpath = "codeberg.org/go-fonts/liberation",
        sum = "h1:SsKoMO1v1OZmzkG2DY+7ZkCL9U+rrWI09niOLfQ5Bo0=",
        version = "v0.5.0",
    )
    go_repository(
        name = "org_codeberg_go_latex_latex",
        importpath = "codeberg.org/go-latex/latex",
        sum = "h1:hoGO86rIbWVyjtlDLzCqZPjNykpWQ9YuTZqAzPcfL3c=",
        version = "v0.1.0",
    )
    go_repository(
        name = "org_codeberg_go_pdf_fpdf",
        importpath = "codeberg.org/go-pdf/fpdf",
        sum = "h1:u+w669foDDx5Ds43mpiiayp40Ov6sZalgcPMDBcZRd4=",
        version = "v0.10.0",
    )
    go_repository(
        name = "org_golang_google_grpc_examples",
        importpath = "google.golang.org/grpc/examples",
        sum = "h1:ExN12ndbJ608cboPYflpTny6mXSzPrDLh0iTaVrRrds=",
        version = "v0.0.0-20250407062114-b368379ef8f6",
    )
    go_repository(
        name = "org_golang_google_grpc_gcp_observability",
        importpath = "google.golang.org/grpc/gcp/observability",
        sum = "h1:2IQ7szW1gobfZaS/sDSAu2uxO0V/aTryMZvlcyqKqQA=",
        version = "v1.0.1",
    )
    go_repository(
        name = "org_golang_google_grpc_security_advancedtls",
        importpath = "google.golang.org/grpc/security/advancedtls",
        sum = "h1:/KQ7VP/1bs53/aopk9QhuPyFAp9Dm9Ejix3lzYkCrDA=",
        version = "v1.0.0",
    )
    go_repository(
        name = "org_golang_google_grpc_stats_opencensus",
        importpath = "google.golang.org/grpc/stats/opencensus",
        sum = "h1:evSYcRZaSToQp+borzWE52+03joezZeXcKJvZDfkUJA=",
        version = "v1.0.0",
    )
    go_repository(
        name = "org_golang_x_arch",
        importpath = "golang.org/x/arch",
        sum = "h1:3wRIsP3pM4yUptoR96otTUOXI367OS0+c9eeRi9doIc=",
        version = "v0.8.0",
    )
    go_repository(
        name = "org_golang_x_tools_go_expect",
        importpath = "golang.org/x/tools/go/expect",
        sum = "h1:jpBZDwmgPhXsKZC6WhL20P4b/wmnpsEAGHaNy0n/rJM=",
        version = "v0.1.1-deprecated",
    )
    go_repository(
        name = "org_golang_x_tools_go_packages_packagestest",
        importpath = "golang.org/x/tools/go/packages/packagestest",
        sum = "h1:1h2MnaIAIXISqTFKdENegdpAgUXz6NrPEsbIeWaBRvM=",
        version = "v0.1.1-deprecated",
    )
    go_repository(
        name = "org_gonum_v1_tools",
        importpath = "gonum.org/v1/tools",
        sum = "h1:cJWOvXtcaFSGXz2F4z2AMM0VV7edDDGrxb5GLQH7ayQ=",
        version = "v0.0.0-20200318103217-c168b003ce8c",
    )
    go_repository(
        name = "org_modernc_sortutil",
        importpath = "modernc.org/sortutil",
        sum = "h1:+xyoGf15mM3NMlPDnFqrteY07klSFxLElE2PVuWIJ7w=",
        version = "v1.2.1",
    )
    go_repository(
        name = "org_uber_go_automaxprocs",
        importpath = "go.uber.org/automaxprocs",
        sum = "h1:O3y2/QNTOdbF+e/dpXNNW7Rx2hZ4sTIPyybbxyNqTUs=",
        version = "v1.6.0",
    )

def go_dependencies():
    go_repository(
        name = "com_github_alecthomas_kingpin_v2",
        importpath = "github.com/alecthomas/kingpin/v2",
        sum = "h1:f48lwail6p8zpO1bC4TxtqACaGqHYA22qkHjHpqDjYY=",
        version = "v2.4.0",
    )
    go_repository(
        name = "com_github_alecthomas_units",
        importpath = "github.com/alecthomas/units",
        sum = "h1:mimo19zliBX/vSQ6PWWSL9lK8qwHozUj03+zLoEB8O0=",
        version = "v0.0.0-20240927000941-0f3dac36c52b",
    )
    go_repository(
        name = "com_github_azure_go_ansiterm",
        importpath = "github.com/Azure/go-ansiterm",
        sum = "h1:udKWzYgxTojEKWjV8V+WSxDXJ4NFATAsZjh8iIbsQIg=",
        version = "v0.0.0-20250102033503-faa5f7b0171c",
    )
    go_repository(
        name = "com_github_jpillora_backoff",
        importpath = "github.com/jpillora/backoff",
        sum = "h1:uvFg412JmmHBHw7iwprIxkPMI+sGQ4kzOWsMeHnm2EA=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_julienschmidt_httprouter",
        importpath = "github.com/julienschmidt/httprouter",
        sum = "h1:U0609e9tgbseu3rBINet9P48AI/D3oJs4dN7jwJOQ1U=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_moby_term",
        importpath = "github.com/moby/term",
        sum = "h1:6qk3FJAFDs6i/q3W/pQ97SX192qKfZgGjCQqfCJkgzQ=",
        version = "v0.5.2",
    )
    go_repository(
        name = "com_github_mwitkow_go_conntrack",
        importpath = "github.com/mwitkow/go-conntrack",
        sum = "h1:KUppIJq7/+SVif2QVs3tOP0zanoHgBEVAwHxUSIzRqU=",
        version = "v0.0.0-20190716064945-2f068394615f",
    )
    go_repository(
        name = "com_github_xhit_go_str2duration_v2",
        importpath = "github.com/xhit/go-str2duration/v2",
        sum = "h1:lxklc02Drh6ynqX+DdPyp5pCKLUQpRT8bp8Ydu2Bstc=",
        version = "v2.1.0",
    )
    go_repository(
        name = "org_golang_x_telemetry",
        importpath = "golang.org/x/telemetry",
        sum = "h1:F5BWKvW126NXR74uxkxuc1jQHhm/rwm/J3rSiFyuRs4=",
        version = "v0.0.0-20260908163034-4bcc4b2ee518",
    )
