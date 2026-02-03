{
  description = "Development environment for da_pipeline using Nix, uv, and .venv";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";

  outputs = { self, nixpkgs }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = [
            pkgs.python312
            pkgs.uv
            pkgs.kubectl
            pkgs.kubernetes-helm
            pkgs.just
            pkgs.jq
            pkgs.curl
            pkgs.openssl
          ];
        };
      });
    };
}
