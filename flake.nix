{
  description = "Development environment for the Data Archive Pipeline (DAP) Orchestrator";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    dap-cli.url = "path:./.cli";
  };

  outputs = { self, nixpkgs, dap-cli }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      packages = forAllSystems (pkgs: {
        dap = dap-cli.packages.${pkgs.stdenv.hostPlatform.system}.dap;
        default = self.packages.${pkgs.stdenv.hostPlatform.system}.dap;
      });

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = [
            dap-cli.packages.${pkgs.stdenv.hostPlatform.system}.dap
            pkgs.python312
            pkgs.uv
            pkgs.go
            pkgs.kubectl
            pkgs.kubernetes-helm
            pkgs.jq
            pkgs.curl
            pkgs.openssl
            dap-cli.gomod2nix.packages.${pkgs.stdenv.hostPlatform.system}.default
          ];

          # Required for pip-installed packages with C++ extensions (e.g., grpcio)
          LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath [
            pkgs.stdenv.cc.cc.lib
          ];
        };
      });
    };
}
