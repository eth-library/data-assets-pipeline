{ buildGoApplication, lib }:

buildGoApplication {
  pname = "dap";
  version = "0.1.0";
  src = lib.cleanSource ./.;
  modules = ./gomod2nix.toml;

  ldflags = [
    "-s" "-w"
    "-X github.com/eth-library/dap/cli/cmd.version=0.1.0"
  ];

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
