{
  description = "Development environment for the Data Archive Pipeline (DAP) Orchestrator";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    dap-cli.url = "path:./cli";
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

      devShells = forAllSystems (pkgs:
        let
          # Package groups
          basePackages = [ pkgs.python312 pkgs.uv ];
          cliPackage = [ dap-cli.packages.${pkgs.stdenv.hostPlatform.system}.dap ];
          k8sPackages = [ pkgs.kubectl pkgs.kubernetes-helm ];
          goPackages = [ pkgs.go dap-cli.gomod2nix.packages.${pkgs.stdenv.hostPlatform.system}.default ];

          # Helper to create shells with common settings
          mkDevShell = packages: pkgs.mkShell {
            inherit packages;
            # Required for pip-installed packages with C++ extensions (e.g., grpcio)
            LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath [ pkgs.stdenv.cc.cc.lib ];
          };

        in {
          # Default: Everything (full development environment)
          default = mkDevShell (basePackages ++ cliPackage ++ k8sPackages ++ goPackages);

          # Minimal: Python + dap CLI (for running pipeline and tests)
          minimal = mkDevShell (basePackages ++ cliPackage);

          # K8s: Python + dap CLI + kubectl/helm (for deployment)
          k8s = mkDevShell (basePackages ++ cliPackage ++ k8sPackages);

          # CLI development: Python + Go (for working on the dap CLI)
          cli-dev = mkDevShell (basePackages ++ goPackages);
        }
      );
    };
}
