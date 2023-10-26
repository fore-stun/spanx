{ lib
, buildGoApplication
}:

buildGoApplication {
  pname = "spanx";
  version = "0.1";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
  subPackages = [ "cmd/caddy" ];
  meta = {
    license = lib.licenses.mit;
    mainProgram = "caddy";
  };
}
