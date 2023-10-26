{ lib
, buildGoApplication
}:

buildGoApplication {
  pname = "spanx";
  version = "0.1";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
  meta = {
    license = lib.licenses.mit;
    mainProgram = "caddy";
  };
}
