{
  description = "A CLI client for Todoist.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    git-hooks-nix = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ self, nixpkgs, ... }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];
      forAllPkgs = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system} system);
    in
    {
      homeModules.default = import ./nix/hm-module.nix self;

      packages = forAllPkgs (
        pkgs: _: {
          todoist-cli = pkgs.callPackage ./nix/package.nix { };
        }
      );

      checks = forAllPkgs (
        pkgs: system: {
          todoist-cli = import ./nix/check.nix self pkgs.nixosTest;

          pre-commit-check = inputs.git-hooks-nix.lib.${system}.run {
            src = ./.;
            hooks = {
              commitizen.enable = true;
              treefmt = {
                enable = true;
                package = self.formatter.${system};
              };
            };
          };
        }
      );

      devShells = forAllPkgs (
        pkgs: system: {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
              gotools
            ];
            CGO_ENABLED = "0";

            inherit (self.checks.${system}.pre-commit-check) shellHook;
          };
        }
      );

      formatter = forAllPkgs (
        pkgs: _:
        pkgs.nixfmt-tree.override {
          runtimeInputs = with pkgs; [ gofumpt ];
          settings.formatter.gofumpt = {
            command = "gofumpt";
            excludes = [ "vendor/*" ];
            includes = [ "*.go" ];
            options = [ "-w" ];
          };
        }
      );
    };
}
