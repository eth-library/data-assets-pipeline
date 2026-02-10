{ buildGoApplication, lib }:

buildGoApplication {
  pname = "dap";
  version = "0.1.0";
  src = lib.cleanSource ./.;
  modules = ./gomod2nix.toml;

  postInstall = ''
    mv $out/bin/cli $out/bin/dap
  '';

  meta = with lib; {
    description = "Development tooling for the DAP Orchestrator";
    homepage = "https://github.com/eth-library/data-assets-pipeline";
    license = licenses.asl20;
    mainProgram = "dap";
  };
}
