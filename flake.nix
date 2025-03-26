{
  description = "A Todoist CLI client";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-parts.url = "github:hercules-ci/flake-parts";

    git-hooks-nix = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    treefmt = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      imports = [
        inputs.git-hooks-nix.flakeModule
        inputs.treefmt.flakeModule
      ];

      perSystem =
        { config, pkgs, ... }:
        {
          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              go
              config.treefmt.build.wrapper
            ];
            CGO_ENABLED = "0";
            shellHook = config.pre-commit.installationScript;
          };

          pre-commit.settings.hooks = {
            commitizen.enable = true;
            treefmt.enable = true;
          };

          treefmt.programs = {
            gofumpt.enable = true;
            nixfmt.enable = true;
            prettier.enable = true;
          };
        };
    };
}
