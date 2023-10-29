{ lib
, buildGoApplication
}:

let
  src = lib.cleanSourceWith {
    filter = name: type:
      ! (type == "directory" && name == ".github")
    ;
    src = lib.cleanSource ./.;
  };
in

buildGoApplication {
  pname = "caddy-extended";
  version = "v2.7.5";
  inherit src;
  modules = ./gomod2nix.toml;
  subPackages = [ "cmd/caddy" ];
  meta = {
    license = lib.licenses.mit;
    mainProgram = "caddy";
  };
}
