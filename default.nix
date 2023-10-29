{ lib
, buildGoApplication
}:

buildGoApplication {
  pname = "caddy-extended";
  version = "v2.7.5";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
  subPackages = [ "cmd/caddy" ];
  meta = {
    license = lib.licenses.mit;
    mainProgram = "caddy";
  };
}
