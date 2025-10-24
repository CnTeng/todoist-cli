{
  description = "A CLI client for Todoist.";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

  outputs =
    { nixpkgs, ... }:
    let
      forEachSystem =
        f:
        nixpkgs.lib.genAttrs [
          "x86_64-linux"
          "aarch64-linux"
        ] (system: f nixpkgs.legacyPackages.${system});
    in
    {
      devShells = forEachSystem (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gotools
          ];
          CGO_ENABLED = "0";
        };
      });

      formatter = forEachSystem (
        pkgs:
        pkgs.nixfmt-tree.override {
          settings.formatter.gofumpt = {
            command = "gofumpt";
            options = [ "-w" ];
            includes = [ "*.go" ];
            excludes = [ "vendor/*" ];
          };
          runtimeInputs = [ pkgs.gofumpt ];
        }
      );
    };
}
