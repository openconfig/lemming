load("@bazel_gazelle//:def.bzl", "go_repository")

def go_repositories():
    go_repository(
        name = "co_honnef_go_tools",
        importpath = "honnef.co/go/tools",
        sum = "h1:UoveltGrhghAA7ePc+e+QYDHXrBps2PqFZiHkGR/xK8=",
        version = "v0.0.1-2020.1.4",
    )
    go_repository(
        name = "com_github_antihax_optional",
        importpath = "github.com/antihax/optional",
        sum = "h1:xK2lYat7ZLaVVcIuj82J8kIro4V6kDe0AUDFboUCwcg=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_antlr_antlr4_runtime_go_antlr",
        importpath = "github.com/antlr/antlr4/runtime/Go/antlr",
        sum = "h1:yL7+Jz0jTC6yykIK/Wh74gnTJnrGr5AyrNMXuA0gves=",
        version = "v1.4.10",
    )
    go_repository(
        name = "com_github_aristanetworks_arista_ceoslab_operator_v2",
        importpath = "github.com/aristanetworks/arista-ceoslab-operator/v2",
        sum = "h1:lOGhvrrxc/bzUMJI4kYmZh/N1a87Vnsk/TENZt5kCoc=",
        version = "v2.0.1",
    )
    go_repository(
        name = "com_github_armon_go_metrics",
        importpath = "github.com/armon/go-metrics",
        sum = "h1:yCQqn7dwca4ITXb+CbubHmedzaQYHhNhrEXLYUeEe8Q=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_armon_go_socks5",
        importpath = "github.com/armon/go-socks5",
        sum = "h1:0CwZNZbxp69SHPdPJAN/hZIm0C4OItdklCFmMRWYpio=",
        version = "v0.0.0-20160902184237-e75332964ef5",
    )
    go_repository(
        name = "com_github_asaskevich_govalidator",
        importpath = "github.com/asaskevich/govalidator",
        sum = "h1:idn718Q4B6AGu/h5Sxe66HYVdqdGu2l9Iebqhi/AEoA=",
        version = "v0.0.0-20190424111038-f61b66f89f4a",
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
        sum = "h1:F3R3q42aWytozkV8ihzcgMO4OA4cuqr3bNlsEuF6//A=",
        version = "v0.11.27",
    )
    go_repository(
        name = "com_github_azure_go_autorest_autorest_adal",
        importpath = "github.com/Azure/go-autorest/autorest/adal",
        sum = "h1:5CHCGJjd9AedaD0iHhMTCRJPypw+hyCrtT9sLk/dD0w=",
        version = "v0.9.19",
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
        name = "com_github_blang_semver_v4",
        importpath = "github.com/blang/semver/v4",
        sum = "h1:1PFHFE6yCCTv8C1TeyNNarDzntLi7wMI5i/pzqYIsAM=",
        version = "v4.0.0",
    )
    go_repository(
        name = "com_github_buger_jsonparser",
        importpath = "github.com/buger/jsonparser",
        sum = "h1:2PnMjfWD7wBILjqQbt530v576A/cAbQvEW9gGIpYMUs=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_burntsushi_toml",
        importpath = "github.com/BurntSushi/toml",
        sum = "h1:9F2/+DoOYIOksmaJFPw1tGFy1eDnIJXg+UHjuD8lTak=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_burntsushi_xgb",
        importpath = "github.com/BurntSushi/xgb",
        sum = "h1:1BDTz0u9nC3//pOCMdNH+CiXJVYJh5UQNCOBG7jbELc=",
        version = "v0.0.0-20160522181843-27f122750802",
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
        sum = "h1:cFAlzYUlVYDysBEH2T5hyJZMh3+5+WCBvSnK6Q8UtC4=",
        version = "v4.1.3",
    )
    go_repository(
        name = "com_github_census_instrumentation_opencensus_proto",
        importpath = "github.com/census-instrumentation/opencensus-proto",
        sum = "h1:iKLQ0xPNFxR/2hzXZMrBo8f1j86j5WHzznCCQxV/b8g=",
        version = "v0.4.1",
    )
    go_repository(
        name = "com_github_cespare_xxhash",
        importpath = "github.com/cespare/xxhash",
        sum = "h1:a6HrQnmkObjyL+Gs60czilIUGqrzKutQD6XZog3p+ko=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_cespare_xxhash_v2",
        importpath = "github.com/cespare/xxhash/v2",
        sum = "h1:DC2CZ1Ep5Y4k3ZQ899DldepgrayRUGE6BBZ/cd9Cj44=",
        version = "v2.2.0",
    )
    go_repository(
        name = "com_github_chzyer_logex",
        importpath = "github.com/chzyer/logex",
        sum = "h1:Swpa1K6QvQznwJRcfTfQJmTE72DqScAa40E+fbHEXEE=",
        version = "v1.1.10",
    )
    go_repository(
        name = "com_github_chzyer_readline",
        importpath = "github.com/chzyer/readline",
        sum = "h1:fY5BOSpyZCqRo5OhCuC+XN+r/bBCmeuuJtjz+bCNIf8=",
        version = "v0.0.0-20180603132655-2972be24d48e",
    )
    go_repository(
        name = "com_github_chzyer_test",
        importpath = "github.com/chzyer/test",
        sum = "h1:q763qf9huN11kDQavWsoZXJNW3xEE4JJyHa5Q25/sd8=",
        version = "v0.0.0-20180213035817-a1ea475d72b1",
    )
    go_repository(
        name = "com_github_client9_misspell",
        importpath = "github.com/client9/misspell",
        sum = "h1:ta993UF76GwbvJcIo3Y68y/M3WxlpEHPWIGDkJYwzJI=",
        version = "v0.3.4",
    )
    go_repository(
        name = "com_github_cncf_udpa_go",
        importpath = "github.com/cncf/udpa/go",
        sum = "h1:QQ3GSy+MqSHxm/d8nCtnAiZdYFd45cYZPs8vOOIYKfk=",
        version = "v0.0.0-20220112060539-c52dc94e7fbe",
    )
    go_repository(
        name = "com_github_cncf_xds_go",
        importpath = "github.com/cncf/xds/go",
        sum = "h1:ACGZRIr7HsgBKHsueQ1yM4WaVaXh21ynwqsF8M8tXhA=",
        version = "v0.0.0-20230105202645-06c439db220b",
    )
    go_repository(
        name = "com_github_containerd_containerd",
        importpath = "github.com/containerd/containerd",
        sum = "h1:doHr6cNxfOLTotWmZs6aZF6LrfJFcjmYFcWlRmQgYPM=",
        version = "v1.5.18",
    )
    go_repository(
        name = "com_github_coreos_go_semver",
        importpath = "github.com/coreos/go-semver",
        sum = "h1:wkHLiw0WNATZnSG7epLsujiMCgPAc9xhjJ4tgnAxmfM=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_coreos_go_systemd_v22",
        importpath = "github.com/coreos/go-systemd/v22",
        sum = "h1:RrqgGjYQKalulkV8NGVIfkXQf6YYmOyiJKk8iXXhfZs=",
        version = "v22.5.0",
    )
    go_repository(
        name = "com_github_cpuguy83_go_md2man_v2",
        importpath = "github.com/cpuguy83/go-md2man/v2",
        sum = "h1:p1EgwI/C7NhT0JmVkwCD2ZBK8j4aeHQX2pMHHBfMQ6w=",
        version = "v2.0.2",
    )
    go_repository(
        name = "com_github_creack_pty",
        importpath = "github.com/creack/pty",
        sum = "h1:n56/Zwd5o6whRC5PMGretI4IdRLlmBXYNjScPaBgsbY=",
        version = "v1.1.18",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_derekparker_trie",
        importpath = "github.com/derekparker/trie",
        sum = "h1:zPGs5GlTifgPVykTSvWB6/+J7EefOmYqH/JuUBI3kf0=",
        version = "v0.0.0-20221221181808-1424fce0c981",
    )
    go_repository(
        name = "com_github_dgryski_go_farm",
        importpath = "github.com/dgryski/go-farm",
        sum = "h1:fAjc9m62+UWV/WAFKLNi6ZS0675eEUC9y3AlwSbQu1Y=",
        version = "v0.0.0-20200201041132-a6ae2369ad13",
    )
    go_repository(
        name = "com_github_docker_distribution",
        importpath = "github.com/docker/distribution",
        sum = "h1:Q50tZOPR6T/hjNsyc9g8/syEs6bk8XXApsHjKukMl68=",
        version = "v2.8.1+incompatible",
    )
    go_repository(
        name = "com_github_docker_docker",
        importpath = "github.com/docker/docker",
        sum = "h1:CEeNmFM0QZIsJCZKMkZx0ZcahTiewkrgiwfYD+dfl1U=",
        version = "v20.10.12+incompatible",
    )
    go_repository(
        name = "com_github_docker_go_connections",
        importpath = "github.com/docker/go-connections",
        sum = "h1:El9xVISelRB7BuFusrZozjnkIM5YnzCViNKohAFqRJQ=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_docker_go_units",
        importpath = "github.com/docker/go-units",
        sum = "h1:3uh0PgVws3nIA0Q+MwDC8yjEPf9zjRfZZWXZYDct3Tw=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_docopt_docopt_go",
        importpath = "github.com/docopt/docopt-go",
        sum = "h1:bWDMxwH3px2JBh6AyO7hdCn/PkvCZXii8TGj7sbtEbQ=",
        version = "v0.0.0-20180111231733-ee0de3bc6815",
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
        name = "com_github_elazarl_goproxy",
        importpath = "github.com/elazarl/goproxy",
        sum = "h1:yUdfgN0XgIJw7foRItutHYUIhlcKzcSf5vDpdhQAKTc=",
        version = "v0.0.0-20180725130230-947c36da3153",
    )
    go_repository(
        name = "com_github_emicklei_go_restful",
        importpath = "github.com/emicklei/go-restful",
        sum = "h1:rgqiKNjTnFQA6kkhFe16D8epTksy9HQ1MyrbDXSdYhM=",
        version = "v2.16.0+incompatible",
    )
    go_repository(
        name = "com_github_emicklei_go_restful_v3",
        importpath = "github.com/emicklei/go-restful/v3",
        sum = "h1:hIovbnmBTLjHXkqEBUz3HGpXZdM7ZrE9fJIZIqlJLqE=",
        version = "v3.10.2",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane",
        importpath = "github.com/envoyproxy/go-control-plane",
        sum = "h1:xdCVXxEe0Y3FQith+0cj2irwZudqGYvecuLB1HtdexY=",
        version = "v0.10.3",
    )
    go_repository(
        name = "com_github_envoyproxy_protoc_gen_validate",
        importpath = "github.com/envoyproxy/protoc-gen-validate",
        sum = "h1:PS7VIOgmSVhWUEeZwTe7z7zouA22Cr590PzXKbZHOVY=",
        version = "v0.9.1",
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
        sum = "h1:b91NhWfaz02IuVxO9faSllyAtNXHMPkC5J8sJCLunww=",
        version = "v5.6.0",
    )
    go_repository(
        name = "com_github_fatih_color",
        importpath = "github.com/fatih/color",
        sum = "h1:8LOYc1KYPPmyKMuN8QV2DNRWNbLo6LZ0iLs8+mlH53w=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_github_felixge_httpsnoop",
        importpath = "github.com/felixge/httpsnoop",
        sum = "h1:s/nj+GCswXYzN5v2DpNMuMQYe+0DDwt5WVCU6CWBdXk=",
        version = "v1.0.3",
    )
    go_repository(
        name = "com_github_flowstack_go_jsonschema",
        importpath = "github.com/flowstack/go-jsonschema",
        sum = "h1:dCrjGJRXIlbDsLAgTJZTjhwUJnnxVWl1OgNyYh5nyDc=",
        version = "v0.1.1",
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
        sum = "h1:FJKSZTDHjyhriyC81FLQ0LY93eSai0ZyR/ZIkd3ZUKE=",
        version = "v1.14.3",
    )
    go_repository(
        name = "com_github_fsnotify_fsnotify",
        importpath = "github.com/fsnotify/fsnotify",
        sum = "h1:n+5WquG0fcWoWp6xPWfHdbskMCQaFnG6PfBrh1Ky4HY=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_ghodss_yaml",
        importpath = "github.com/ghodss/yaml",
        sum = "h1:wQHKEahhL6wmXdzwWG11gIVCkOv05bNOh+Rxn0yngAk=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_go_gl_glfw",
        importpath = "github.com/go-gl/glfw",
        sum = "h1:QbL/5oDUmRBzO9/Z7Seo6zf912W/a6Sr4Eu0G/3Jho0=",
        version = "v0.0.0-20190409004039-e6da0acd62b1",
    )
    go_repository(
        name = "com_github_go_gl_glfw_v3_3_glfw",
        importpath = "github.com/go-gl/glfw/v3.3/glfw",
        sum = "h1:WtGNWLvXpe6ZudgnXrq0barxBImvnnJoMEhXAzcbM0I=",
        version = "v0.0.0-20200222043503-6f7a984d4dc4",
    )
    go_repository(
        name = "com_github_go_kit_log",
        importpath = "github.com/go-kit/log",
        sum = "h1:MRVx0/zhvdseW+Gza6N9rVzU/IVzaeE1SFI4raAhmBU=",
        version = "v0.2.1",
    )
    go_repository(
        name = "com_github_go_logfmt_logfmt",
        importpath = "github.com/go-logfmt/logfmt",
        sum = "h1:otpy5pqBCBZ1ng9RQ0dPu4PN7ba75Y/aA+UpowDyNVA=",
        version = "v0.5.1",
    )
    go_repository(
        name = "com_github_go_logr_logr",
        importpath = "github.com/go-logr/logr",
        sum = "h1:2DntVwHkVopvECVRSlL5PSo9eG+cAkDCuckLubN+rq0=",
        version = "v1.2.3",
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
        sum = "h1:a9vnzlIBPQBBkeaR9IuMUfmVOrQlkoC4YfPoFkX3T7A=",
        version = "v1.2.3",
    )
    go_repository(
        name = "com_github_go_openapi_jsonpointer",
        importpath = "github.com/go-openapi/jsonpointer",
        sum = "h1:eCs3fxoIi3Wh6vtgmLTOjdhSpiqphQ+DaPn38N2ZdrE=",
        version = "v0.19.6",
    )
    go_repository(
        name = "com_github_go_openapi_jsonreference",
        importpath = "github.com/go-openapi/jsonreference",
        sum = "h1:3sVjiK66+uXK/6oQ8xgcRKcFgQ5KXa2KvnJRumpMGbE=",
        version = "v0.20.2",
    )
    go_repository(
        name = "com_github_go_openapi_swag",
        importpath = "github.com/go-openapi/swag",
        sum = "h1:yMBqmnQ0gyZvEb/+KzuWZOXgllrXT4SADYbvDaXHv/g=",
        version = "v0.22.3",
    )
    go_repository(
        name = "com_github_go_test_deep",
        importpath = "github.com/go-test/deep",
        sum = "h1:WOcxcdHcvdgThNXjw0t76K42FXTU7HpNQWHpA2HHNlg=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_gogo_protobuf",
        importpath = "github.com/gogo/protobuf",
        sum = "h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=",
        version = "v1.3.2",
    )
    go_repository(
        name = "com_github_golang_glog",
        importpath = "github.com/golang/glog",
        sum = "h1:/d3pCKDPWNnvIWe0vVUpNP32qc8U3PDVxySP/y360qE=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_golang_groupcache",
        importpath = "github.com/golang/groupcache",
        sum = "h1:oI5xCqsCo564l8iNU+DwB5epxmsaqB+rhGL0m5jtYqE=",
        version = "v0.0.0-20210331224755-41bb18bfe9da",
    )
    go_repository(
        name = "com_github_golang_jwt_jwt_v4",
        importpath = "github.com/golang-jwt/jwt/v4",
        sum = "h1:pC5DB52sCeK48Wlb9oPcdhnjkz1TKt1D/P7WKJ0kUcQ=",
        version = "v4.4.1",
    )
    go_repository(
        name = "com_github_golang_mock",
        importpath = "github.com/golang/mock",
        sum = "h1:ErTB+efbowRARo13NNdxyJji2egdxLGQhRaY+DUumQc=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=",
        version = "v1.5.3",
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
        sum = "h1:gK4Kx5IaGY9CD5sPJ36FHiBJ6ZXl0kilRiiCj+jdYp4=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_google_cel_go",
        importpath = "github.com/google/cel-go",
        sum = "h1:kjeKudqV0OygrAqA9fX6J55S8gj+Jre2tckIm5RoG4M=",
        version = "v0.12.6",
    )
    go_repository(
        name = "com_github_google_gnostic",
        build_file_proto_mode = "disable",
        importpath = "github.com/google/gnostic",
        sum = "h1:ZK/5VhkoX835RikCHpSUJV9a+S3e1zLh59YnyWeBW+0=",
        version = "v0.6.9",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:O2Tfq5qg4qc4AmwVlvv0oLiVAGB7enBSJ2x2DqQFi38=",
        version = "v0.5.9",
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
        name = "com_github_google_martian",
        importpath = "github.com/google/martian",
        sum = "h1:/CP5g8u/VJHijgedC/Legn3BAbAaWPgecwXBIDzw5no=",
        version = "v2.1.0+incompatible",
    )
    go_repository(
        name = "com_github_google_martian_v3",
        importpath = "github.com/google/martian/v3",
        sum = "h1:IqNFLAmvJOgVlpdEBiQbDc2EwKW77amAycfTuWKdfvw=",
        version = "v3.3.2",
    )
    go_repository(
        name = "com_github_google_pprof",
        importpath = "github.com/google/pprof",
        sum = "h1:K6RDEckDVWvDI9JAJYCmNdQXq6neHJOYx3V6jnqNEec=",
        version = "v0.0.0-20210720184732-4bb14d4b1be1",
    )
    go_repository(
        name = "com_github_google_protobuf",
        importpath = "github.com/google/protobuf",
        sum = "h1:J+nrmmGmKWdQobsM0Vom92fe/UGLowFjFl9jX+Oe/S4=",
        version = "v3.11.4+incompatible",
    )
    go_repository(
        name = "com_github_google_renameio",
        importpath = "github.com/google/renameio",
        sum = "h1:GOZbcHa3HfsPKPlmyPyN2KEohoMXOhdMbHrvbpl2QaA=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_google_uuid",
        importpath = "github.com/google/uuid",
        sum = "h1:t6JiXgmwXMjEs8VusXIJk2BXHsn+wx8BZdTaoZ5fu7I=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_googleapis_enterprise_certificate_proxy",
        importpath = "github.com/googleapis/enterprise-certificate-proxy",
        sum = "h1:yk9/cqRKtT9wXZSsRH9aurXEpJX+U6FLtpYTdC3R06k=",
        version = "v0.2.3",
    )
    go_repository(
        name = "com_github_googleapis_gax_go_v2",
        importpath = "github.com/googleapis/gax-go/v2",
        sum = "h1:IcsPKeInNvYi7eqSaDjiZqDDKu5rsmunY0Y1YupQSSQ=",
        version = "v2.7.0",
    )
    go_repository(
        name = "com_github_googleapis_gnostic",
        importpath = "github.com/googleapis/gnostic",
        sum = "h1:9fHAtK0uDfpveeqqo1hkEZJcFvYXAiCN3UutL8F9xHw=",
        version = "v0.5.5",
    )
    go_repository(
        name = "com_github_googleapis_go_type_adapters",
        importpath = "github.com/googleapis/go-type-adapters",
        sum = "h1:9XdMn+d/G57qq1s8dNc5IesGCXHf6V2HZ2JwRxfA2tA=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_googleapis_google_cloud_go_testing",
        importpath = "github.com/googleapis/google-cloud-go-testing",
        sum = "h1:tlyzajkF3030q6M8SvmJSemC9DTHL/xaMa18b65+JM4=",
        version = "v0.0.0-20200911160855-bcd43fbb19e8",
    )
    go_repository(
        name = "com_github_gorilla_mux",
        importpath = "github.com/gorilla/mux",
        sum = "h1:i40aqfkR1h2SlN9hojwV5ZA91wcXFOvkdNIeFDP5koI=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_github_gorilla_websocket",
        importpath = "github.com/gorilla/websocket",
        sum = "h1:+/TMaTYc4QFitKJxsQ7Yye35DkWvkdLcvGKqM+x0Ufc=",
        version = "v1.4.2",
    )
    go_repository(
        name = "com_github_gregjones_httpcache",
        importpath = "github.com/gregjones/httpcache",
        sum = "h1:pdN6V1QBWetyv/0+wjACpqVH+eVULgEjkurDLq3goeM=",
        version = "v0.0.0-20180305231024-9cad4c3443a7",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_go_grpc_middleware",
        importpath = "github.com/grpc-ecosystem/go-grpc-middleware",
        sum = "h1:+9834+KizmvFV7pXQGSXQTsaWhq2GjuNUt0aUU0YBYw=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_go_grpc_prometheus",
        importpath = "github.com/grpc-ecosystem/go-grpc-prometheus",
        sum = "h1:Ovs26xHkKqVztRpIrF/92BcuyuQ/YW4NSIpoGtfXNho=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_grpc_gateway",
        importpath = "github.com/grpc-ecosystem/grpc-gateway",
        sum = "h1:gmcG1KaJ57LophUzW0Hy8NmPhnMZb4M0+kPpLofRdBo=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_grpc_gateway_v2",
        importpath = "github.com/grpc-ecosystem/grpc-gateway/v2",
        sum = "h1:lLT7ZLSzGLI08vc9cpd+tYmNWjdKDqyr/2L+f6U12Fk=",
        version = "v2.11.3",
    )
    go_repository(
        name = "com_github_h_fam_errdiff",
        importpath = "github.com/h-fam/errdiff",
        sum = "h1:rPsW4ob2fMOIulwTEoZXaaUIuud7XUudw5SLKTZj3Ss=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_hashicorp_consul_api",
        importpath = "github.com/hashicorp/consul/api",
        sum = "h1:R7PPNzTCeN6VuQNDwwhZWJvzCtGSrNpJqfb22h3yH9g=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_cleanhttp",
        importpath = "github.com/hashicorp/go-cleanhttp",
        sum = "h1:035FKYIWjmULyFRBKPs8TBQoi0x6d9G4xc9neXJWAZQ=",
        version = "v0.5.2",
    )
    go_repository(
        name = "com_github_hashicorp_go_hclog",
        importpath = "github.com/hashicorp/go-hclog",
        sum = "h1:La19f8d7WIlm4ogzNHB0JGqs5AUDAZ2UfCY4sJXcJdM=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_immutable_radix",
        importpath = "github.com/hashicorp/go-immutable-radix",
        sum = "h1:DKHmCUm2hRBK510BaiZlwvpD40f8bJFeZnpfm2KLowc=",
        version = "v1.3.1",
    )
    go_repository(
        name = "com_github_hashicorp_go_rootcerts",
        importpath = "github.com/hashicorp/go-rootcerts",
        sum = "h1:jzhAVGtqPKbwpyCPELlgNWhE1znq+qwJtW5Oi2viEzc=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_hashicorp_go_version",
        importpath = "github.com/hashicorp/go-version",
        sum = "h1:aAQzgqIrRKRa7w75CKpbBxYsmUoPjzVm1W59ca1L0J4=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_hashicorp_golang_lru",
        importpath = "github.com/hashicorp/golang-lru",
        sum = "h1:YDjusn29QI/Das2iO9M0BHnIbxPeyuCHsjMW+lJfyTc=",
        version = "v0.5.4",
    )
    go_repository(
        name = "com_github_hashicorp_hcl",
        importpath = "github.com/hashicorp/hcl",
        sum = "h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_hashicorp_serf",
        importpath = "github.com/hashicorp/serf",
        sum = "h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=",
        version = "v0.10.1",
    )
    go_repository(
        name = "com_github_iancoleman_strcase",
        importpath = "github.com/iancoleman/strcase",
        sum = "h1:05I4QRnGpI0m37iZQRuskXh+w77mr6Z41lwQzuHLwW0=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_ianlancetaylor_demangle",
        importpath = "github.com/ianlancetaylor/demangle",
        sum = "h1:mV02weKRL81bEnm8A0HT1/CAelMQDBuQIfLw8n+d6xI=",
        version = "v0.0.0-20200824232613-28f6c0f3b639",
    )
    go_repository(
        name = "com_github_imdario_mergo",
        importpath = "github.com/imdario/mergo",
        sum = "h1:fOqeC1+nCuuk6PKQdg9YmosXX7Y7mHX6R/0ZldI9iHo=",
        version = "v0.3.14",
    )
    go_repository(
        name = "com_github_inconshreveable_mousetrap",
        importpath = "github.com/inconshreveable/mousetrap",
        sum = "h1:wN+x4NVGpMsO7ErUn/mUI3vEoE6Jt13X2s0bqwp9tc8=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_jessevdk_go_flags",
        importpath = "github.com/jessevdk/go-flags",
        sum = "h1:1jKYvbxEjfUl0fmqTCOfonvskHHXMjBySTLW4y9LFvc=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_jonboulle_clockwork",
        importpath = "github.com/jonboulle/clockwork",
        sum = "h1:UOGuzwb1PwsrDAObMuhUnj0p5ULPj8V/xJ7Kx9qUBdQ=",
        version = "v0.2.2",
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
        name = "com_github_jstemmer_go_junit_report",
        importpath = "github.com/jstemmer/go-junit-report",
        sum = "h1:6QPYqodiu3GuPL+7mfx+NwDdp2eTkp9IfEUpgAwUN0o=",
        version = "v0.9.1",
    )
    go_repository(
        name = "com_github_jstemmer_go_junit_report_v2",
        importpath = "github.com/jstemmer/go-junit-report/v2",
        sum = "h1:BVBb1o0TfOuRCMykVAYJ1r2yoZ+ByE0f19QNF4ngQ0M=",
        version = "v2.0.1-0.20220823220451-7b10b4285462",
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
        sum = "h1:WZcp8V8GQhsya0bMZuXktEH/Wz+aBlhiMle4tExkj6M=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_kisielk_errcheck",
        importpath = "github.com/kisielk/errcheck",
        sum = "h1:e8esj/e4R+SAOwFwN+n3zr0nYeCyeweozKfO23MvHzY=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_kisielk_gotool",
        importpath = "github.com/kisielk/gotool",
        sum = "h1:AV2c/EiW3KqPNT9ZKl07ehoAGi4C5/01Cfbblndcapg=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_kr_fs",
        importpath = "github.com/kr/fs",
        sum = "h1:Jskdu9ieNAYnjxsi0LbQp1ulIKZV1LAFgK1tWhpZgl8=",
        version = "v0.1.0",
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
        name = "com_github_lyft_protoc_gen_star",
        importpath = "github.com/lyft/protoc-gen-star",
        sum = "h1:erE0rdztuaDq3bpGifD95wfoPrSZc95nGA6tbiNYh6M=",
        version = "v0.6.1",
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
        sum = "h1:3MEsd0SM6jqZojhjLWWeBY+Kcjy9i6MQAeY7YgDP83g=",
        version = "v3.2.0",
    )
    go_repository(
        name = "com_github_mattn_go_colorable",
        importpath = "github.com/mattn/go-colorable",
        sum = "h1:jF+Du6AlPIjs2BiUiQlKOX0rt3SujHxPnksPKZbaA40=",
        version = "v0.1.12",
    )
    go_repository(
        name = "com_github_mattn_go_isatty",
        importpath = "github.com/mattn/go-isatty",
        sum = "h1:yVuAays6BHfxijgZPzw+3Zlu5yQgKGP2/hcQbHb7S9Y=",
        version = "v0.0.14",
    )
    go_repository(
        name = "com_github_matttproud_golang_protobuf_extensions",
        importpath = "github.com/matttproud/golang_protobuf_extensions",
        sum = "h1:mmDVorXM7PCGKw94cs5zkfA9PSy5pEvNWRP0ET0TIVo=",
        version = "v1.0.4",
    )
    go_repository(
        name = "com_github_microsoft_go_winio",
        importpath = "github.com/Microsoft/go-winio",
        sum = "h1:a9IhgEQBCUEk6QCdml9CiJGhAws+YwffDHEMp1VMrpA=",
        version = "v0.5.2",
    )
    go_repository(
        name = "com_github_mitchellh_go_homedir",
        importpath = "github.com/mitchellh/go-homedir",
        sum = "h1:lukF9ziXFxDFPkA1vsr5zpc1XuPDn/wFntq5mG+4E0Y=",
        version = "v1.1.0",
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
        sum = "h1:cjW1zVyyoiM0T7b6UoySUFqzXMoqRckQtXwGPiBhOM8=",
        version = "v0.2.0",
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
        sum = "h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_morikuni_aec",
        importpath = "github.com/morikuni/aec",
        sum = "h1:nP9CBfwrvYnBRgY6qfDQkygYDmYwOilePFkwzv4dU8A=",
        version = "v1.0.0",
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
        name = "com_github_niemeyer_pretty",
        importpath = "github.com/niemeyer/pretty",
        sum = "h1:fD57ERR4JtEqsWbfPhv4DMiApHyliiK5xCTNVSPiaAs=",
        version = "v0.0.0-20200227124842-a10e7caefd8e",
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
        name = "com_github_oneofone_xxhash",
        importpath = "github.com/OneOfOne/xxhash",
        sum = "h1:KMrpdQIwFcEqXDklaen+P1axHaj9BSKzvpUUfnHldSE=",
        version = "v1.2.2",
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
        sum = "h1:9t9b9vRUbFq3C4qKFCGkVuq/fIHji802N1nrtkh1mNc=",
        version = "v2.6.0",
    )
    go_repository(
        name = "com_github_onsi_gomega",
        importpath = "github.com/onsi/gomega",
        sum = "h1:KORJXNNTzJXzu4ScJWssJfJMnJ+2QJqhoQSRwNlze9E=",
        version = "v1.24.1",
    )
    go_repository(
        name = "com_github_open_traffic_generator_ixia_c_operator",
        importpath = "github.com/open-traffic-generator/ixia-c-operator",
        sum = "h1:DFINZYbM9MLNnKvi2+BgQgrwXZvXFG3Y7jRm+Qv/4mQ=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_open_traffic_generator_snappi_gosnappi",
        build_file_proto_mode = "disable",
        importpath = "github.com/open-traffic-generator/snappi/gosnappi",
        sum = "h1:s73SOj9BzQLaX8jRXYfz2tGvLZfMDpNqhteTcmYXowU=",
        version = "v0.10.11",
    )
    go_repository(
        name = "com_github_openconfig_gnmi",
        build_directives = ["gazelle:proto_import_prefix github.com/openconfig/gnmi"],
        build_file_generation = "on",
        importpath = "github.com/openconfig/gnmi",
        sum = "h1:hVOdLTaRjdy68oCGJbkf2vrmnUoQ5xbINqBOAMix4xM=",
        version = "v0.9.1",
    )
    go_repository(
        name = "com_github_openconfig_gnoi",
        importpath = "github.com/openconfig/gnoi",
        patch_args = ["-p1"],
        patches = ["//patches:gnoi.patch"],
        sum = "h1:nyVpgoCaVqng9lHRFlqbdvL+z/h8dUeKmoqmdHX5WPk=",
        version = "v0.0.0-20230331232029-cc419f3696d3",
    )
    go_repository(
        name = "com_github_openconfig_gnsi",
        importpath = "github.com/openconfig/gnsi",
        patch_args = ["-p1"],
        patches = ["//patches:gnsi.patch"],
        sum = "h1:oes3klKmvdxvORgGKhJMOD7Oc3bkwPnGzjAxyBezdI8=",
        version = "v0.2.1-0.20230406234913-6acb1b496c20",
    )
    go_repository(
        name = "com_github_openconfig_gocloser",
        importpath = "github.com/openconfig/gocloser",
        sum = "h1:NSYuxdlOWLldNpid1dThR6Dci96juXioUguMho6aliI=",
        version = "v0.0.0-20220310182203-c6c950ed3b0b",
    )
    go_repository(
        name = "com_github_openconfig_goyang",
        importpath = "github.com/openconfig/goyang",
        sum = "h1:mChUZvp1kCWq6Q00wVCtOToddFzEsGlMGG+V+wNXva8=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_openconfig_gribi",
        build_file_proto_mode = "disable",
        importpath = "github.com/openconfig/gribi",
        sum = "h1:xMwEg0mBD+21mOxuFOw0d9dBKuIPwJEhMUUeUulZdLg=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_openconfig_gribigo",
        importpath = "github.com/openconfig/gribigo",
        sum = "h1:N+aVXDlRwVLmI55woEs7wxzKwh7JafGc6xGtlucuBl4=",
        version = "v0.0.0-20230207233343-ef59db57c4fc",
    )
    go_repository(
        name = "com_github_openconfig_grpctunnel",
        importpath = "github.com/openconfig/grpctunnel",
        sum = "h1:t6SvvdfWCMlw0XPlsdxO8EgO+q/fXnTevDjdYREKFwU=",
        version = "v0.0.0-20220819142823-6f5422b8ca70",
    )
    go_repository(
        name = "com_github_openconfig_kne",
        importpath = "github.com/openconfig/kne",
        sum = "h1:mXYY96iWZ5D4C4boNhuhnNAFd/F3BzfHKVJGCVrCPAw=",
        version = "v0.1.8-0.20230309181728-8fc46a2e57e2",
    )
    go_repository(
        name = "com_github_openconfig_lemming_operator",
        importpath = "github.com/openconfig/lemming/operator",
        sum = "h1:Ux+NDwc/QLmYENlESHld/47o1HYcLlxWIAQbDK6UC70=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_openconfig_ondatra",
        build_directives = [
            "gazelle:resolve go github.com/openconfig/gnsi/pathz @com_github_openconfig_gnsi//pathz:pathz",
            "gazelle:resolve go github.com/openconfig/gnoi/bgp @com_github_openconfig_gnoi//bgp:bgp",
            "gazelle:resolve go github.com/openconfig/gnoi/interface @com_github_openconfig_gnoi//interface:interface",
            "gazelle:resolve go github.com/openconfig/gnoi/healthz @com_github_openconfig_gnoi//healthz:healthz",
            "gazelle:resolve go github.com/openconfig/gnoi/cert @com_github_openconfig_gnoi//cert:cert",
            "gazelle:resolve go github.com/openconfig/gnoi/diag @com_github_openconfig_gnoi//diag:diag",
            "gazelle:resolve go github.com/openconfig/gnoi/file @com_github_openconfig_gnoi//file:file",
            "gazelle:resolve go github.com/openconfig/gnoi/types @com_github_openconfig_gnoi//types:types",
            "gazelle:resolve go github.com/openconfig/gnoi/packet_link_qualification @com_github_openconfig_gnoi//packet_link_qualification:linkqual",
            "gazelle:resolve go github.com/openconfig/gnoi/factory_reset @com_github_openconfig_gnoi//factory_reset:factory_reset",
            "gazelle:resolve go github.com/openconfig/gnoi/layer2 @com_github_openconfig_gnoi//layer2:layer2",
            "gazelle:resolve go github.com/openconfig/gnoi/mpls @com_github_openconfig_gnoi//mpls:mpls",
            "gazelle:resolve go github.com/openconfig/gnoi/system @com_github_openconfig_gnoi//system:system",
            "gazelle:resolve go github.com/openconfig/gnoi/os @com_github_openconfig_gnoi//os:os",
            "gazelle:resolve go github.com/openconfig/gnoi/otdr @com_github_openconfig_gnoi//otdr:otdr",
            "gazelle:resolve go github.com/openconfig/gnoi/wavelength_router @com_github_openconfig_gnoi//wavelength_router:wavelength_router",
            "gazelle:resolve go github.com/p4lang/p4runtime/go/p4/v1 @com_github_p4lang_p4runtime//:p4runtime_go_proto",
        ],
        importpath = "github.com/openconfig/ondatra",
        sum = "h1:TVUz99aoOl7Fe5u7bHlxTbV2cEHJ8iGnI+njxBGU5Bw=",
        version = "v0.1.7",
    )
    go_repository(
        name = "com_github_openconfig_testt",
        importpath = "github.com/openconfig/testt",
        sum = "h1:X631iD/B0ximGFb5P9LY5wHju4SiedxUhc5UZEo7VSw=",
        version = "v0.0.0-20220311054427-efbb1a32ec07",
    )
    go_repository(
        name = "com_github_openconfig_ygnmi",
        importpath = "github.com/openconfig/ygnmi",
        sum = "h1:37EJcxZq6LgszCBHd4otNZkFoYzt8s4ENi9w5Ho3d9k=",
        version = "v0.7.9",
    )
    go_repository(
        name = "com_github_openconfig_ygot",
        build_file_proto_mode = "disable",
        importpath = "github.com/openconfig/ygot",
        sum = "h1:B2f6BMp4v0Lex5qzcsqODv+HdQWqcc748pQ4eDF+ty4=",
        version = "v0.27.0",
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
        sum = "h1:9yCKha/T5XdGtO0q9Q9a6T5NUCsTn/DrBg0D7ufOcFM=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_osrg_gobgp_v3",
        importpath = "github.com/osrg/gobgp/v3",
        sum = "h1:zfXhM8+QSfAPUkzWYTbH8pMkvxkcUfvQEKnBU7ikk4g=",
        version = "v3.14.0",
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
        sum = "h1:YtFkrqsMEj7YqpIhRteVxJxCeC3jJBieuLr0d4C4rSA=",
        version = "v0.0.0-20190409184431-ee0cd42419d3",
    )
    go_repository(
        name = "com_github_pborman_uuid",
        importpath = "github.com/pborman/uuid",
        sum = "h1:+ZZIw58t/ozdjRaXh/3awHfmWRbzYxJoAdNJxe/3pvw=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_pelletier_go_toml",
        importpath = "github.com/pelletier/go-toml",
        sum = "h1:4yBQzkHv+7BHq2PQUZF3Mx0IYxG7LsP222s7Agd3ve8=",
        version = "v1.9.5",
    )
    go_repository(
        name = "com_github_pelletier_go_toml_v2",
        importpath = "github.com/pelletier/go-toml/v2",
        sum = "h1:muncTPStnKRos5dpVKULv2FVd4bMOhNePj9CjgDb8Us=",
        version = "v2.0.7",
    )
    go_repository(
        name = "com_github_peterbourgon_diskv",
        importpath = "github.com/peterbourgon/diskv",
        sum = "h1:UBdAOUP5p4RWqPBg048CAvpKN+vxiaj6gdUUzhl4XmI=",
        version = "v2.0.1+incompatible",
    )
    go_repository(
        name = "com_github_pkg_errors",
        importpath = "github.com/pkg/errors",
        sum = "h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
        version = "v0.9.1",
    )
    go_repository(
        name = "com_github_pkg_sftp",
        importpath = "github.com/pkg/sftp",
        sum = "h1:I2qBYMChEhIjOgazfJmV3/mZM256btk6wkCDRmW7JYs=",
        version = "v1.13.1",
    )
    go_repository(
        name = "com_github_pmezard_go_difflib",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_prometheus_client_golang",
        importpath = "github.com/prometheus/client_golang",
        sum = "h1:nJdhIvne2eSX/XRAFV9PcvFFRbrjbcTUj0VP62TMhnw=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_github_prometheus_client_model",
        build_file_proto_mode = "disable",
        importpath = "github.com/prometheus/client_model",
        sum = "h1:UBgGFHqYdG/TPFD1B1ogZywDqEkwp3fBMvqdiQ7Xew4=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_prometheus_common",
        importpath = "github.com/prometheus/common",
        sum = "h1:VTQitp6mXTdUoCmDMugDVOJ1opi6ADftKfp/yeqTR/E=",
        version = "v0.38.0",
    )
    go_repository(
        name = "com_github_prometheus_procfs",
        importpath = "github.com/prometheus/procfs",
        sum = "h1:ODq8ZFEaYeCaZOJlZZdJA2AbQR98dSHSM1KW/You5mo=",
        version = "v0.8.0",
    )
    go_repository(
        name = "com_github_protocolbuffers_txtpbfmt",
        importpath = "github.com/protocolbuffers/txtpbfmt",
        sum = "h1:AKJY61V2SQtJ2a2PdeswKk0NM1qF77X+julRNYRxPOk=",
        version = "v0.0.0-20220608084003-fc78c767cd6a",
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
        name = "com_github_remyoudompheng_bigfft",
        importpath = "github.com/remyoudompheng/bigfft",
        sum = "h1:OdAsTTz6OkFY5QxjkYwrChwuRruF69c169dPK26NUlk=",
        version = "v0.0.0-20200410134404-eec4a21b6bb0",
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
        sum = "h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_github_russross_blackfriday_v2",
        importpath = "github.com/russross/blackfriday/v2",
        sum = "h1:JIOH55/0cWyOuilr9/qlrm0BSXldqnqwMsf35Ld67mk=",
        version = "v2.1.0",
    )
    go_repository(
        name = "com_github_sagikazarmark_crypt",
        importpath = "github.com/sagikazarmark/crypt",
        sum = "h1:fipzMFW34hFUEc4D7fsLQFtE7yElkpgyS2zruedRdZk=",
        version = "v0.9.0",
    )
    go_repository(
        name = "com_github_scrapli_scrapligo",
        importpath = "github.com/scrapli/scrapligo",
        sum = "h1:V4RmedNOa5gtET8/wee7i+DVQjuF0I0DJDp22d0BI90=",
        version = "v1.1.6",
    )
    go_repository(
        name = "com_github_scrapli_scrapligocfg",
        importpath = "github.com/scrapli/scrapligocfg",
        sum = "h1:540SuGqqM6rKN87SLCfR54IageQ6s3a/ZOycGRgbbak=",
        version = "v1.0.0",
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
        sum = "h1:trlNQbNUG3OdDrDil03MCb1H2o9nJ1x4/5LYw7byDE0=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_github_soheilhy_cmux",
        importpath = "github.com/soheilhy/cmux",
        sum = "h1:jjzc5WVemNEDTLwv9tlmemhC73tI08BNOIGwBOo10Js=",
        version = "v0.1.5",
    )
    go_repository(
        name = "com_github_spaolacci_murmur3",
        importpath = "github.com/spaolacci/murmur3",
        sum = "h1:qLC7fQah7D6K1B0ujays3HV9gkFtllcxhzImRR7ArPQ=",
        version = "v0.0.0-20180118202830-f09979ecbc72",
    )
    go_repository(
        name = "com_github_spf13_afero",
        importpath = "github.com/spf13/afero",
        sum = "h1:stMpOSZFs//0Lv29HduCmli3GUfpFoF3Y1Q/aXj/wVM=",
        version = "v1.9.5",
    )
    go_repository(
        name = "com_github_spf13_cast",
        importpath = "github.com/spf13/cast",
        sum = "h1:rj3WzYc11XZaIZMPKmwP96zkFEnnAmV8s6XbB2aY32w=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_spf13_cobra",
        importpath = "github.com/spf13/cobra",
        sum = "h1:o94oiPyS4KD1mPy2fmcYYHHfCxLqYjJOhGsCHFZtEzA=",
        version = "v1.6.1",
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
        sum = "h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=",
        version = "v1.0.5",
    )
    go_repository(
        name = "com_github_spf13_viper",
        importpath = "github.com/spf13/viper",
        sum = "h1:js3yy885G8xwJa6iOISGFwd+qlUo5AvyXb7CiihdtiU=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_github_srl_labs_srl_controller",
        importpath = "github.com/srl-labs/srl-controller",
        sum = "h1:Jy4iqo/t6skTg3O2mhcz2quCw6eWYyQT2nd7ey5mobE=",
        version = "v0.4.6",
    )
    go_repository(
        name = "com_github_srl_labs_srlinux_scrapli",
        importpath = "github.com/srl-labs/srlinux-scrapli",
        sum = "h1:2vsHljvAv7MHmkcB3yW1Qk+i8kYE12OYiSbb+k28qEY=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_stoewer_go_strcase",
        importpath = "github.com/stoewer/go-strcase",
        sum = "h1:Z2iHWqGXH00XYgqDmNgQbIBxf3wrNq0F3feEy0ainaU=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_stretchr_objx",
        importpath = "github.com/stretchr/objx",
        sum = "h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        sum = "h1:+h33VjcLVPDHtOdpUCuF+7gSuG3yGIftsP1YvFihtJ8=",
        version = "v1.8.2",
    )
    go_repository(
        name = "com_github_subosito_gotenv",
        importpath = "github.com/subosito/gotenv",
        sum = "h1:X1TuBLAMDFbaTAChgCBLu3DU3UPyELpnF2jjJ2cz/S8=",
        version = "v1.4.2",
    )
    go_repository(
        name = "com_github_tmc_grpc_websocket_proxy",
        importpath = "github.com/tmc/grpc-websocket-proxy",
        sum = "h1:uruHq4dN7GR16kFc5fp3d1RIYzJW5onx8Ybykw2YQFA=",
        version = "v0.0.0-20201229170055-e5319fda7802",
    )
    go_repository(
        name = "com_github_vishvananda_netlink",
        importpath = "github.com/vishvananda/netlink",
        sum = "h1:Llsql0lnQEbHj0I1OuKyp8otXp0r3q0mPkuhwHfStVs=",
        version = "v1.2.1-beta.2",
    )
    go_repository(
        name = "com_github_vishvananda_netns",
        importpath = "github.com/vishvananda/netns",
        sum = "h1:Oeaw1EM2JMxD51g9uhtC0D7erkIjgmj8+JZc26m1YX8=",
        version = "v0.0.4",
    )
    go_repository(
        name = "com_github_wenovus_gobgp_v3",
        build_file_proto_mode = "disable",
        importpath = "github.com/wenovus/gobgp/v3",
        sum = "h1:O3EW2hZ0IQDlT4VKEmXujLL/9YVslMIz3vG/wMQW9TI=",
        version = "v3.0.0-20230512225508-639b46f0d89c",
    )
    go_repository(
        name = "com_github_xeipuuv_gojsonpointer",
        importpath = "github.com/xeipuuv/gojsonpointer",
        sum = "h1:J9EGpcZtP0E/raorCMxlFGSTBrsSlaDGf3jU/qvAE2c=",
        version = "v0.0.0-20180127040702-4e3ac2762d5f",
    )
    go_repository(
        name = "com_github_xeipuuv_gojsonreference",
        importpath = "github.com/xeipuuv/gojsonreference",
        sum = "h1:EzJWgHovont7NscjpAxXsDA8S8BMYve8Y5+7cuRE7R0=",
        version = "v0.0.0-20180127040603-bd5ef7bd5415",
    )
    go_repository(
        name = "com_github_xeipuuv_gojsonschema",
        importpath = "github.com/xeipuuv/gojsonschema",
        sum = "h1:LhYJRs+L4fBtjZUfuSZIKGeVu0QRy8e5Xi7D17UxZ74=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_xiang90_probing",
        importpath = "github.com/xiang90/probing",
        sum = "h1:eY9dn8+vbi4tKz5Qo6v2eYzo7kUS51QINcR5jNpbZS8=",
        version = "v0.0.0-20190116061207-43a291ad63a2",
    )
    go_repository(
        name = "com_github_yuin_goldmark",
        importpath = "github.com/yuin/goldmark",
        sum = "h1:fVcFKWvrslecOb/tg+Cc05dkeYx540o0FuFt3nUVDoE=",
        version = "v1.4.13",
    )
    go_repository(
        name = "com_google_cloud_go",
        importpath = "cloud.google.com/go",
        sum = "h1:Zc8gqp3+a9/Eyph2KDmcGaPtbKRIoqq4YTlL4NMD0Ys=",
        version = "v0.110.0",
    )
    go_repository(
        name = "com_google_cloud_go_accessapproval",
        importpath = "cloud.google.com/go/accessapproval",
        sum = "h1:x0cEHro/JFPd7eS4BlEWNTMecIj2HdXjOVB5BtvwER0=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_accesscontextmanager",
        importpath = "cloud.google.com/go/accesscontextmanager",
        sum = "h1:r7DpDlWkCMtH/w+gu6Yq//EeYgNWSUbR1+n8ZYr4YWk=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_aiplatform",
        importpath = "cloud.google.com/go/aiplatform",
        sum = "h1:8frB0cIswlhVnYnGrMr+JjZaNC7DHZahvoGHpU9n+RY=",
        version = "v1.35.0",
    )
    go_repository(
        name = "com_google_cloud_go_analytics",
        importpath = "cloud.google.com/go/analytics",
        sum = "h1:uN80RHQeT2jGA3uAFDZSBnKdful4bFw0IHJV6t3EkqU=",
        version = "v0.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigateway",
        importpath = "cloud.google.com/go/apigateway",
        sum = "h1:ZI9mVO7x3E9RK/BURm2p1aw9YTBSCQe3klmyP1WxWEg=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeconnect",
        importpath = "cloud.google.com/go/apigeeconnect",
        sum = "h1:sWOmgDyAsi1AZ48XRHcATC0tsi9SkPT7DA/+VCfkaeA=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeregistry",
        importpath = "cloud.google.com/go/apigeeregistry",
        sum = "h1:BwTPDPTBlYIoQGiwtRUsNFRDZ24cT/02Xb3yFH614YQ=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_apikeys",
        importpath = "cloud.google.com/go/apikeys",
        sum = "h1:+77+/BhFuU476/s78kYiWHObxaYBHsC6Us+Gd7W9pJ4=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_appengine",
        importpath = "cloud.google.com/go/appengine",
        sum = "h1:uTDtjzuHpig1lrf8lycxNSKrthiTDgXnadu+WxYEKxQ=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_area120",
        importpath = "cloud.google.com/go/area120",
        sum = "h1:ugckkFh4XkHJMPhTIx0CyvdoBxmOpMe8rNs4Ok8GAag=",
        version = "v0.7.1",
    )
    go_repository(
        name = "com_google_cloud_go_artifactregistry",
        importpath = "cloud.google.com/go/artifactregistry",
        sum = "h1:G9kjfHsDto5AdKK93hkHWHsY9Oe+6Nv66i7o/KgUO8E=",
        version = "v1.11.2",
    )
    go_repository(
        name = "com_google_cloud_go_asset",
        importpath = "cloud.google.com/go/asset",
        sum = "h1:yObuRcVfexhYQuIWbjNt+9PVPikXIRhERXZxga7qAAY=",
        version = "v1.11.1",
    )
    go_repository(
        name = "com_google_cloud_go_assuredworkloads",
        importpath = "cloud.google.com/go/assuredworkloads",
        sum = "h1:VLGnVFta+N4WM+ASHbhc14ZOItOabDLH1MSoDv+Xuag=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_automl",
        importpath = "cloud.google.com/go/automl",
        sum = "h1:50VugllC+U4IGl3tDNcZaWvApHBTrn/TvyHDJ0wM+Uw=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_baremetalsolution",
        importpath = "cloud.google.com/go/baremetalsolution",
        sum = "h1:2AipdYXL0VxMboelTTw8c1UJ7gYu35LZYUbuRv9Q28s=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_batch",
        importpath = "cloud.google.com/go/batch",
        sum = "h1:YbMt0E6BtqeD5FvSv1d56jbVsWEzlGm55lYte+M6Mzs=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_beyondcorp",
        importpath = "cloud.google.com/go/beyondcorp",
        sum = "h1:qwXDVYf4fQ9DrKci8/40X1zaKYxzYK07vSdPeI9mEQw=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_bigquery",
        importpath = "cloud.google.com/go/bigquery",
        sum = "h1:u+fhS1jJOkPO9vdM84M8HO5VznTfVUicBeoXNKD26ho=",
        version = "v1.48.0",
    )
    go_repository(
        name = "com_google_cloud_go_billing",
        importpath = "cloud.google.com/go/billing",
        sum = "h1:k8pngyiI8uAFhVAhH5+iXSa3Me406XW17LYWZ/3Fr84=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_binaryauthorization",
        importpath = "cloud.google.com/go/binaryauthorization",
        sum = "h1:d3pMDBCCNivxt5a4eaV7FwL7cSH0H7RrEnFrTb1QKWs=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_certificatemanager",
        importpath = "cloud.google.com/go/certificatemanager",
        sum = "h1:5C5UWeSt8Jkgp7OWn2rCkLmYurar/vIWIoSQ2+LaTOc=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_channel",
        importpath = "cloud.google.com/go/channel",
        sum = "h1:/ToBJYu+7wATtd3h8T7hpc4+5NfzlJMDRZjPLIm4EZk=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_cloudbuild",
        build_directives = [
            "gazelle:resolve go google.golang.org/genproto/googleapis/longrunning @org_golang_google_genproto//googleapis/longrunning",  # keep
            "gazelle:resolve go google.golang.org/genproto/googleapis/logging/v2 @org_golang_google_genproto//googleapis/logging/v2:logging",  # keep
        ],
        build_file_proto_mode = "disable",
        importpath = "cloud.google.com/go/cloudbuild",
        sum = "h1:osBOHQJqLPqNfHfkRQXz6sCKAIEKRrupA9NaAGiLN4s=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_clouddms",
        importpath = "cloud.google.com/go/clouddms",
        sum = "h1:E7v4TpDGUyEm1C/4KIrpVSOCTm0P6vWdHT0I4mostRA=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_cloudtasks",
        importpath = "cloud.google.com/go/cloudtasks",
        sum = "h1:Cc2/20hMhGLV2pBGk/i6zNY+eTT9IsV3mrK6TKBu3gs=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_compute",
        importpath = "cloud.google.com/go/compute",
        sum = "h1:FEigFqoDbys2cvFkZ9Fjq4gnHBP55anJ0yQyau2f9oY=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_compute_metadata",
        importpath = "cloud.google.com/go/compute/metadata",
        sum = "h1:mg4jlk7mCAj6xXp9UJ4fjI9VUI5rubuGBW5aJ7UnBMY=",
        version = "v0.2.3",
    )
    go_repository(
        name = "com_google_cloud_go_contactcenterinsights",
        importpath = "cloud.google.com/go/contactcenterinsights",
        sum = "h1:jXIpfcH/VYSE1SYcPzO0n1VVb+sAamiLOgCw45JbOQk=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_container",
        importpath = "cloud.google.com/go/container",
        sum = "h1:q8lTpyAsjcJZQCjGI8JJfcOG4ixl998vwe6TAgQROcM=",
        version = "v1.13.1",
    )
    go_repository(
        name = "com_google_cloud_go_containeranalysis",
        importpath = "cloud.google.com/go/containeranalysis",
        sum = "h1:kw0dDRJPIN8L50Nwm8qa5VuGKPrbVup5lM3ULrvuWrg=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_datacatalog",
        importpath = "cloud.google.com/go/datacatalog",
        sum = "h1:3uaYULZRLByPdbuUvacGeqneudztEM4xqKQsBcxbDnY=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataflow",
        importpath = "cloud.google.com/go/dataflow",
        sum = "h1:eYyD9o/8Nm6EttsKZaEGD84xC17bNgSKCu0ZxwqUbpg=",
        version = "v0.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataform",
        importpath = "cloud.google.com/go/dataform",
        sum = "h1:HBegGOzStIXPWo49FaVTzJOD4EPo8BndPFBUfsuoYe0=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_datafusion",
        importpath = "cloud.google.com/go/datafusion",
        sum = "h1:sZjRnS3TWkGsu1LjYPFD/fHeMLZNXDK6PDHi2s2s/bk=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_datalabeling",
        importpath = "cloud.google.com/go/datalabeling",
        sum = "h1:ch4qA2yvddGRUrlfwrNJCr79qLqhS9QBwofPHfFlDIk=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataplex",
        importpath = "cloud.google.com/go/dataplex",
        sum = "h1:uSkmPwbgOWp3IFtCVEM0Xew80dczVyhNXkvAtTapRn8=",
        version = "v1.5.2",
    )
    go_repository(
        name = "com_google_cloud_go_dataproc",
        importpath = "cloud.google.com/go/dataproc",
        sum = "h1:W47qHL3W4BPkAIbk4SWmIERwsWBaNnWm0P2sdx3YgGU=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataqna",
        importpath = "cloud.google.com/go/dataqna",
        sum = "h1:yFzi/YU4YAdjyo7pXkBE2FeHbgz5OQQBVDdbErEHmVQ=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_datastore",
        importpath = "cloud.google.com/go/datastore",
        sum = "h1:4siQRf4zTiAVt/oeH4GureGkApgb2vtPQAtOmhpqQwE=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_datastream",
        importpath = "cloud.google.com/go/datastream",
        sum = "h1:v6j8C4p0TfXA9Wcea3iH7ZUm05Cx4BiPsH4vEkH7A9g=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_deploy",
        importpath = "cloud.google.com/go/deploy",
        sum = "h1:hdXxUdVw+NOrCQeqg9eQPB3hF1mFEchoS3h+K4IAU9s=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_dialogflow",
        importpath = "cloud.google.com/go/dialogflow",
        sum = "h1:TwmxDsdFcQdExfShoLRlTtdPTor8qSxNu9KZ13o+TUQ=",
        version = "v1.31.0",
    )
    go_repository(
        name = "com_google_cloud_go_dlp",
        importpath = "cloud.google.com/go/dlp",
        sum = "h1:1JoJqezlgu6NWCroBxr4rOZnwNFILXr4cB9dMaSKO4A=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_documentai",
        importpath = "cloud.google.com/go/documentai",
        sum = "h1:tHZA9dB2xo3VaCP4JPxs5jHRntJnmg38kZ0UxlT/u90=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_domains",
        importpath = "cloud.google.com/go/domains",
        sum = "h1:2ti/o9tlWL4N+wIuWUNH+LbfgpwxPr8J1sv9RHA4bYQ=",
        version = "v0.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_edgecontainer",
        importpath = "cloud.google.com/go/edgecontainer",
        sum = "h1:i57Q4zg9j8h4UQoKTD7buXbLCvofmmV8+8owwSmM3ew=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_errorreporting",
        importpath = "cloud.google.com/go/errorreporting",
        sum = "h1:kj1XEWMu8P0qlLhm3FwcaFsUvXChV/OraZwA70trRR0=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_essentialcontacts",
        importpath = "cloud.google.com/go/essentialcontacts",
        sum = "h1:gIzEhCoOT7bi+6QZqZIzX1Erj4SswMPIteNvYVlu+pM=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_eventarc",
        importpath = "cloud.google.com/go/eventarc",
        sum = "h1:4cELkxrOYntz1VRNi2deLRkOr+R6u175kF4hUyd/4Ms=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_filestore",
        importpath = "cloud.google.com/go/filestore",
        sum = "h1:M/iQpbNJw+ELfEvFAW2mAhcHOn1HQQzIkzqmA4njTwg=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_firestore",
        importpath = "cloud.google.com/go/firestore",
        sum = "h1:IBlRyxgGySXu5VuW0RgGFlTtLukSnNkpDiEOMkQkmpA=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_functions",
        importpath = "cloud.google.com/go/functions",
        sum = "h1:WC0JiI5ZBTPSgjzFccqZ8TMkhoPRpDClN99KXhHJp6I=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_gaming",
        importpath = "cloud.google.com/go/gaming",
        sum = "h1:7vEhFnZmd931Mo7sZ6pJy7uQPDxF7m7v8xtBheG08tc=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkebackup",
        importpath = "cloud.google.com/go/gkebackup",
        sum = "h1:za3QZvw6ujR0uyqkhomKKKNoXDyqYGPJies3voUK8DA=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkeconnect",
        importpath = "cloud.google.com/go/gkeconnect",
        sum = "h1:gXYKciHS/Lgq0GJ5Kc9SzPA35NGc3yqu6SkjonpEr2Q=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkehub",
        importpath = "cloud.google.com/go/gkehub",
        sum = "h1:C4p1ZboBOexyCgZSCq+QdP+xfta9+puxgHFy8cjbgYI=",
        version = "v0.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkemulticloud",
        importpath = "cloud.google.com/go/gkemulticloud",
        sum = "h1:8I84Q4vl02rJRsFiinBxl7WCozfdLlUVBQuSrqr9Wtk=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_grafeas",
        importpath = "cloud.google.com/go/grafeas",
        sum = "h1:CYjC+xzdPvbV65gi6Dr4YowKcmLo045pm18L0DhdELM=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_google_cloud_go_gsuiteaddons",
        importpath = "cloud.google.com/go/gsuiteaddons",
        sum = "h1:1mvhXqJzV0Vg5Fa95QwckljODJJfDFXV4pn+iL50zzA=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_iam",
        importpath = "cloud.google.com/go/iam",
        sum = "h1:DRtTY29b75ciH6Ov1PHb4/iat2CLCvrOm40Q0a6DFpE=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_iap",
        importpath = "cloud.google.com/go/iap",
        sum = "h1:a6Heb3z12tUHJqXvmYqLhr7cWz3zzl566xtlbavD5Q0=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_ids",
        importpath = "cloud.google.com/go/ids",
        sum = "h1:fodnCDtOXuMmS8LTC2y3h8t24U8F3eKWfhi+3LY6Qf0=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_iot",
        importpath = "cloud.google.com/go/iot",
        sum = "h1:so1XASBu64OWGylrv5xjvsi6U+/CIR2KiRuZt+WLyKk=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_kms",
        importpath = "cloud.google.com/go/kms",
        sum = "h1:b0votJQa/9DSsxgHwN33/tTLA7ZHVzfWhDCrfiXijSo=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_language",
        importpath = "cloud.google.com/go/language",
        sum = "h1:7Ulo2mDk9huBoBi8zCE3ONOoBrL6UXfAI71CLQ9GEIM=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_lifesciences",
        importpath = "cloud.google.com/go/lifesciences",
        sum = "h1:uWrMjWTsGjLZpCTWEAzYvyXj+7fhiZST45u9AgasasI=",
        version = "v0.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_logging",
        importpath = "cloud.google.com/go/logging",
        sum = "h1:CJYxlNNNNAMkHp9em/YEXcfJg+rPDg7YfwoRpMU+t5I=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_longrunning",
        build_file_proto_mode = "disable",
        importpath = "cloud.google.com/go/longrunning",
        sum = "h1:v+yFJOfKC3yZdY6ZUI933pIYdhyhV8S3NpWrXWmg7jM=",
        version = "v0.4.1",
    )
    go_repository(
        name = "com_google_cloud_go_managedidentities",
        importpath = "cloud.google.com/go/managedidentities",
        sum = "h1:ZRQ4k21/jAhrHBVKl/AY7SjgzeJwG1iZa+mJ82P+VNg=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_maps",
        importpath = "cloud.google.com/go/maps",
        sum = "h1:soPzd0NABgCOGZavyZCAKrJ9L1JAwg3To6n5kuMCm98=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_mediatranslation",
        importpath = "cloud.google.com/go/mediatranslation",
        sum = "h1:anPxH+/WWt8Yc3EdoEJhPMBRF7EhIdz426A+tuoA0OU=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_memcache",
        importpath = "cloud.google.com/go/memcache",
        sum = "h1:8/VEmWCpnETCrBwS3z4MhT+tIdKgR1Z4Tr2tvYH32rg=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_metastore",
        importpath = "cloud.google.com/go/metastore",
        sum = "h1:QCFhZVe2289KDBQ7WxaHV2rAmPrmRAdLC6gbjUd3HPo=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_monitoring",
        importpath = "cloud.google.com/go/monitoring",
        sum = "h1:+X79DyOP/Ny23XIqSIb37AvFWSxDN15w/ktklVvPLso=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_networkconnectivity",
        importpath = "cloud.google.com/go/networkconnectivity",
        sum = "h1:DJwVcr97sd9XPc9rei0z1vUI2ExJyXpA11DSi+Yh7h4=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_networkmanagement",
        importpath = "cloud.google.com/go/networkmanagement",
        sum = "h1:8KWEUNGcpSX9WwZXq7FtciuNGPdPdPN/ruDm769yAEM=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_networksecurity",
        importpath = "cloud.google.com/go/networksecurity",
        sum = "h1:sAKgrzvEslukcwezyEIoXocU2vxWR1Zn7xMTp4uLR0E=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_notebooks",
        importpath = "cloud.google.com/go/notebooks",
        sum = "h1:mMI+/ETVBmCZjdiSYYkN6VFgFTR68kh3frJ8zWvg6go=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_optimization",
        importpath = "cloud.google.com/go/optimization",
        sum = "h1:dj8O4VOJRB4CUwZXdmwNViH1OtI0WtWL867/lnYH248=",
        version = "v1.3.1",
    )
    go_repository(
        name = "com_google_cloud_go_orchestration",
        importpath = "cloud.google.com/go/orchestration",
        sum = "h1:Vw+CEXo8M/FZ1rb4EjcLv0gJqqw89b7+g+C/EmniTb8=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_orgpolicy",
        importpath = "cloud.google.com/go/orgpolicy",
        sum = "h1:XDriMWug7sd0kYT1QKofRpRHzjad0bK8Q8uA9q+XrU4=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_osconfig",
        importpath = "cloud.google.com/go/osconfig",
        sum = "h1:PkSQx4OHit5xz2bNyr11KGcaFccL5oqglFPdTboyqwQ=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_oslogin",
        importpath = "cloud.google.com/go/oslogin",
        sum = "h1:whP7vhpmc+ufZa90eVpkfbgzJRK/Xomjz+XCD4aGwWw=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_phishingprotection",
        importpath = "cloud.google.com/go/phishingprotection",
        sum = "h1:l6tDkT7qAEV49MNEJkEJTB6vOO/onbSOcNtAT09HPuA=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_policytroubleshooter",
        importpath = "cloud.google.com/go/policytroubleshooter",
        sum = "h1:/fRzv4eqv9PDCEL7nBgJiA1EZxhdKMQ4/JIfheCdUZI=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_privatecatalog",
        importpath = "cloud.google.com/go/privatecatalog",
        sum = "h1:7d0gcifTV9As6zzBQo34ZsFiRRlENjD3kw0o3uHn+fY=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub",
        importpath = "cloud.google.com/go/pubsub",
        sum = "h1:XzabfdPx/+eNrsVVGLFgeUnQQKPGkMb8klRCeYK52is=",
        version = "v1.28.0",
    )
    go_repository(
        name = "com_google_cloud_go_pubsublite",
        importpath = "cloud.google.com/go/pubsublite",
        sum = "h1:qh04RCSOnQDVHYmzT74ANu8WR9czAXG3Jl3TV4iR5no=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_recaptchaenterprise",
        importpath = "cloud.google.com/go/recaptchaenterprise",
        sum = "h1:u6EznTGzIdsyOsvm+Xkw0aSuKFXQlyjGE9a4exk6iNQ=",
        version = "v1.3.1",
    )
    go_repository(
        name = "com_google_cloud_go_recaptchaenterprise_v2",
        importpath = "cloud.google.com/go/recaptchaenterprise/v2",
        sum = "h1:E9VgcQxj9M3HS945E3Jb53qd14xcpHBaEG1LgQhnxW8=",
        version = "v2.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_recommendationengine",
        importpath = "cloud.google.com/go/recommendationengine",
        sum = "h1:VibRFCwWXrFebEWKHfZAt2kta6pS7Tlimsnms0fjv7k=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_recommender",
        importpath = "cloud.google.com/go/recommender",
        sum = "h1:ZnFRY5R6zOVk2IDS1Jbv5Bw+DExCI5rFumsTnMXiu/A=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_redis",
        importpath = "cloud.google.com/go/redis",
        sum = "h1:JoAd3SkeDt3rLFAAxEvw6wV4t+8y4ZzfZcZmddqphQ8=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_resourcemanager",
        importpath = "cloud.google.com/go/resourcemanager",
        sum = "h1:m2RQU8UzBCIO+wsdwoehpuyAaF1i7ahFhj7TLocxuJE=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_resourcesettings",
        importpath = "cloud.google.com/go/resourcesettings",
        sum = "h1:8Dua37kQt27CCWHm4h/Q1XqCF6ByD7Ouu49xg95qJzI=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_retail",
        importpath = "cloud.google.com/go/retail",
        sum = "h1:1Dda2OpFNzIb4qWgFZjYlpP7sxX3aLeypKG6A3H4Yys=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_run",
        importpath = "cloud.google.com/go/run",
        sum = "h1:monNAz/FXgo8A31aR9sbrsv+bEbqy6H/arSgLOfA2Fk=",
        version = "v0.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_scheduler",
        importpath = "cloud.google.com/go/scheduler",
        sum = "h1:NRzIXqVxpyoiyonpYOKJmVJ9iif/Acw36Jri+cVHZ9U=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_secretmanager",
        importpath = "cloud.google.com/go/secretmanager",
        sum = "h1:pu03bha7ukxF8otyPKTFdDz+rr9sE3YauS5PliDXK60=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_security",
        importpath = "cloud.google.com/go/security",
        sum = "h1:WIyVxhrdex1geaAV0pC/4yXy/sZdurjHXLzMopcjers=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_securitycenter",
        importpath = "cloud.google.com/go/securitycenter",
        sum = "h1:DRUo2MFSq3Kt0a4hWRysdMHcu2obPwnSQNgHfOuwR4Q=",
        version = "v1.18.1",
    )
    go_repository(
        name = "com_google_cloud_go_servicecontrol",
        importpath = "cloud.google.com/go/servicecontrol",
        sum = "h1:iEiMJgD1bzRL9Zu4JYDQUWfqZ+kRLX8wWZSCMBK8Qzs=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_servicedirectory",
        importpath = "cloud.google.com/go/servicedirectory",
        sum = "h1:DPvPdb6O/lg7xK+BFKlzZN+w6upeJ/bbfcUnnqU66b8=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_servicemanagement",
        importpath = "cloud.google.com/go/servicemanagement",
        sum = "h1:flWoX0eJy21+34I/7HPUbpr6xTHPVzws1xnecLFlUm0=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_serviceusage",
        importpath = "cloud.google.com/go/serviceusage",
        sum = "h1:fl1AGgOx7E2eyBmH5ofDXT9w8xGvEaEnHYyNYGkxaqg=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_shell",
        importpath = "cloud.google.com/go/shell",
        sum = "h1:wT0Uw7ib7+AgZST9eCDygwTJn4+bHMDtZo5fh7kGWDU=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_spanner",
        importpath = "cloud.google.com/go/spanner",
        sum = "h1:fba7k2apz4aI0BE59/kbeaJ78dPOXSz2PSuBIfe7SBM=",
        version = "v1.44.0",
    )
    go_repository(
        name = "com_google_cloud_go_speech",
        importpath = "cloud.google.com/go/speech",
        sum = "h1:x4ZJWhop/sLtnIP97IMmPtD6ZF003eD8hykJ0lOgEtw=",
        version = "v1.14.1",
    )
    go_repository(
        name = "com_google_cloud_go_storage",
        importpath = "cloud.google.com/go/storage",
        sum = "h1:F5QDG5ChchaAVQhINh24U99OWHURqrW8OmQcGKXcbgI=",
        version = "v1.28.1",
    )
    go_repository(
        name = "com_google_cloud_go_storagetransfer",
        importpath = "cloud.google.com/go/storagetransfer",
        sum = "h1:doREJk5f36gq7yJDJ2HVGaYTuQ8Nh6JWm+6tPjdfh+g=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_talent",
        importpath = "cloud.google.com/go/talent",
        sum = "h1:nI9sVZPjMKiO2q3Uu0KhTDVov3Xrlpt63fghP9XjyEM=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_texttospeech",
        importpath = "cloud.google.com/go/texttospeech",
        sum = "h1:H4g1ULStsbVtalbZGktyzXzw6jP26RjVGYx9RaYjBzc=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_tpu",
        importpath = "cloud.google.com/go/tpu",
        sum = "h1:/34T6CbSi+kTv5E19Q9zbU/ix8IviInZpzwz3rsFE+A=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_trace",
        importpath = "cloud.google.com/go/trace",
        sum = "h1:GFPLxbp5/FzdgTzor3nlNYNxMd6hLmzkE7sA9F0qQcA=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_translate",
        importpath = "cloud.google.com/go/translate",
        sum = "h1:oBW4KVgcUq4OAXGdKEdyV7lqWiA3keQ3+8FKreAQv4g=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_video",
        importpath = "cloud.google.com/go/video",
        sum = "h1:FL+xG+4vgZASVIxcWACxneKPhFOnOX75GJhhTP7yUkQ=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_videointelligence",
        importpath = "cloud.google.com/go/videointelligence",
        sum = "h1:Uh5BdoET8XXqXX2uXIahGb+wTKbLkGH7s4GXR58RrG8=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_vision",
        importpath = "cloud.google.com/go/vision",
        sum = "h1:/CsSTkbmO9HC8iQpxbK8ATms3OQaX3YQUeTMGCxlaK4=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_google_cloud_go_vision_v2",
        importpath = "cloud.google.com/go/vision/v2",
        sum = "h1:WKt7VNhMLKaT9NmdisWnU2LVO5CaHvisssTaAqfV3dg=",
        version = "v2.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_vmmigration",
        importpath = "cloud.google.com/go/vmmigration",
        sum = "h1:+2zAH2Di1FB02kAv8L9In2chYRP2Mw0bl41MiWwF+Fc=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_vmwareengine",
        importpath = "cloud.google.com/go/vmwareengine",
        sum = "h1:ZM35wN4xuxDZSpKFypLMTsB02M+NEIZ2wr7/VpT3osw=",
        version = "v0.2.2",
    )
    go_repository(
        name = "com_google_cloud_go_vpcaccess",
        importpath = "cloud.google.com/go/vpcaccess",
        sum = "h1:FOe6CuiQD3BhHJWt7E8QlbBcaIzVRddupwJlp7eqmn4=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_webrisk",
        importpath = "cloud.google.com/go/webrisk",
        sum = "h1:IY+L2+UwxcVm2zayMAtBhZleecdIFLiC+QJMzgb0kT0=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_websecurityscanner",
        importpath = "cloud.google.com/go/websecurityscanner",
        sum = "h1:AHC1xmaNMOZtNqxI9Rmm87IJEyPaRkOxeI0gpAacXGk=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_workflows",
        importpath = "cloud.google.com/go/workflows",
        sum = "h1:FfGp9w0cYnaKZJhUOMqCOJCYT/WlvYBfTQhFWV3sRKI=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_lukechampine_uint128",
        importpath = "lukechampine.com/uint128",
        sum = "h1:pnxCASz787iMf+02ssImqk6OLt+Z5QHMoZyUXR4z6JU=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_shuralyov_dmitri_gpu_mtl",
        importpath = "dmitri.shuralyov.com/gpu/mtl",
        sum = "h1:VpgP7xuJadIUuKccphEpTJnWhS2jkQyMt6Y7pJCD7fY=",
        version = "v0.0.0-20190408044501-666a987793e9",
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
        sum = "h1:1Lc07Kr7qY4U2YPouBjpCLxpiyxIVoxqXgkXLknAOE8=",
        version = "v2.0.0",
    )
    go_repository(
        name = "in_gopkg_tomb_v1",
        importpath = "gopkg.in/tomb.v1",
        sum = "h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=",
        version = "v1.0.0-20141024135613-dd632973f1e7",
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
        sum = "h1:/ecaJf0sk1l4l6V4awd65v2C3ILy7MSj+s/x1ADCIMU=",
        version = "v1.3.6",
    )
    go_repository(
        name = "io_etcd_go_etcd_api_v3",
        importpath = "go.etcd.io/etcd/api/v3",
        sum = "h1:Cy2qx3npLcYqTKqGJzMypnMv2tiRyifZJ17BlWIWA7A=",
        version = "v3.5.6",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_pkg_v3",
        importpath = "go.etcd.io/etcd/client/pkg/v3",
        sum = "h1:TXQWYceBKqLp4sa87rcPs11SXxUA/mHwH975v+BDvLU=",
        version = "v3.5.6",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_v2",
        importpath = "go.etcd.io/etcd/client/v2",
        sum = "h1:fIDR0p4KMjw01MJMfUIDWdQbjo06PD6CeYM5z4EHLi0=",
        version = "v2.305.6",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_v3",
        importpath = "go.etcd.io/etcd/client/v3",
        sum = "h1:coLs69PWCXE9G4FKquzNaSHrRyMCAXwF+IX1tAPVO8E=",
        version = "v3.5.6",
    )
    go_repository(
        name = "io_etcd_go_etcd_pkg_v3",
        importpath = "go.etcd.io/etcd/pkg/v3",
        sum = "h1:Ablg7T7OkR+AeeeU32kdVhw/AGDsitkKPl7aW73ssjU=",
        version = "v3.5.5",
    )
    go_repository(
        name = "io_etcd_go_etcd_raft_v3",
        importpath = "go.etcd.io/etcd/raft/v3",
        sum = "h1:Ibz6XyZ60OYyRopu73lLM/P+qco3YtlZMOhnXNS051I=",
        version = "v3.5.5",
    )
    go_repository(
        name = "io_etcd_go_etcd_server_v3",
        importpath = "go.etcd.io/etcd/server/v3",
        sum = "h1:jNjYm/9s+f9A9r6+SC4RvNaz6AqixpOvhrFdT0PvIj0=",
        version = "v3.5.5",
    )
    go_repository(
        name = "io_k8s_api",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/api",
        sum = "h1:dM3cinp3PGB6asOySalOZxEG4CZ0IAdJsrYZXE/ovGQ=",
        version = "v0.26.2",
    )
    go_repository(
        name = "io_k8s_apiextensions_apiserver",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/apiextensions-apiserver",
        sum = "h1:/yTG2B9jGY2Q70iGskMf41qTLhL9XeNN2KhI0uDgwko=",
        version = "v0.26.2",
    )
    go_repository(
        name = "io_k8s_apimachinery",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/apimachinery",
        sum = "h1:da1u3D5wfR5u2RpLhE/ZtZS2P7QvDgLZTi9wrNZl/tQ=",
        version = "v0.26.2",
    )
    go_repository(
        name = "io_k8s_apiserver",
        importpath = "k8s.io/apiserver",
        sum = "h1:Pk8lmX4G14hYqJd1poHGC08G03nIHVqdJMR0SD3IH3o=",
        version = "v0.26.2",
    )
    go_repository(
        name = "io_k8s_client_go",
        importpath = "k8s.io/client-go",
        sum = "h1:s1WkVujHX3kTp4Zn4yGNFK+dlDXy1bAAkIl+cFAiuYI=",
        version = "v0.26.2",
    )
    go_repository(
        name = "io_k8s_code_generator",
        importpath = "k8s.io/code-generator",
        sum = "h1:QMgN5oXUgQe27uMaqpbT0hg6ti+rvgCWaHEDMHVhox8=",
        version = "v0.26.2",
    )
    go_repository(
        name = "io_k8s_component_base",
        importpath = "k8s.io/component-base",
        sum = "h1:IfWgCGUDzrD6wLLgXEstJKYZKAFS2kO+rBRi0p3LqcI=",
        version = "v0.26.2",
    )
    go_repository(
        name = "io_k8s_gengo",
        importpath = "k8s.io/gengo",
        sum = "h1:U9tB195lKdzwqicbJvyJeOXV7Klv+wNAWENRnXEGi08=",
        version = "v0.0.0-20220902162205-c0856e24416d",
    )
    go_repository(
        name = "io_k8s_klog_v2",
        importpath = "k8s.io/klog/v2",
        sum = "h1:m4bYOKall2MmOiRaR1J+We67Do7vm9KiQVlT96lnHUw=",
        version = "v2.90.1",
    )
    go_repository(
        name = "io_k8s_kms",
        importpath = "k8s.io/kms",
        sum = "h1:GM1gg3tFK3OUU/QQFi93yGjG3lJT8s8l3Wkn2+VxBLM=",
        version = "v0.26.2",
    )
    go_repository(
        name = "io_k8s_kube_openapi",
        importpath = "k8s.io/kube-openapi",
        sum = "h1:gmovKNur38vgoWfGtP5QOGNOA7ki4n6qNYoFAgMlNvg=",
        version = "v0.0.0-20230308215209-15aac26d736a",
    )
    go_repository(
        name = "io_k8s_sigs_apiserver_network_proxy_konnectivity_client",
        importpath = "sigs.k8s.io/apiserver-network-proxy/konnectivity-client",
        sum = "h1:+xBL5uTc+BkPBwmMi3vYfUJjq+N3K+H6PXeETwf5cPI=",
        version = "v0.0.35",
    )
    go_repository(
        name = "io_k8s_sigs_controller_runtime",
        importpath = "sigs.k8s.io/controller-runtime",
        sum = "h1:6xaWFqzT5KuAQ9ufgUaj1G/+C4Y1GRkhrxl+BJ9i+5s=",
        version = "v0.14.5",
    )
    go_repository(
        name = "io_k8s_sigs_json",
        importpath = "sigs.k8s.io/json",
        sum = "h1:EDPBXCAspyGV4jQlpZSudPeMmr1bNJefnuqLsRAsHZo=",
        version = "v0.0.0-20221116044647-bc3834ca7abd",
    )
    go_repository(
        name = "io_k8s_sigs_structured_merge_diff_v4",
        importpath = "sigs.k8s.io/structured-merge-diff/v4",
        sum = "h1:PRbqxJClWWYMNV1dhaG4NsibJbArud9kFxnAMREiWFE=",
        version = "v4.2.3",
    )
    go_repository(
        name = "io_k8s_sigs_yaml",
        importpath = "sigs.k8s.io/yaml",
        sum = "h1:a2VclLzOGrwOHDiV8EfBGhvjHvP46CtW5j6POvhYGGo=",
        version = "v1.3.0",
    )
    go_repository(
        name = "io_k8s_utils",
        importpath = "k8s.io/utils",
        sum = "h1:xMMXJlJbsU8w3V5N2FLDQ8YgU8s1EoULdbQBcAeNJkY=",
        version = "v0.0.0-20230313181309-38a27ef9d749",
    )
    go_repository(
        name = "io_opencensus_go",
        importpath = "go.opencensus.io",
        sum = "h1:y73uSU6J157QMP2kn2r30vwW1A2W2WFwSCGnAVxeaD0=",
        version = "v0.24.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc",
        importpath = "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc",
        sum = "h1:xFSRQBbXF6VvYRf2lqMJXxoB72XI1K/azav8TekHHSw=",
        version = "v0.35.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp",
        importpath = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
        sum = "h1:Ajldaqhxqw/gNzQA45IKFWLdG7jZuXX/wBW1d5qvbUI=",
        version = "v0.35.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel",
        importpath = "go.opentelemetry.io/otel",
        sum = "h1:Y7DTJMR6zs1xkS/upamJYk0SxxN4C9AqRd77jmZnyY4=",
        version = "v1.10.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_internal_retry",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/internal/retry",
        sum = "h1:TaB+1rQhddO1sF71MpZOZAuSPW1klK2M8XxfrBMfK7Y=",
        version = "v1.10.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace",
        sum = "h1:pDDYmo0QadUPal5fwXoY1pmMpFcdyhXOmL5drCrI3vU=",
        version = "v1.10.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace_otlptracegrpc",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc",
        sum = "h1:KtiUEhQmj/Pa874bVYKGNVdq8NPKiacPbaRRtgXi+t4=",
        version = "v1.10.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:6SiklT+gfWAwWUR0meEMxQBtihpiEs4c+vL9spDTqUs=",
        version = "v0.31.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk",
        importpath = "go.opentelemetry.io/otel/sdk",
        sum = "h1:jZ6K7sVn04kk/3DNUdJ4mqRlGDiXAVuIG+MMENpTNdY=",
        version = "v1.10.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:npQMbR8o7mum8uF95yFbOEJffhs1sbCOfDh8zAJiH5E=",
        version = "v1.10.0",
    )
    go_repository(
        name = "io_opentelemetry_go_proto_otlp",
        importpath = "go.opentelemetry.io/proto/otlp",
        sum = "h1:IVN6GR+mhC4s5yfcTbmzHYODqvWAp3ZedA2SJPI1Nnw=",
        version = "v0.19.0",
    )
    go_repository(
        name = "io_rsc_binaryregexp",
        importpath = "rsc.io/binaryregexp",
        sum = "h1:HfqmD5MEmC0zvwBuF187nq9mdnXjXsSivRiXN7SmRkE=",
        version = "v0.2.0",
    )
    go_repository(
        name = "io_rsc_quote_v3",
        importpath = "rsc.io/quote/v3",
        sum = "h1:9JKUTTIUgS6kzR9mK1YuGKv6Nl+DijDNIc0ghT58FaY=",
        version = "v3.1.0",
    )
    go_repository(
        name = "io_rsc_sampler",
        importpath = "rsc.io/sampler",
        sum = "h1:7uVkIFmeBqHfdjD+gZwtXXI+RODJ2Wc4O7MPEh/QiW4=",
        version = "v1.3.0",
    )
    go_repository(
        name = "org_golang_google_api",
        importpath = "google.golang.org/api",
        sum = "h1:l+rh0KYUooe9JGbGVx71tbFo4SMbMTXK3I3ia2QSEeU=",
        version = "v0.110.0",
    )
    go_repository(
        name = "org_golang_google_appengine",
        importpath = "google.golang.org/appengine",
        sum = "h1:FZR1q0exgwxzPzp/aF+VccGrSfxfPpkBqjIIEq3ru6c=",
        version = "v1.6.7",
    )
    go_repository(
        name = "org_golang_google_genproto",
        build_file_proto_mode = "disable",
        importpath = "google.golang.org/genproto",
        sum = "h1:khxVcsk/FhnzxMKOyD+TDGwjbEOpcPuIpmafPGFmhMA=",
        version = "v0.0.0-20230320184635-7606e756e683",
    )
    go_repository(
        name = "org_golang_google_grpc",
        build_file_proto_mode = "disable",
        importpath = "google.golang.org/grpc",
        sum = "h1:LAv2ds7cmFV/XTS3XG1NneeENYrXGmorPxsBbptIjNc=",
        version = "v1.53.0",
    )
    go_repository(
        name = "org_golang_google_grpc_cmd_protoc_gen_go_grpc",
        importpath = "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
        sum = "h1:M1YKkFIboKNieVO5DLUEVzQfGwJD30Nv2jfUgzb5UcE=",
        version = "v1.1.0",
    )
    go_repository(
        name = "org_golang_google_protobuf",
        importpath = "google.golang.org/protobuf",
        sum = "h1:kPPoIgf3TsEvrm0PFe15JQ+570QVxYzEvvHqChK+cng=",
        version = "v1.30.0",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:U/0M97KRkSFvyD/3FSmdP5W5swImpNgle/EHFhOsQPE=",
        version = "v0.5.0",
    )
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:Fq7F/w7MAa1KJ5bt2aJ62ihqp9HDcRuyILskkpIAurw=",
        version = "v0.0.0-20221106115401-f9659909a136",
    )
    go_repository(
        name = "org_golang_x_image",
        importpath = "golang.org/x/image",
        sum = "h1:+qEpEAPhDZ1o0x3tHzZTQDArnOixOzGD9HUJfcg0mb4=",
        version = "v0.0.0-20190802002840-cff245a6509b",
    )
    go_repository(
        name = "org_golang_x_lint",
        importpath = "golang.org/x/lint",
        sum = "h1:VLliZ0d+/avPrXXH+OakdXhpJuEoBZuwh1m2j7U6Iug=",
        version = "v0.0.0-20210508222113-6edffad5e616",
    )
    go_repository(
        name = "org_golang_x_mobile",
        importpath = "golang.org/x/mobile",
        sum = "h1:4+4C/Iv2U4fMZBiMCc98MG1In4gJY5YRhtpDNeDeHWs=",
        version = "v0.0.0-20190719004257-d2bd2a29d028",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:LUYupSeNrTNCGzR/hVBk2NHZO4hXcVaW1k4Qx7rjPx8=",
        version = "v0.8.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:Zrh2ngAOFYneWTAIAPethzeaQLuHwhuBkuV6ZiRnUaQ=",
        version = "v0.8.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:Lh8GPgSKBfWSwFvtuWOfeI3aAAnbXTSutYxJiOJFgIw=",
        version = "v0.6.0",
    )
    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:wsuoTGHzEhffawBOhz5CYhcrV4IdKZbEyZjBMuTp12o=",
        version = "v0.1.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:MVltZSvRTcU2ljQOhs94SXPftV6DCNnZViHeQps87pQ=",
        version = "v0.6.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:clScbb1cHjoCkyRbWwBEUZ5H/tIFu5TAXIqaZD0Gcjw=",
        version = "v0.6.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:57P1ETyNKtuIjB4SRd15iJxuhj8Gc416Y78H3qgMh68=",
        version = "v0.8.0",
    )
    go_repository(
        name = "org_golang_x_time",
        importpath = "golang.org/x/time",
        sum = "h1:rg5rLMjNzMS1RkNLzCG38eapWhnYLFYXDXj2gOlr8j4=",
        version = "v0.3.0",
    )
    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:BOw41kyTf3PuCW1pVQf8+Cyg8pMlkYB1oo9iJ6D/lKM=",
        version = "v0.6.0",
    )
    go_repository(
        name = "org_golang_x_xerrors",
        importpath = "golang.org/x/xerrors",
        sum = "h1:H2TDz8ibqkAF6YGhCdN3jS9O0/s90v0rJh3X/OLHEUk=",
        version = "v0.0.0-20220907171357-04be3eba64a2",
    )
    go_repository(
        name = "org_modernc_cc_v4",
        importpath = "modernc.org/cc/v4",
        sum = "h1:3avQiH2LBr2L2SYwjKHdVIG+buRnxxoR7ZFoifW/OaQ=",
        version = "v4.1.3",
    )
    go_repository(
        name = "org_modernc_ccorpus2",
        importpath = "modernc.org/ccorpus2",
        sum = "h1:I+dskokQYHl+wL+Y2q9dS9B1o7JbrxBTyYvESaqmBgQ=",
        version = "v1.2.0",
    )
    go_repository(
        name = "org_modernc_mathutil",
        importpath = "modernc.org/mathutil",
        sum = "h1:rV0Ko/6SfM+8G+yKiyI830l3Wuz1zRutdslNoQ0kfiQ=",
        version = "v1.5.0",
    )
    go_repository(
        name = "org_modernc_opt",
        importpath = "modernc.org/opt",
        sum = "h1:3XOZf2yznlhC+ibLltsDGzABUGVx8J6pnFMS3E4dcq4=",
        version = "v0.1.3",
    )
    go_repository(
        name = "org_modernc_strutil",
        importpath = "modernc.org/strutil",
        sum = "h1:fNMm+oJklMGYfU9Ylcywl0CO5O6nTfaowNsh2wpPjzY=",
        version = "v1.1.3",
    )
    go_repository(
        name = "org_modernc_token",
        importpath = "modernc.org/token",
        sum = "h1:Xl7Ap9dKaEs5kLoOQeQmPWevfnk/DM5qcLcYlA8ys6Y=",
        version = "v1.1.0",
    )

    go_repository(
        name = "org_uber_go_atomic",
        importpath = "go.uber.org/atomic",
        sum = "h1:9qC72Qh0+3MqyJbAn8YU5xVq1frD8bn3JtD2oXtafVQ=",
        version = "v1.10.0",
    )
    go_repository(
        name = "org_uber_go_goleak",
        importpath = "go.uber.org/goleak",
        sum = "h1:xqgm/S+aQvhWFTtR0XK3Jvg7z8kGV8P4X14IzwN3Eqk=",
        version = "v1.2.0",
    )
    go_repository(
        name = "org_uber_go_multierr",
        importpath = "go.uber.org/multierr",
        sum = "h1:S0h4aNzvfcFsC3dRF1jLoaov7oRaKqRGC/pUEJ2yvPQ=",
        version = "v1.10.0",
    )
    go_repository(
        name = "org_uber_go_zap",
        importpath = "go.uber.org/zap",
        sum = "h1:FiJd5l1UOLj0wCgbSE0rwwXHzEdAZS6hiiSnxJN/D60=",
        version = "v1.24.0",
    )
    go_repository(
        name = "tf_universe_go_metallb",
        importpath = "go.universe.tf/metallb",
        sum = "h1:SYaiJDn5tNRUeM7VFwA6zdb8KvYaTf7pGFwdq2lhuGs=",
        version = "v0.13.5",
    )
    go_repository(
        name = "xyz_gomodules_jsonpatch_v2",
        importpath = "gomodules.xyz/jsonpatch/v2",
        sum = "h1:4pT439QV83L+G9FkcCriY6EkpcK6r6bK+A5FBUMI7qY=",
        version = "v2.2.0",
    )
