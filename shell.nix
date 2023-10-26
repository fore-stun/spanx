{ mkGoEnv
, gomod2nix
, mkShell
}:

let
  goEnv = mkGoEnv { pwd = ./.; };
in
mkShell {
  packages = [
    goEnv
    gomod2nix
  ];
}
