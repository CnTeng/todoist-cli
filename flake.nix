{
  description = "A CLI client for Todoist.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-parts.url = "github:hercules-ci/flake-parts";

    git-hooks-nix = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ self, flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      imports = [ inputs.git-hooks-nix.flakeModule ];

      perSystem =
        { config, pkgs, ... }:
        {
          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              go
              gotools
            ];
            CGO_ENABLED = "0";
            shellHook = config.pre-commit.installationScript;
          };

          formatter = pkgs.nixfmt-tree.override {
            settings.formatter.gofumpt = {
              command = "gofumpt";
              excludes = [ "vendor/*" ];
              includes = [ "*.go" ];
              options = [ "-w" ];
            };
            runtimeInputs = [ pkgs.gofumpt ];
          };

          pre-commit.settings.hooks = {
            treefmt = {
              enable = true;
              package = config.formatter;
            };
            commitizen.enable = true;
          };
        };
    };
}
