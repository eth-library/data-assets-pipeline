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
    description = "Developer tools for the Data Archive Pipeline (DAP) Orchestrator";
    homepage = "https://github.com/eth-library/dap";
    license = licenses.mit;
    mainProgram = "dap";
  };
}
