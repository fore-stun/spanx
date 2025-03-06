{ lib
, dex
, gomod2nix
, go-mod-graph-chart
, hostPlatform
, mkGoEnv
, mkShell
, xcaddy
, xdg-utils
, zsh
}:

let
  name = "go";

  goEnv = mkGoEnv { pwd = ./.; };

  BROWSER =
    if hostPlatform.isLinux then ''
      ${lib.getExe dex} -d "/run/current-system/sw/share/applications/$(
        ${xdg-utils}/bin/xdg-settings get default-web-browser
      )" \
        | sed -E -n -e 's/^Executing command: //p'
    ''
    else "echo open";

  customPackages = [
    xcaddy
  ];
in
mkShell {
  inherit name;

  packages = [
    goEnv
    gomod2nix
    go-mod-graph-chart
  ] ++ customPackages;

  shellHook = ''
    export NIX_SHELL_NAME="${name}"
    export BROWSER="$(${BROWSER})"
    RPROMPT='%F{magenta}${name}%f %1(j.«%j» .)%*'
    ${zsh}/bin/zsh
    exit 0
  '';
}
