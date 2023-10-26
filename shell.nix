{ lib
, gomod2nix
, mkGoEnv
, mkShell
, zsh
}:

let
  name = "go";

  goEnv = mkGoEnv { pwd = ./.; };
in
mkShell {
  inherit name;

  packages = [
    goEnv
    gomod2nix
  ];

  shellHook = ''
    export NIX_SHELL_NAME="${name}"
    RPROMPT='%F{magenta}${name}%f %1(j.«%j» .)%*'
    ${zsh}/bin/zsh
    exit 0
  '';
}
